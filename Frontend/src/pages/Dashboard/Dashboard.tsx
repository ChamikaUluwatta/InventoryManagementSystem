import { Link } from 'react-router-dom'
import {
  Package,
  DollarSign,
  AlertTriangle,
  MapPin,
  Truck,
  PackageCheck,
  Boxes,
  Building2,
  FolderTree
} from 'lucide-react'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'

const statCards = [
  {
    label: 'Total Products',
    value: '1,234',
    icon: Package,
    color: 'bg-blue-500',
  },
  {
    label: 'Stock Value',
    value: '$45,678',
    icon: DollarSign,
    color: 'bg-green-500',
  },
  {
    label: 'Low Stock',
    value: '23',
    icon: AlertTriangle,
    color: 'bg-amber-500',
  },
  {
    label: 'Locations',
    value: '8',
    icon: MapPin,
    color: 'bg-purple-500',
  },
]

const modules = [
  {
    title: 'Products',
    description: 'Manage your product catalog',
    icon: Package,
    href: '/products',
  },
  {
    title: 'Inventory',
    description: 'Track stock levels',
    icon: Boxes,
    href: '/inventory',
  },
  {
    title: 'Locations',
    description: 'Manage warehouses',
    icon: MapPin,
    href: '/locations',
  },
  {
    title: 'Categories',
    description: 'Organize products',
    icon: FolderTree,
    href: '/categories',
  },
  {
    title: 'Companies',
    description: 'Supplier management',
    icon: Building2,
    href: '/companies',
  },
]

const recentActivity = [
  {
    primary: 'New product added: Widget Pro X',
    timestamp: '2 minutes ago',
    initials: 'JD',
  },
  {
    primary: 'Stock updated for: Cable USB-C',
    timestamp: '15 minutes ago',
    initials: 'MK',
  },
  {
    primary: 'Location "Warehouse A" updated',
    timestamp: '1 hour ago',
    initials: 'AL',
  },
  {
    primary: 'New supplier added: Acme Corp',
    timestamp: '3 hours ago',
    initials: 'TB',
  },
  {
    primary: 'Low stock alert: Adapter 5W',
    timestamp: '5 hours ago',
    initials: 'SY',
  },
]

const alerts = [
  {
    type: 'Low Stock',
    message: '15 items below minimum threshold',
    icon: AlertTriangle,
    iconColor: 'text-amber-500',
    bgColor: 'bg-amber-50',
  },
  {
    type: 'Incoming Shipment',
    message: '3 orders arriving today',
    icon: Truck,
    iconColor: 'text-blue-500',
    bgColor: 'bg-blue-50',
  },
  {
    type: 'Restocked',
    message: '8 items refilled',
    icon: PackageCheck,
    iconColor: 'text-green-500',
    bgColor: 'bg-green-50',
  },
]

export default function Dashboard() {
  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
        <p className="text-muted-foreground">Overview of your inventory operations</p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {statCards.map((stat) => (
          <Card key={stat.label}>
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-sm font-medium">{stat.label}</CardTitle>
              <Avatar className={stat.color}>
                <AvatarFallback className="text-white">
                  <stat.icon className="size-4" />
                </AvatarFallback>
              </Avatar>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stat.value}</div>
            </CardContent>
          </Card>
        ))}
      </div>

      <div className="grid gap-4 md:grid-cols-3 lg:grid-cols-5">
        {modules.map((module) => (
          <Card key={module.title} className="hover:bg-muted/50 transition-colors">
            <Link to={module.href}>
              <CardHeader className="flex flex-row items-center gap-4 p-4">
                <Avatar className="bg-accent">
                  <AvatarFallback className="bg-accent text-accent-foreground">
                    <module.icon className="size-5" />
                  </AvatarFallback>
                </Avatar>
                <CardTitle className="text-base">{module.title}</CardTitle>
              </CardHeader>
              <CardContent>
                <CardDescription>{module.description}</CardDescription>
              </CardContent>
            </Link>
          </Card>
        ))}
      </div>

      <div className="grid gap-4 md:grid-cols-12">
        <Card className="md:col-span-7">
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-base font-medium">Recent Activity</CardTitle>
            <Badge variant="success">Live</Badge>
          </CardHeader>
          <CardContent>
            <div className="flex flex-col gap-4">
              {recentActivity.map((activity, index) => (
                <div key={index} className="flex items-center gap-4">
                  <Avatar>
                    <AvatarFallback className="bg-muted">{activity.initials}</AvatarFallback>
                  </Avatar>
                  <div className="flex flex-col gap-1">
                    <p className="text-sm font-medium">{activity.primary}</p>
                    <p className="text-xs text-muted-foreground">{activity.timestamp}</p>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        <Card className="md:col-span-5">
          <CardHeader>
            <CardTitle className="text-base font-medium">Alerts</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex flex-col gap-3">
              {alerts.map((alert, index) => (
                <div
                  key={index}
                  className={`flex items-center gap-3 rounded-lg p-3 ${alert.bgColor}`}
                >
                  <div className={alert.iconColor}>
                    <alert.icon className="size-5" />
                  </div>
                  <div className="flex flex-col gap-0.5">
                    <p className="text-sm font-medium">{alert.type}</p>
                    <p className="text-xs text-muted-foreground">{alert.message}</p>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}