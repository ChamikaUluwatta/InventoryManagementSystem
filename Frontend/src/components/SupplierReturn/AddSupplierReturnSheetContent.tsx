import { useEffect, useState } from 'react'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm, useFieldArray } from 'react-hook-form'
import { z } from 'zod'
import type { Company } from '@/types/company'
import type { Product } from '@/types/product'
import { createSupplierReturn } from '@/services/supplierReturnService'
import { getAllCompanies } from '@/services/companyService'
import { getAllProducts } from '@/services/productService'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Textarea } from '@/components/ui/textarea'
import { X, Package, Plus, Trash2 } from 'lucide-react'
import { SectionLabel, EditLabel, EditCell } from '../ui/sheet-label'

type Props = {
  onClose: () => void
  onSuccess: () => void
}

const itemSchema = z.object({
  product_id: z.string().min(1, 'Required'),
  quantity: z.number().int().positive('Quantity must be positive'),
  unit_cost: z.number().nonnegative('Unit cost must be non-negative'),
})

const formSchema = z.object({
  company_id: z.string().min(1, 'Company is required'),
  return_no: z.string().min(1, 'Return number is required'),
  reason: z.string().optional(),
  notes: z.string().optional(),
  items: z.array(itemSchema).min(1, 'At least one item is required'),
})

type FormData = z.infer<typeof formSchema>

export default function AddSupplierReturnSheetContent({ onClose, onSuccess }: Props) {
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [companies, setCompanies] = useState<Company[]>([])
  const [products, setProducts] = useState<Product[]>([])

  const form = useForm<FormData>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      company_id: '',
      return_no: '',
      reason: '',
      notes: '',
      items: [{ product_id: '', quantity: 1, unit_cost: 0 }],
    },
  })

  const { fields, append, remove } = useFieldArray({
    control: form.control,
    name: 'items',
  })

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [companiesData, productsData] = await Promise.all([
          getAllCompanies(),
          getAllProducts(),
        ])
        setCompanies(companiesData)
        setProducts(productsData)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load data')
      } finally {
        setLoading(false)
      }
    }
    fetchData()
  }, [])

  async function onSubmit(data: FormData) {
    setSaving(true)
    try {
      await createSupplierReturn({
        company_id: data.company_id,
        return_no: data.return_no,
        reason: data.reason || null,
        notes: data.notes || null,
        items: data.items,
      })
      onSuccess()
      onClose()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create supplier return')
    } finally {
      setSaving(false)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center gap-4 justify-center h-full">
        <div className="w-12 h-12 border-2 border-primary border-t-transparent rounded-full animate-spin" />
        <p>Loading...</p>
      </div>
    )
  }

  if (error && !companies.length) {
    return (
      <div className="m-4 p-3 border border-red-200 bg-red-50 dark:bg-red-950/20 dark:border-red-900 rounded text-sm text-red-600 dark:text-red-400 font-mono">
        ERR: {error}
      </div>
    )
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
            <h2 className="text-sm font-semibold leading-none">New Supplier Return</h2>
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
        <form id="supplier-return-form" onSubmit={form.handleSubmit(onSubmit)} className="space-y-5">
          {error && (
            <div className="p-3 border border-red-200 bg-red-50 dark:bg-red-950/20 dark:border-red-900 rounded text-sm text-red-600 dark:text-red-400 font-mono">
              {error}
            </div>
          )}

          <div>
            <SectionLabel>Return Details</SectionLabel>
            <div className="grid grid-cols-2 border border-border rounded-md overflow-hidden">
              <EditCell>
                <EditLabel>RETURN NO</EditLabel>
                <Input
                  id="return_no"
                  className="h-8 text-sm font-mono"
                  placeholder="e.g. SR-001"
                  {...form.register('return_no')}
                />
                {form.formState.errors.return_no && (
                  <p className="text-xs text-destructive font-mono mt-1">
                    {form.formState.errors.return_no.message}
                  </p>
                )}
              </EditCell>
              <EditCell bordered>
                <EditLabel>COMPANY</EditLabel>
                <Select
                  value={form.watch('company_id')}
                  onValueChange={(value) => form.setValue('company_id', value)}
                >
                  <SelectTrigger className="h-8 text-sm font-mono">
                    <SelectValue placeholder="Select company" />
                  </SelectTrigger>
                  <SelectContent position="popper">
                    {companies.map((comp) => (
                      <SelectItem
                        key={comp.company_id}
                        value={comp.company_id.toString()}
                        className="font-mono text-sm"
                      >
                        {comp.company_name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                {form.formState.errors.company_id && (
                  <p className="text-xs text-destructive font-mono mt-1">
                    {form.formState.errors.company_id.message}
                  </p>
                )}
              </EditCell>
              <EditCell topBorder className="col-span-2">
                <EditLabel>REASON</EditLabel>
                <Input
                  id="reason"
                  className="h-8 text-sm font-mono"
                  placeholder="e.g. Defective items, Wrong shipment"
                  {...form.register('reason')}
                />
              </EditCell>
            </div>
          </div>

          <div>
            <SectionLabel>Notes</SectionLabel>
            <Textarea
              id="notes"
              className="min-h-20 text-sm font-mono resize-none"
              placeholder="Additional notes..."
              {...form.register('notes')}
            />
          </div>

          <div>
            <div className="flex items-center justify-between mb-2">
              <SectionLabel>Return Items</SectionLabel>
              <Button
                type="button"
                variant="outline"
                size="sm"
                className="h-7 text-xs font-mono gap-1"
                onClick={() => append({ product_id: '', quantity: 1, unit_cost: 0 })}
              >
                <Plus className="h-3 w-3" />
                ADD ITEM
              </Button>
            </div>
            <div className="space-y-2">
              {fields.map((field, index) => (
                <div key={field.id} className="flex items-start gap-2 p-3 border border-border rounded-md">
                  <div className="flex-1 grid grid-cols-3 gap-2">
                    <div>
                      <EditLabel>PRODUCT</EditLabel>
                      <Select
                        value={form.watch(`items.${index}.product_id`).toString()}
                        onValueChange={(value) =>
                          form.setValue(`items.${index}.product_id`, value)
                        }
                      >
                        <SelectTrigger className="h-8 text-sm font-mono">
                          <SelectValue placeholder="Select" />
                        </SelectTrigger>
                        <SelectContent position="popper">
                          {products.map((prod) => (
                            <SelectItem
                              key={prod.product_id}
                              value={prod.product_id.toString()}
                              className="font-mono text-sm"
                            >
                              {prod.product_name}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      {form.formState.errors.items?.[index]?.product_id && (
                        <p className="text-xs text-destructive font-mono mt-1">
                          {form.formState.errors.items[index]?.product_id?.message}
                        </p>
                      )}
                    </div>
                    <div>
                      <EditLabel>QUANTITY</EditLabel>
                      <Input
                        type="number"
                        min={1}
                        className="h-8 text-sm font-mono"
                        {...form.register(`items.${index}.quantity`, { valueAsNumber: true })}
                      />
                      {form.formState.errors.items?.[index]?.quantity && (
                        <p className="text-xs text-destructive font-mono mt-1">
                          {form.formState.errors.items[index]?.quantity?.message}
                        </p>
                      )}
                    </div>
                    <div>
                      <EditLabel>UNIT COST</EditLabel>
                      <Input
                        type="number"
                        step="0.01"
                        min={0}
                        className="h-8 text-sm font-mono"
                        {...form.register(`items.${index}.unit_cost`, { valueAsNumber: true })}
                      />
                      {form.formState.errors.items?.[index]?.unit_cost && (
                        <p className="text-xs text-destructive font-mono mt-1">
                          {form.formState.errors.items[index]?.unit_cost?.message}
                        </p>
                      )}
                    </div>
                  </div>
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8 mt-5 text-muted-foreground hover:text-destructive"
                    onClick={() => remove(index)}
                    disabled={fields.length === 1}
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              ))}
              {form.formState.errors.items?.root && (
                <p className="text-xs text-destructive font-mono">
                  {form.formState.errors.items.root.message}
                </p>
              )}
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
          form="supplier-return-form"
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
            'Create Return'
          )}
        </Button>
      </div>
    </div>
  )
}
