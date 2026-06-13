ALTER TABLE "products"
ADD COLUMN IF NOT EXISTS "location_id" TEXT NOT NULL DEFAULT 'unassigned';