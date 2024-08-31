-- +goose Up
CREATE TABLE IF NOT EXISTS users ( 
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY, 
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(), 
    first_name text NOT NULL, 
    last_name text NOT NULL,
    date_of_birth DATE NOT NULL,
    email citext UNIQUE NOT NULL, 
    password_hash bytea NOT NULL, 
    activated bool NOT NULL, 
    isAdmin bool NOT NULL DEFAULT FALSE, 
    version integer NOT NULL DEFAULT 1 );
    
-- +goose Down
DROP TABLE IF EXISTS users;