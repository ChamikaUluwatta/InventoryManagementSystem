import type { Category } from '@/types/category'

const API_BASE_URL = import.meta.env.VITE_API_URL;

if (!API_BASE_URL && import.meta.env.MODE === 'production') {
  throw new Error('VITE_API_URL environment variable is required for production build');
}

const API_BASE = API_BASE_URL || 'http://localhost:8080/api/v1'

export const getAllCategories = async (): Promise<Category[]> => {
  const response = await fetch(`${API_BASE}/categories`)
  if (!response.ok) {
    throw new Error('Failed to fetch categories')
  }
  return response.json()
}

export const getCategoryById = async (id: string): Promise<Category> => {
  const response = await fetch(`${API_BASE}/categories/${id}`)
  if (!response.ok) {
    throw new Error('Failed to fetch category')
  }
  return response.json()
}

export const createCategory = async (
  category: Omit<Category, 'category_id'>,
): Promise<Category> => {
  const response = await fetch(`${API_BASE}/categories`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(category),
  })
  if (!response.ok) {
    throw new Error('Failed to create category')
  }
  return response.json()
}

export const updateCategory = async (
  id: string,
  category: Partial<Category>,
): Promise<Category> => {
  const response = await fetch(`${API_BASE}/categories/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(category),
  })
  if (!response.ok) {
    throw new Error('Failed to update category')
  }
  return response.json()
}

export const deleteCategory = async (id: string): Promise<void> => {
  const response = await fetch(`${API_BASE}/categories/${id}`, {
    method: 'DELETE',
  })
  if (!response.ok) {
    throw new Error('Failed to delete category')
  }
}
