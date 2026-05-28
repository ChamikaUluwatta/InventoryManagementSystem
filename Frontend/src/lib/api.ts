import { getAccessToken, attemptRefresh } from './tokenStore'

const API_BASE_URL = import.meta.env.VITE_API_URL

if (!API_BASE_URL && import.meta.env.MODE === 'production') {
  console.warn('VITE_API_URL environment variable is not set')
}

const API_BASE = API_BASE_URL || 'http://localhost:8080/api/v1'

export async function apiFetch<T>(
  path: string,
  options: RequestInit = {},
): Promise<T> {
  const controller = new AbortController()
  const timeoutId = setTimeout(() => controller.abort(), 10000)

  try {
    const headers = new Headers(options.headers)

    const token = getAccessToken()
    if (token) {
      headers.set('Authorization', `Bearer ${token}`)
    }

    const response = await fetch(`${API_BASE}${path}`, {
      ...options,
      headers,
      signal: options.signal || controller.signal,
    })

    if (response.status === 401) {
      const newToken = await attemptRefresh()
      if (newToken) {
        headers.set('Authorization', `Bearer ${newToken}`)
        const retryResponse = await fetch(`${API_BASE}${path}`, {
          ...options,
          headers,
          signal: options.signal || controller.signal,
        })

        if (!retryResponse.ok) {
          throw new Error('Something went wrong')
        }

        if (retryResponse.status === 204) return undefined as T
        return retryResponse.json() as Promise<T>
      }
      throw new Error('Your session has expired')
    }

    if (!response.ok) {
      if (response.status === 400) {
        const errorBody = await response.json()
        throw new Error(errorBody.message || 'Please check your input')
      }
      throw new Error('Something went wrong')
    }

    if (response.status === 204) return undefined as T

    return response.json() as Promise<T>
  } finally {
    clearTimeout(timeoutId)
  }
}
