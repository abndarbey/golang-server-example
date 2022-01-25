BEGIN;
-- Permissions
CREATE TABLE "permissions" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT NOW(),
  "updated_at" timestamptz NOT NULL DEFAULT NOW()
);
-- Organizations
CREATE TABLE "organizations" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "code" varchar NOT NULL,
  "name" varchar NOT NULL,
  "website" varchar,
  "is_archived" boolean NOT NULL DEFAULT FALSE,
  "created_at" timestamptz NOT NULL DEFAULT NOW(),
  "updated_at" timestamptz NOT NULL DEFAULT NOW()
);
CREATE TABLE "roles" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "code" varchar NOT NULL,
  "name" varchar NOT NULL,
  "permissions" text[],
  "is_org_admin" boolean NOT NULL DEFAULT FALSE,
  "is_archived" boolean NOT NULL DEFAULT FALSE,
  "organization_id" bigint NOT NULL REFERENCES organizations (id),
  "created_at" timestamptz NOT NULL DEFAULT NOW(),
  "updated_at" timestamptz NOT NULL DEFAULT NOW()
);
-- Users
CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "phone" varchar UNIQUE NOT NULL,
  "is_admin" boolean NOT NULL DEFAULT FALSE,
  "is_member" boolean NOT NULL DEFAULT FALSE, -- defines if user is affilated to an org
  "is_customer" boolean NOT NULL DEFAULT FALSE,
  "password_hash" varchar NOT NULL,
  "organization_id" bigint REFERENCES organizations (id),
  "role_id" bigint REFERENCES roles (id),
  "created_at" timestamptz NOT NULL DEFAULT NOW(),
  "updated_at" timestamptz NOT NULL DEFAULT NOW()
);
CREATE TABLE "profiles" (
  "user_id" bigserial PRIMARY KEY NOT NULL REFERENCES users (id),
  "date_of_birth" varchar UNIQUE,
  "referral_code" varchar UNIQUE,
  "wallet_points" bigserial NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT NOW(),
  "updated_at" timestamptz NOT NULL DEFAULT NOW()
);

-- Inventory
CREATE TABLE "containers" (
  "id" bigserial NOT NULL PRIMARY KEY,
  "uid" uuid UNIQUE NOT NULL,
  "code" text UNIQUE NOT NULL,
  "description" text NOT NULL DEFAULT '',
  "is_archived" boolean NOT NULL DEFAULT FALSE,
  "organization_id" bigint REFERENCES organizations (id),
  "created_by_id" bigint NOT NULL REFERENCES users (id),
  "created_at" timestamptz NOT NULL DEFAULT NOW(),
  "updated_at" timestamptz NOT NULL DEFAULT NOW()
);
CREATE TABLE "pallets" (
  "id" bigserial NOT NULL PRIMARY KEY,
  "uid" uuid UNIQUE NOT NULL,
  "code" text UNIQUE NOT NULL,
  "description" text NOT NULL DEFAULT '',
  "container_id" bigint REFERENCES containers (id),
  "is_archived" boolean NOT NULL DEFAULT FALSE,
  "organization_id" bigint REFERENCES organizations (id),
  "created_by_id" bigint NOT NULL REFERENCES users (id),
  "created_at" timestamptz NOT NULL DEFAULT NOW(),
  "updated_at" timestamptz NOT NULL DEFAULT NOW()
);

COMMIT;