CREATE TABLE IF NOT EXISTS "inventories" (
    "inventory_id" INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "product_id"   UUID NOT NULL,
    "location_id"  VARCHAR(255) NOT NULL,
    "stock"        INT NOT NULL DEFAULT 0,
    CONSTRAINT "fk_inventory_product" FOREIGN KEY ("product_id")
        REFERENCES "products" ("product_id"),
    CONSTRAINT "fk_inventory_location" FOREIGN KEY ("location_id")
        REFERENCES "locations" ("location_id"),
    CONSTRAINT "uq_inventory_product_location" UNIQUE ("product_id", "location_id")
);