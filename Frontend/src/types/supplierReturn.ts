export type ReturnStatus =
  | 'draft'
  | 'approved'
  | 'sent'
  | 'credited'
  | 'cancelled'
  | 'rejected'
  | 'completed'

export interface SupplierReturnItemSnapshot {
  supplier_return_item_id: number
  supplier_return_id: number
  product_id: string | null
  location_id: string | null
  quantity: number
  unit_cost: number
  product_name_snapshot: string
  location_snapshot: string
}

export interface SupplierReturn {
  supplier_return_id: number
  return_no: string
  company_id: string
  status: ReturnStatus
  reason?: string | null
  notes?: string | null
  created_at: string
  approved_at?: string | null
  completed_at?: string | null
  items?: SupplierReturnItemSnapshot[]
}

export interface SupplierReturnCreateRequest {
  company_id: string
  return_no: string
  reason?: string | null
  notes?: string | null
  items: SupplierReturnItemCreate[]
}

export interface SupplierReturnItemCreate {
  product_id: string
  location_id: string
  quantity: number
  unit_cost: number
}

export interface SupplierReturnItem {
  product_id: string
  quantity: number
  unit_cost: number
}
