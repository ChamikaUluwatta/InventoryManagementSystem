CREATE TABLE IF NOT EXISTS "companies" (
    "company_id"   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "company_name" VARCHAR(255) NOT NULL
);