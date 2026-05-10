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
import { Search, Plus, Package, X, Trash2 } from 'lucide-react'
import { Sheet, SheetContent } from '@/components/ui/sheet'
import { ErrorMessage } from '@/components/ui/error-message'
import { getAllCompanies } from '@/services/companyService'
import { getAllSupplierReturns, createSupplierReturn } from '@/services/supplierReturnService'
import type { ReturnStatus, SupplierReturn } from '@/types/supplierReturn'
import type { Company } from '@/types/company'
import type { Product } from '@/types/product'
import SelectCompanySheet from '@/components/SupplierReturn/SelectCompanySheet'
import SelectProductsSheet from '@/components/SupplierReturn/SelectProductsSheet'
import SupplierReturnSheetContent from '@/components/SupplierReturn/SupplierReturnSheetContent'

type SupplierReturnView = SupplierReturn & {
  company_name: string
}

type SelectedReturnItem = {
  product: Product
  quantity: number
  unit_cost: number
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

function mergeSelectedItems(existingItems: SelectedReturnItem[], newItems: SelectedReturnItem[]) {
  const mergedItems = existingItems.map((item) => ({ ...item }))

  for (const newItem of newItems) {
    if (newItem.quantity <= 0) {
      continue
    }

    const existingIndex = mergedItems.findIndex(
      (item) =>
        item.product.product_id === newItem.product.product_id &&
        item.unit_cost === newItem.unit_cost,
    )

    if (existingIndex >= 0) {
      mergedItems[existingIndex] = {
        ...mergedItems[existingIndex],
        quantity: mergedItems[existingIndex].quantity + newItem.quantity,
      }
      continue
    }

    mergedItems.push({ ...newItem })
  }

  return mergedItems
}

export default function SupplierReturns() {
  const [returns, setReturns] = useState<SupplierReturn[]>([])
  const [companyNameByID, setCompanyNameByID] = useState<Map<string, string>>(new Map())
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [sorting, setSorting] = useState<SortingState>([])
  const [globalFilter, setGlobalFilter] = useState('')
  const [statusFilter, setStatusFilter] = useState<'all' | ReturnStatus>('all')
  const [addSheetOpen, setAddSheetOpen] = useState(false)
  const [selectedReturnId, setSelectedReturnId] = useState<number | null>(null)
  const [selectCompanySheetOpen, setSelectCompanySheetOpen] = useState(false)
  const [selectProductsSheetOpen, setSelectProductsSheetOpen] = useState(false)
  const [selectedCompany, setSelectedCompany] = useState<Company | null>(null)
  const [selectedItems, setSelectedItems] = useState<SelectedReturnItem[]>([])
  const [returnNo, setReturnNo] = useState('')
  const [reason, setReason] = useState('')
  const [notes, setNotes] = useState('')

  const resetDraft = () => {
    setAddSheetOpen(false)
    setSelectCompanySheetOpen(false)
    setSelectProductsSheetOpen(false)
    setSelectedCompany(null)
    setSelectedItems([])
    setReturnNo('')
    setReason('')
    setNotes('')
  }

  const generateReturnNo = () => {
    const now = new Date()
    const year = now.getFullYear()
    const month = String(now.getMonth() + 1).padStart(2, '0')
    const random = Math.floor(Math.random() * 1000).toString().padStart(3, '0')
    return `SR-${year}${month}-${random}`
  }

  const startNewReturn = () => {
    resetDraft()
    setReturnNo(generateReturnNo())
    setSelectCompanySheetOpen(true)
  }

  const handleAddSuccess = () => {
    getAllSupplierReturns().then(setReturns).catch(console.error)
  }

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

  if (error) return (
    <div className="flex items-center justify-center h-full">
      <ErrorMessage message={error} className="max-w-md text-center" />
    </div>
  )

  const pageIndex = table.getState().pagination.pageIndex
  const pageSize = table.getState().pagination.pageSize
  const totalRows = table.getFilteredRowModel().rows.length
  const startRow = totalRows === 0 ? 0 : pageIndex * pageSize + 1
  const endRow = totalRows === 0 ? 0 : Math.min((pageIndex + 1) * pageSize, totalRows)

  return (
    <div className="h-full min-w-0 flex flex-col">
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

          <Button
            variant="outline"
            size="sm"
            className="gap-2 font-mono text-xs"
            onClick={startNewReturn}
          >
            <Plus className="h-4 w-4" />
            NEW RETURN
          </Button>
          <Sheet open={selectCompanySheetOpen} onOpenChange={(open) => {
            setSelectCompanySheetOpen(open)
            if (!open) {
              setSelectedCompany(null)
              setSelectedItems([])
            }
          }}>
            <SheetContent side="right" className="sm:max-w-lg w-full">
              <SelectCompanySheet
                onClose={() => setSelectCompanySheetOpen(false)}
                onSelectCompany={(company) => {
                  setSelectedCompany(company)
                  setSelectCompanySheetOpen(false)
                  setSelectProductsSheetOpen(true)
                }}
              />
            </SheetContent>
          </Sheet>
          <Sheet open={selectProductsSheetOpen} onOpenChange={setSelectProductsSheetOpen}>
            <SheetContent side="left" size="3xl" className="w-full">
              {selectedCompany && (
                <SelectProductsSheet
                  open={selectProductsSheetOpen}
                  company={selectedCompany}
                  initialItems={selectedItems}
                  onBack={() => {
                    setSelectProductsSheetOpen(false)
                    setSelectCompanySheetOpen(true)
                  }}
                  onConfirm={(items) => {
                    setSelectedItems((previousItems) => mergeSelectedItems(previousItems, items))
                    setSelectProductsSheetOpen(false)
                    setAddSheetOpen(true)
                  }}
                />
              )}
            </SheetContent>
          </Sheet>
          <Sheet
            open={addSheetOpen}
            onOpenChange={(open) => {
              if (!open) {
                resetDraft()
                return
              }
              setAddSheetOpen(open)
            }}
          >
            <SheetContent size='md' side="right" className="w-full">
              <div className="flex flex-col h-full bg-background">
                <div className="flex items-center justify-between px-5 py-4 border-b shrink-0">
                  <div className="flex items-center gap-3">
                    <div className="flex items-center justify-center w-8 h-8 rounded border border-border">
                      <Package className="h-4 w-4" />
                    </div>
                    <div>
                      <p className="text-[10px] font-mono uppercase tracking-widest text-muted-foreground leading-none mb-0.5">
                        Creating
                      </p>
                      <h2 className="text-sm font-semibold leading-none">New Supplier Return</h2>
                    </div>
                  </div>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-7 w-7 rounded-sm text-muted-foreground hover:text-foreground hover:bg-muted"
                    onClick={resetDraft}
                  >
                    <X className="h-3.5 w-3.5" />
                  </Button>
                </div>

                <div className="flex-1 overflow-y-auto p-5 space-y-5">
                  <div className="flex items-baseline justify-between px-4 py-3 rounded-md border border-border">
                    <span className="text-[10px] font-mono uppercase tracking-widest text-muted-foreground">
                      Return No
                    </span>
                    <Input
                      id="return_no"
                      className="w-32 h-8 text-sm font-mono text-right"
                      placeholder="SR-000"
                      value={returnNo}
                      onChange={(e) => setReturnNo(e.target.value)}
                    />
                  </div>

                  {selectedCompany && (
                    <div>
                      <p className="text-[10px] font-mono uppercase tracking-widest text-muted-foreground mb-2">
                        Company
                      </p>
                      <div className="px-4 py-3 border border-border rounded-md">
                        <p className="text-sm font-semibold">{selectedCompany.company_name}</p>
                      </div>
                    </div>
                  )}

                  {selectedItems.length > 0 && (
                    <div>
                      <p className="text-[10px] font-mono uppercase tracking-widest text-muted-foreground mb-2">
                        Return Items ({selectedItems.length})
                      </p>
                      <div className="border border-border rounded-md overflow-hidden">
                        <Table>
                          <TableHeader>
                            <TableRow className="bg-muted/50">
                              <TableHead className="text-xs">PRODUCT</TableHead>
                              <TableHead className="text-xs text-right">QTY</TableHead>
                              <TableHead className="text-xs text-right">UNIT COST</TableHead>
                              <TableHead className="text-xs text-right">TOTAL</TableHead>
                              <TableHead className="text-xs text-right">ACTION</TableHead>
                            </TableRow>
                          </TableHeader>
                          <TableBody>
                            {selectedItems.map((item, idx) => (
                              <TableRow key={idx}>
                                <TableCell className="text-xs font-mono">{item.product.product_name}</TableCell>
                                <TableCell className="text-xs font-mono text-right">{item.quantity}</TableCell>
                                <TableCell className="text-xs font-mono text-right">${item.unit_cost}</TableCell>
                                <TableCell className="text-xs font-mono text-right">${(item.quantity * item.unit_cost).toFixed(2)}</TableCell>
                                <TableCell className="text-right">
                                  <Button
                                    type="button"
                                    variant="ghost"
                                    size="icon"
                                    className="h-7 w-7 text-muted-foreground hover:text-destructive"
                                    onClick={() => {
                                      setSelectedItems((previousItems) => previousItems.filter((_, itemIndex) => itemIndex !== idx))
                                    }}
                                  >
                                    <Trash2 className="h-4 w-4" />
                                  </Button>
                                </TableCell>
                              </TableRow>
                            ))}
                          </TableBody>
                        </Table>
                      </div>
                    </div>
                  )}

                  <Button
                    variant="outline"
                    size="sm"
                    className="w-full text-xs"
                    onClick={() => {
                      setSelectProductsSheetOpen(true)
                    }}
                  >
                    <Plus className="h-3 w-3 mr-1" />
                    Add More Items
                  </Button>
                </div>

                <div className="px-5 py-4 border-t shrink-0 flex gap-2 justify-end items-center">
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="h-8 px-3 text-xs text-muted-foreground"
                    onClick={resetDraft}
                  >
                    Cancel
                  </Button>
                  <Button
                    type="button"
                    size="sm"
                    className="h-8 px-4 text-xs"
                    disabled={!returnNo || !selectedCompany || selectedItems.length === 0}
                    onClick={async () => {
                      try {
                        await createSupplierReturn({
                          company_id: selectedCompany!.company_id.toString(),
                          return_no: returnNo,
                          reason: reason || null,
                          notes: notes || null,
                          items: selectedItems.map(item => ({
                            product_id: item.product.product_id,
                            location_id: item.product.location_id,
                            quantity: item.quantity,
                            unit_cost: item.unit_cost,
                          })),
                        })
                        handleAddSuccess()
                        resetDraft()
                      } catch (err) {
                        console.error('Failed to create supplier return:', err)
                      }
                    }}
                  >
                    Create Return
                  </Button>
                </div>
              </div>
            </SheetContent>
          </Sheet>
        </div>
      </div>

      <div className="min-w-0 flex-1 overflow-y-auto">
        <Table className="table-industrial table-fixed [&_th]:whitespace-normal [&_th]:wrap-break-word [&_td]:whitespace-normal [&_td]:wrap-break-word">
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
                <TableRow
                  key={row.id}
                  className="cursor-pointer"
                  onClick={() => setSelectedReturnId(row.original.supplier_return_id)}
                >
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

      <Sheet open={selectedReturnId !== null} onOpenChange={(open) => {
        if (!open) setSelectedReturnId(null)
      }}>
        <SheetContent className="w-100 sm:w-125 md:w-150 lg:w-175 xl:w-200 max-w-[90vw]">
          {selectedReturnId !== null && (
            <SupplierReturnSheetContent
              supplierReturnId={selectedReturnId}
              onClose={() => setSelectedReturnId(null)}
              onSuccess={() => {
                getAllSupplierReturns().then(setReturns).catch(console.error)
              }}
            />
          )}
        </SheetContent>
      </Sheet>
    </div>
  )
}
