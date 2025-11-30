export class ApiError extends Error {
  status: number

  constructor(message: string, status: number) {
    super(message)
    this.status = status
  }
}

export type ApiClient = {
  subscribe(
    topic: string,
    id: string,
    subscription: PushSubscriptionJSON
  ): Promise<void>

  unsubscribe(topic: string, id: string): Promise<void>

  subscriptionExists(topic: string, id: string): Promise<boolean>
}
