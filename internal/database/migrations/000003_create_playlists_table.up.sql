CREATE TABLE public.playlists (
	user_uuid uuid NOT NULL,
	playlist_uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
	playlist_id text NOT NULL,
	"name" citext NOT NULL,
	created_at timestamptz DEFAULT now() NOT NULL,
	updated_at timestamptz DEFAULT now() NOT NULL,
	CONSTRAINT playlists_pk PRIMARY KEY (playlist_uuid),
	CONSTRAINT playlists_name_unique UNIQUE (user_uuid, name),
	CONSTRAINT playlists_id_unique UNIQUE (playlist_id),
	CONSTRAINT playlists_users_fk FOREIGN KEY (user_uuid) REFERENCES public.users(user_uuid) ON DELETE CASCADE ON UPDATE CASCADE
);