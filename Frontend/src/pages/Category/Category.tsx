import { useEffect, useState, useMemo, useCallback } from 'react'
import type { Category } from '@/types/category'
import { getAllCategories } from '@/services/categoryService'
import { Spinner } from '@/components/ui/spinner'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Sheet, SheetContent } from '@/components/ui/sheet'
import AddCategorySheetContent from '@/components/Category/AddCategorySheetContent'
import CategorySheetContent from '@/components/Category/CategorySheetContent'
import { IconPlus, IconMinus } from '@tabler/icons-react'
import { Search, Plus } from 'lucide-react'
import { ErrorMessage } from '@/components/ui/error-message'

interface CategoryTreeNode extends Category {
  children: CategoryTreeNode[]
}

function buildTree(categories: Category[]): CategoryTreeNode[] {
  const map = new Map<number, CategoryTreeNode>()
  const roots: CategoryTreeNode[] = []

  for (const cat of categories) {
    map.set(cat.category_id, { ...cat, children: [] })
  }

  for (const node of map.values()) {
    if (node.parent_id !== undefined && map.has(node.parent_id)) {
      map.get(node.parent_id)!.children.push(node)
    } else {
      roots.push(node)
    }
  }

  const sortByName = (nodes: CategoryTreeNode[]) => {
    nodes.sort((a, b) => a.category_name.localeCompare(b.category_name))
    for (const node of nodes) sortByName(node.children)
  }
  sortByName(roots)

  return roots
}

function filterTree(nodes: CategoryTreeNode[], search: string): CategoryTreeNode[] {
  if (!search) return nodes
  const term = search.toLowerCase()

  const result: CategoryTreeNode[] = []
  for (const node of nodes) {
    const match = node.category_name.toLowerCase().includes(term)
    const filteredChildren = filterTree(node.children, search)
    if (match || filteredChildren.length > 0) {
      result.push({ ...node, children: filteredChildren })
    }
  }
  return result
}

interface FlatRow {
  node: CategoryTreeNode
  depth: number
}

function flattenTree(nodes: CategoryTreeNode[], expanded: Set<number>, depth = 0): FlatRow[] {
  const result: FlatRow[] = []
  for (const node of nodes) {
    result.push({ node, depth })
    if (node.children.length > 0 && expanded.has(node.category_id)) {
      result.push(...flattenTree(node.children, expanded, depth + 1))
    }
  }
  return result
}

export default function Category() {
  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')
  const [expanded, setExpanded] = useState<Set<number>>(new Set())
  const [addSheetOpen, setAddSheetOpen] = useState(false)
  const [editCategory, setEditCategory] = useState<Category | null>(null)

  useEffect(() => {
    const fetchCategories = async () => {
      try {
        const data = await getAllCategories()
        setCategories(data)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to fetch categories')
      } finally {
        setLoading(false)
      }
    }
    fetchCategories()
  }, [])

  const tree = useMemo(() => buildTree(categories), [categories])

  useEffect(() => {
    const parents = new Set<number>()
    const collectParents = (nodes: CategoryTreeNode[]) => {
      for (const node of nodes) {
        if (node.children.length > 0) {
          parents.add(node.category_id)
          collectParents(node.children)
        }
      }
    }
    collectParents(tree)

    setExpanded(prev => {
      const next = new Set(prev)
      for (const id of parents) next.add(id)
      return next
    })
  }, [tree])

  const handleAddSuccess = useCallback(() => {
    getAllCategories().then(setCategories).catch(console.error)
  }, [])

  const toggleExpand = useCallback((id: number) => {
    setExpanded(prev => {
      const next = new Set(prev)
      if (next.has(id)) next.delete(id)
      else next.add(id)
      return next
    })
  }, [])

  const filteredTree = useMemo(() => filterTree(tree, search), [tree, search])
  const flatRows = useMemo(() => flattenTree(filteredTree, expanded), [filteredTree, expanded])

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
              value={search}
              onChange={(e) => setSearch(e.target.value)}
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
        <table className="table-industrial">
          <thead>
            <tr>
              <th>CATEGORY NAME</th>
            </tr>
          </thead>
          <tbody>
            {flatRows.length === 0 ? (
              <tr>
                <td className="h-24 text-center text-muted-foreground">
                  No categories found.
                </td>
              </tr>
            ) : (
              flatRows.map((row) => (
                <tr
                  key={row.node.category_id}
                  className="cursor-pointer"
                  onClick={() => setEditCategory(row.node)}
                >
                  <td>
                    <div className="flex items-center gap-1" style={{ paddingLeft: `${row.depth * 20}px` }}>
                      {row.node.children.length > 0 ? (
                        <button
                          onClick={(e) => { e.stopPropagation(); toggleExpand(row.node.category_id) }}
                          className="text-muted-foreground hover:text-foreground"
                        >
                          {expanded.has(row.node.category_id) ? <IconMinus className="size-4" /> : <IconPlus className="size-4" />}
                        </button>
                      ) : (
                        <span className="size-4" />
                      )}
                      <span className="font-medium">{row.node.category_name}</span>
                    </div>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {flatRows.length > 0 && (
        <div className="border-t border-border p-3 flex items-center shrink-0 text-xs text-muted-foreground font-mono">
          {flatRows.length} categor{flatRows.length === 1 ? 'y' : 'ies'}
        </div>
      )}

      <Sheet open={addSheetOpen} onOpenChange={(open) => {
        if (!open) setAddSheetOpen(false)
      }}>
        <SheetContent className="w-100 sm:w-125 md:w-150 lg:w-175 xl:w-200 max-w-[90vw]">
          <AddCategorySheetContent
            onClose={() => setAddSheetOpen(false)}
            onSuccess={handleAddSuccess}
          />
        </SheetContent>
      </Sheet>

      <Sheet open={editCategory !== null} onOpenChange={(open) => {
        if (!open) setEditCategory(null)
      }}>
        <SheetContent className="w-100 sm:w-125 md:w-150 lg:w-175 xl:w-200 max-w-[90vw]">
          {editCategory && (
            <CategorySheetContent
              category={editCategory}
              onClose={() => setEditCategory(null)}
              onSuccess={handleAddSuccess}
            />
          )}
        </SheetContent>
      </Sheet>
    </div>
  )
}
