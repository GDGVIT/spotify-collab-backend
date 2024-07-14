CREATE TABLE IF NOT EXISTS public.events (
	user_uuid uuid NOT NULL,
	event_uuid uuid DEFAULT uuid_generate_v4() NOT NULL,
	created_at timestamptz DEFAULT now() NOT NULL,
	updated_at timestamptz DEFAULT now() NOT NULL,
	name citext NOT NULL,
	event_code citext NOT NULL,
	CONSTRAINT events_pk PRIMARY KEY (event_uuid),
	CONSTRAINT events_name_unique UNIQUE (event_uuid, name),
	CONSTRAINT events_code_unique UNIQUE (event_code),
	CONSTRAINT events_users_fk FOREIGN KEY (user_uuid) REFERENCES public.users(user_uuid) ON UPDATE CASCADE
);