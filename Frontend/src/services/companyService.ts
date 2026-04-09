import type { Company } from '@/types/company'

const API_BASE_URL = import.meta.env.VITE_API_URL;

if (!API_BASE_URL && import.meta.env.MODE === 'production') {
  throw new Error('VITE_API_URL environment variable is required for production build');
}

const API_BASE = API_BASE_URL || 'http://localhost:8080/api/v1'

export const getAllCompanies = async (): Promise<Company[]> => {
  const response = await fetch(`${API_BASE}/companies`)
  if (!response.ok) {
    throw new Error('Failed to fetch companies')
  }
  return response.json()
}

export const getCompanyById = async (id: string): Promise<Company> => {
  const response = await fetch(`${API_BASE}/companies/${id}`)
  if (!response.ok) {
    throw new Error('Failed to fetch company')
  }
  return response.json()
}

export const createCompany = async (company: Omit<Company, 'company_id'>): Promise<Company> => {
  const response = await fetch(`${API_BASE}/companies`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(company),
  })
  if (!response.ok) {
    throw new Error('Failed to create company')
  }
  return response.json()
}

export const updateCompany = async (id: string, company: Partial<Company>): Promise<Company> => {
  const response = await fetch(`${API_BASE}/companies/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(company),
  })
  if (!response.ok) {
    throw new Error('Failed to update company')
  }
  return response.json()
}

export const deleteCompany = async (id: string): Promise<void> => {
  const response = await fetch(`${API_BASE}/companies/${id}`, {
    method: 'DELETE',
  })
  if (!response.ok) {
    throw new Error('Failed to delete company')
  }
}
