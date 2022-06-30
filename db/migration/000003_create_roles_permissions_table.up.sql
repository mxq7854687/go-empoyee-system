CREATE TABLE "roles" (
    "id" bigserial PRIMARY KEY,
    "role" varchar NOT NULL,
    "privileges" jsonb NOT NULL, 
    "updated_at" timestamptz NOT NULL DEFAULT (now()),
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TYPE "privilege" AS ENUM (
  'CreateAndUpdateJobs',
  'CreateAndUpdateDepartments',
  'DeleteJobs',
  'DeleteDepartments',
  'CreateAndUpdateEmployees',
  'DelteEmployees',
  'ReadAllEmployees'
);
