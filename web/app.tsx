import type { JSX } from 'react'
import { useIsStandalone } from './lib/pwa'
import { StandalonePage } from './pages/standalone'
import { WebPage } from './pages/web'

export function App(): JSX.Element {
  const isStandalone = useIsStandalone()

  return isStandalone ? <StandalonePage /> : <WebPage />
}
