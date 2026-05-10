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
    const response = await fetch(`${API_BASE}${path}`, {
      ...options,
      signal: options.signal || controller.signal,
    })

    if (!response.ok) {
      let message: string
      try {
        const errorBody = await response.json()
        message = errorBody.error || errorBody.message || `Request failed with status ${response.status}`
      } catch {
        message = `Request failed with status ${response.status}`
      }
      throw new Error(message)
    }

    if (response.status === 204) return undefined as T

    return response.json() as Promise<T>
  } finally {
    clearTimeout(timeoutId)
  }
}
