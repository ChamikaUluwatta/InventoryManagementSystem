import type { Category } from '@/types/category'
import { apiFetch } from '@/lib/api'

export const getAllCategories = async (): Promise<Category[]> => {
  return apiFetch<Category[]>('/categories')
}

export const getCategoryById = async (id: string): Promise<Category> => {
  return apiFetch<Category>(`/categories/${id}`)
}

export const createCategory = async (
  category: Omit<Category, 'category_id'>,
): Promise<Category> => {
  return apiFetch<Category>('/categories', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(category),
  })
}

export const updateCategory = async (
  id: string,
  category: Partial<Category>,
): Promise<Category> => {
  return apiFetch<Category>(`/categories/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(category),
  })
}

export const deleteCategory = async (id: string): Promise<void> => {
  return apiFetch<void>(`/categories/${id}`, {
    method: 'DELETE',
  })
}
