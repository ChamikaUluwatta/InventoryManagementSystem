import { useEffect, useState } from 'react'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import type { Company } from '@/types/company'
import type { Product } from '@/types/product'
import type { SupplierReturn } from '@/types/supplierReturn'
import { updateCompany, deleteCompany, getCompanyDependencies } from '@/services/companyService'
import { getProductsByCompany } from '@/services/productService'
import { getSupplierReturnsByCompany } from '@/services/supplierReturnService'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Button } from '@/components/ui/button'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogMedia,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import { Pencil, Trash2, X, Package, TriangleAlert, ChevronDown, ChevronRight, List } from 'lucide-react'
import { SectionLabel, DataCell } from '@/components/ui/sheet-label'
import { ErrorMessage } from '@/components/ui/error-message'

type Props = {
  company: Company
  onClose: () => void
  onSuccess: () => void
}

const formSchema = z.object({
  company_name: z.string().min(1, 'Required'),
  description: z.string().optional(),
})

type FormData = z.infer<typeof formSchema>

export default function CompanySheetContent({ company, onClose, onSuccess }: Props) {
  const [mode, setMode] = useState<'view' | 'edit'>('view')
  const [localCompany, setLocalCompany] = useState<Company>(company)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [showBlockDialog, setShowBlockDialog] = useState(false)
  const [blockingDeps, setBlockingDeps] = useState<{ product_count: number; supplier_count: number } | null>(null)
  const [showProducts, setShowProducts] = useState(false)
  const [showSupplierReturns, setShowSupplierReturns] = useState(false)
  const [products, setProducts] = useState<Product[]>([])
  const [supplierReturns, setSupplierReturns] = useState<SupplierReturn[]>([])
  const [loadingProducts, setLoadingProducts] = useState(false)
  const [loadingReturns, setLoadingReturns] = useState(false)

  const form = useForm<FormData>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      company_name: company.company_name,
      description: company.description || '',
    },
  })

  useEffect(() => {
    if (mode === 'edit') {
      form.reset({
        company_name: localCompany.company_name,
        description: localCompany.description || '',
      })
    }
  }, [mode, localCompany, form])

  async function onSubmit(data: FormData) {
    setSaving(true)
    try {
      await updateCompany(company.company_id, { company_name: data.company_name, description: data.description || '' })
      setLocalCompany({ ...localCompany, company_name: data.company_name, description: data.description || '' })
      setMode('view')
      onSuccess()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update company')
    } finally {
      setSaving(false)
    }
  }

  async function handleDelete() {
    try {
      await deleteCompany(company.company_id)
      onSuccess()
      onClose()
    } catch {
      try {
        const deps = await getCompanyDependencies(company.company_id)
        if (deps.product_count > 0 || deps.supplier_count > 0) {
          setBlockingDeps(deps)
          setShowBlockDialog(true)
          return
        }
      } catch {
        /* fall through */
      }
      setError('Failed to delete company')
    }
  }

  async function handleShowProducts() {
    setLoadingProducts(true)
    try {
      const data = await getProductsByCompany(company.company_id)
      setProducts(data)
      setShowProducts(true)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load products')
    } finally {
      setLoadingProducts(false)
    }
  }

  async function handleShowSupplierReturns() {
    setLoadingReturns(true)
    try {
      const data = await getSupplierReturnsByCompany(company.company_id)
      setSupplierReturns(data)
      setShowSupplierReturns(true)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load supplier returns')
    } finally {
      setLoadingReturns(false)
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
              {mode === 'view' ? 'Company Details' : 'Editing'}
            </p>
            <h2 className="text-sm font-semibold leading-none">{localCompany.company_name}</h2>
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
        {error && (
          <ErrorMessage message={error} />
        )}

        {mode === 'view' ? (
          <div>
            <SectionLabel>Details</SectionLabel>
            <div className="grid grid-cols-1 border border-border rounded-md overflow-hidden">
              <DataCell label="COMPANY NAME" value={localCompany.company_name} />
            </div>
            <div className="mt-3">
              <SectionLabel>Description</SectionLabel>
              <div className="px-4 py-3 border border-border rounded-md">
                {localCompany.description ? (
                  <p className="text-sm text-muted-foreground leading-relaxed">{localCompany.description}</p>
                ) : (
                  <p className="text-xs font-mono italic text-muted-foreground/50">No description</p>
                )}
              </div>
            </div>
          </div>
        ) : (
          <form id="company-form" onSubmit={form.handleSubmit(onSubmit)} className="space-y-5">
            <div>
              <SectionLabel>Company Name</SectionLabel>
              <Input
                id="company_name"
                className="h-9 text-sm font-mono"
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
                  <AlertDialogTitle>Delete this company?</AlertDialogTitle>
                  <AlertDialogDescription>
                    This will permanently remove <strong>{localCompany.company_name}</strong>. This action cannot be undone.
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
              onClick={() => {
                window.setTimeout(() => setMode('edit'), 0)
              }}
            >
              <Pencil className="h-3.5 w-3.5" />
              Edit Company
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
              form="company-form"
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

      <AlertDialog open={showBlockDialog} onOpenChange={setShowBlockDialog}>
        <AlertDialogContent className="sm:max-w-lg">
          <AlertDialogHeader>
            <AlertDialogMedia>
              <TriangleAlert className="h-6 w-6 text-amber-500" />
            </AlertDialogMedia>
            <AlertDialogTitle>Cannot delete this company</AlertDialogTitle>
            <AlertDialogDescription>
              <strong>{localCompany.company_name}</strong> still has:
            </AlertDialogDescription>
          </AlertDialogHeader>

          <div className="px-6 space-y-3">
            <ul className="text-sm space-y-1 list-disc list-inside">
              {blockingDeps && blockingDeps.product_count > 0 && (
                <li>{blockingDeps.product_count} product{blockingDeps.product_count !== 1 ? 's' : ''}</li>
              )}
              {blockingDeps && blockingDeps.supplier_count > 0 && (
                <li>{blockingDeps.supplier_count} supplier return{blockingDeps.supplier_count !== 1 ? 's' : ''}</li>
              )}
            </ul>

            <p className="text-xs text-muted-foreground">
              Remove these before deleting the company.
            </p>

            {blockingDeps && blockingDeps.product_count > 0 && (
              <div className="border border-border rounded-md overflow-hidden">
                <button
                  type="button"
                  onClick={handleShowProducts}
                  disabled={loadingProducts}
                  className="w-full flex items-center justify-between px-3 py-2 text-xs font-mono font-medium text-left hover:bg-muted transition-colors disabled:opacity-50"
                >
                  <span className="flex items-center gap-2">
                    <List className="h-3 w-3" />
                    Show Products ({blockingDeps.product_count})
                  </span>
                  {loadingProducts ? (
                    <span className="w-3 h-3 border border-current border-t-transparent rounded-full animate-spin" />
                  ) : showProducts ? (
                    <ChevronDown className="h-3 w-3" />
                  ) : (
                    <ChevronRight className="h-3 w-3" />
                  )}
                </button>
                {showProducts && products.length > 0 && (
                  <div className="border-t border-border max-h-48 overflow-y-auto">
                    {products.slice(0, 20).map((p) => (
                      <div
                        key={p.product_id}
                        className="flex items-center justify-between px-3 py-1.5 text-xs border-b border-border last:border-b-0"
                      >
                        <span className="font-mono truncate min-w-0">{p.product_name}</span>
                        <span className="text-muted-foreground shrink-0 ml-2">${p.price}</span>
                      </div>
                    ))}
                    {products.length > 20 && (
                      <div className="px-3 py-2 text-xs text-muted-foreground font-mono italic">
                        ...and {products.length - 20} more
                      </div>
                    )}
                  </div>
                )}
              </div>
            )}

            {blockingDeps && blockingDeps.supplier_count > 0 && (
              <div className="border border-border rounded-md overflow-hidden">
                <button
                  type="button"
                  onClick={handleShowSupplierReturns}
                  disabled={loadingReturns}
                  className="w-full flex items-center justify-between px-3 py-2 text-xs font-mono font-medium text-left hover:bg-muted transition-colors disabled:opacity-50"
                >
                  <span className="flex items-center gap-2">
                    <List className="h-3 w-3" />
                    Show Supplier Returns ({blockingDeps.supplier_count})
                  </span>
                  {loadingReturns ? (
                    <span className="w-3 h-3 border border-current border-t-transparent rounded-full animate-spin" />
                  ) : showSupplierReturns ? (
                    <ChevronDown className="h-3 w-3" />
                  ) : (
                    <ChevronRight className="h-3 w-3" />
                  )}
                </button>
                {showSupplierReturns && supplierReturns.length > 0 && (
                  <div className="border-t border-border max-h-48 overflow-y-auto">
                    {supplierReturns.slice(0, 20).map((sr) => (
                      <div
                        key={sr.supplier_return_id}
                        className="flex items-center justify-between px-3 py-1.5 text-xs border-b border-border last:border-b-0"
                      >
                        <span className="font-mono truncate min-w-0">{sr.return_no}</span>
                        <span className="text-muted-foreground shrink-0 ml-2 capitalize">{sr.status}</span>
                      </div>
                    ))}
                    {supplierReturns.length > 20 && (
                      <div className="px-3 py-2 text-xs text-muted-foreground font-mono italic">
                        ...and {supplierReturns.length - 20} more
                      </div>
                    )}
                  </div>
                )}
              </div>
            )}
          </div>

          <AlertDialogFooter>
            <AlertDialogCancel>Got it</AlertDialogCancel>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  )
}
