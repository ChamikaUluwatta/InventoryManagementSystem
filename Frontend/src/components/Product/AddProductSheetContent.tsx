import { useEffect, useState } from 'react'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import type { Category } from '@/types/category'
import type { Company } from '@/types/company'
import type { Location } from '@/types/location'
import { createProduct } from '@/services/productService'
import { getAllCategories } from '@/services/categoryService'
import { getAllCompanies } from '@/services/companyService'
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
import { Textarea } from '@/components/ui/textarea'
import { Package, X } from 'lucide-react'
import { Spinner } from '../ui/spinner'
import { SectionLabel, EditLabel, EditCell } from '../ui/sheet-label'
import { ErrorMessage } from '@/components/ui/error-message'
import { createInventory } from '@/services/inventoryService'
import type { InventoryCreateRequest } from '@/types/inventory'

type Props = {
  onClose: () => void
  onSuccess: () => void
}

const formSchema = z
  .object({
    product_name: z.string().min(1, 'Required'),
    product_description: z.string().optional(),
    diameter: z.number(),
    width: z.number(),
    price: z.number(),
    category_id: z.number().int().positive().optional(),
    company_id: z.string().min(1, 'Required'),
    location_id: z.string().optional(),
    stock: z.number().int().nonnegative(),
  })
  .refine((data) => data.diameter > 0, {
    message: 'Must be positive',
    path: ['diameter'],
  })
  .refine((data) => data.width > 0, {
    message: 'Must be positive',
    path: ['width'],
  })
  .refine((data) => data.price > 0, {
    message: 'Must be positive',
    path: ['price'],
  }).refine((data) => data.stock >= 0, {
    message: 'Must be non-negative',
    path: ['stock'],
  })

type FormData = z.infer<typeof formSchema>

export default function AddProductSheetContent({ onClose, onSuccess }: Props) {
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [categories, setCategories] = useState<Category[]>([])
  const [companies, setCompanies] = useState<Company[]>([])
  const [locations, setLocations] = useState<Location[]>([])

  const form = useForm<FormData>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      product_name: '',
      product_description: '',
      diameter: 0,
      width: 0,
      price: 0,
      category_id: undefined,
      company_id: '',
      location_id: '',
      stock: 0,
    },
  })

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [categoriesData, companiesData, locationsData] = await Promise.all([
          getAllCategories(),
          getAllCompanies(),
          getAllLocations(),
        ])
        setCategories(categoriesData)
        setCompanies(companiesData)
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
      const locationValue = data.location_id !== '' ? data.location_id : 'unassigned'
      const createdProduct = {
        product_name: data.product_name,
        product_description: data.product_description || '',
        diameter: data.diameter,
        width: data.width,
        price: data.price,
        category_id: data.category_id || 0,
        company_id: data.company_id,
        location_id: locationValue as string,
      }

      const newProduct = await createProduct(createdProduct)

      const createdInventory: InventoryCreateRequest = {
        location_id: locationValue?.toString() || 'unassigned',
        stock: data.stock,
        product_id: newProduct.product_id,
      }

      await createInventory(createdInventory)
      onSuccess()
      onClose()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create product')
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

  if (error && !categories.length) {
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
            <h2 className="text-sm font-semibold leading-none">New Product</h2>
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
        <form id="product-form" onSubmit={form.handleSubmit(onSubmit)} className="space-y-5">
          <div className="flex items-baseline justify-between px-4 py-3 rounded-md border border-border">
            <span className="text-[10px] font-mono uppercase tracking-widest text-muted-foreground">
              Unit Price
            </span>
            <Input
              id="price"
              type="number"
              step="0.01"
              className="w-24 h-8 text-xl font-bold font-mono text-right"
              placeholder="0.00"
              {...form.register('price', { valueAsNumber: true })}
            />
          </div>

          <div >
            <SectionLabel>Product Name</SectionLabel>
            <Input
              id="product_name"
              className="h-9 text-sm font-mono"
              placeholder="Enter product name"
              {...form.register('product_name')}
            />
            {form.formState.errors.product_name && (
              <p className="text-xs text-destructive font-mono mt-1">
                {form.formState.errors.product_name.message}
              </p>
            )}
          </div>

          <div>
            <SectionLabel>Description</SectionLabel>
            <div className="px-4 py-3 border border-border rounded-md">
              <Textarea
                id="product_description"
                className="min-h-20 text-sm font-mono resize-none border-0 p-3 focus-visible:ring-0"
                placeholder="Enter description"
                {...form.register('product_description')}
              />
            </div>
          </div>

          <div>
            <SectionLabel>Specifications</SectionLabel>
            <div className="grid grid-cols-2 border border-border rounded-md overflow-hidden">
              <EditCell>
                <EditLabel>DIAMETER</EditLabel>
                <Input
                  id="diameter"
                  type="number"
                  step="0.01"
                  className="h-8 text-sm font-mono"
                  placeholder="0.00"
                  {...form.register('diameter', { valueAsNumber: true })}
                />
                {form.formState.errors.diameter && (
                  <p className="text-xs text-destructive font-mono mt-1">
                    {form.formState.errors.diameter.message}
                  </p>
                )}
              </EditCell>
              <EditCell bordered>
                <EditLabel>WIDTH</EditLabel>
                <Input
                  id="width"
                  type="number"
                  step="0.01"
                  className="h-8 text-sm font-mono"
                  placeholder="0.00"
                  {...form.register('width', { valueAsNumber: true })}
                />
                {form.formState.errors.width && (
                  <p className="text-xs text-destructive font-mono mt-1">
                    {form.formState.errors.width.message}
                  </p>
                )}
              </EditCell>
              <EditCell topBorder className="col-span-2 flex justify-between items-center">
                <EditLabel>CATEGORY</EditLabel>
                <Select
                  value={form.watch('category_id')?.toString()}
                  onValueChange={(value) =>
                    form.setValue('category_id', value ? parseInt(value) : undefined)
                  }
                >
                  <SelectTrigger className="h-8 text-sm font-mono">
                    <SelectValue placeholder="Select" />
                  </SelectTrigger>
                  <SelectContent position="popper">
                    {categories.map((cat) => (
                      <SelectItem
                        key={cat.category_id}
                        value={cat.category_id.toString()}
                        className="font-mono text-sm"
                      >
                        {cat.category_name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </EditCell>
              <EditCell topBorder className="col-span-2 flex justify-between">
                <EditLabel>COMPANY</EditLabel>
                <Select
                  value={form.watch('company_id')}
                  onValueChange={(value) => form.setValue('company_id', value)}
                >
                  <SelectTrigger className="h-8 text-sm font-mono">
                    <SelectValue placeholder="Select" />
                  </SelectTrigger>
                  <SelectContent position="popper">
                    {companies.map((comp) => (
                      <SelectItem
                        key={comp.company_id}
                        value={comp.company_id.toString()}
                        className="font-mono text-sm text-wrap"
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
            </div>
          </div>

          <div>
            <SectionLabel>Location</SectionLabel>
            <div className="flex items-center gap-3 px-4 py-3 border border-border rounded-md ">
              <div className="flex justify-between w-full items-center">
                <EditLabel>STORAGE UNIT</EditLabel>
                <Select
                  value={form.watch('location_id')}
                  onValueChange={(value) => form.setValue('location_id', value)}
                >
                  <SelectTrigger className="h-8 text-sm font-mono">
                    <SelectValue placeholder="Select" />
                  </SelectTrigger>
                  <SelectContent position="popper">
                    {locations.map((loc) => (
                      <SelectItem
                        key={loc.location_id}
                        value={loc.location_id.toString()}
                        className="font-mono text-sm"
                      >
                        {loc.location_id}
                      </SelectItem>
                    ))}
                    <SelectItem value="unassigned" className="font-mono text-sm">
                      Unassigned
                    </SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>
          </div>
          <div>
            <SectionLabel>Inventory</SectionLabel>
            <div className="flex items-center gap-3 px-4 py-3 border border-border rounded-md">
              <div className="flex justify-between w-full items-center">
                <span className="text-[10px] font-mono uppercase tracking-widest text-muted-foreground">
                  Stock
                </span>
                <Input
                  id="stock"
                  type="number"
                  step="1"
                  className="w-24 h-8 text-xl font-bold font-mono text-right"
                  placeholder="0.00"
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
          form="product-form"
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
            'Create Product'
          )}
        </Button>
      </div>
    </div>
  )
}
