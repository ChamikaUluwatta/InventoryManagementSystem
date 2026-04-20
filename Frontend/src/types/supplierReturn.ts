export type ReturnStatus =
  | 'draft'
  | 'approved'
  | 'sent'
  | 'credited'
  | 'cancelled'
  | 'rejected'
  | 'completed'

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
}

