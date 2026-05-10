import { Link, useLocation } from 'react-router-dom'
import {
  Sidebar, SidebarContent, SidebarGroup, SidebarHeader,
  SidebarMenu, SidebarMenuButton, SidebarMenuItem, SidebarRail,
} from '@/components/ui/sidebar'
import {
  PackageIcon, HomeIcon, WarehouseIcon,
  RotateCcwIcon, MapPinIcon, TagsIcon, Building2Icon,
} from 'lucide-react'

const navItems = [
  { title: 'Home', url: '/', icon: HomeIcon },
  { title: 'Products', url: '/products', icon: PackageIcon },
  { title: 'Stock', url: '/stock', icon: WarehouseIcon },
  { title: 'Returns', url: '/returns', icon: RotateCcwIcon },
  { title: 'Locations', url: '/locations', icon: MapPinIcon },
  { title: 'Categories', url: '/categories', icon: TagsIcon },
  { title: 'Companies', url: '/companies', icon: Building2Icon },
]

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const location = useLocation()
  const isActive = (url: string) => location.pathname === url

  return (
    <Sidebar {...props}>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild>
              <Link to="/">
                <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                  <PackageIcon className="size-4" />
                </div>
                <div className="flex flex-col gap-0.5 leading-none">
                  <span className="font-medium">IMS</span>
                  <span className="">v1.0.0</span>
                </div>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>

      <SidebarContent>
        <SidebarGroup>
          <SidebarMenu>
            {navItems.map((item) => {
              const Icon = item.icon
              return (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton asChild isActive={isActive(item.url)}>
                    <Link to={item.url}>
                      {Icon && <Icon className="size-4" />}
                      {item.title}
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              )
            })}
          </SidebarMenu>
        </SidebarGroup>
      </SidebarContent>
      <SidebarRail />
    </Sidebar>
  )
}