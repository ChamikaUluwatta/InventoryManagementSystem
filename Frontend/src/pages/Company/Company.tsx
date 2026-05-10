import { useEffect, useState } from 'react'
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
import type { Company } from '@/types/company'
import { getAllCompanies } from '@/services/companyService'
import { Spinner } from '@/components/ui/spinner'
import { Sheet, SheetContent } from '@/components/ui/sheet'
import { ErrorMessage } from '@/components/ui/error-message'
import AddCompanySheetContent from '@/components/Company/AddCompanySheetContent'
import CompanySheetContent from '@/components/Company/CompanySheetContent'
import { Search, Plus } from 'lucide-react'

const columns: ColumnDef<Company>[] = [
  {
    accessorKey: 'company_name',
    header: 'COMPANY NAME',
    cell: ({ row }) => <span className="font-medium">{row.getValue('company_name')}</span>,
  },
  {
    accessorKey: 'description',
    header: 'DESCRIPTION',
    cell: ({ row }) => {
      const desc = row.getValue('description') as string | null | undefined
      if (!desc) {
        return <span className="italic text-muted-foreground/50 text-xs font-mono">No description</span>
      }
      return <span className="font-mono text-xs text-muted-foreground">{desc}</span>
    },
  },
]

export default function Company() {
  const [companies, setCompanies] = useState<Company[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [sorting, setSorting] = useState<SortingState>([])
  const [globalFilter, setGlobalFilter] = useState('')
  const [addSheetOpen, setAddSheetOpen] = useState(false)
  const [selectedCompany, setSelectedCompany] = useState<Company | null>(null)
  const [detailSheetOpen, setDetailSheetOpen] = useState(false)

  useEffect(() => {
    const fetchCompanies = async () => {
      try {
        const data = await getAllCompanies()
        setCompanies(data)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to fetch companies')
      } finally {
        setLoading(false)
      }
    }
    fetchCompanies()
  }, [])

  const handleAddSuccess = () => {
    getAllCompanies().then(setCompanies).catch(console.error)
  }

  const handleRowClick = (company: Company) => {
    setSelectedCompany(company)
    setDetailSheetOpen(true)
  }

  const handleDetailClose = () => {
    setDetailSheetOpen(false)
    setSelectedCompany(null)
  }

  const handleDetailSuccess = () => {
    getAllCompanies().then(setCompanies).catch(console.error)
  }

  const table = useReactTable({
    data: companies,
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
          <Button variant="outline" size="sm" className="gap-2 font-mono text-xs" onClick={() => setAddSheetOpen(true)}>
            <Plus className="h-4 w-4" />
            ADD
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
                  No companies found.
                </TableCell>
              </TableRow>
            ) : (
              table.getRowModel().rows.map((row) => (
                <TableRow
                  key={row.id}
                  className="cursor-pointer"
                  onClick={() => handleRowClick(row.original)}
                >
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id}>
                      {flexRender(cell.column.columnDef.cell, cell.getContext())}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      <div className="border-t border-border p-3 flex items-center justify-between shrink-0 text-xs text-muted-foreground font-mono">
        <div>
          SHOWING {table.getState().pagination.pageIndex * table.getState().pagination.pageSize + 1}-
          {Math.min(
            (table.getState().pagination.pageIndex + 1) * table.getState().pagination.pageSize,
            companies.length,
          )}{' '}
          OF {companies.length}
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

      <Sheet open={addSheetOpen} onOpenChange={(open) => {
        if (!open) setAddSheetOpen(false)
      }}>
        <SheetContent className="w-100 sm:w-125 md:w-150 lg:w-175 xl:w-200 max-w-[90vw]">
          <AddCompanySheetContent
            onClose={() => setAddSheetOpen(false)}
            onSuccess={handleAddSuccess}
          />
        </SheetContent>
      </Sheet>

      <Sheet open={detailSheetOpen} onOpenChange={(open) => {
        if (!open) handleDetailClose()
      }}>
        <SheetContent className="w-100 sm:w-125 md:w-150 lg:w-175 xl:w-200 max-w-[90vw]">
          {selectedCompany && (
            <CompanySheetContent
              company={selectedCompany}
              onClose={handleDetailClose}
              onSuccess={handleDetailSuccess}
            />
          )}
        </SheetContent>
      </Sheet>
    </div>
  )
}
