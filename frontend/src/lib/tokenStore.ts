let accessToken: string | null = null
let refreshFn: (() => Promise<string | null>) | null = null

export function setAccessToken(token: string | null) {
  accessToken = token
}

export function getAccessToken(): string | null {
  return accessToken
}

export function setRefreshFunction(fn: () => Promise<string | null>) {
  refreshFn = fn
}

export async function attemptRefresh(): Promise<string | null> {
  if (!refreshFn) return null
  return refreshFn()
}
