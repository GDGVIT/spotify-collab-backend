CREATE TABLE IF NOT EXISTS public.songs (
	song_uri text NOT NULL,
	playlist_id text NOT NULL,
	count integer DEFAULT 1 NOT NULL,
	CONSTRAINT songs_pk PRIMARY KEY (playlist_id, song_uri),
	CONSTRAINT songs_playlists_fk FOREIGN KEY (playlist_id) REFERENCES public.playlists(playlist_id) ON DELETE CASCADE ON UPDATE CASCADE
);
