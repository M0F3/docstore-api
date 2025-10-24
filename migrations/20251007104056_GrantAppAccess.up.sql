GRANT USAGE ON SCHEMA base TO app;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA base TO app;
ALTER TABLE base.tenants ENABLE ROW LEVEL SECURITY;
ALTER TABLE base.users ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_id_tenant ON base.tenants TO app
    USING (id = current_setting('app.current_tenant_id')::uuid);
CREATE POLICY tenant_id_users ON base.users TO app
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid);