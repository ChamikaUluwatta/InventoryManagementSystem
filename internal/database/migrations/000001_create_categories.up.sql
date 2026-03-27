CREATE TABLE IF NOT EXISTS "categories" (
    "category_id"   INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "category_name" VARCHAR(255) NOT NULL,
    "parent_id"     INT,
    CONSTRAINT "fk_category_parent" FOREIGN KEY ("parent_id")
        REFERENCES "categories" ("category_id")
        DEFERRABLE INITIALLY IMMEDIATE
);