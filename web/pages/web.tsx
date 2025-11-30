import { type JSX, useState } from 'react'
import { DeviconGrafana } from '../components/icons/DeviconGrafana'
import { OcticonChevronRight24 } from '../components/icons/OcticonChevronRight24'
import { OcticonInbox16 } from '../components/icons/OcticonInbox16'
import { SfSymbolsMenubarDockRectangle24 } from '../components/icons/SfSymbolsMenubarDockRectangle24'
import { SfSymbolsPlusApp24 } from '../components/icons/SfSymbolsPlusApp24'
import { SfSymbolsShare24 } from '../components/icons/SfSymbolsShare'
import { useLocationPathPattern } from '../lib/routing'

export function WebPage(): JSX.Element {
  const pathPatternMatch = useLocationPathPattern('/topics/:topic', 'topic')
  const topic = pathPatternMatch?.topic

  return (
    <div className="flex justify-center px-2 py-10">
      <div className="flex flex-col gap-y-2 flex-grow max-w-[600px]">
        <h1>Grapevine</h1>
        <h2>Get started</h2>
        {topic ? (
          <>
            <div className="card items-center gap-y-2">
              <p>Add Grapevine to the Dock or home screen</p>
              <p className="text-foreground-1-alt text-center">
                Grapevine works by adding a web app to your device. Once added,
                open the app to complete the installation. You can add Grapevine
                multiple times, once for each available topic.
                <ul className="flex flex-col gap-y-2 mt-2">
                  <li className="flex flex-shrink justify-center items-center gap-x-2">
                    iOS: <SfSymbolsShare24 />
                    <SfSymbolsPlusApp24 /> Add to home screen
                  </li>
                  <li className="flex flex-shrink justify-center items-center gap-x-2">
                    macOS: <SfSymbolsShare24 />
                    <SfSymbolsMenubarDockRectangle24 />
                    Add to Dock
                  </li>
                </ul>
              </p>
            </div>
          </>
        ) : (
          <>
            <div className="card items-center gap-y-2">
              <p>Selecting a topic</p>
              <p className="text-foreground-1-alt text-center">
                Grapevine groups notifications into topics. You need to select a
                topic before continuing.
              </p>
            </div>
            <ul className="card">
              <a href="/topics/default" className="hover:bg-surface-1-hover">
                <li className="flex gap-x-2 items-center">
                  <div className="flex items-center justify-center w-[30px] h-[30px] rounded bg-[#40d663]">
                    <OcticonInbox16 />
                  </div>
                  <p className="flex-grow">Default</p>
                  <OcticonChevronRight24 className="text-foreground-1-alt" />
                </li>
              </a>
              <a href="/topics/grafana" className="hover:bg-surface-1-hover">
                <li className="flex gap-x-2 items-center">
                  <div className="flex items-center justify-center w-[30px] h-[30px] rounded bg-surface-2">
                    <DeviconGrafana />
                  </div>
                  <p className="flex-grow">Grafana</p>
                  <OcticonChevronRight24 className="text-foreground-1-alt" />
                </li>
              </a>
            </ul>
          </>
        )}
      </div>
    </div>
  )
}
