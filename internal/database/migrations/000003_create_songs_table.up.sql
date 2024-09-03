CREATE TABLE IF NOT EXISTS public.songs (
	song_uri text NOT NULL,
	playlist_uuid uuid NOT NULL,
	count integer DEFAULT 1 NOT NULL,
	CONSTRAINT songs_pk PRIMARY KEY (song_uri, playlist_uuid),
	CONSTRAINT songs_playlists_fk FOREIGN KEY (playlist_uuid) REFERENCES public.playlists(playlist_uuid) ON DELETE CASCADE ON UPDATE CASCADE
);
