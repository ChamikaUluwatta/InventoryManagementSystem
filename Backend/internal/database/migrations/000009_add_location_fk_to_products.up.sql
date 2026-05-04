ALTER TABLE "products" DROP CONSTRAINT IF EXISTS "fk_products_location";

-- add unassigned to make sure unassigned exist in location table --
INSERT INTO "locations" ("location_id") VALUES ('unassigned') ON CONFLICT DO NOTHING;

--update existing products with no reference to location to unassigned --
UPDATE "products" SET "location_id" = 'unassigned' WHERE "location_id" IS NULL or "location_id" NOT IN (SELECT "location_id" FROM "locations");

--add the foreign key constraint to products table --
ALTER TABLE "products" ADD CONSTRAINT "fk_location" FOREIGN KEY ("location_id") REFERENCES "locations"("location_id") ON DELETE SET DEFAULT;