CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE schema IF NOT EXISTS base;
CREATE TABLE IF NOT EXISTS base.tenants (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
ALTER TABLE base.tenants ADD CONSTRAINT "tenant_name_unique" UNIQUE ("name");