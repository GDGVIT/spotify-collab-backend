CREATE TABLE public.events_playlists_mapping (
	event_uuid uuid NOT NULL,
	playlist_uuid uuid NOT NULL,
	CONSTRAINT events_playlists_mapping_events_fk FOREIGN KEY (event_uuid) REFERENCES public.events(event_uuid) ON DELETE CASCADE ON UPDATE CASCADE,
	CONSTRAINT events_playlists_mapping_playlists_fk FOREIGN KEY (playlist_uuid) REFERENCES public.playlists(playlist_uuid) ON DELETE CASCADE ON UPDATE CASCADE
);