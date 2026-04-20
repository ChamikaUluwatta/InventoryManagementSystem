import { zodResolver } from '@hookform/resolvers/zod'
import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import type { Product } from '@/types/product'
import type { Category } from '@/types/category'
import type { Company } from '@/types/company'
import type { Location } from '@/types/location'
import { updateProduct, getProductById, deleteProduct } from '@/services/productService'
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
import { Spinner } from '../ui/spinner'
import { SectionLabel, EditLabel, EditCell, DataCell } from '../ui/sheet-label'

type Props = {
  product: Product
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
  })

type FormData = z.infer<typeof formSchema>

export default function ProductSheetContent({ product, onClose, onSuccess }: Props) {
  const [mode, setMode] = useState<'view' | 'edit'>('view')
  const [loading, setLoading] = useState(false)
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
    const fetchLookupData = async () => {
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
        console.error('Failed to load lookup data:', err)
      }
    }
    fetchLookupData()
  }, [])

  useEffect(() => {
    const fetchProductData = async () => {
      try {
        setLoading(true)
        const productData = await getProductById(product.product_id)
        form.reset({
          product_name: productData.product_name,
          product_description: productData.product_description || '',
          diameter: productData.diameter,
          width: productData.width,
          price: productData.price,
          category_id: productData.category_id,
          company_id: productData.company_id.toString(),
          location_id: productData.location_id,
          stock: productData.stock,
        })
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load product data')
      } finally {
        setLoading(false)
      }
    }

    if (mode === 'edit') {
      fetchProductData()
    }
  }, [mode, product.product_id])

  const getCategoryName = (id: number | undefined) => {
    if (!id) return '—'
    const cat = categories.find((c) => c.category_id === id)
    return cat?.category_name || '—'
  }

  const getCompanyName = (id: string) => {
    if (!id) return '—'
    const comp = companies.find((c) => c.company_id.toString() === id.toString())
    return comp?.company_name || '—'
  }

  const getLocationName = (id: string | undefined) => {
    if (!id) return '—'
    const loc = locations.find((l) => l.location_id === id)
    return loc?.location_id || '—'
  }

  async function onSubmit(data: FormData) {
    setSaving(true)
    try {
      const updatedProduct: Partial<Product> = {
        product_name: data.product_name,
        product_description: data.product_description,
        diameter: data.diameter,
        width: data.width,
        price: data.price,
        category_id: data.category_id,
        company_id: data.company_id,
        location_id: data.location_id,
      }
      await updateProduct(product.product_id, updatedProduct)
      setMode('view')
      onSuccess()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update product')
    } finally {
      setSaving(false)
    }
  }

  async function handleDelete(id: string) {
    try {
      await deleteProduct(id)
      onSuccess()
      onClose()
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to delete product')
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

  if (error) {
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
              {mode === 'view' ? 'Product Details' : 'Editing'}
            </p>
            <h2 className="text-sm font-semibold leading-none">{product.product_name}</h2>
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
        {mode === 'view' ? (
          <>
            <div className="flex items-baseline justify-between px-4 py-3 rounded-md border border-border">
              <span className="text-[10px] font-mono uppercase tracking-widest text-muted-foreground">
                Unit Price
              </span>
              <span className="text-2xl font-bold font-mono tabular-nums">${product.price}</span>
            </div>
            <div>
              <SectionLabel>Description</SectionLabel>
              <div className="px-4 py-3 border border-border rounded-md">
                <p className="text-sm text-muted-foreground leading-relaxed">
                  {product.product_description ? product.product_description : '—'}
                </p>
              </div>
            </div>
            <div>
              <SectionLabel>Specifications</SectionLabel>
              <div className="grid grid-cols-2 border border-border rounded-md overflow-hidden">
                <DataCell label="DIAMETER" value={`${product.diameter} mm`} bordered />
                <DataCell label="WIDTH" value={`${product.width} mm`} />
                <DataCell
                  label="CATEGORY"
                  value={getCategoryName(product.category_id)}
                  bordered
                  topBorder
                />
                <DataCell label="COMPANY" value={getCompanyName(product.company_id)} topBorder />
              </div>
            </div>

            <div>
              <SectionLabel>Location</SectionLabel>
              <div className="flex items-center gap-3 px-4 py-3 border border-border rounded-md">
                <div className="flex-1">
                  <p className="text-[10px] font-mono text-muted-foreground uppercase tracking-wider mb-0.5">
                    Storage Unit
                  </p>
                  <p className="text-sm  font-mono">
                    {getLocationName(product.location_id)}
                  </p>
                </div>
              </div>
            </div>
          </>
        ) : (
          //Edit section
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
                {...form.register('price', { valueAsNumber: true })}
              />
            </div>
            <div>
              <SectionLabel>Product Name</SectionLabel>
              <Input
                id="product_name"
                className="h-9 text-sm font-medium"
                {...form.register('product_name')}
              />
            </div>
            <div>
              <SectionLabel>Description</SectionLabel>
              <div className="px-4 py-3 border border-border rounded-md">
                <Textarea
                  id="product_description"
                  className="min-h-20 text-sm font-mono resize-none border-0 p-3"
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
                    {...form.register('diameter', { valueAsNumber: true })}
                  />
                </EditCell>
                <EditCell bordered>
                  <EditLabel>WIDTH</EditLabel>
                  <Input
                    id="width"
                    type="number"
                    step="0.01"
                    className="h-8 text-sm font-mono"
                    {...form.register('width', { valueAsNumber: true })}
                  />
                </EditCell>
                <EditCell topBorder className="col-span-2">
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
                <EditCell topBorder className="col-span-2">
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
                </EditCell>
              </div>
            </div>

            <div>
              <SectionLabel>Location</SectionLabel>
              <div className="flex items-center gap-3 px-4 py-3 border border-border rounded-md">
                <div className="flex-1">
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
                <div className="flex-1">
                  <DataCell
                    className="p-4 flex justify-between"
                    label="IN STOCK"
                    bordered
                    value={form.watch('stock').toString()}
                  />
                </div>
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
                  <AlertDialogTitle>Delete this product?</AlertDialogTitle>
                  <AlertDialogDescription>
                    This will permanently remove <strong>{product.product_name}</strong> from the
                    inventory. This action cannot be undone.
                  </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel>Cancel</AlertDialogCancel>
                  <AlertDialogAction
                    className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                    onClick={() => handleDelete(product.product_id.toString())}
                  >
                    Delete
                  </AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>

            <Button size="sm" className="h-8 px-4 text-xs gap-1.5" onClick={() => setMode('edit')}>
              <Pencil className="h-3.5 w-3.5" />
              Edit Product
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
              form="product-form"
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




