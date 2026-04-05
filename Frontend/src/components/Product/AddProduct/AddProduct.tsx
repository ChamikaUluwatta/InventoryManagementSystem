import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
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

export const formSchema = z
  .object({
    product_name: z.string().min(1, 'Product name is required'),
    product_description: z.string().optional(),
    diameter: z.number(),
    width: z.number(),
    price: z.number(),
    category_id: z.number().int().positive().optional(),
    company_id: z.string().min(1, 'Company is required'),
    location_id: z.string().optional(),
  })
  .refine((data) => data.diameter > 0, {
    message: 'Diameter must be a positive number',
    path: ['diameter'],
  })
  .refine((data) => data.width > 0, {
    message: 'Width must be a positive number',
    path: ['width'],
  })
  .refine((data) => data.price > 0, {
    message: 'Price must be a positive number',
    path: ['price'],
  })

type FormData = z.infer<typeof formSchema>

export default function AddProduct() {
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
    return <div className="p-4">Loading...</div>
  }

  if (error) {
    return <div className="p-4 text-red-500">Error: {error}</div>
  }

  return (
    <div className="container mx-auto py-10">
      <Card>
        <CardHeader className="flex flex-row justify-center">
          <CardTitle>Add Product</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <div className='flex flex-row gap-6'>
              <div className="space-y-2">
                <Label htmlFor="product_name">Product Name</Label>
                <Input id="product_name" placeholder='Type your product name' {...form.register('product_name')} />
                {form.formState.errors.product_name && (
                  <p className="text-sm text-red-500">
                    {form.formState.errors.product_name.message}
                  </p>
                )}
              </div>

              <div className="space-y-2 grow-3">
                <Label htmlFor="product_description">Description</Label>
                <Textarea id="product_description" className='min-h-25 resize-y' placeholder='Type your product description' {...form.register('product_description')} />
              </div>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="diameter">Diameter</Label>
                <Input
                  id="diameter"
                  type="number"
                  step="0.01"
                  {...form.register('diameter', { valueAsNumber: true })}
                />
                {form.formState.errors.diameter && (
                  <p className="text-sm text-red-500">{form.formState.errors.diameter.message}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="width">Width</Label>
                <Input
                  id="width"
                  type="number"
                  step="0.01"
                  {...form.register('width', { valueAsNumber: true })}
                />
                {form.formState.errors.width && (
                  <p className="text-sm text-red-500">{form.formState.errors.width.message}</p>
                )}
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="price">Price</Label>
              <Input
                id="price"
                type="number"
                step="0.01"
                {...form.register('price', { valueAsNumber: true })}
              />
              {form.formState.errors.price && (
                <p className="text-sm text-red-500">{form.formState.errors.price.message}</p>
              )}
            </div>

            <div className="grid grid-cols-3 gap-4">
              <div className="flex flex-col items-center gap-3">
                <Label htmlFor="category_id">Category</Label>
                <Select
                  value={form.watch('category_id')?.toString()}
                  onValueChange={(value) =>
                    form.setValue('category_id', value ? parseInt(value) : undefined)
                  }
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select category" />
                  </SelectTrigger>
                  <SelectContent position="popper">
                    {categories.map((cat) => (
                      <SelectItem key={cat.category_id} value={cat.category_id.toString()}>
                        {cat.category_name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                {form.formState.errors.category_id && (
                  <p className="text-sm text-red-500">
                    {form.formState.errors.category_id.message}
                  </p>
                )}
              </div>

              <div className="flex flex-col items-center gap-3">
                <Label htmlFor="company_id">Company</Label>
                <Select
                  value={form.watch('company_id')}
                  onValueChange={(value) => form.setValue('company_id', value)}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select company" />
                  </SelectTrigger>
                  <SelectContent position="popper">
                    {companies.map((comp) => (
                      <SelectItem key={comp.company_id} value={comp.company_id.toString()}>
                        {comp.company_name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                {form.formState.errors.company_id && (
                  <p className="text-sm text-red-500">{form.formState.errors.company_id.message}</p>
                )}
              </div>

              <div className="flex flex-col items-center gap-3">
                <Label htmlFor="location_id">Location</Label>
                <Select
                  value={form.watch('location_id')}
                  onValueChange={(value) => form.setValue('location_id', value)}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select location" />
                  </SelectTrigger>
                  <SelectContent position="popper">
                    {locations.map((loc) => (
                      <SelectItem key={loc.location_id} value={loc.location_id.toString()}>
                        {loc.location_id}
                      </SelectItem>
                    ))}
                    <SelectItem value="unassigned">Unassigned</SelectItem>
                  </SelectContent>
                </Select>
                {form.formState.errors.location_id && (
                  <p className="text-sm text-red-500">
                    {form.formState.errors.location_id.message}
                  </p>
                )}
              </div>
            </div>
            <CardFooter>
              <div className="flex gap-4 pt-4">
                <Button type="submit" disabled={saving}>
                  {saving ? 'Saving...' : 'Save Changes'}
                </Button>
                <Button type="button" variant="outline" onClick={() => navigate('/products')}>
                  Cancel
                </Button>
              </div>
            </CardFooter>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
