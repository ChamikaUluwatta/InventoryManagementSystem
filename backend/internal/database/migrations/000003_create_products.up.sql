CREATE TABLE IF NOT EXISTS "products" (
    "product_id"   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "product_name" VARCHAR(255) NOT NULL,
    "diameter"     NUMERIC(10, 2) NOT NULL DEFAULT 0,
    "width"        NUMERIC(10, 3) NOT NULL DEFAULT 0,
    "company_id"   UUID NOT NULL,
    "price"        NUMERIC(10, 2) NOT NULL DEFAULT 0,
    "category_id"  INT NOT NULL,
    CONSTRAINT "fk_product_category" FOREIGN KEY ("category_id")
        REFERENCES "categories" ("category_id"),
    CONSTRAINT "fk_product_company" FOREIGN KEY ("company_id")
        REFERENCES "companies" ("company_id"),
    CONSTRAINT "uq_product_identity" UNIQUE ("product_name", "diameter", "width", "company_id")
);