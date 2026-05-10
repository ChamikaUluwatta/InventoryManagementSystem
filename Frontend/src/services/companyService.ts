import type { Company } from '@/types/company'
import { apiFetch } from '@/lib/api'

export const getAllCompanies = async (): Promise<Company[]> => {
  return apiFetch<Company[]>('/companies')
}

export const getCompanyById = async (id: string): Promise<Company> => {
  return apiFetch<Company>(`/companies/${id}`)
}

export const createCompany = async (company: Omit<Company, 'company_id'>): Promise<Company> => {
  return apiFetch<Company>('/companies', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(company),
  })
}

export const updateCompany = async (id: string, company: Partial<Company>): Promise<Company> => {
  return apiFetch<Company>(`/companies/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(company),
  })
}

export const deleteCompany = async (id: string): Promise<void> => {
  return apiFetch<void>(`/companies/${id}`, {
    method: 'DELETE',
  })
}

export const getCompanyDependencies = async (id: string): Promise<{ product_count: number; supplier_count: number }> => {
  return apiFetch<{ product_count: number; supplier_count: number }>(`/companies/${id}/dependencies`)
}
