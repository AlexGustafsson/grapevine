import {
  createContext,
  type JSX,
  type PropsWithChildren,
  useCallback,
  useContext,
  useEffect,
  useState,
} from 'react'

interface WebPushContextType {
  applicationServerKey: string
}

const WebPushContext = createContext<WebPushContextType>(
  {} as WebPushContextType
)

export function WebPushProvider({
  children,
  applicationServerKey,
}: PropsWithChildren<{ applicationServerKey: string }>): JSX.Element {
  return (
    <WebPushContext value={{ applicationServerKey }}>{children}</WebPushContext>
  )
}

export function useApplicationServerKey() {
  const context = useContext(WebPushContext)
  return context.applicationServerKey
}

export function useWebPushSubscription(): [
  PushSubscription | undefined,
  () => Promise<void>,
  () => Promise<void>,
] {
  const [subscription, setSubscription] = useState<PushSubscription>()

  const applicationServerKey = useApplicationServerKey()

  useEffect(() => {
    window.pushManager
      .getSubscription()
      .then((subscription) => {
        setSubscription(subscription || undefined)
      })
      .catch((error) => {
        console.error('Failed to get subscription', error)
      })
  }, [])

  const subscribe = useCallback(async () => {
    const subscription = await window.pushManager.subscribe({
      // MUST be true for declerative web push
      userVisibleOnly: true,
      applicationServerKey,
    })

    setSubscription(subscription)
  }, [applicationServerKey])

  const unsubscribe = useCallback(async () => {
    if (!subscription) {
      return
    }

    const ok = await subscription.unsubscribe()
    if (ok) {
      setSubscription(undefined)
    } else {
      throw new Error('failed to unsubscribe')
    }
  }, [subscription])

  return [subscription, subscribe, unsubscribe]
}
