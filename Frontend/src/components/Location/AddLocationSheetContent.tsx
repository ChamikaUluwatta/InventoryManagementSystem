import { useState } from 'react'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { createLocation } from '@/services/locationService'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Package, X } from 'lucide-react'
import { SectionLabel } from '@/components/ui/sheet-label'
import { ErrorMessage } from '@/components/ui/error-message'

type Props = {
  onClose: () => void
  onSuccess: () => void
}

const formSchema = z.object({
  location_id: z.string().min(1, 'Required'),
  image: z.string().optional(),
})

type FormData = z.infer<typeof formSchema>

export default function AddLocationSheetContent({ onClose, onSuccess }: Props) {
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const form = useForm<FormData>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      location_id: '',
      image: '',
    },
  })

  async function onSubmit(data: FormData) {
    setSaving(true)
    try {
      await createLocation({
        location_id: data.location_id,
        image: data.image || null,
      })
      onSuccess()
      onClose()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create location')
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="flex flex-col h-full bg-background">
      <div className="flex items-center justify-between px-5 py-4 border-b shrink-0">
        <div className="flex items-center gap-3">
          <div className="flex items-center justify-center w-8 h-8 rounded border border-border">
            <Package className="h-4 w-4" />
          </div>
          <div>
            <p className="text-[10px] font-mono uppercase tracking-widest text-muted-foreground leading-none mb-0.5">
              Creating
            </p>
            <h2 className="text-sm font-semibold leading-none">New Location</h2>
          </div>
        </div>
        <Button
          variant="ghost"
          size="icon"
          className="h-7 w-7 rounded-sm text-muted-foreground hover:text-foreground hover:bg-muted"
          onClick={onClose}
        >
          <X className="h-3.5 w-3.5" />
        </Button>
      </div>

      <div className="flex-1 overflow-y-auto p-5 space-y-5">
        {error && (
          <ErrorMessage message={error} />
        )}

        <form id="location-form" onSubmit={form.handleSubmit(onSubmit)} className="space-y-5">
          <div>
            <SectionLabel>Location ID</SectionLabel>
            <Input
              id="location_id"
              className="h-9 text-sm font-mono"
              placeholder="e.g. A1-B2-C3"
              {...form.register('location_id')}
            />
            {form.formState.errors.location_id && (
              <p className="text-xs text-destructive font-mono mt-1">
                {form.formState.errors.location_id.message}
              </p>
            )}
          </div>

          <div>
            <SectionLabel>Image URL (optional)</SectionLabel>
            <Input
              id="image"
              className="h-9 text-sm font-mono"
              placeholder="https://example.com/image.jpg"
              {...form.register('image')}
            />
          </div>
        </form>
      </div>

      <div className="px-5 py-4 border-t shrink-0 flex gap-2 justify-end items-center">
        <Button
          type="button"
          variant="ghost"
          size="sm"
          className="h-8 px-3 text-xs text-muted-foreground"
          onClick={onClose}
        >
          Cancel
        </Button>
        <Button
          type="submit"
          form="location-form"
          size="sm"
          className="h-8 px-4 text-xs"
          disabled={saving}
        >
          {saving ? (
            <span className="flex items-center gap-2">
              <span className="w-3 h-3 border border-current border-t-transparent rounded-full animate-spin" />
              Creating...
            </span>
          ) : (
            'Create Location'
          )}
        </Button>
      </div>
    </div>
  )
}
