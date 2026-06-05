import { createContext, useContext, useState, useEffect, useCallback, useRef, type ReactNode } from 'react'
import * as authService from '@/services/authService'
import { setAccessToken, setRefreshFunction } from '@/lib/tokenStore'

interface User {
  email: string
  permissions: string[]
}

interface AuthContextType {
  user: User | null
  isLoading: boolean
  isAuthDialogOpen: boolean
  openAuthDialog: () => void
  closeAuthDialog: () => void
  login: (email: string, password: string) => Promise<void>
  register: (email: string, password: string) => Promise<void>
  loginAsGuest: () => Promise<{ email: string; password: string }>
  logout: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | null>(null)

let refreshPromise: Promise<string | null> | null = null

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [isAuthDialogOpen, setIsAuthDialogOpen] = useState(false)
  const autoOpenTimerRef = useRef<ReturnType<typeof setTimeout> | undefined>(undefined)

  const openAuthDialog = useCallback(() => setIsAuthDialogOpen(true), [])
  const closeAuthDialog = useCallback(() => setIsAuthDialogOpen(false), [])

  const refreshAccessToken = useCallback(async (): Promise<string | null> => {
    if (refreshPromise) {
      return refreshPromise
    }

    refreshPromise = (async () => {
      try {
        const response = await authService.refresh()
        setAccessToken(response.accessToken)
        setUser({
          email: response.email,
          permissions: response.permissions,
        })
        return response.accessToken
      } catch {
        setAccessToken(null)
        setUser(null)
        return null
      } finally {
        refreshPromise = null
      }
    })()

    return refreshPromise
  }, [])

  useEffect(() => {
    setRefreshFunction(refreshAccessToken)

    const initAuth = async () => {
      await refreshAccessToken()
      setIsLoading(false)
    }
    initAuth()
  }, [refreshAccessToken])

  useEffect(() => {
    if (!isLoading && !user) {
      autoOpenTimerRef.current = setTimeout(() => openAuthDialog(), 150)
    }
    return () => {
      if (autoOpenTimerRef.current) {
        clearTimeout(autoOpenTimerRef.current)
      }
    }
  }, [isLoading, user, openAuthDialog])

  const login = async (email: string, password: string) => {
    const response = await authService.login(email, password)
    setAccessToken(response.accessToken)
    setUser({
      email: response.email,
      permissions: response.permissions,
    })
  }

  const register = async (email: string, password: string) => {
    await authService.register(email, password)
    await login(email, password)
  }

  const loginAsGuest = async (): Promise<{ email: string; password: string }> => {
    const response = await authService.guestLogin()
    setAccessToken(response.accessToken)
    setUser({
      email: response.email,
      permissions: response.permissions,
    })
    return { email: response.email, password: response.password }
  }

  const logout = async () => {
    try {
      await authService.logout()
    } finally {
      setAccessToken(null)
      setUser(null)
    }
  }

  return (
    <AuthContext.Provider value={{ user, isLoading, isAuthDialogOpen, openAuthDialog, closeAuthDialog, login, register, loginAsGuest, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
