import {
  createContext,
  type JSX,
  type PropsWithChildren,
  use,
  useCallback,
  useEffect,
  useRef,
  useState,
} from 'react'
import { useWebPushSubscription } from '../WebPushProvider'
import type { ApiClient } from './client'

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

export function useSubscription() {
  // const client = useApiClient()
  // const [subscription, setSubscription] = useState<PushSubscription>()
  // const subscriptionId = useState<string>()
  // // Get initial state of the local subscription
  // useEffect(() => {
  //   setSubscription(some)
  // }, [])
  // // Keep subscription id up-to-date (for end-user verification?)
  // useEffect(() => {
  // }, [subscription])
  // // Upsert the subscription
  // useEffect(() => {
  //   if (subscription) {
  //     deriveSubscriptionId(subscription).then((subscriptionId) => client.createSubscription(window.grapevine.topic, subscriptionId, subscription.toJSON()))
  //   }
  // }, [subscription])
  // useEffect(() => {
  //   if (subscription) {
  //     const id = await crypto.subtle
  //               .digest(
  //                 "SHA-256",
  //                 new TextEncoder().encode(subscription.endpoint)
  //               )
  //               .then((x) =>
  //                 new Uint8Array(x).toBase64({
  //                   alphabet: "base64url",
  //                   omitPadding: true,
  //                 })
  //               );
  //     client
  //       .subscribe(window.grapevine.topic, id, subscription.toJSON())
  //       .then(() => {})
  //       .catch((error) => {
  //         console.log('Failed to subscribe', error)
  //       })
  //   } else {
  //   }
  // }, [subscription])
  // const subscribe = useCallback(async () => {
  //   setStatus('subscribing')
  //   try {
  //     const subscription = await window.pushManager.subscribe({
  //       // MUST be true for declerative web push
  //       userVisibleOnly: true,
  //       applicationServerKey: window.grapevine.applicationServerKey,
  //     })
  //     subscriptionRef.current = subscription
  //     await client.subscribe(window.grapevine.topic, subscription.toJSON())
  //     setStatus('subscribed')
  //   } catch (error) {
  //     console.error('Failed to subscribe', error)
  //     if (subscriptionRef.current) {
  //       try {
  //         const ok = await subscriptionRef.current.unsubscribe()
  //         if (ok) {
  //           setStatus('unsubscribed')
  //         }
  //       } catch (error) {
  //         console.error('Failed to unsubscribe', error)
  //       }
  //     } else {
  //       setStatus('unsubscribed')
  //     }
  //   }
  //   window.pushManager
  //     .subscribe({
  //       // MUST be true for declerative web push
  //       userVisibleOnly: true,
  //       applicationServerKey: window.grapevine.applicationServerKey,
  //     })
  //     .then((subscription) => {
  //       subscriptionRef.current = subscription
  //       return client.subscribe(window.grapevine.topic, subscription.toJSON())
  //     })
  //     .then(() => {
  //       setStatus('subscribed')
  //     })
  //     .catch((error) => {
  //       console.error('Failed to subscribe', error)
  //       if (subscriptionRef.current) {
  //         subscriptionRef.current
  //           .unsubscribe()
  //           .then((ok) => {
  //             if (!ok) {
  //               throw new Error('not ok')
  //             }
  //             subscriptionRef.current = null
  //           })
  //           .catch((error) => {
  //             console.error('Failed to unsubscribe', error)
  //           })
  //           .finally(() => {
  //             setStatus('unsubscribed')
  //           })
  //       } else {
  //         setStatus('unsubscribed')
  //       }
  //     })
  // }, [client])
  // const unsubscribe = useCallback(async () => {})
  // useEffect(() => {
  //   if (subscription) {
  //     oldSubscriptionRef.current = subscription
  //   }
  // }, [subscription])
  // useEffect(() => {
  //   if (subscription) {
  //   } else {
  //   }
  // }, [client, subscription])
}
