export interface Product {
  product_id: string
  product_name: string
  product_description?: string
  diameter: number
  width: number
  company_id: string
  price: number
  category_id: number
  location_id: string
  stock: number
}

export type CreateProductRequest = Omit<Product, 'product_id' | 'stock'>
