import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import * as z from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import type { Category } from '@/types/category'
import type { Company } from '@/types/company'
import type { Location } from '@/types/location'
import { createProduct } from '@/services/productService'
import { getAllCategories } from '@/services/categoryService'
import { getAllCompanies } from '@/services/companyService'
import { getAllLocations } from '@/services/locationService'
import { Textarea } from '@/components/ui/textarea'
import { Spinner } from '@/components/ui/spinner'
import { ArrowLeft } from 'lucide-react'

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

function FieldError({ message }: { message?: string }) {
  if (!message) return null
  return <p className="text-xs text-destructive font-mono mt-0.5">{message}</p>
}

export default function ViewAddProduct() {
  const navigate = useNavigate()
  const [loading, setLoading] = useState(true)
  const [categories, setCategories] = useState<Category[]>([])
  const [companies, setCompanies] = useState<Company[]>([])
  const [locations, setLocations] = useState<Location[]>([])
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)

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

      await createProduct(createdProduct)
      navigate('/products')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create product')
    } finally {
      setSaving(false)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64 gap-3">
        <Spinner className="size-6" />
        <span className="text-xs font-mono uppercase tracking-widest text-muted-foreground">Loading</span>
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
    <div className="max-w-2xl mx-auto">
      {/* Header */}
      <div className="border-b border-border p-4 flex items-center gap-4">
        <Button variant="ghost" size="icon" className="h-8 w-8" onClick={() => navigate('/products')}>
          <ArrowLeft className="h-4 w-4" />
        </Button>
        <h1 className="text-lg font-semibold tracking-tight">NEW PRODUCT</h1>
      </div>

      {/* Form */}
      <form id="add-product-form" onSubmit={form.handleSubmit(onSubmit)} className="divide-y divide-border">
        {/* Section: Basic Info */}
        <div className="p-5">
          <p className="section-header">BASIC INFORMATION</p>
          <div className="space-y-4 mt-3">
            <div className="space-y-1">
              <Label htmlFor="product_name" className="text-xs">PRODUCT NAME *</Label>
              <Input
                id="product_name"
                className="font-mono"
                placeholder="Enter product name"
                {...form.register('product_name')}
              />
              <FieldError message={form.formState.errors.product_name?.message} />
            </div>

            <div className="space-y-1">
              <Label htmlFor="product_description" className="text-xs">DESCRIPTION</Label>
              <Textarea
                id="product_description"
                className="font-mono min-h-24"
                placeholder="Enter description"
                {...form.register('product_description')}
              />
            </div>

            <div className="grid grid-cols-3 gap-4">
              <div className="space-y-1">
                <Label htmlFor="diameter" className="text-xs">DIAMETER *</Label>
                <Input
                  id="diameter"
                  type="number"
                  step="0.01"
                  className="font-mono"
                  placeholder="0.00"
                  {...form.register('diameter', { valueAsNumber: true })}
                />
                <FieldError message={form.formState.errors.diameter?.message} />
              </div>

              <div className="space-y-1">
                <Label htmlFor="width" className="text-xs">WIDTH *</Label>
                <Input
                  id="width"
                  type="number"
                  step="0.01"
                  className="font-mono"
                  placeholder="0.00"
                  {...form.register('width', { valueAsNumber: true })}
                />
                <FieldError message={form.formState.errors.width?.message} />
              </div>

              <div className="space-y-1">
                <Label htmlFor="price" className="text-xs">PRICE *</Label>
                <Input
                  id="price"
                  type="number"
                  step="0.01"
                  className="font-mono"
                  placeholder="0.00"
                  {...form.register('price', { valueAsNumber: true })}
                />
                <FieldError message={form.formState.errors.price?.message} />
              </div>
            </div>
          </div>
        </div>

        {/* Section: Relationships */}
        <div className="p-5">
          <p className="section-header">RELATIONSHIPS</p>
          <div className="grid grid-cols-3 gap-4 mt-3">
            <div className="space-y-1">
              <Label htmlFor="category_id" className="text-xs">CATEGORY*</Label>
              <Select
                value={form.watch('category_id')?.toString()}
                onValueChange={(value) =>
                  form.setValue('category_id', value ? parseInt(value) : undefined)
                }
              >
                <SelectTrigger className="font-mono">
                  <SelectValue placeholder="SELECT" />
                </SelectTrigger>
                <SelectContent position="popper">
                  {categories.map((cat) => (
                    <SelectItem key={cat.category_id} value={cat.category_id.toString()} className="font-mono">
                      {cat.category_name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-1">
              <Label htmlFor="company_id" className="text-xs">COMPANY *</Label>
              <Select
                value={form.watch('company_id')}
                onValueChange={(value) => form.setValue('company_id', value)}
              >
                <SelectTrigger className="font-mono">
                  <SelectValue placeholder="SELECT" />
                </SelectTrigger>
                <SelectContent position="popper">
                  {companies.map((comp) => (
                    <SelectItem key={comp.company_id} value={comp.company_id.toString()} className="font-mono">
                      {comp.company_name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <FieldError message={form.formState.errors.company_id?.message} />
            </div>

            <div className="space-y-1">
              <Label htmlFor="location_id" className="text-xs">LOCATION</Label>
              <Select
                value={form.watch('location_id')}
                onValueChange={(value) => form.setValue('location_id', value)}
              >
                <SelectTrigger className="font-mono">
                  <SelectValue placeholder="SELECT" />
                </SelectTrigger>
                <SelectContent position="popper">
                  {locations.map((loc) => (
                    <SelectItem key={loc.location_id} value={loc.location_id.toString()} className="font-mono">
                      {loc.location_id}
                    </SelectItem>
                  ))}
                  <SelectItem value="unassigned" className="font-mono">UNASSIGNED</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        </div>

        {/* Footer */}
        <div className="p-5 flex gap-3 justify-end">
          <Button
            type="button"
            variant="outline"
            size="sm"
            className="font-mono text-xs"
            onClick={() => navigate('/products')}
          >
            CANCEL
          </Button>
          <Button
            type="submit"
            size="sm"
            className="font-mono text-xs"
            disabled={saving}
          >
            {saving ? 'SAVING...' : 'CREATE PRODUCT'}
          </Button>
        </div>
      </form>
    </div>
  )
}