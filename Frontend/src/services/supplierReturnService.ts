import type { ReturnStatus, SupplierReturn } from '@/types/supplierReturn'
import { apiFetch } from '@/lib/api'

export const getSupplierReturnById = async (id: number): Promise<SupplierReturn> => {
  return apiFetch<SupplierReturn>(`/supplier-returns/${id}`)
}

export const getAllSupplierReturns = async (): Promise<SupplierReturn[]> => {
  return apiFetch<SupplierReturn[]>('/supplier-returns')
}

export const getSupplierReturnsByCompany = async (companyId: string): Promise<SupplierReturn[]> => {
  return apiFetch<SupplierReturn[]>(`/supplier-returns?company=${encodeURIComponent(companyId)}`)
}

export interface CreateSupplierReturnRequest {
  company_id: string
  return_no: string
  reason?: string | null
  notes?: string | null
  items: { product_id: string; location_id: string; quantity: number; unit_cost: number }[]
}

export const createSupplierReturn = async (
  data: CreateSupplierReturnRequest,
): Promise<SupplierReturn> => {
  return apiFetch<SupplierReturn>('/supplier-returns', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
}

export const updateSupplierReturnStatus = async (
  id: number,
  status: ReturnStatus,
): Promise<SupplierReturn> => {
  return apiFetch<SupplierReturn>(`/supplier-returns/${id}/status`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ status }),
  })
}
