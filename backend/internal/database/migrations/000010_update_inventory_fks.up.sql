BEGIN;

-- Ensure 'unassigned' location exists for SET DEFAULT
INSERT INTO "locations" ("location_id") VALUES ('unassigned') ON CONFLICT DO NOTHING;

-- Add default so ON DELETE SET DEFAULT has a value to use
ALTER TABLE "inventories"
  ALTER COLUMN "location_id" SET DEFAULT 'unassigned';

-- Drop existing FK constraints (no ON DELETE action)
ALTER TABLE "inventories"
  DROP CONSTRAINT IF EXISTS "fk_inventory_product",
  DROP CONSTRAINT IF EXISTS "fk_inventory_location";

-- Re-add with ON DELETE CASCADE for product (delete product → delete its inventory)
ALTER TABLE "inventories"
  ADD CONSTRAINT "fk_inventory_product"
    FOREIGN KEY ("product_id")
    REFERENCES "products" ("product_id")
    ON DELETE CASCADE;

-- Re-add with ON DELETE SET DEFAULT for location (delete location → set to 'unassigned')
ALTER TABLE "inventories"
  ADD CONSTRAINT "fk_inventory_location"
    FOREIGN KEY ("location_id")
    REFERENCES "locations" ("location_id")
    ON DELETE SET DEFAULT;

COMMIT;
