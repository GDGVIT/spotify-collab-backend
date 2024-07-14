CREATE TABLE IF NOT EXISTS public.event_configs (
	event_uuid uuid NOT NULL,
	explicit bool DEFAULT False NOT NULL,
	require_approval bool DEFAULT true NOT NULL,
	max_song integer DEFAULT 5 NOT NULL,
    CONSTRAINT event_configs_events_fk FOREIGN KEY (event_uuid) REFERENCES public.events(event_uuid) ON DELETE CASCADE ON UPDATE CASCADE
);