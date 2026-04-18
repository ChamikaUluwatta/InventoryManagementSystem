import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { AppLayout } from '@/components/Layout/AppLayout'
import ManageProducts from '@/pages/Products/ManageProducts'
import AddProduct from '@/pages/Products/AddProduct'
import Category from '@/pages/Category/Category'
import Inventory from '@/pages/Inventory/Inventory'
import Location from '@/pages/Location/Location'
import { Button } from './components/ui/button'
import Dashboard from './pages/Dashboard/Dashboard'



function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route
          path="/products"
          element={
            <AppLayout>
              <ManageProducts />
            </AppLayout>
          }
        />
        <Route
          path="/"
          element={
            <AppLayout>
              <Dashboard />
            </AppLayout>
          }
        />
        <Route
          path="*"
          element={
            <AppLayout>
              <div className="p-4 h-full flex items-center justify-center flex-col gap-2">
                <h1 className="text-2xl font-bold">Page Not Found</h1>
                <Button variant="outline" className="ml-4" onClick={() => window.history.back()}>
                  Go Back
                </Button>
              </div>
            </AppLayout>
          }
        />
        <Route
          path="/products/new"
          element={
            <AppLayout>
              <AddProduct />
            </AppLayout>
          }
        />
        <Route
          path="/inventory"
          element={
            <AppLayout>
              <Inventory />
            </AppLayout>
          }
        />
        <Route
          path="/locations"
          element={
            <AppLayout>
              <Location />
            </AppLayout>
          }
        />
        <Route
          path="/categories"
          element={
            <AppLayout>
              <Category />
            </AppLayout>
          }
        />
        <Route
          path="/companies"
          element={
            <AppLayout>
              <div className="p-4 h-full flex items-center justify-center">
                <h1 className="text-2xl font-bold">Companies - Coming Soon</h1>
              </div>
            </AppLayout>
          }
        />
      </Routes>
    </BrowserRouter>
  )
}

export default App