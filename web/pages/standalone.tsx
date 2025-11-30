import { type JSX, useState } from 'react'
import { Notification } from '../components/Notification'
import { useLocationPathPattern } from '../lib/routing'

export function StandalonePage(): JSX.Element {
  const pathPatternMatch = useLocationPathPattern('/topics/:topic', 'topic')
  const topic = pathPatternMatch?.topic

  const [subscribed, setSubscribed] = useState(false)
  const [notifications, setNotifications] = useState([])

  if (!topic) {
    return (
      <div className="flex justify-center px-2 py-10">
        <div className="flex flex-col gap-y-2 flex-grow max-w-[600px]">
          <h1>Grapevine</h1>
          <div className="card">
            <p className="text-center">Invalid topic</p>
            <p className="text-center text-foreground-1-alt">
              You've added the web app to your device without specifying a
              topic. Please delete the app, visit your Grapevine instance and
              follow the instructions.
            </p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="flex justify-center px-2 py-10">
      <div className="flex flex-col gap-y-2 flex-grow max-w-[600px]">
        <h1>Grapevine</h1>

        {subscribed ? (
          <>
            <h2>Recent</h2>
            <div className="card">
              <p className="text-center">
                You're all set up and ready to receive notifications!
              </p>
              <p className="text-center text-foreground-1-alt">
                Your recent notifications will show up here. Try to send a
                notification to the '{topic}' topic.
              </p>
            </div>
            {notifications.length > 0 && (
              <>
                <ul className="card">
                  <a href="/" className="hover:bg-surface-1-hover">
                    <li>
                      <Notification />
                    </li>
                  </a>
                  <a href="/" className="hover:bg-surface-1-hover">
                    <li>
                      <Notification />
                    </li>
                  </a>
                </ul>
                <button type="button" className="big">
                  View all
                </button>
              </>
            )}
            <h2>Settings</h2>
            <div className="card items-center gap-y-2">
              <button
                type="button"
                className="big w-full danger"
                onClick={() => setSubscribed(false)}
              >
                Unsubscribe
              </button>
            </div>
          </>
        ) : (
          <>
            <h2>Get started</h2>
            <div className="card items-center gap-y-2">
              <p>One last step</p>
              <p className="text-foreground-1-alt text-center">
                Click the subscribe button to allow Grapevine to send you push
                notifications.
              </p>
              <button
                type="button"
                className="big w-full primary"
                onClick={() => setSubscribed(true)}
              >
                Subscribe
              </button>
            </div>
          </>
        )}
      </div>
    </div>
  )
}
