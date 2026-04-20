import { useEffect, useMemo, useState } from 'react'
import {
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  useReactTable,
} from '@tanstack/react-table'
import type { ColumnDef, SortingState } from '@tanstack/react-table'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'
import { Spinner } from '@/components/ui/spinner'
import { Search, Plus } from 'lucide-react'
import { getAllCompanies } from '@/services/companyService'
import { getAllSupplierReturns } from '@/services/supplierReturnService'
import type { ReturnStatus, SupplierReturn } from '@/types/supplierReturn'

type SupplierReturnView = SupplierReturn & {
  company_name: string
}

const statusOptions: Array<'all' | ReturnStatus> = [
  'all',
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

function formatDate(value?: string | null) {
  if (!value) return '—'

  const parsed = new Date(value)
  if (Number.isNaN(parsed.getTime())) return '—'

  return parsed.toLocaleString()
}

function formatStatus(status: ReturnStatus) {
  return status.replace('_', ' ').toUpperCase()
}

const columns: ColumnDef<SupplierReturnView>[] = [
  {
    accessorKey: 'return_no',
    header: 'RETURN NO',
    cell: ({ row }) => <span className="font-mono">{row.getValue('return_no')}</span>,
  },
  {
    accessorKey: 'company_name',
    header: 'COMPANY',
    cell: ({ row }) => <span className="font-mono">{row.getValue('company_name')}</span>,
  },
  {
    accessorKey: 'status',
    header: 'STATUS',
    cell: ({ row }) => {
      const status = row.getValue('status') as ReturnStatus
      return <Badge variant={statusVariantMap[status]}>{formatStatus(status)}</Badge>
    },
  },
  {
    accessorKey: 'created_at',
    header: 'CREATED',
    cell: ({ row }) => (
      <span className="font-mono text-xs text-muted-foreground">{formatDate(row.getValue('created_at'))}</span>
    ),
  },
  {
    accessorKey: 'completed_at',
    header: 'COMPLETED',
    cell: ({ row }) => (
      <span className="font-mono text-xs text-muted-foreground">{formatDate(row.getValue('completed_at'))}</span>
    ),
  },
]

export default function SupplierReturns() {
  const [returns, setReturns] = useState<SupplierReturn[]>([])
  const [companyNameByID, setCompanyNameByID] = useState<Map<string, string>>(new Map())
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [sorting, setSorting] = useState<SortingState>([])
  const [globalFilter, setGlobalFilter] = useState('')
  const [statusFilter, setStatusFilter] = useState<'all' | ReturnStatus>('all')

  useEffect(() => {
    const fetchReturns = async () => {
      try {
        const [returnData, companyData] = await Promise.all([
          getAllSupplierReturns(),
          getAllCompanies(),
        ])

        setReturns(returnData)
        setCompanyNameByID(new Map(companyData.map((company) => [company.company_id, company.company_name])))
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to fetch supplier returns and companies')
      } finally {
        setLoading(false)
      }
    }

    fetchReturns()
  }, [])

  const filteredReturns = useMemo(() => {
    if (statusFilter === 'all') return returns
    return returns.filter((item) => item.status === statusFilter)
  }, [returns, statusFilter])

  const tableData = useMemo<SupplierReturnView[]>(() => {
    return filteredReturns.map((item) => ({
      ...item,
      company_name: companyNameByID.get(item.company_id) ?? 'Unknown Company',
    }))
  }, [filteredReturns, companyNameByID])

  const table = useReactTable({
    data: tableData,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    onSortingChange: setSorting,
    onGlobalFilterChange: setGlobalFilter,
    state: {
      sorting,
      globalFilter,
    },
  })

  if (loading)
    return (
      <div className="flex items-center gap-4 justify-center h-full">
        <Spinner className="size-12" />
        <p>Loading...</p>
      </div>
    )

  if (error) return <div className="p-4 text-red-500">Error: {error}</div>

  const pageIndex = table.getState().pagination.pageIndex
  const pageSize = table.getState().pagination.pageSize
  const totalRows = table.getFilteredRowModel().rows.length
  const startRow = totalRows === 0 ? 0 : pageIndex * pageSize + 1
  const endRow = totalRows === 0 ? 0 : Math.min((pageIndex + 1) * pageSize, totalRows)

  return (
    <div className="h-full flex flex-col">
      <div className="border-b border-border p-4 flex items-center justify-end shrink-0">
        <div className="flex items-center gap-3">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="SEARCH..."
              value={globalFilter ?? ''}
              onChange={(e) => setGlobalFilter(e.target.value)}
              className="pl-9 w-48 font-mono text-xs uppercase"
            />
          </div>

          <Select value={statusFilter} onValueChange={(value) => setStatusFilter(value as 'all' | ReturnStatus)}>
            <SelectTrigger className="w-40 h-9 font-mono text-xs">
              <SelectValue placeholder="FILTER STATUS" />
            </SelectTrigger>
            <SelectContent>
              {statusOptions.map((status) => (
                <SelectItem key={status} value={status}>
                  {status === 'all' ? 'ALL STATUSES' : formatStatus(status)}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <Button variant="outline" size="sm" className="gap-2 font-mono text-xs" disabled>
            <Plus className="h-4 w-4" />
            NEW RETURN
          </Button>
        </div>
      </div>

      <div className="flex-1 overflow-auto">
        <Table className="table-industrial">
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => (
                  <TableHead
                    key={header.id}
                    className="cursor-pointer select-none"
                    onClick={header.column.getToggleSortingHandler()}
                  >
                    <div className="flex items-center gap-2">
                      {flexRender(header.column.columnDef.header, header.getContext())}
                      {{
                        asc: '↑',
                        desc: '↓',
                      }[header.column.getIsSorted() as string] ?? null}
                    </div>
                  </TableHead>
                ))}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows.length === 0 ? (
              <TableRow>
                <TableCell colSpan={columns.length} className="h-24 text-center text-muted-foreground">
                  No supplier returns found.
                </TableCell>
              </TableRow>
            ) : (
              table.getRowModel().rows.map((row) => (
                <TableRow key={row.id}>
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id}>{flexRender(cell.column.columnDef.cell, cell.getContext())}</TableCell>
                  ))}
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      <div className="border-t border-border p-3 flex items-center justify-between shrink-0 text-xs text-muted-foreground font-mono">
        <div>
          SHOWING {startRow}-{endRow} OF {totalRows}
        </div>
        <div className="flex items-center gap-2">
          <Select
            value={table.getState().pagination.pageSize.toString()}
            onValueChange={(value) => table.setPageSize(Number(value))}
          >
            <SelectTrigger className="w-16 h-8">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {[10, 25, 50].map((size) => (
                <SelectItem key={size} value={size.toString()}>
                  {size}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <Button
            variant="outline"
            size="sm"
            className="h-8 font-mono text-xs"
            onClick={() => table.previousPage()}
            disabled={!table.getCanPreviousPage()}
          >
            PREV
          </Button>
          <Button
            variant="outline"
            size="sm"
            className="h-8 font-mono text-xs"
            onClick={() => table.nextPage()}
            disabled={!table.getCanNextPage()}
          >
            NEXT
          </Button>
        </div>
      </div>
    </div>
  )
}
