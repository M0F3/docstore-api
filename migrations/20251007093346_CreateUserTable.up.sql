CREATE TABLE IF NOT EXISTS base.users (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(), 
    first_name text NOT NULL,
    last_name text NOT NULL,
    email text NOT NULL,
    tenant_id uuid references base.tenants(id)
);
ALTER TABLE base.users ADD CONSTRAINT "user_email_unique" UNIQUE ("email");