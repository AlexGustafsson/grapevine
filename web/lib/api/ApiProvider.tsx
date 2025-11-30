import {
  createContext,
  type JSX,
  type PropsWithChildren,
  use,
  useCallback,
  useEffect,
  useState,
} from 'react'
import { ApiError, type ApiClient } from './client'

const ApiContext = createContext<ApiClient>({} as ApiClient)

export function ApiProvider({
  children,
  client,
}: PropsWithChildren<{ client: ApiClient }>): JSX.Element {
  return <ApiContext value={client}>{children}</ApiContext>
}

export function useApiClient(): ApiClient {
  const client = use(ApiContext)
  return client
}

async function deriveSubscriptionId(
  subscription: PushSubscription
): Promise<string> {
  return await crypto.subtle
    .digest('SHA-256', new TextEncoder().encode(subscription.endpoint))
    .then((x) =>
      new Uint8Array(x).toBase64({
        alphabet: 'base64url',
        omitPadding: true,
      })
    )
}

export function useSubscription(): [
  string | undefined,
  () => Promise<void>,
  () => Promise<void>,
] {
  const client = useApiClient()

  const [subscription, setSubscription] = useState<PushSubscription>()
  const [subscriptionId, setSubscriptionId] = useState<string>()
  const [serverHasSubscription, setServerHasSubscription] = useState(false)

  // Get initial state of the local subscription
  useEffect(() => {
    window.pushManager
      .getSubscription()
      .then((subscription) => {
        if (subscription) {
          setSubscription(subscription)
          deriveSubscriptionId(subscription)
            .then((subscriptionId) => {
              setSubscriptionId(subscriptionId)

              client
                .subscriptionExists(window.grapevine.topic, subscriptionId)
                .then(setServerHasSubscription)
                .catch((error) => {
                  console.error('Failed to check if subscription exists', error)
                })
            })
            .catch((error) => {
              console.error('Failed to derive subscription id', error)
            })
        }
      })
      .catch((error) => {
        console.error('Failed to identify existing subscription', error)
      })
  }, [client])

  // TODO: Error handling
  const subscribe = useCallback(async () => {
    const subscription = await window.pushManager.subscribe({
      // MUST be true for declerative web push
      userVisibleOnly: true,
      applicationServerKey: window.grapevine.applicationServerKey,
    })
    setSubscription(subscription)

    const subscriptionId = await deriveSubscriptionId(subscription)
    setSubscriptionId(subscriptionId)

    await client.subscribe(
      window.grapevine.topic,
      subscriptionId,
      subscription.toJSON()
    )
    setServerHasSubscription(true)
  }, [client])

  // TODO: Error handling
  const unsubscribe = useCallback(async () => {
    if (serverHasSubscription && subscriptionId) {
      try {
        await client.unsubscribe(window.grapevine.topic, subscriptionId)
        setServerHasSubscription(false)
      } catch (error) {
        if (error instanceof ApiError && error.status === 404) {
          // Assume the subscription is already removed
          setServerHasSubscription(false)
        } else {
          console.error('Failed to unsubscribe', error)
          return
        }
      }
    }

    if (subscription) {
      const ok = await subscription.unsubscribe()
      if (!ok) {
        throw new Error('not ok')
      }

      setSubscription(undefined)
      setSubscriptionId(undefined)
    }
  }, [client, serverHasSubscription, subscriptionId, subscription])

  return [subscriptionId, subscribe, unsubscribe]
}
