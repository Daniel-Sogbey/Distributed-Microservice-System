CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE payment_status as ENUM ('unspecified', 'pending', 'paid', 'failed');

CREATE TABLE IF NOT EXISTS orders (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    status payment_status DEFAULT 'unspecified',
    amount_cents bigint NOT NULL,
    created_unix bigint NOT NULL
);
