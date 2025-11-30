import type { ApiClient as IApiClient } from './client'

export const DEFAULT_API_ENDPOINT = import.meta.env.VITE_API_ENDPOINT

export class ApiClient implements IApiClient {
  #endpoint: string

  constructor(endpoint: string) {
    this.#endpoint = endpoint
  }
}
