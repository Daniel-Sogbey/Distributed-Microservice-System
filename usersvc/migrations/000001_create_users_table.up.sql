CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
email text UNIQUE NOT NULL,
name text NOT NULL,
created_unix bigint NOT NULL
);
