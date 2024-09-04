-- CREATE TABLE IF NOT EXISTS tokens (
--     hash bytea PRIMARY KEY,
--     user_uuid uuid NOT NULL,
--     expiry timestamp(0) with time zone NOT NULL,
--     scope text NOT NULL,
--     CONSTRAINT tokens_users_fk FOREIGN KEY (user_uuid) REFERENCES public.users(user_uuid) ON UPDATE CASCADE
-- );

CREATE TABLE IF NOT EXISTS tokens (
    user_uuid uuid NOT NULL,
    refresh bytea PRIMARY KEY,
    access bytea NOT NULL, 
    expiry timestamp(0) with time zone NOT NULL,
    CONSTRAINT tokens_users_fk FOREIGN KEY (user_uuid) REFERENCES public.users(user_uuid) ON UPDATE CASCADE ON DELETE CASCADE
)