import { useState, useEffect, useCallback, useRef } from 'react'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import { useAuth } from '@/contexts/AuthContext'
import { apiFetch } from '@/lib/api'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Field, FieldDescription, FieldGroup, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'

const loginSchema = z.object({
  email: z.string().email('Please enter a valid email'),
  password: z.string().min(1, 'Password is required'),
})

const signupSchema = z.object({
  email: z.string().email('Please enter a valid email'),
  password: z.string().min(6, 'Password must be at least 6 characters'),
  confirmPassword: z.string().min(1, 'Please confirm your password'),
}).refine(data => data.password === data.confirmPassword, {
  message: 'Passwords do not match',
  path: ['confirmPassword'],
})

const SEED_STEPS = [
  'Creating companies',
  'Creating categories',
  'Creating locations',
  'Creating products',
  'Creating inventories',
]

type DialogContent = 'login' | 'signup' | 'guest_credentials' | 'seeding'

export function AuthDialog({ open, onClose }: { open: boolean; onClose: () => void }) {
  const [mode, setMode] = useState<'login' | 'signup'>('login')
  const [error, setError] = useState<string | null>(null)
  const [content, setContent] = useState<DialogContent>('login')
  const [isGuestLoading, setIsGuestLoading] = useState(false)
  const [guestEmail, setGuestEmail] = useState('')
  const [guestPassword, setGuestPassword] = useState('')
  const [guestExpiresAt, setGuestExpiresAt] = useState<Date | null>(null)
  const [copiedField, setCopiedField] = useState<'email' | 'password' | null>(null)
  const [seedStep, setSeedStep] = useState(0)
  const [mounted, setMounted] = useState(false)
  const [visible, setVisible] = useState(false)
  const exitTimerRef = useRef<ReturnType<typeof setTimeout> | undefined>(undefined)
  const cancelledRef = useRef(false)
  const { login, register, loginAsGuest } = useAuth()

  const loginForm = useForm({ resolver: zodResolver(loginSchema) })
  const signupForm = useForm({ resolver: zodResolver(signupSchema) })

  const isSubmitting = mode === 'login'
    ? loginForm.formState.isSubmitting
    : signupForm.formState.isSubmitting

  const currentErrors = mode === 'login' ? loginForm.formState.errors : signupForm.formState.errors

  const emailError = currentErrors.email?.message as string | undefined
  const passwordError = currentErrors.password?.message as string | undefined
  const confirmPasswordError = mode === 'signup' ? signupForm.formState.errors.confirmPassword?.message as string | undefined : undefined

  const resetToLogin = useCallback(() => {
    setContent('login')
    setMode('login')
    setGuestEmail('')
    setGuestPassword('')
    setGuestExpiresAt(null)
    setCopiedField(null)
    setSeedStep(0)
    setError(null)
  }, [])

  useEffect(() => {
    if (open) {
      setMounted(true)
      requestAnimationFrame(() => {
        requestAnimationFrame(() => setVisible(true))
      })
    } else {
      setVisible(false)
      exitTimerRef.current = setTimeout(() => {
        setMounted(false)
        resetToLogin()
      }, 200)
    }
    return () => {
      if (exitTimerRef.current) {
        clearTimeout(exitTimerRef.current)
      }
    }
  }, [open, resetToLogin])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)

    try {
      if (mode === 'login') {
        const data = loginForm.getValues()
        const result = loginSchema.safeParse(data)
        if (!result.success) {
          loginForm.setError('email', { message: result.error.issues.find(e => e.path[0] === 'email')?.message })
          loginForm.setError('password', { message: result.error.issues.find(e => e.path[0] === 'password')?.message })
          return
        }
        await login(result.data.email, result.data.password)
      } else {
        const data = signupForm.getValues()
        const result = signupSchema.safeParse(data)
        if (!result.success) {
          signupForm.setError('email', { message: result.error.issues.find(e => e.path[0] === 'email')?.message })
          signupForm.setError('password', { message: result.error.issues.find(e => e.path[0] === 'password')?.message })
          signupForm.setError('confirmPassword', { message: result.error.issues.find(e => e.path[0] === 'confirmPassword')?.message })
          return
        }
        await register(result.data.email, result.data.password)
      }
      onClose()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Something went wrong')
    }
  }

  const handleGuestSignIn = async () => {
    setIsGuestLoading(true)
    setError(null)
    try {
      const { email, password } = await loginAsGuest()
      setGuestEmail(email)
      setGuestPassword(password)
      setGuestExpiresAt(new Date(Date.now() + 6 * 60 * 60 * 1000))
      setContent('guest_credentials')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create guest account')
    } finally {
      setIsGuestLoading(false)
    }
  }

  const handleGuestContinue = async () => {
    setContent('seeding')
    setSeedStep(0)
    setError(null)

    const STEP_DURATION_MS = 1000
    cancelledRef.current = false

    const runAnimation = (step: number) => {
      if (cancelledRef.current) return
      setSeedStep(step)
      if (step >= SEED_STEPS.length) return
      setTimeout(() => runAnimation(step + 1), STEP_DURATION_MS)
    }
    runAnimation(0)

    try {
      await apiFetch('/seed', { method: 'POST' })
    } catch {
      // seed error surfaced elsewhere; animation continues until last step
    } finally {
      cancelledRef.current = true
      setSeedStep(SEED_STEPS.length)
      setTimeout(() => onClose(), 1000)
    }
  }

  const handleCopy = async (text: string, field: 'email' | 'password') => {
    try {
      await navigator.clipboard.writeText(text)
    } catch {
      const textarea = document.createElement('textarea')
      textarea.value = text
      textarea.style.position = 'fixed'
      textarea.style.opacity = '0'
      document.body.appendChild(textarea)
      textarea.select()
      document.execCommand('copy')
      document.body.removeChild(textarea)
    }
    setCopiedField(field)
    setTimeout(() => setCopiedField(null), 2000)
  }

  if (!mounted) return null

  const renderLoginForm = () => (
    <form onSubmit={handleSubmit}>
      {error && (
        <div className="mb-4 rounded-md bg-destructive/15 px-3 py-2 text-sm text-destructive">
          {error}
        </div>
      )}
      <FieldGroup>
        <Field>
          <FieldLabel htmlFor="auth-email">Email</FieldLabel>
          <Input
            id="auth-email"
            type="email"
            placeholder="m@example.com"
            {...(mode === 'login' ? loginForm.register('email') : signupForm.register('email'))}
            aria-invalid={!!emailError}
          />
          {emailError && <p className="text-sm text-destructive">{emailError}</p>}
        </Field>

        <Field>
          <FieldLabel htmlFor="auth-password">Password</FieldLabel>
          <Input
            id="auth-password"
            type="password"
            placeholder="Enter your password"
            {...(mode === 'login' ? loginForm.register('password') : signupForm.register('password'))}
            aria-invalid={!!passwordError}
          />
          {passwordError && <p className="text-sm text-destructive">{passwordError}</p>}
        </Field>

        <div
          className="grid transition-[grid-template-rows] duration-300"
          style={{ gridTemplateRows: mode === 'signup' ? '1fr' : '0fr' }}
        >
          <div className="overflow-hidden">
            <Field>
              <FieldLabel htmlFor="auth-confirm">Confirm Password</FieldLabel>
              <Input
                id="auth-confirm"
                type="password"
                placeholder="Confirm your password"
                {...signupForm.register('confirmPassword')}
                aria-invalid={!!confirmPasswordError}
                tabIndex={mode === 'signup' ? undefined : -1}
              />
              {confirmPasswordError && <p className="text-sm text-destructive">{confirmPasswordError}</p>}
            </Field>
          </div>
        </div>

        <Field>
          <Button type="submit" disabled={isSubmitting} className="w-full">
            {isSubmitting ? 'Please wait...' : mode === 'login' ? 'Login' : 'Sign Up'}
          </Button>

          {mode === 'login' && (
            <>
              <div className="relative my-2">
                <div className="absolute inset-0 flex items-center">
                  <span className="w-full border-t" />
                </div>
                <div className="relative flex justify-center text-xs uppercase">
                  <span className="bg-card px-2 text-muted-foreground">or</span>
                </div>
              </div>

              <Button
                type="button"
                variant="outline"
                className="w-full"
                disabled={isGuestLoading}
                onClick={handleGuestSignIn}
              >
                {isGuestLoading ? 'Creating guest account...' : 'Sign in as Guest'}
              </Button>
            </>
          )}

          <FieldDescription className="text-center">
            {mode === 'login' ? (
              <>
                Don&apos;t have an account?{' '}
                <button type="button" onClick={() => { setMode('signup'); setError(null); loginForm.clearErrors(); signupForm.clearErrors() }} className="underline underline-offset-4 hover:text-foreground">
                  Sign up
                </button>
              </>
            ) : (
              <>
                Already have an account?{' '}
                <button type="button" onClick={() => { setMode('login'); setError(null); loginForm.clearErrors(); signupForm.clearErrors() }} className="underline underline-offset-4 hover:text-foreground">
                  Login
                </button>
              </>
            )}
          </FieldDescription>
        </Field>
      </FieldGroup>
    </form>
  )

  const renderGuestCredentials = () => (
    <div className="flex flex-col gap-4">
      <div className="rounded-lg border bg-muted/30 p-4 space-y-3">
        <div className="flex items-center justify-between gap-2">
          <div className="min-w-0 flex-1">
            <p className="text-xs text-muted-foreground mb-0.5">Email</p>
            <p className="text-sm font-medium truncate">{guestEmail}</p>
          </div>
          <Button
            type="button"
            variant="outline"
            size="sm"
            className="shrink-0 h-8"
            onClick={() => handleCopy(guestEmail, 'email')}
          >
            {copiedField === 'email' ? 'Copied!' : 'Copy'}
          </Button>
        </div>
        <div className="flex items-center justify-between gap-2">
          <div className="min-w-0 flex-1">
            <p className="text-xs text-muted-foreground mb-0.5">Password</p>
            <p className="text-sm font-mono font-medium">{guestPassword}</p>
          </div>
          <Button
            type="button"
            variant="outline"
            size="sm"
            className="shrink-0 h-8 font-mono"
            onClick={() => handleCopy(guestPassword, 'password')}
          >
            {copiedField === 'password' ? 'Copied!' : 'Copy'}
          </Button>
        </div>
      </div>

      <p className="text-xs text-muted-foreground text-center">
        Account expires at{' '}
        {guestExpiresAt?.toLocaleString(undefined, {
          year: 'numeric',
          month: 'short',
          day: 'numeric',
          hour: '2-digit',
          minute: '2-digit',
        })}
      </p>

      <Button type="button" className="w-full" onClick={handleGuestContinue}>
        Continue
      </Button>
    </div>
  )

  const renderSeedAnimation = () => (
    <div className="flex flex-col gap-3 py-2">
      {SEED_STEPS.map((label, i) => {
        const isDone = i < seedStep
        const isActive = i === seedStep
        const isPending = i > seedStep

        return (
          <div key={label} className="flex items-center gap-3">
            {isDone ? (
              <span
                key="done"
                className="text-green-600 text-sm font-bold shrink-0"
              >
                &#10003;
              </span>
            ) : isActive ? (
              <span
                key="active"
                className="text-primary text-sm shrink-0 animate-pulse"
              >
                &#9679;
              </span>
            ) : (
              <span
                key="pending"
                className="text-muted-foreground/30 text-sm shrink-0"
              >
                &#9679;
              </span>
            )}
            <span
              className={cn(
                'text-sm transition-colors duration-300',
                isPending ? 'text-muted-foreground/40' : 'text-foreground'
              )}
            >
              {label}
            </span>
          </div>
        )
      })}
    </div>
  )

  const getTitle = () => {
    if (content === 'guest_credentials') return 'Your Guest Account'
    if (content === 'seeding') return 'Setting up your workspace'
    return mode === 'login' ? 'Login to your account' : 'Sign Up'
  }

  const getDescription = () => {
    if (content === 'guest_credentials') return 'Save your credentials before continuing'
    if (content === 'seeding') return 'Please wait while we prepare your workspace'
    return mode === 'login'
      ? 'Enter your email below to login to your account'
      : 'Create an account to get started'
  }

  const overlayClass = visible ? 'animate-auth-overlay-in' : 'animate-auth-overlay-out'
  const cardClass = visible ? 'animate-auth-card-in' : 'animate-auth-card-out'

  return (
    <div className={cn('fixed inset-0 z-60 flex items-center justify-center bg-black/30 backdrop-blur-sm p-4', overlayClass)}>
      <div className={cn('w-full max-w-sm', cardClass)}>
        <div className="flex flex-col gap-6">
          <Card>
            <CardHeader>
              <CardTitle>{getTitle()}</CardTitle>
              <CardDescription>{getDescription()}</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="relative overflow-hidden">
                <div
                  className={cn(
                    'transition-all duration-300',
                    content === 'login' || content === 'signup'
                      ? 'opacity-100 translate-x-0 relative'
                      : 'opacity-0 translate-x-4 absolute inset-0 pointer-events-none'
                  )}
                >
                  {renderLoginForm()}
                </div>
                <div
                  className={cn(
                    'transition-all duration-300',
                    content === 'guest_credentials'
                      ? 'opacity-100 translate-x-0 relative'
                      : 'opacity-0 -translate-x-4 absolute inset-0 pointer-events-none'
                  )}
                >
                  {renderGuestCredentials()}
                </div>
                <div
                  className={cn(
                    'transition-all duration-300',
                    content === 'seeding'
                      ? 'opacity-100 translate-x-0 relative'
                      : 'opacity-0 translate-x-4 absolute inset-0 pointer-events-none'
                  )}
                >
                  {renderSeedAnimation()}
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}
