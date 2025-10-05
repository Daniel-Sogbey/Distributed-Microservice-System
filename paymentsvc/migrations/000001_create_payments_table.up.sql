CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS payments(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id uuid NOT NULL,
    status text NOT NULL CHECK (status IN ('unspecified', 'pending', 'paid', 'failed')),
    amount_cents bigint NOT NULL,
    created_unix bigint NOT NULL
);
