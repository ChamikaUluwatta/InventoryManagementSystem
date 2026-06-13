import { useEffect, useState } from 'react'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import type { Product } from '@/types/product'
import type { Location } from '@/types/location'
import { createInventory } from '@/services/inventoryService'
import { getAllProducts } from '@/services/productService'
import { getAllLocations } from '@/services/locationService'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Package, X } from 'lucide-react'
import { Spinner } from '@/components/ui/spinner'
import { SectionLabel, EditLabel } from '@/components/ui/sheet-label'
import { ErrorMessage } from '@/components/ui/error-message'

type Props = {
  onClose: () => void
  onSuccess: () => void
}

const formSchema = z.object({
  product_id: z.string().min(1, 'Required'),
  location_id: z.string().min(1, 'Required'),
  stock: z.number().int().nonnegative(),
})

type FormData = z.infer<typeof formSchema>

export default function AddInventorySheetContent({ onClose, onSuccess }: Props) {
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [products, setProducts] = useState<Product[]>([])
  const [locations, setLocations] = useState<Location[]>([])

  const form = useForm<FormData>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      product_id: '',
      location_id: '',
      stock: 0,
    },
  })

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [productsData, locationsData] = await Promise.all([
          getAllProducts(),
          getAllLocations(),
        ])
        setProducts(productsData)
        setLocations(locationsData)
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
      await createInventory({
        product_id: data.product_id,
        location_id: data.location_id,
        stock: data.stock,
      })
      onSuccess()
      onClose()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create inventory')
    } finally {
      setSaving(false)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center gap-4 justify-center h-full">
        <Spinner className="size-12" />
        <p>Loading...</p>
      </div>
    )
  }

  if (error && !products.length) {
    return (
      <ErrorMessage message={error} />
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
            <h2 className="text-sm font-semibold leading-none">New Inventory</h2>
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

        <form id="inventory-form" onSubmit={form.handleSubmit(onSubmit)} className="space-y-5">
          <div>
            <SectionLabel>Product</SectionLabel>
            <div className="border border-border rounded-md overflow-hidden">
              <div className="px-4 py-3">
                <EditLabel>PRODUCT</EditLabel>
                <Select
                  value={form.watch('product_id')}
                  onValueChange={(value) => form.setValue('product_id', value)}
                >
                  <SelectTrigger className="h-8 text-sm font-mono">
                    <SelectValue placeholder="Select product" />
                  </SelectTrigger>
                  <SelectContent position="popper">
                    {products.map((prod) => (
                      <SelectItem
                        key={prod.product_id}
                        value={prod.product_id}
                        className="font-mono text-sm"
                      >
                        {prod.product_name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                {form.formState.errors.product_id && (
                  <p className="text-xs text-destructive font-mono mt-1">
                    {form.formState.errors.product_id.message}
                  </p>
                )}
              </div>
            </div>
          </div>

          <div>
            <SectionLabel>Location</SectionLabel>
            <div className="border border-border rounded-md overflow-hidden">
              <div className="px-4 py-3">
                <EditLabel>STORAGE UNIT</EditLabel>
                <Select
                  value={form.watch('location_id')}
                  onValueChange={(value) => form.setValue('location_id', value)}
                >
                  <SelectTrigger className="h-8 text-sm font-mono">
                    <SelectValue placeholder="Select location" />
                  </SelectTrigger>
                  <SelectContent position="popper">
                    {locations.map((loc) => (
                      <SelectItem
                        key={loc.location_id}
                        value={loc.location_id}
                        className="font-mono text-sm"
                      >
                        {loc.location_id}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                {form.formState.errors.location_id && (
                  <p className="text-xs text-destructive font-mono mt-1">
                    {form.formState.errors.location_id.message}
                  </p>
                )}
              </div>
            </div>
          </div>

          <div>
            <SectionLabel>Stock</SectionLabel>
            <div className="flex items-center gap-3 px-4 py-3 border border-border rounded-md">
              <div className="flex justify-between w-full items-center">
                <span className="text-[10px] font-mono uppercase tracking-widest text-muted-foreground">
                  Quantity
                </span>
                <Input
                  id="stock"
                  type="number"
                  step="1"
                  className="w-24 h-8 text-xl font-bold font-mono text-right"
                  placeholder="0"
                  {...form.register('stock', { valueAsNumber: true })}
                />
              </div>
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
          form="inventory-form"
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
            'Create Inventory'
          )}
        </Button>
      </div>
    </div>
  )
}
