import * as React from 'react'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { useSidebar } from '@/components/ui/sidebar'
import { PanelLeftIcon } from 'lucide-react'
import { Link, useLocation } from 'react-router-dom'

const breadcrumbConfig: Record<string, { label: string; parent?: { label: string; url: string } }> =
  {
    '/': { label: '' },
    '/products': { label: 'Products' },
    '/products/add': { label: 'Add Product' },
    '/stock': { label: 'Stock' },
    '/returns': { label: 'Returns' },
    '/locations': { label: 'Locations' },
    '/categories': { label: 'Categories' },
  }

function generateBreadcrumbs(pathname: string) {
  const exactMatch = breadcrumbConfig[pathname]
  if (exactMatch) {
    return exactMatch.parent
      ? [
          { label: 'Home', url: '/', isPage: false },
          { ...exactMatch.parent, isPage: false },
          { label: exactMatch.label, isPage: true },
        ]
      : [
          { label: 'Home', url: '/', isPage: false },
          { label: exactMatch.label, isPage: true },
        ]
  }

  if (pathname.startsWith('/products/') && pathname.includes('/edit')) {
    return [
      { label: 'Home', url: '/', isPage: false },
      { label: 'Products', url: '/products/manage', isPage: false },
      { label: 'Edit Product', isPage: true },
    ]
  }

  return [
    { label: 'Home', url: '/', isPage: false },
    { label: 'Page Not Found', isPage: true },
  ]
}

export function SiteHeader() {
  const { toggleSidebar } = useSidebar()
  const location = useLocation()
  const breadcrumbs = generateBreadcrumbs(location.pathname)

  return (
    <header className="sticky top-0 z-50 flex w-full items-center border-b bg-background">
      <div className="flex h-(--header-height) w-full items-center gap-2 px-4">
        <Button className="h-8 w-8" variant="ghost" size="icon" onClick={toggleSidebar}>
          <PanelLeftIcon />
        </Button>
        <Separator
          orientation="vertical"
          className="mr-2 data-vertical:h-4 data-vertical:self-auto"
        />
        <Breadcrumb className="hidden sm:block">
          <BreadcrumbList>
            {breadcrumbs.map((item, index) => (
              <React.Fragment key={index}>
                <BreadcrumbItem>
                  {item.isPage ? (
                    <BreadcrumbPage>{item.label}</BreadcrumbPage>
                  ) : (
                    <BreadcrumbLink asChild>
                      <Link to={item.url || '#'}>{item.label}</Link>
                    </BreadcrumbLink>
                  )}
                </BreadcrumbItem>
                {index < breadcrumbs.length - 1 && <BreadcrumbSeparator />}
              </React.Fragment>
            ))}
          </BreadcrumbList>
        </Breadcrumb>
        {/* <SearchForm className="w-full sm:ml-auto sm:w-auto" /> */}
      </div>
    </header>
  )
}
