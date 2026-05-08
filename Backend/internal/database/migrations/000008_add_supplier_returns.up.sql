CREATE TABLE IF NOT EXISTS "supplier_returns" (
    "supplier_return_id" INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "return_no" VARCHAR(50) NOT NULL UNIQUE,
    "company_id" UUID NOT NULL,
    "status" VARCHAR(30) NOT NULL DEFAULT 'draft',
    "reason" TEXT,
    "notes" TEXT,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "approved_at" TIMESTAMPTZ,
    "completed_at" TIMESTAMPTZ,
    CONSTRAINT "fk_supplier_return_company" FOREIGN KEY ("company_id")
        REFERENCES "companies" ("company_id"),
    CONSTRAINT "chk_supplier_return_status" CHECK (
        "status" IN ('draft', 'approved', 'sent', 'credited', 'cancelled','rejected','completed')
    )
);

CREATE TABLE IF NOT EXISTS "supplier_return_items" (
    "supplier_return_item_id" INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "supplier_return_id" INT NOT NULL,
    "product_id" UUID,
    "location_id" VARCHAR(255),
    "quantity" INT NOT NULL CHECK ("quantity" > 0),
    "unit_cost" NUMERIC(10, 2) NOT NULL DEFAULT 0 CHECK ("unit_cost" >= 0),

    "product_name_snapshot" VARCHAR(255) NOT NULL,
    "location_snapshot" VARCHAR(255) NOT NULL,

    CONSTRAINT "fk_supplier_return_item_return" FOREIGN KEY ("supplier_return_id")
        REFERENCES "supplier_returns" ("supplier_return_id") ON DELETE CASCADE,
    CONSTRAINT "fk_supplier_return_item_product" FOREIGN KEY ("product_id")
        REFERENCES "products" ("product_id") ON DELETE SET NULL,
    CONSTRAINT "fk_supplier_return_item_location" FOREIGN KEY ("location_id")
        REFERENCES "locations" ("location_id") ON DELETE SET NULL
);