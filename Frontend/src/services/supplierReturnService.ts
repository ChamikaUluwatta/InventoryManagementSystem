import type { ReturnStatus, SupplierReturn } from '@/types/supplierReturn'

const API_BASE_URL = import.meta.env.VITE_API_URL

if (!API_BASE_URL && import.meta.env.MODE === 'production') {
  throw new Error('VITE_API_URL environment variable is not set')
}

const API_BASE = API_BASE_URL || 'http://localhost:8080/api/v1'

export const getAllSupplierReturns = async (): Promise<SupplierReturn[]> => {
  const response = await fetch(`${API_BASE}/supplier-returns`)
  if (!response.ok) {
    throw new Error('Failed to fetch supplier returns')
  }
  return response.json()
}

export const createSupplierReturn = async (
  data: Omit<SupplierReturn, 'supplier_return_id' | 'created_at' | 'approved_at' | 'completed_at'>,
): Promise<SupplierReturn> => {
  const response = await fetch(`${API_BASE}/supplier-returns`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  if (!response.ok) {
    throw new Error('Failed to create supplier return')
  }
  return response.json()
}

export const updateSupplierReturnStatus = async (
  id: number,
  status: ReturnStatus,
): Promise<SupplierReturn> => {
  const response = await fetch(`${API_BASE}/supplier-returns/${id}/status`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ status }),
  })
  if (!response.ok) {
    throw new Error('Failed to update supplier return status')
  }
  return response.json()
}
