CREATE TABLE IF NOT EXISTS public.playlist_configs (
	playlist_uuid uuid NOT NULL,
	explicit bool DEFAULT False NOT NULL,
	require_approval bool DEFAULT true NOT NULL,
	max_song integer DEFAULT 5 NOT NULL,
    CONSTRAINT playlist_configs_playlist_fk FOREIGN KEY (playlist_uuid) REFERENCES public.playlists(playlist_uuid) ON DELETE CASCADE ON UPDATE CASCADE
);