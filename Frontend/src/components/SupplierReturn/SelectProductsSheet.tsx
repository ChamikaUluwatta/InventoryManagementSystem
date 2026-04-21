import { useEffect, useMemo, useState } from 'react'
import { X, Package, Check, ChevronLeft } from 'lucide-react'
import type { Company } from '@/types/company'
import type { Product } from '@/types/product'
import { getProductsByCompany } from '@/services/productService'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Checkbox } from '@/components/ui/checkbox'
import { Spinner } from '../ui/spinner'

type SelectedItem = {
  product: Product
  quantity: number
  unit_cost: number
}

type Props = {
  open: boolean
  company: Company
  initialItems?: SelectedItem[]
  onBack: () => void
  onConfirm: (items: SelectedItem[]) => void
}

export default function SelectProductsSheet({ open, company, initialItems = [], onBack, onConfirm }: Props) {
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [selectedProductIds, setSelectedProductIds] = useState<Set<string>>(new Set())
  const [quantities, setQuantities] = useState<Record<string, number>>({})
  const [unitCosts, setUnitCosts] = useState<Record<string, number>>({})

  const initiallyAddedProductIds = useMemo(
    () => new Set(initialItems.map((item) => item.product.product_id)),
    [initialItems],
  )

  useEffect(() => {
    if (!open) {
      setSelectedProductIds(new Set())
      return
    }

    const fetchProducts = async () => {
      setLoading(true)
      setError(null)
      try {
        const data = await getProductsByCompany(company.company_id.toString())
        setProducts(data)
        const initialQuantities: Record<string, number> = {}
        const initialUnitCosts: Record<string, number> = {}

        data.forEach(p => {
          initialQuantities[p.product_id] = 1
          initialUnitCosts[p.product_id] = p.price || 0
        })

        setSelectedProductIds(new Set())
        setQuantities(initialQuantities)
        setUnitCosts(initialUnitCosts)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load products')
      } finally {
        setLoading(false)
      }
    }
    fetchProducts()
  }, [open, company.company_id])

  const handleToggleProduct = (productId: string) => {
    const newSelected = new Set(selectedProductIds)
    if (newSelected.has(productId)) {
      newSelected.delete(productId)
    } else {
      newSelected.add(productId)
    }
    setSelectedProductIds(newSelected)
  }

  const handleQuantityChange = (productId: string, value: string) => {
    const qty = parseInt(value) || 0
    setQuantities(prev => ({ ...prev, [productId]: qty }))
  }

  const handleUnitCostChange = (productId: string, value: string) => {
    const cost = parseFloat(value) || 0
    setUnitCosts(prev => ({ ...prev, [productId]: cost }))
  }

  const handleConfirm = () => {
    const selectedItems: SelectedItem[] = []
    selectedProductIds.forEach(id => {
      const product = products.find(p => p.product_id === id)
      if (product) {
        selectedItems.push({
          product,
          quantity: quantities[id] || 1,
          unit_cost: unitCosts[id] || 0
        })
      }
    })
    onConfirm(selectedItems)
  }

  const allSelected = products.length > 0 && selectedProductIds.size === products.length

  return (
    <div className="flex flex-col h-full bg-background">
      <div className="flex items-center justify-between px-5 py-4 border-b shrink-0">
        <div className="flex items-center gap-3">
          <Button
            variant="ghost"
            size="icon"
            className="h-7 w-7"
            onClick={onBack}
          >
            <ChevronLeft className="h-4 w-4" />
          </Button>
          <div className="flex items-center justify-center w-8 h-8 rounded border border-border">
            <Package className="h-4 w-4" />
          </div>
          <div>
            <p className="text-[10px] font-mono uppercase tracking-widest text-muted-foreground leading-none mb-0.5">
              Step 2 of 2
            </p>
            <h2 className="text-sm font-semibold leading-none">Select Products</h2>
          </div>
        </div>
        <Button
          variant="ghost"
          size="icon"
          className="h-7 w-7 rounded-sm text-muted-foreground hover:text-foreground hover:bg-muted"
          onClick={onBack}
        >
          <X className="h-3.5 w-3.5" />
        </Button>
      </div>

      <div className="flex-1 overflow-y-auto">
        {loading ? (
          <div className="flex items-center gap-4 justify-center h-40">
            <Spinner className="size-8" />
            <p>Loading products for {company.company_name}...</p>
          </div>
        ) : error ? (
          <div className="m-4 p-3 border border-red-200 bg-red-50 dark:bg-red-950/20 dark:border-red-900 rounded text-sm text-red-600 dark:text-red-400 font-mono">
            ERR: {error}
          </div>
        ) : products.length === 0 ? (
          <div className="p-8 text-center text-muted-foreground">
            <Package className="h-12 w-12 mx-auto mb-4 opacity-50" />
            <p>No products found for this company</p>
          </div>
        ) : (
          <>
            <div className="p-4 border-b border-border bg-muted/30">
              <p className="text-[10px] font-mono uppercase tracking-widest text-muted-foreground mb-1">
                Company
              </p>
              <p className="text-sm font-semibold">{company.company_name}</p>
            </div>

            <div className="overflow-visible">
              <Table className="table-industrial">
                <TableHeader>
                  <TableRow className="bg-muted/50">
                  <TableHead className="w-10">
                    <Checkbox
                      checked={allSelected}
                      onCheckedChange={(checked) => {
                        if (checked) {
                          setSelectedProductIds(new Set(products.map(p => p.product_id)))
                        } else {
                          setSelectedProductIds(new Set())
                        }
                      }}
                    />
                  </TableHead>
                  <TableHead className="flex-1">PRODUCT</TableHead>
                  <TableHead className="w-16 text-right">STOCK</TableHead>
                  <TableHead className="w-20 text-right">PRICE</TableHead>
                  <TableHead className="w-24 text-right">QTY</TableHead>
                  <TableHead className="w-28 text-right">UNIT COST</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {products.map((product) => {
                  const isSelected = selectedProductIds.has(product.product_id)
                  return (
                    <TableRow
                      key={product.product_id}
                      className={isSelected ? 'bg-muted/30' : undefined}
                    >
                      <TableCell>
                        <Checkbox
                          checked={isSelected}
                          onCheckedChange={() => handleToggleProduct(product.product_id)}
                        />
                      </TableCell>
                      <TableCell className="font-mono text-xs">
                        <div className="flex items-center gap-2">
                          <span>{product.product_name}</span>
                          {initiallyAddedProductIds.has(product.product_id) ? (
                            <span className="text-[10px] uppercase tracking-wide text-muted-foreground">Added</span>
                          ) : null}
                        </div>
                      </TableCell>
                      <TableCell className="text-right font-mono text-xs">
                        {product.stock ?? 0}
                      </TableCell>
                      <TableCell className="text-right font-mono text-xs">
                        ${typeof product.price === 'number' ? product.price.toFixed(2) : Number(product.price).toFixed(2)}
                      </TableCell>
                      <TableCell>
                        <Input
                          type="number"
                          min={1}
                          className="h-8 text-xs font-mono text-right"
                          value={quantities[product.product_id] || 1}
                          onChange={(e) => handleQuantityChange(product.product_id, e.target.value)}
                          disabled={!isSelected}
                        />
                      </TableCell>
                      <TableCell>
                        <Input
                          type="number"
                          step="0.01"
                          min={0}
                          className="h-8 text-xs font-mono text-right"
                          value={unitCosts[product.product_id] || 0}
                          onChange={(e) => handleUnitCostChange(product.product_id, e.target.value)}
                          disabled={!isSelected}
                        />
                      </TableCell>
                    </TableRow>
                  )
                })}
              </TableBody>
              </Table>
            </div>
          </>
        )}
      </div>

      <div className="px-5 py-4 border-t shrink-0 flex items-center justify-between">
        <div className="text-xs text-muted-foreground font-mono">
          {selectedProductIds.size} of {products.length} selected
        </div>
        <div className="flex gap-2">
          <Button
            type="button"
            variant="ghost"
            size="sm"
            className="h-8 px-3 text-xs text-muted-foreground"
            onClick={onBack}
          >
            Back
          </Button>
          <Button
            type="button"
            size="sm"
            className="h-8 px-4 text-xs"
            disabled={selectedProductIds.size === 0}
            onClick={handleConfirm}
          >
            <Check className="h-3 w-3 mr-1" />
            Add Selected
          </Button>
        </div>
      </div>
    </div>
  )
}