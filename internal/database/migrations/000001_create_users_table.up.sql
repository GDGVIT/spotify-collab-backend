CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS public.users (
	id bigserial NOT NULL,
	user_uuid uuid DEFAULT uuid_generate_v4() NOT NULL,
	spotify_id TEXT NOT NULL,
	created_at timestamp with time zone DEFAULT Now() NOT NULL,
	updated_at timestamp with time zone DEFAULT Now() NOT NULL,
	name text,
	email public."citext" NOT NULL,
	password_hash bytea NOT NULL,
	activated boolean DEFAULT false NOT NULL,
	"version" integer DEFAULT 1 NOT NULL,
	CONSTRAINT users_email_unique UNIQUE (email),
	CONSTRAINT users_spotify_id_unique UNIQUE (spotify_id),
	CONSTRAINT users_pk PRIMARY KEY (user_uuid)
);
