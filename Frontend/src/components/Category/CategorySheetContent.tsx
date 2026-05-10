import { useEffect, useState, useMemo } from 'react'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import type { Category } from '@/types/category'
import { updateCategory, deleteCategory, getAllCategories } from '@/services/categoryService'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import { Package, X, Pencil, Trash2 } from 'lucide-react'
import { SectionLabel, EditLabel, EditCell, DataCell } from '@/components/ui/sheet-label'
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
  category: Category
  onClose: () => void
  onSuccess: () => void
}

const formSchema = z.object({
  category_name: z.string().min(1, 'Required'),
  parent_id: z.string(),
})

type FormData = z.infer<typeof formSchema>

export default function CategorySheetContent({ category, onClose, onSuccess }: Props) {
  const [mode, setMode] = useState<'view' | 'edit'>('view')
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [localCategory, setLocalCategory] = useState<Category>(category)
  const [categories, setCategories] = useState<Category[]>([])

  const form = useForm<FormData>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      category_name: localCategory.category_name,
      parent_id: localCategory.parent_id?.toString() ?? 'none',
    },
  })

  useEffect(() => {
    getAllCategories().then(setCategories).catch(() => {})
  }, [])

  useEffect(() => {
    if (mode === 'edit') {
      form.reset({
        category_name: localCategory.category_name,
        parent_id: localCategory.parent_id?.toString() ?? 'none',
      })
    }
  }, [mode, localCategory, form])

  const tree = useMemo(() => buildTree(categories), [categories])

  const filteredTree = useMemo(() => {
    const excludeNode = (nodes: CategoryTreeNode[], id: number): CategoryTreeNode[] =>
      nodes.filter(n => n.category_id !== id).map(n => ({
        ...n,
        children: excludeNode(n.children, id),
      }))
    return excludeNode(tree, localCategory.category_id)
  }, [tree, localCategory.category_id])

  const flatOptions = useMemo(() => flattenOptions(filteredTree), [filteredTree])

  const parentName = localCategory.parent_id
    ? categories.find(c => c.category_id === localCategory.parent_id)?.category_name ?? 'Unknown'
    : 'None'

  async function onSubmit(data: FormData) {
    setSaving(true)
    try {
      const parentId = data.parent_id && data.parent_id !== 'none' ? parseInt(data.parent_id) : null
      await updateCategory(localCategory.category_id.toString(), {
        category_name: data.category_name,
        parent_id: parentId,
      } as any)
      setLocalCategory(prev => ({ ...prev, category_name: data.category_name, parent_id: parentId ?? undefined }))
      setMode('view')
      onSuccess()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update category')
    } finally {
      setSaving(false)
    }
  }

  async function handleDelete() {
    try {
      await deleteCategory(localCategory.category_id.toString())
      onSuccess()
      onClose()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete category')
    }
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
              {mode === 'view' ? 'Category Details' : 'Editing'}
            </p>
            <h2 className="text-sm font-semibold leading-none">{localCategory.category_name}</h2>
          </div>
        </div>
        <Button
          variant="ghost"
          size="icon"
          className="h-7 w-7 rounded-sm text-muted-foreground hover:text-foreground hover:bg-muted"
          onClick={mode === 'view' ? onClose : () => setMode('view')}
        >
          <X className="h-3.5 w-3.5" />
        </Button>
      </div>

      <div className="flex-1 overflow-y-auto p-5 space-y-5">
        {error && (
          <ErrorMessage message={error} />
        )}

        {mode === 'view' ? (
          <div>
            <SectionLabel>Details</SectionLabel>
            <div className="border border-border rounded-md overflow-hidden">
              <DataCell label="CATEGORY NAME" value={localCategory.category_name} />
              <div className="border-t border-border">
                <DataCell label="PARENT CATEGORY" value={parentName} />
              </div>
            </div>
          </div>
        ) : (
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
                    value={form.watch('parent_id')}
                    onValueChange={(value) => form.setValue('parent_id', value)}
                  >
                    <SelectTrigger className="h-8 text-sm font-mono">
                      <SelectValue placeholder="None" />
                    </SelectTrigger>
                    <SelectContent position="popper">
                      <SelectItem value="none" className="font-mono text-sm italic text-muted-foreground">
                        None
                      </SelectItem>
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
        )}
      </div>

      <div className="px-5 py-4 border-t shrink-0 flex gap-2 justify-between items-center">
        {mode === 'view' ? (
          <>
            <AlertDialog>
              <AlertDialogTrigger asChild>
                <Button
                  variant="ghost"
                  size="sm"
                  className="text-muted-foreground hover:text-destructive hover:bg-destructive/10 h-8 px-3 text-xs gap-1.5"
                >
                  <Trash2 className="h-3.5 w-3.5" />
                  Delete
                </Button>
              </AlertDialogTrigger>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>Delete this category?</AlertDialogTitle>
                  <AlertDialogDescription>
                    This will permanently remove <strong>{localCategory.category_name}</strong>. This action cannot be undone.
                  </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel>Cancel</AlertDialogCancel>
                  <AlertDialogAction
                    className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                    onClick={handleDelete}
                  >
                    Delete
                  </AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>

            <Button
              type="button"
              size="sm"
              className="h-8 px-4 text-xs gap-1.5"
              onClick={() => window.setTimeout(() => setMode('edit'), 0)}
            >
              <Pencil className="h-3.5 w-3.5" />
              Edit Category
            </Button>
          </>
        ) : (
          <>
            <Button
              type="button"
              variant="ghost"
              size="sm"
              className="h-8 px-3 text-xs text-muted-foreground"
              onClick={() => setMode('view')}
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
                  Saving...
                </span>
              ) : (
                'Save Changes'
              )}
            </Button>
          </>
        )}
      </div>
    </div>
  )
}
