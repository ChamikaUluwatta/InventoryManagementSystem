BEGIN;

ALTER TABLE "inventories"
  DROP CONSTRAINT IF EXISTS "fk_inventory_product",
  DROP CONSTRAINT IF EXISTS "fk_inventory_location";

ALTER TABLE "inventories"
  ADD CONSTRAINT "fk_inventory_product"
    FOREIGN KEY ("product_id")
    REFERENCES "products" ("product_id");

ALTER TABLE "inventories"
  ADD CONSTRAINT "fk_inventory_location"
    FOREIGN KEY ("location_id")
    REFERENCES "locations" ("location_id");

ALTER TABLE "inventories"
  ALTER COLUMN "location_id" DROP DEFAULT;

COMMIT;
