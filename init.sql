CREATE USER niflheim WITH PASSWORD 'niflguard';
ALTER ROLE niflheim WITH CREATEDB;

DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'iot_db') THEN
        CREATE DATABASE iot_db;
    END IF;
END $$;

\c iot_db



CREATE TABLE IF NOT EXISTS "users" (
    "id" serial PRIMARY KEY,
    "email" varchar(255) UNIQUE NOT NULL,
    "password" varchar(255) NOT NULL,
    "key_card" varchar(64) UNIQUE,
    "access_level" INTEGER NOT NULL DEFAULT 0,
    "last_entered" Timestamp WITH TIME ZONE NOT NULL DEFAULT now(),
    "created_at" Timestamp WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "entry_logs" (
    "id" serial PRIMARY KEY,
    "key_card" varchar(64),
    "status" varchar(64) NOT NULL,
    "message" text,
    "created_at" Timestamp WITH TIME ZONE NOT NULL DEFAULT now(),
    FOREIGN KEY("key_card") REFERENCES users("key_card")
);

INSERT INTO users(email, password, access_level) VALUES('admin', '\x2432612431302441624866564a414239385776397258524674544c504f54696c6e6a75326932552e323868496f4f776931327261393738336e556c57', 10);
