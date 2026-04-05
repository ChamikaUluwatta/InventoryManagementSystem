import { BrowserRouter, Routes, Route, useParams } from "react-router-dom"
import { AppLayout } from "@/components/Layout/AppLayout"
import ProductList from "@/components/Product/ProductListView/ProductList"
import { Button } from "./components/ui/button"
import EditProduct from "./components/Product/EditProduct/EditProduct"

function EditProductWrapper() {
  const { productId } = useParams<{ productId: string }>();
  return <EditProduct uuid={productId || ""} />;
}

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route
          path="/products"
          element={
            <AppLayout>
              <ProductList />
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
          path="/products/:productId/edit"
          element={
            <AppLayout>
              <EditProductWrapper />
            </AppLayout>
          }
        />
      </Routes>
    </BrowserRouter>
  )
}

export default App