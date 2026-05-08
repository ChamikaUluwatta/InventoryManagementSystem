export interface Inventory {
  inventory_id: number;
  product_id: string;
  location_id: string;
  stock: number;
}

export interface InventoryView {
  inventory_id: number;
  product_id: string;
  product_name: string;
  location_id: string;
  stock: number;
}

export interface InventoryCreateRequest {
  product_id: string;
  location_id: string;
  stock: number;
}