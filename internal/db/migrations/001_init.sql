CREATE TABLE IF NOT EXISTS schema_versions (
  version INT PRIMARY KEY,
  applied_at TIMESTAMPTZ DEFAULT NOW()
INSERT INTO schema_versions (version) VALUES (1) 
ON CONFLICT (version) DO NOTHING;
);

INSERT INTO schema_versions (version) VALUES (1) 
ON CONFLICT (version) DO NOTHING;

--  !!ВЫПОЛНЯТЬ ВРУЧНУЮ!!