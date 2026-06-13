import { cn } from '@/lib/utils'

interface ErrorMessageProps {
  message: string
  className?: string
}

export function ErrorMessage({ message, className }: ErrorMessageProps) {
  return (
    <div
      className={cn(
        'p-3 border border-red-200 bg-red-50 dark:bg-red-950/20 dark:border-red-900 rounded text-sm text-red-600 dark:text-red-400 font-mono',
        className,
      )}
    >
      {message}
    </div>
  )
}
