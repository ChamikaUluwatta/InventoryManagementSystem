ALTER TABLE "categories" DROP CONSTRAINT "fk_category_parent";
ALTER TABLE "categories" ADD CONSTRAINT "fk_category_parent" FOREIGN KEY ("parent_id") REFERENCES "categories" ("category_id");