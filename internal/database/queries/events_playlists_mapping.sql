-- name: SetPlaylistForEvent :exec
INSERT INTO events_playlists_mapping (event_uuid, playlist_uuid) 
VALUES ($1, $2)
ON CONFLICT (event_uuid)
DO UPDATE SET playlist_uuid = $2;