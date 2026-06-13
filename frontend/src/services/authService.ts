export interface TokenResponse {
  accessToken: string
  tokenType: string
  expiresIn: number
  email: string
  permissions: string[]
}

export interface User {
  email: string
  permissions: string[]
}

const AUTH_BASE = '/api/v1/auth'

export async function login(email: string, password: string): Promise<TokenResponse> {
  const response = await fetch(`${AUTH_BASE}/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ email, password }),
  })

  if (!response.ok) {
    if (response.status === 401) {
      throw new Error('Invalid email or password')
    }
    if (response.status === 400) {
      const body = await response.json()
      throw new Error(body.message || 'Please check your input')
    }
    throw new Error('Something went wrong, please try again later')
  }

  return response.json()
}

export async function refresh(): Promise<TokenResponse> {
  const response = await fetch(`${AUTH_BASE}/refresh`, {
    method: 'POST',
    credentials: 'include',
  })

  if (!response.ok) {
    throw new Error('Session expired')
  }

  return response.json()
}

export async function register(email: string, password: string): Promise<void> {
  const response = await fetch(`${AUTH_BASE}/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  })

  if (!response.ok) {
    if (response.status === 409) {
      throw new Error('An account with this email already exists')
    }
    throw new Error('Something went wrong, please try again later')
  }
}

export interface GuestLoginResponse extends TokenResponse {
  password: string
}

export async function guestLogin(): Promise<GuestLoginResponse> {
  const response = await fetch(`${AUTH_BASE}/register?type=guest`, {
    method: 'POST',
    credentials: 'include',
  })

  if (!response.ok) {
    throw new Error('Failed to create guest account')
  }

  return response.json()
}

export async function logout(): Promise<void> {
  await fetch(`${AUTH_BASE}/logout`, {
    method: 'POST',
    credentials: 'include',
  })
}
