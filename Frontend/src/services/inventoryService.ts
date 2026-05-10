import type { Inventory, InventoryView } from '@/types/inventory';
import type { Product } from '@/types/product';
import { apiFetch } from '@/lib/api';

export const getAllInventories = async (): Promise<Inventory[]> => {
  return apiFetch<Inventory[]>('/inventories');
};

export const getAllProducts = async (): Promise<Product[]> => {
  return apiFetch<Product[]>('/products');
};

export const getInventoryWithProductDetails = async (): Promise<InventoryView[]> => {
  const [inventories, products] = await Promise.all([
    getAllInventories(),
    getAllProducts(),
  ]);

  const productMap = new Map(products.map(p => [p.product_id, p]));

  return inventories.map(inv => {
    const product = productMap.get(inv.product_id);
    return {
      inventory_id: inv.inventory_id,
      product_id: inv.product_id,
      product_name: product?.product_name || 'Unknown',
      location_id: inv.location_id,
      stock: inv.stock,
    };
  });
};

export const createInventory = async (inventory: Omit<Inventory, 'inventory_id'>): Promise<Inventory> => {
  return apiFetch<Inventory>('/inventories', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(inventory),
  });
};

export const updateInventory = async (id: number, inventory: { product_id: string; location_id: string; stock: number }): Promise<Inventory> => {
  return apiFetch<Inventory>(`/inventories/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(inventory),
  })
}

export const deleteInventory = async (id: number): Promise<void> => {
  return apiFetch<void>(`/inventories/${id}`, {
    method: 'DELETE',
  });
};
