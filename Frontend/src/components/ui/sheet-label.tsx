function SectionLabel({ children }: { children: React.ReactNode }) {
  return (
    <p className="text-[10px] font-mono uppercase tracking-widest text-muted-foreground mb-2 flex items-center gap-2">
      <span className="inline-block w-3 h-px bg-muted-foreground/40" />
      {children}
    </p>
  )
}

function EditLabel({ children }: { children: React.ReactNode }) {
  return (
    <p className="text-[10px] font-mono uppercase tracking-wider text-muted-foreground mb-1">
      {children}
    </p>
  )
}

function EditCell({
  children,
  bordered,
  topBorder,
  className,
}: {
  children: React.ReactNode
  bordered?: boolean
  topBorder?: boolean
  className?: string
}) {
  return (
    <div
      className={[
        'px-4 py-3',
        bordered ? 'border-r border-border' : '',
        topBorder ? 'border-t border-border' : '',
        className || '',
      ].join(' ')}
    >
      {children}
    </div>
  )
}

function DataCell({
  label,
  value,
  bordered,
  topBorder,
  className,
}: {
  label: string
  value: string
  bordered?: boolean
  topBorder?: boolean
  className?: string
}) {
  return (
    <div
      className={[
        'px-4 py-3',
        bordered ? 'border-r border-border' : '',
        topBorder ? 'border-t border-border' : '',
        className ? className : '',
      ].join(' ')}
    >
      <p className="text-[10px] font-mono uppercase tracking-wider text-muted-foreground mb-0.5">
        {label}
      </p>
      <p className="text-sm  font-mono">{value}</p>
    </div>
  )
}

export { SectionLabel, EditLabel, EditCell, DataCell }