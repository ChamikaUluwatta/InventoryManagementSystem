import type { Product } from '@/types/product'
import { apiFetch } from '@/lib/api'

export const getAllProducts = async (): Promise<Product[]> => {
  return apiFetch<Product[]>('/products')
}

export const getProductById = async (id: string): Promise<Product> => {
  return apiFetch<Product>(`/products/${id}`)
}

export const createProduct = async (product: Omit<Product, 'product_id' | 'stock'>): Promise<Product> => {
  return apiFetch<Product>('/products', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(product),
  })
}

export const updateProduct = async (id: string, product: Partial<Product>): Promise<Product> => {
  return apiFetch<Product>(`/products/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(product),
  })
}

export const deleteProduct = async (id: string): Promise<void> => {
  return apiFetch<void>(`/products/${id}`, {
    method: 'DELETE',
  })
}

export const getProductsByCompany = async (companyId: string): Promise<Product[]> => {
  return apiFetch<Product[]>(`/products?company=${encodeURIComponent(companyId)}`)
}
