ALTER TABLE "users" ADD COLUMN "role_id" bigserial;
ALTER TABLE "users" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");