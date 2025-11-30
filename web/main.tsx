import React from 'react'
import { createRoot } from 'react-dom/client'

import './main.css'
import { App } from './app'
import { ApiProvider } from './lib/api/ApiProvider'
import { ApiClient, DEFAULT_API_ENDPOINT } from './lib/api/api-client'

const apiClient = new ApiClient(DEFAULT_API_ENDPOINT)

const root = document.getElementById('root')
if (root) {
  createRoot(root).render(
    <React.StrictMode>
      <ApiProvider client={apiClient}>
        <App />
      </ApiProvider>
    </React.StrictMode>
  )
}
