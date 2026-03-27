CREATE TABLE "Products" (
  "product_id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "product_name" varchar(255) NOT NULL,
  "diameter" numeric(10, 2) NOT NULL DEFAULT 0, 
  "width" numeric(10, 3) NOT NULL DEFAULT 0, 
  "company_id" uuid NOT NULL, 
  "price" numeric(10, 2) DEFAULT 0,
  "category_id" int NOT NULL, 
  CONSTRAINT "uq_product_identity" UNIQUE ("product_name", "diameter", "width", "company_id")
);

CREATE TABLE "Categories" (
  "category_id" INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  "category_name" varchar(255),
  "parent_id" int
);

CREATE TABLE "Companies" (
  "company_id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "company_name" varchar(255)
);

CREATE TABLE "Locations" (
  "location_id" varchar(255) PRIMARY KEY,
  "image" text
);

CREATE TABLE "Inventories" (
  "inventory_id" INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  "product_id" uuid,
  "location_id" varchar(255),
  "stock" int
);


ALTER TABLE "Products" ADD FOREIGN KEY ("category_id") REFERENCES "Categories" ("category_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "Products" ADD FOREIGN KEY ("company_id") REFERENCES "Companies" ("company_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "Categories" ADD FOREIGN KEY ("parent_id") REFERENCES "Categories" ("category_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "Inventories" ADD FOREIGN KEY ("product_id") REFERENCES "Products" ("product_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "Inventories" ADD FOREIGN KEY ("location_id") REFERENCES "Locations" ("location_id") DEFERRABLE INITIALLY IMMEDIATE;
