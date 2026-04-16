import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { AppLayout } from '@/components/Layout/AppLayout'
import ViewManageProducts from '@/components/Product/ViewManageProducts/ViewManageProducts'
import { Button } from './components/ui/button'
import ViewAddProduct from '@/components/Product/ViewAddProduct/ViewAddProduct'
import ViewInventory from '@/components/Inventory/ViewInventory/ViewInventory'
import ViewLocation from '@/components/Location/ViewLocation/ViewLocation'
import ViewCategory from '@/components/Category/ViewCategory/ViewCategory'


function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route
          path="/products"
          element={
            <AppLayout>
              <ViewManageProducts />
            </AppLayout>
          }
        />
        <Route
          path="/"
          element={
            <AppLayout>
              <div className="p-4 h-full flex items-center justify-center">
                <h1 className="text-2xl font-bold">Welcome to Inventory Management system</h1>
              </div>
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
              <ViewAddProduct />
            </AppLayout>
          }
        />
        <Route
          path="/inventory"
          element={
            <AppLayout>
              <ViewInventory />
            </AppLayout>
          }
        />
        <Route
          path="/locations"
          element={
            <AppLayout>
              <ViewLocation />
            </AppLayout>
          }
        />
        <Route
          path="/categories"
          element={
            <AppLayout>
              <ViewCategory />
            </AppLayout>
          }
        />
      </Routes>
    </BrowserRouter>
  )
}

export default App