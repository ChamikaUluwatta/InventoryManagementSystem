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
import type { Product } from '@/types/product'
import { getAllProducts } from '@/services/productService'
import type { Category } from '@/types/category'
import { getAllCategories } from '@/services/categoryService'
import { Spinner } from '@/components/ui/spinner'
import { Plus, Search } from 'lucide-react'
import { Sheet, SheetContent } from '@/components/ui/sheet'
import ProductSheetContent from '@/components/Product/ProductSheetContent'
import AddProductSheetContent from '@/components/Product/AddProductSheetContent'

export default function ManageProducts() {
  const [products, setProducts] = useState<Product[]>([])
  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [sorting, setSorting] = useState<SortingState>([])
  const [globalFilter, setGlobalFilter] = useState('')

  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null)
  const [sheetOpen, setSheetOpen] = useState(false)
  const [addSheetOpen, setAddSheetOpen] = useState(false)

  useEffect(() => {
    const fetchProducts = async () => {
      try {
        const [productData, categoriesData] = await Promise.all([
          getAllProducts(),
          getAllCategories(),
        ])
        setProducts(productData)
        setCategories(categoriesData)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to fetch products')
      } finally {
        setLoading(false)
      }
    }

    fetchProducts()
  }, [])

  const handleRowClick = (product: Product) => {
    setSelectedProduct(product)
    setSheetOpen(true)
  }

  const handleAddSuccess = () => {
    getAllProducts().then(setProducts).catch(console.error)
  }

  const handleSuccess = () => {
    getAllProducts().then(setProducts).catch(console.error)
  }

  const handleClose = () => {
    setSheetOpen(false)
    setSelectedProduct(null)
  }

  const columns: ColumnDef<Product>[] = [
    {
      accessorKey: 'product_name',
      header: 'PRODUCT',
      cell: ({ row }) => <span className="font-mono">{row.getValue('product_name')}</span>,
    },
    {
      accessorKey: 'diameter',
      header: 'DIAMETER',
      cell: ({ row }) => <span className="font-data">{row.getValue('diameter')}</span>,
    },
    {
      accessorKey: 'width',
      header: 'WIDTH',
      cell: ({ row }) => <span className="font-data">{row.getValue('width')}</span>,
    },
    {
      accessorKey: 'price',
      header: 'PRICE',
      cell: ({ row }) => <span className="font-data">${row.getValue('price')}</span>,
    },
    {
      accessorKey: 'category_id',
      header: 'CATEGORY',
      cell: ({ row }) => {
        const catId = row.getValue('category_id') as number
        const category = categories.find((c) => c.category_id === catId)
        return category?.category_name || '-'
      },
    },
    {
      accessorKey: 'location_id',
      header: 'LOCATION',
      cell: ({ row }) => <span className="font-data">{row.getValue('location_id') || '-'}</span>,
    },
  ]

  const table = useReactTable({
    data: products,
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
            ADD PRODUCT
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
                  No products found.
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
            products.length,
          )}{' '}
          OF {products.length}
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
      
      <Sheet open={sheetOpen} onOpenChange={(open) => {
        if (!open) handleClose()
      }}>
        <SheetContent className="w-100 sm:w-125 md:w-150 lg:w-175 xl:w-200 max-w-[90vw]">
          {selectedProduct && (
            <ProductSheetContent
              product={selectedProduct}
              onClose={handleClose}
              onSuccess={handleSuccess}
            />
          )}
        </SheetContent>
      </Sheet>

      <Sheet open={addSheetOpen} onOpenChange={(open) => {
        if (!open) setAddSheetOpen(false)
      }}>
        <SheetContent className="w-100 sm:w-125 md:w-150 lg:w-175 xl:w-200 max-w-[90vw]">
          <AddProductSheetContent
            onClose={() => setAddSheetOpen(false)}
            onSuccess={handleAddSuccess}
          />
        </SheetContent>
      </Sheet>
    </div>
  )
}