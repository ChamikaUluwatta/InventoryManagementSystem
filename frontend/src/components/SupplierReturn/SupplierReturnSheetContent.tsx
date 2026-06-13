import { useEffect, useState } from 'react'
import type { ReturnStatus, SupplierReturn } from '@/types/supplierReturn'
import { getSupplierReturnById, updateSupplierReturnStatus } from '@/services/supplierReturnService'
import { getAllCompanies } from '@/services/companyService'
import { Button } from '@/components/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { X, Package } from 'lucide-react'
import { Spinner } from '@/components/ui/spinner'
import { Badge } from '@/components/ui/badge'
import { SectionLabel, DataCell } from '@/components/ui/sheet-label'
import { ErrorMessage } from '@/components/ui/error-message'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'

type Props = {
  supplierReturnId: number
  onClose: () => void
  onSuccess: () => void
}

const statusOptions: ReturnStatus[] = [
  'draft',
  'approved',
  'sent',
  'credited',
  'cancelled',
  'rejected',
  'completed',
]

const statusVariantMap: Record<ReturnStatus, 'default' | 'secondary' | 'success' | 'warning' | 'destructive' | 'outline'> = {
  draft: 'secondary',
  approved: 'default',
  sent: 'warning',
  credited: 'success',
  cancelled: 'destructive',
  rejected: 'destructive',
  completed: 'success',
}

function formatStatus(status: ReturnStatus) {
  return status.replace('_', ' ').toUpperCase()
}

function formatDate(value?: string | null) {
  if (!value) return '—'
  const parsed = new Date(value)
  if (Number.isNaN(parsed.getTime())) return '—'
  return parsed.toLocaleString()
}

export default function SupplierReturnSheetContent({ supplierReturnId, onClose, onSuccess }: Props) {
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [data, setData] = useState<SupplierReturn | null>(null)
  const [companyName, setCompanyName] = useState<string>('—')

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [returnData, companies] = await Promise.all([
          getSupplierReturnById(supplierReturnId),
          getAllCompanies(),
        ])
        setData(returnData)
        const company = companies.find(c => c.company_id.toString() === returnData.company_id.toString())
        setCompanyName(company?.company_name || 'Unknown')
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load supplier return')
      } finally {
        setLoading(false)
      }
    }
    fetchData()
  }, [supplierReturnId])

  async function handleStatusChange(newStatus: ReturnStatus) {
    if (!data || data.status === newStatus) return
    setSaving(true)
    try {
      const updated = await updateSupplierReturnStatus(supplierReturnId, newStatus)
      setData(updated)
      onSuccess()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update status')
    } finally {
      setSaving(false)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center gap-4 justify-center h-full">
        <Spinner className="size-12" />
        <p>Loading...</p>
      </div>
    )
  }

  if (error && !data) {
    return (
      <ErrorMessage message={error} />
    )
  }

  if (!data) return null

  const totalCost = data.items?.reduce((sum, item) => sum + item.quantity * item.unit_cost, 0) ?? 0

  return (
    <div className="flex flex-col h-full bg-background">
      <div className="flex items-center justify-between px-5 py-4 border-b shrink-0">
        <div className="flex items-center gap-3">
          <div className="flex items-center justify-center w-8 h-8 rounded border border-border">
            <Package className="h-4 w-4" />
          </div>
          <div>
            <p className="text-[10px] font-mono uppercase tracking-widest text-muted-foreground leading-none mb-0.5">
              Supplier Return
            </p>
            <h2 className="text-sm font-semibold leading-none">{data.return_no}</h2>
          </div>
        </div>
        <Button
          variant="ghost"
          size="icon"
          className="h-7 w-7 rounded-sm text-muted-foreground hover:text-foreground hover:bg-muted"
          onClick={onClose}
        >
          <X className="h-3.5 w-3.5" />
        </Button>
      </div>

      <div className="flex-1 overflow-y-auto p-5 space-y-5">
        {error && (
          <ErrorMessage message={error} />
        )}

        <div className="flex items-baseline justify-between px-4 py-3 rounded-md border border-border">
          <span className="text-[10px] font-mono uppercase tracking-widest text-muted-foreground">
            Total Value
          </span>
          <span className="text-2xl font-bold font-mono tabular-nums">
            ${totalCost.toFixed(2)}
          </span>
        </div>

        <div>
          <SectionLabel>Details</SectionLabel>
          <div className="grid grid-cols-2 border border-border rounded-md overflow-hidden">
            <DataCell label="RETURN NO" value={data.return_no} bordered />
            <DataCell label="COMPANY" value={companyName} />
            <DataCell label="CREATED" value={formatDate(data.created_at)} bordered topBorder />
            <DataCell label="COMPLETED" value={formatDate(data.completed_at)} topBorder />
          </div>
        </div>

        <div>
          <SectionLabel>Status</SectionLabel>
          <div className="flex items-center gap-3 px-4 py-3 border border-border rounded-md">
            <div className="flex justify-between w-full items-center">
              <div className="flex items-center gap-2">
                <Badge variant={statusVariantMap[data.status]}>
                  {formatStatus(data.status)}
                </Badge>
              </div>
              <Select
                value={data.status}
                onValueChange={(value) => handleStatusChange(value as ReturnStatus)}
                disabled={saving}
              >
                <SelectTrigger className="h-8 text-xs font-mono w-36">
                  <SelectValue placeholder="Change status" />
                </SelectTrigger>
                <SelectContent position="popper">
                  {statusOptions
                    .filter(s => s !== data.status)
                    .map((status) => (
                      <SelectItem key={status} value={status} className="font-mono text-xs">
                        {formatStatus(status)}
                      </SelectItem>
                    ))}
                </SelectContent>
              </Select>
            </div>
          </div>
        </div>

        {data.reason && (
          <div>
            <SectionLabel>Reason</SectionLabel>
            <div className="px-4 py-3 border border-border rounded-md">
              <p className="text-sm text-muted-foreground">{data.reason}</p>
            </div>
          </div>
        )}

        {data.notes && (
          <div>
            <SectionLabel>Notes</SectionLabel>
            <div className="px-4 py-3 border border-border rounded-md">
              <p className="text-sm text-muted-foreground">{data.notes}</p>
            </div>
          </div>
        )}

        {data.items && data.items.length > 0 && (
          <div>
            <SectionLabel>Items ({data.items.length})</SectionLabel>
            <div className="border border-border rounded-md overflow-hidden">
              <Table className="table-industrial">
                <TableHeader>
                  <TableRow className="bg-muted/50">
                    <TableHead className="text-xs">PRODUCT</TableHead>
                    <TableHead className="text-xs">LOCATION</TableHead>
                    <TableHead className="text-xs text-right">QTY</TableHead>
                    <TableHead className="text-xs text-right">UNIT COST</TableHead>
                    <TableHead className="text-xs text-right">TOTAL</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {data.items.map((item) => (
                    <TableRow key={item.supplier_return_item_id}>
                      <TableCell className="text-xs font-mono">{item.product_name_snapshot}</TableCell>
                      <TableCell className="text-xs font-mono">{item.location_snapshot || '—'}</TableCell>
                      <TableCell className="text-xs font-mono text-right">{item.quantity}</TableCell>
                      <TableCell className="text-xs font-mono text-right">${item.unit_cost}</TableCell>
                      <TableCell className="text-xs font-mono text-right">${(item.quantity * item.unit_cost).toFixed(2)}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          </div>
        )}
      </div>

      <div className="px-5 py-4 border-t shrink-0 flex gap-2 justify-end items-center">
        <Button
          type="button"
          variant="ghost"
          size="sm"
          className="h-8 px-3 text-xs text-muted-foreground"
          onClick={onClose}
        >
          Close
        </Button>
      </div>
    </div>
  )
}
