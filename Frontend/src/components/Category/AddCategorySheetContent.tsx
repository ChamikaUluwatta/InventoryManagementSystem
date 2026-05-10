import { useEffect, useState, useMemo } from 'react'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import type { Category } from '@/types/category'
import { createCategory } from '@/services/categoryService'
import { getAllCategories } from '@/services/categoryService'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Package, X } from 'lucide-react'
import { Spinner } from '@/components/ui/spinner'
import { SectionLabel, EditLabel, EditCell } from '@/components/ui/sheet-label'
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

interface FlatOption {
  node: CategoryTreeNode
  depth: number
}

function flattenOptions(nodes: CategoryTreeNode[], depth = 0): FlatOption[] {
  const result: FlatOption[] = []
  for (const node of nodes) {
    result.push({ node, depth })
    if (node.children.length > 0) {
      result.push(...flattenOptions(node.children, depth + 1))
    }
  }
  return result
}

type Props = {
  onClose: () => void
  onSuccess: () => void
}

const formSchema = z.object({
  category_name: z.string().min(1, 'Required'),
  parent_id: z.number().int().positive().optional(),
})

type FormData = z.infer<typeof formSchema>

export default function AddCategorySheetContent({ onClose, onSuccess }: Props) {
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [categories, setCategories] = useState<Category[]>([])

  const form = useForm<FormData>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      category_name: '',
      parent_id: undefined,
    },
  })

  useEffect(() => {
    const fetchData = async () => {
      try {
        const data = await getAllCategories()
        setCategories(data)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load categories')
      } finally {
        setLoading(false)
      }
    }
    fetchData()
  }, [])

  const tree = useMemo(() => buildTree(categories), [categories])
  const flatOptions = useMemo(() => flattenOptions(tree), [tree])

  async function onSubmit(data: FormData) {
    setSaving(true)
    try {
      await createCategory({
        category_name: data.category_name,
        parent_id: data.parent_id || undefined,
      })
      onSuccess()
      onClose()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create category')
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

  if (error && !categories.length) {
    return (
      <ErrorMessage message={error} />
    )
  }

  return (
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
            <h2 className="text-sm font-semibold leading-none">New Category</h2>
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

        <form id="category-form" onSubmit={form.handleSubmit(onSubmit)} className="space-y-5">
          <div>
            <SectionLabel>Category Name</SectionLabel>
            <Input
              id="category_name"
              className="h-9 text-sm font-mono"
              placeholder="Enter category name"
              {...form.register('category_name')}
            />
            {form.formState.errors.category_name && (
              <p className="text-xs text-destructive font-mono mt-1">
                {form.formState.errors.category_name.message}
              </p>
            )}
          </div>

          <div>
            <SectionLabel>Parent Category (optional)</SectionLabel>
            <div className="border border-border rounded-md overflow-hidden">
              <EditCell className="col-span-2 flex justify-between items-center">
                <EditLabel>PARENT</EditLabel>
                <Select
                  value={form.watch('parent_id')?.toString()}
                  onValueChange={(value) =>
                    form.setValue('parent_id', value ? parseInt(value) : undefined)
                  }
                >
                  <SelectTrigger className="h-8 text-sm font-mono">
                    <SelectValue placeholder="None" />
                  </SelectTrigger>
                  <SelectContent position="popper">
                    {flatOptions.map(({ node, depth }) => (
                      <SelectItem
                        key={node.category_id}
                        value={node.category_id.toString()}
                        className="font-mono text-sm"
                      >
                        <span style={{ paddingLeft: `${depth * 16}px` }}>
                          {node.category_name}
                        </span>
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </EditCell>
            </div>
          </div>
        </form>
      </div>

      <div className="px-5 py-4 border-t shrink-0 flex gap-2 justify-end items-center">
        <Button
          type="button"
          variant="ghost"
          size="sm"
          className="h-8 px-3 text-xs text-muted-foreground"
          onClick={onClose}
        >
          Cancel
        </Button>
        <Button
          type="submit"
          form="category-form"
          size="sm"
          className="h-8 px-4 text-xs"
          disabled={saving}
        >
          {saving ? (
            <span className="flex items-center gap-2">
              <span className="w-3 h-3 border border-current border-t-transparent rounded-full animate-spin" />
              Creating...
            </span>
          ) : (
            'Create Category'
          )}
        </Button>
      </div>
    </div>
  )
}
