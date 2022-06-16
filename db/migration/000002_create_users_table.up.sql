CREATE TYPE user_status AS ENUM ('pending', 'activated', 'deactivated');

CREATE TABLE "users" (
    -- "id" bigserial PRIMARY KEY,
    "email" varchar(100) PRIMARY KEY,
    "status" user_status DEFAULT 'pending',
    "hashed_password" varchar NOT NULL,
    "updated_at" timestamptz NOT NULL DEFAULT (now()),
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "employees" ADD FOREIGN KEY ("email") REFERENCES "users" ("email");