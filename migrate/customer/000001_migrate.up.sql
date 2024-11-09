CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE customers (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  first_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255) NOT NULL,
  age INTEGER NOT NULL,
  membership VARCHAR(255) NOT NULL,
  orders JSONB NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO customers (id, first_name, last_name, age, membership, orders) VALUES
('4abf6d8a-35f9-4989-ab4e-ed9e9b414c98','Jane', 'Doe', 30, 'premium', '[{"name": "Headphones", "description": "Blue headphones"}, {"name": "Mouse", "description": "Black mouse"}, {"name": "Monitor", "description": "Black monitor"}]');
