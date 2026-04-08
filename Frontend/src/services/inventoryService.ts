import type { Inventory, InventoryView } from '@/types/inventory';
import type { Product } from '@/types/product';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

export const getAllInventories = async (): Promise<Inventory[]> => {
  const response = await fetch(`${API_BASE_URL}/inventories`);
  if (!response.ok) {
    throw new Error('Failed to fetch inventories');
  }
  return response.json();
};

export const getAllProducts = async (): Promise<Product[]> => {
  const response = await fetch(`${API_BASE_URL}/products`);
  if (!response.ok) {
    throw new Error('Failed to fetch products');
  }
  return response.json();
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