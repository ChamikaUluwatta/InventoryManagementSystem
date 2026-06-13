import { useState } from 'react'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import type { Location } from '@/types/location'
import { updateLocation, deleteLocation } from '@/services/locationService'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import { Pencil, Trash2, X, Package } from 'lucide-react'
import { SectionLabel, DataCell } from '@/components/ui/sheet-label'
import { ErrorMessage } from '@/components/ui/error-message'

type Props = {
  location: Location
  onClose: () => void
  onSuccess: () => void
}

const formSchema = z.object({
  image: z.string().optional(),
})

type FormData = z.infer<typeof formSchema>

export default function LocationSheetContent({ location, onClose, onSuccess }: Props) {
  const [mode, setMode] = useState<'view' | 'edit'>('view')
  const [localLocation, setLocalLocation] = useState<Location>(location)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const form = useForm<FormData>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      image: location.image || '',
    },
  })

  async function onSubmit(data: FormData) {
    setSaving(true)
    try {
      await updateLocation(localLocation.location_id, {
        image: data.image || null,
      })
      setLocalLocation({ ...localLocation, image: data.image || null })
      setMode('view')
      onSuccess()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update location')
    } finally {
      setSaving(false)
    }
  }

  async function handleDelete() {
    try {
      await deleteLocation(localLocation.location_id)
      onSuccess()
      onClose()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete location')
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
              {mode === 'view' ? 'Location Details' : 'Editing'}
            </p>
            <h2 className="text-sm font-semibold leading-none">{localLocation.location_id}</h2>
          </div>
        </div>
        <Button
          variant="ghost"
          size="icon"
          className="h-7 w-7 rounded-sm text-muted-foreground hover:text-foreground hover:bg-muted"
          onClick={mode === 'view' ? onClose : () => setMode('view')}
        >
          <X className="h-3.5 w-3.5" />
        </Button>
      </div>

      <div className="flex-1 overflow-y-auto p-5 space-y-5">
        {error && <ErrorMessage message={error} />}

        {mode === 'view' ? (
          <div>
            <SectionLabel>Details</SectionLabel>
            <div className="grid grid-cols-1 border border-border rounded-md overflow-hidden">
              <DataCell label="LOCATION ID" value={localLocation.location_id} />
            </div>
            <div className="mt-3">
              <SectionLabel>Image</SectionLabel>
              <div className="px-4 py-3 border border-border rounded-md">
                {localLocation.image ? (
                  <p className="text-sm text-muted-foreground leading-relaxed break-all font-mono">{localLocation.image}</p>
                ) : (
                  <p className="text-xs font-mono italic text-muted-foreground/50">No image</p>
                )}
              </div>
            </div>
          </div>
        ) : (
          <form id="location-form" onSubmit={form.handleSubmit(onSubmit)} className="space-y-5">
            <div>
              <SectionLabel>Location ID</SectionLabel>
              <div className="px-4 py-3 border border-border rounded-md">
                <p className="text-sm font-mono text-muted-foreground">{localLocation.location_id}</p>
              </div>
            </div>

            <div>
              <SectionLabel>Image URL</SectionLabel>
              <Input
                id="image"
                className="h-9 text-sm font-mono"
                placeholder="https://example.com/image.jpg"
                {...form.register('image')}
              />
            </div>
          </form>
        )}
      </div>

      <div className="px-5 py-4 border-t shrink-0 flex gap-2 justify-between items-center">
        {mode === 'view' ? (
          <>
            <AlertDialog>
              <AlertDialogTrigger asChild>
                <Button
                  variant="ghost"
                  size="sm"
                  className="text-muted-foreground hover:text-destructive hover:bg-destructive/10 h-8 px-3 text-xs gap-1.5"
                >
                  <Trash2 className="h-3.5 w-3.5" />
                  Delete
                </Button>
              </AlertDialogTrigger>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>Delete this location?</AlertDialogTitle>
                  <AlertDialogDescription>
                    This will permanently remove <strong>{localLocation.location_id}</strong>. This action cannot be undone.
                  </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel>Cancel</AlertDialogCancel>
                  <AlertDialogAction
                    className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                    onClick={handleDelete}
                  >
                    Delete
                  </AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>

            <Button
              type="button"
              size="sm"
              className="h-8 px-4 text-xs gap-1.5"
              onClick={() => window.setTimeout(() => setMode('edit'), 0)}
            >
              <Pencil className="h-3.5 w-3.5" />
              Edit Location
            </Button>
          </>
        ) : (
          <>
            <Button
              type="button"
              variant="ghost"
              size="sm"
              className="h-8 px-3 text-xs text-muted-foreground"
              onClick={() => setMode('view')}
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
                  Saving…
                </span>
              ) : (
                'Save Changes'
              )}
            </Button>
          </>
        )}
      </div>
    </div>
  )
}
