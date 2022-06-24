-- ALTER TABLE "employees" DROP COLUMN IF EXISTS "user_id";
ALTER TABLE IF EXISTS "employees" DROP CONSTRAINT IF EXISTS "employees_email_fkey";
DROP TABLE IF EXISTS "users";
DROP TYPE IF EXISTS "user_status";