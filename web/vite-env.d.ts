/// <reference types="vite/client" />

declare interface ViteTypeOptions {
  strictImportMetaEnv: unknown
}

interface ImportMetaEnv {
  readonly VITE_API_ENDPOINT: string
}

declare interface ImportMeta {
  readonly env: ImportMetaEnv
}

declare interface Window {
  grapevine: {
    applicationServerKey: string
    topic: string
  }
  pushManager: PushManager
}

declare interface Navigator {
  /** Available on Safari. */
  standalone?: boolean
}

declare interface Uint8Array {
  toBase64(options?: { alphabet: 'base64url'; omitPadding: true }): string
}
