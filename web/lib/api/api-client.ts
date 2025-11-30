import { ApiError, type ApiClient as IApiClient } from './client'

export const DEFAULT_API_ENDPOINT = import.meta.env.VITE_API_ENDPOINT

export class ApiClient implements IApiClient {
  #endpoint: string

  constructor(endpoint: string) {
    this.#endpoint = endpoint
  }

  async subscribe(
    topic: string,
    id: string,
    subscription: PushSubscriptionJSON
  ): Promise<void> {
    const res = await fetch(
      `${this.#endpoint}/subscriptions/${encodeURIComponent(topic)}/${encodeURIComponent(id)}`,
      {
        method: 'post',
        headers: {
          'content-type': 'application/json',
        },
        body: JSON.stringify(subscription),
      }
    )

    if (res.status !== 201) {
      throw new ApiError('unexpected status code', res.status)
    }
  }

  async unsubscribe(topic: string, id: string): Promise<void> {
    const res = await fetch(
      `${this.#endpoint}/subscriptions/${encodeURIComponent(topic)}/${encodeURIComponent(id)}`,
      {
        method: 'delete',
      }
    )

    if (res.status !== 204) {
      throw new ApiError('unexpected status code', res.status)
    }
  }

  async subscriptionExists(topic: string, id: string): Promise<boolean> {
    const res = await fetch(
      `${this.#endpoint}/subscriptions/${encodeURIComponent(topic)}/${encodeURIComponent(id)}`,
      {
        method: 'head',
      }
    )

    if (res.status !== 200 && res.status !== 404) {
      throw new ApiError('unexpected status code', res.status)
    }

    return res.status === 200
  }
}
