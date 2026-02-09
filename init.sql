CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE service (
	service_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	service_name TEXT NOT NULL,
	service_url TEXT NOT NULL,
	service_created_at TIMESTAMP DEFAULT now()
);
