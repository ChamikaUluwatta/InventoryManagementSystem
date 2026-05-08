import * as React from 'react'
import { Link, useLocation } from 'react-router-dom'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible'
import {
  Sidebar, SidebarContent, SidebarGroup, SidebarHeader,
  SidebarMenu, SidebarMenuButton, SidebarMenuItem,
  SidebarMenuSub, SidebarMenuSubButton, SidebarMenuSubItem, SidebarRail,
} from '@/components/ui/sidebar'
import { PackageIcon, PlusIcon, MinusIcon, HomeIcon, LayoutGridIcon } from 'lucide-react'

type NavItem =
  | { type: 'link'; title: string; url: string; icon?: React.ElementType }
  | { type: 'group'; title: string; icon?: React.ElementType; items: { title: string; url: string; disabled?: boolean }[] }

const data: { navMain: NavItem[] } = {
  navMain: [
    {
      type: 'link',
      title: 'Home',
      url: '/',
      icon: HomeIcon,
    },
    {
      type: 'group',
      title: 'Inventory',
      icon: LayoutGridIcon,
      items: [
        { title: 'Products',   url: '/products' },
        { title: 'Stock',  url: '/stock' },
        { title: 'Returns', url: '/returns' },
        { title: 'Locations',  url: '/locations' },
        { title: 'Categories', url: '/categories' },
        { title: 'Companies',  url: '/companies' },
      ],
    },
  ],
}

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
            {data.navMain.map((item) => {
              if (item.type === 'link') {
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
              }

              const Icon = item.icon
              return (
                <Collapsible key={item.title} className="group/collapsible">
                  <SidebarMenuItem>
                    <CollapsibleTrigger asChild>
                      <SidebarMenuButton>
                        {Icon && <Icon className="size-4" />}
                        {item.title}
                        <PlusIcon  className="ml-auto group-data-[state=open]/collapsible:hidden" />
                        <MinusIcon className="ml-auto group-data-[state=closed]/collapsible:hidden" />
                      </SidebarMenuButton>
                    </CollapsibleTrigger>
                    {item.items?.length ? (
                      <CollapsibleContent>
                        <SidebarMenuSub>
                          {item.items.map((subItem) => (
                            <SidebarMenuSubItem key={subItem.title}>
                              {subItem.disabled ? (
                                <SidebarMenuSubButton asChild isActive={false}>
                                  <span aria-disabled="true" className="cursor-not-allowed opacity-50" title="Coming soon">
                                    {subItem.title}
                                  </span>
                                </SidebarMenuSubButton>
                              ) : (
                                <SidebarMenuSubButton asChild isActive={isActive(subItem.url)}>
                                  <Link to={subItem.url}>{subItem.title}</Link>
                                </SidebarMenuSubButton>
                              )}
                            </SidebarMenuSubItem>
                          ))}
                        </SidebarMenuSub>
                      </CollapsibleContent>
                    ) : null}
                  </SidebarMenuItem>
                </Collapsible>
              )
            })}
          </SidebarMenu>
        </SidebarGroup>
      </SidebarContent>
      <SidebarRail />
    </Sidebar>
  )
}