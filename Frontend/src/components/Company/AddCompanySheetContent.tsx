import { useState } from 'react'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { createCompany } from '@/services/companyService'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Button } from '@/components/ui/button'
import { Package, X } from 'lucide-react'
import { SectionLabel } from '@/components/ui/sheet-label'
import { ErrorMessage } from '@/components/ui/error-message'

type Props = {
  onClose: () => void
  onSuccess: () => void
}

const formSchema = z.object({
  company_name: z.string().min(1, 'Required'),
  description: z.string().optional(),
})

type FormData = z.infer<typeof formSchema>

export default function AddCompanySheetContent({ onClose, onSuccess }: Props) {
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const form = useForm<FormData>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      company_name: '',
      description: '',
    },
  })

  async function onSubmit(data: FormData) {
    setSaving(true)
    try {
      await createCompany({ company_name: data.company_name, description: data.description || '' })
      onSuccess()
      onClose()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create company')
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
            <h2 className="text-sm font-semibold leading-none">New Company</h2>
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

        <form id="company-form" onSubmit={form.handleSubmit(onSubmit)} className="space-y-5">
          <div>
            <SectionLabel>Company Name</SectionLabel>
            <Input
              id="company_name"
              className="h-9 text-sm font-mono"
              placeholder="Enter company name"
              {...form.register('company_name')}
            />
            {form.formState.errors.company_name && (
              <p className="text-xs text-destructive font-mono mt-1">
                {form.formState.errors.company_name.message}
              </p>
            )}
          </div>

          <div>
            <SectionLabel>Description (optional)</SectionLabel>
            <div className="px-4 py-3 border border-border rounded-md">
              <Textarea
                id="description"
                className="min-h-20 text-sm font-mono resize-none border-0 p-3 focus-visible:ring-0"
                placeholder="Enter company description"
                {...form.register('description')}
              />
            </div>
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
          form="company-form"
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
            'Create Company'
          )}
        </Button>
      </div>
    </div>
  )
}
