
import { cn } from '@/lib/utils'

export function FieldGroup({ className, ...props }: React.ComponentProps<'div'>) {
  return <div className={cn('flex flex-col gap-4', className)} {...props} />
}

export function Field({ className, ...props }: React.ComponentProps<'div'>) {
  return <div className={cn('flex flex-col gap-2', className)} {...props} />
}

export function FieldLabel({ className, ...props }: React.ComponentProps<'label'>) {
  return <label className={cn('text-sm font-medium', className)} {...props} />
}

export function FieldDescription({ className, ...props }: React.ComponentProps<'p'>) {
  return (
    <p
      className={cn('text-sm text-muted-foreground text-balance', className)}
      {...props}
    />
  )
}
