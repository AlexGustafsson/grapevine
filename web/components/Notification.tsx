import { JSX } from 'react'

export function Notification(): JSX.Element {
  return (
    <div className="cursor-pointer flex gap-x-2">
      <div className="flex justify-center w-8 gap-y-1">
        <p className="text-[#8D909D]">x</p>
      </div>
      <div className="flex flex-col flex-grow gap-y-1">
        <p className="text-[#8D909D]">AlexGustafsson / cupdate #451</p>
        <p>Server Bug: Insufficient scope on Caddy</p>
        <p className="text-[#8D909D]">I don't think this has to do with C...</p>
      </div>
      <div className="flex justify-center w-8 gap-y-1">
        <p className="text-[#8D909D]">20h</p>
      </div>
    </div>
  )
}
