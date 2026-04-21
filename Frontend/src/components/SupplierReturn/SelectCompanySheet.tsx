import { useEffect, useState } from 'react'
import { X, Package } from 'lucide-react'
import type { Company } from '@/types/company'
import { getAllCompanies } from '@/services/companyService'
import { Button } from '@/components/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { SectionLabel, EditCell } from '../ui/sheet-label'

type Props = {
  onClose: () => void
  onSelectCompany: (company: Company) => void
}

export default function SelectCompanySheet({ onClose, onSelectCompany }: Props) {
  const [companies, setCompanies] = useState<Company[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [selectedCompanyId, setSelectedCompanyId] = useState<string>('')

  useEffect(() => {
    const fetchCompanies = async () => {
      try {
        const data = await getAllCompanies()
        setCompanies(data)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load companies')
      } finally {
        setLoading(false)
      }
    }
    fetchCompanies()
  }, [])

  const handleContinue = () => {
    if (!selectedCompanyId) return
    const company = companies.find(c => c.company_id.toString() === selectedCompanyId)
    if (company) {
      onSelectCompany(company)
    }
  }

  const selectedCompany = companies.find(c => c.company_id.toString() === selectedCompanyId)

  return (
    <div className="flex flex-col h-full bg-background">
      <div className="flex items-center justify-between px-5 py-4 border-b shrink-0">
        <div className="flex items-center gap-3">
          <div className="flex items-center justify-center w-8 h-8 rounded border border-border">
            <Package className="h-4 w-4" />
          </div>
          <div>
            <p className="text-[10px] font-mono uppercase tracking-widest text-muted-foreground leading-none mb-0.5">
              Step 1 of 2
            </p>
            <h2 className="text-sm font-semibold leading-none">Select Company</h2>
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
        {loading ? (
          <div className="flex items-center gap-4 justify-center h-40">
            <div className="w-8 h-8 border-2 border-primary border-t-transparent rounded-full animate-spin" />
            <p>Loading...</p>
          </div>
        ) : error ? (
          <div className="m-4 p-3 border border-red-200 bg-red-50 dark:bg-red-950/20 dark:border-red-900 rounded text-sm text-red-600 dark:text-red-400 font-mono">
            ERR: {error}
          </div>
        ) : (
          <>
            <div>
              <SectionLabel>Company</SectionLabel>
              <div className="border border-border rounded-md overflow-hidden">
                <EditCell>
                  <Select
                    value={selectedCompanyId}
                    onValueChange={(value) => setSelectedCompanyId(value)}
                  >
                    <SelectTrigger className="h-10 text-sm font-mono">
                      <SelectValue placeholder="Select a company" />
                    </SelectTrigger>
                    <SelectContent position="popper" className="max-h-60">
                      {companies.map((comp) => (
                        <SelectItem
                          key={comp.company_id}
                          value={comp.company_id.toString()}
                          className="font-mono text-sm"
                        >
                          {comp.company_name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </EditCell>
              </div>
            </div>

            {selectedCompany && (
              <div className="p-4 border border-border rounded-md bg-muted/30">
                <p className="text-[10px] font-mono uppercase tracking-widest text-muted-foreground mb-1">
                  Selected Company
                </p>
                <p className="text-sm font-semibold">{selectedCompany.company_name}</p>
              </div>
            )}
          </>
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
          Cancel
        </Button>
        <Button
          type="button"
          size="sm"
          className="h-8 px-4 text-xs"
          disabled={!selectedCompanyId}
          onClick={handleContinue}
        >
          Continue
        </Button>
      </div>
    </div>
  )
}
