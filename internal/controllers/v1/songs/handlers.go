package songs

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"spotify-collab/internal/controllers/v1/auth"
	"spotify-collab/internal/database"
	"spotify-collab/internal/merrors"
	"spotify-collab/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

type SongHandler struct {
	db          *pgxpool.Pool
	spotifyauth *spotifyauth.Authenticator
}

func Handler(db *pgxpool.Pool, spotifyAuth *spotifyauth.Authenticator) *SongHandler {
	return &SongHandler{
		db:          db,
		spotifyauth: spotifyAuth,
	}
}

// Participant adds song
func (s *SongHandler) AddSongToDB(c *gin.Context) {
	req, err := validateAddSongToDBReq(c)
	if err != nil {
		merrors.Validation(c, err.Error())
		return
	}

	tx, err := s.db.Begin(c)
	if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}
	defer tx.Rollback(c)
	qtx := database.New(s.db).WithTx(tx)

	playlist, err := qtx.GetPlaylistUUIDByCode(c, req.PlaylistCode)
	if errors.Is(err, pgx.ErrNoRows) {
		merrors.NotFound(c, "no playlist found")
		return
	} else if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	song, err := qtx.AddSong(c, database.AddSongParams{
		SongUri:      req.SongURI,
		PlaylistUuid: playlist,
	})
	if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	err = tx.Commit(c)
	if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BaseResponse{
		Success:    true,
		Message:    "Song successfully added!",
		Data:       song,
		StatusCode: http.StatusOK,
	})
}

// Adding a song through spotify api to the playlist
func (s *SongHandler) AddSongToPlaylist(c *gin.Context) {
	req, err := validateAddSongToPlaylistReq(c)
	if err != nil {
		merrors.Validation(c, err.Error())
		return
	}

	u, ok := c.Get("user")
	if !ok {
		panic(" user failed to set in context ")
	}
	user := u.(*auth.ContextUser)
	if user == auth.AnonymousUser {
		merrors.Unauthorized(c, "This action is forbidden.")
		return
	}

	tx, err := s.db.Begin(c)
	if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}
	defer tx.Rollback(c)
	qtx := database.New(s.db).WithTx(tx)

	var message string
	message = "song rejected"

	tokenChanged := false
	var oauthToken *oauth2.Token

	if req.Option == "accept" {
		message = "song successfully added"

		token, err := qtx.GetOAuthToken(c, user.UserUUID)
		if err != nil {
			merrors.InternalServer(c, err.Error())
			return
		}

		oauthToken = &oauth2.Token{
			AccessToken:  string(token.Access),
			RefreshToken: string(token.Refresh),
			Expiry:       token.Expiry.Time,
		}

		if !oauthToken.Valid() {
			oauthToken, err = s.spotifyauth.RefreshToken(c, oauthToken)
			tokenChanged = true
			if err != nil {
				merrors.InternalServer(c, fmt.Sprintf("Couldn't get access token %s", err))
				return
			}

			_, err := qtx.UpdateToken(c, database.UpdateTokenParams{
				Refresh:  []byte(oauthToken.RefreshToken),
				Access:   []byte(oauthToken.AccessToken),
				Expiry:   pgtype.Timestamptz{Time: oauthToken.Expiry, Valid: true},
				UserUuid: user.UserUUID,
			})
			if err != nil {
				merrors.InternalServer(c, err.Error())
				return
			}
		}

		playlist, err := qtx.GetPlaylistIDByUUID(c, req.PlaylistUUID)
		if errors.Is(err, pgx.ErrNoRows) {
			merrors.NotFound(c, "no playlist found")
			return
		} else if err != nil {
			merrors.InternalServer(c, err.Error())
			return
		}

		client := spotify.New(s.spotifyauth.Client(c, oauthToken))
		_, err = client.AddTracksToPlaylist(c, spotify.ID(playlist), spotify.ID(req.SongURI))
		if err != nil {
			merrors.InternalServer(c, fmt.Sprintf("Error while adding to playlist: %s", err.Error()))
			return
		}
	}

	err = qtx.AddSongToPlaylist(c, database.AddSongToPlaylistParams{
		SongUri:      req.SongURI,
		PlaylistUuid: req.PlaylistUUID,
	})
	if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	err = tx.Commit(c)
	if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	if tokenChanged == true {
		c.JSON(http.StatusOK, utils.BaseResponse{
			Success:    true,
			Message:    message,
			MetaData:   oauthToken,
			StatusCode: http.StatusOK,
		})
	} else {
		c.JSON(http.StatusOK, utils.BaseResponse{
			Success:    true,
			Message:    message,
			StatusCode: http.StatusOK,
		})
	}
}

func (s *SongHandler) GetAllSongs(c *gin.Context) {
	req, err := validateGetAllSongsReq(c)
	if err != nil {
		merrors.Validation(c, err.Error())
		return
	}

	u, ok := c.Get("user")
	if !ok {
		panic(" user failed to set in context ")
	}
	user := u.(*auth.ContextUser)
	if user == auth.AnonymousUser {
		merrors.Unauthorized(c, "This action is forbidden.")
		return
	}

	tx, err := s.db.Begin(c)
	if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}
	defer tx.Rollback(c)
	qtx := database.New(s.db).WithTx(tx)

	token, err := qtx.GetOAuthToken(c, user.UserUUID)
	if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	oauthToken := &oauth2.Token{
		AccessToken:  string(token.Access),
		RefreshToken: string(token.Refresh),
		Expiry:       token.Expiry.Time,
	}

	tokenChanged := false
	if !oauthToken.Valid() {
		oauthToken, err = s.spotifyauth.RefreshToken(c, oauthToken)
		tokenChanged = true
		if err != nil {
			merrors.InternalServer(c, fmt.Sprintf("Couldn't get access token %s", err))
			return
		}

		_, err := qtx.UpdateToken(c, database.UpdateTokenParams{
			Refresh:  []byte(oauthToken.RefreshToken),
			Access:   []byte(oauthToken.AccessToken),
			Expiry:   pgtype.Timestamptz{Time: oauthToken.Expiry, Valid: true},
			UserUuid: user.UserUUID,
		})
		if err != nil {
			merrors.InternalServer(c, err.Error())
			return
		}
	}

	songs, err := qtx.GetAllSongs(c, req.PlaylistUUID)
	log.Println(errors.Is(err, pgx.ErrNoRows))
	log.Println(err)
	log.Println(songs)
	if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	playlist, err := qtx.GetPlaylistIDByUUID(c, req.PlaylistUUID)
	if errors.Is(err, pgx.ErrNoRows) {
		merrors.NotFound(c, "no playlist found")
		return
	} else if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	err = tx.Commit(c)
	if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	client := spotify.New(s.spotifyauth.Client(c, oauthToken))
	songIDs := []spotify.ID{}
	for _, v := range songs {
		songIDs = append(songIDs, spotify.ID(v.SongUri))
	}
	log.Println(songs)
	log.Println(songIDs)
	var tracks []*spotify.FullTrack
	if len(songIDs) != 0 {
		tracks, err = client.GetTracks(c, songIDs)
		if err != nil {
			merrors.InternalServer(c, err.Error())
			return
		}
	}
	log.Println(playlist)

	// offset := 0
	limit := 100
	var playlist_tracks []spotify.PlaylistItem

	// new_tracks, err := client.GetPlaylistItems(c, spotify.ID(playlist), spotify.Limit(limit), spotify.Fields("next,items(track(name,artists(name)))"))
	// new_tracks, err := client.GetPlaylistItems(c, spotify.ID(playlist), spotify.Limit(limit), spotify.Fields("next,items(track(name,artists(name)))"))
	// new_tracks, err := client.GetPlaylistItems(c, spotify.ID(playlist), spotify.Limit(limit), spotify.Fields("tracks.items(track(name,href,album(!name,href)))"))
	new_tracks, err := client.GetPlaylistItems(c, spotify.ID(playlist), spotify.Limit(limit))
	log.Println(new_tracks)
	if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}
	playlist_tracks = append(playlist_tracks, new_tracks.Items...)

	// for new_tracks.Next != "" {
	// 	offset += limit
	// 	new_tracks, err := client.GetPlaylistItems(c, spotify.ID(playlist), spotify.Limit(limit), spotify.Offset(offset), spotify.Fields("next,items(track(name,artists(name)))"))
	// 	if err != nil {
	// 		merrors.InternalServer(c, err.Error())
	// 		return
	// 	}
	// 	playlist_tracks = append(playlist_tracks, new_tracks.Items...)
	// }

	if tokenChanged == true {
		c.JSON(http.StatusOK, utils.BaseResponse{
			Success: true,
			Message: "Songs successfully retrieved",
			Data: gin.H{
				"submitted": tracks,
				"accepted":  playlist_tracks,
			},
			MetaData:   oauthToken,
			StatusCode: http.StatusOK,
		})
	} else {
		c.JSON(http.StatusOK, utils.BaseResponse{
			Success: true,
			Message: "Songs successfully retrieved",
			Data: gin.H{
				"submitted": tracks,
				"accepted":  playlist_tracks,
			},
			StatusCode: http.StatusOK,
		})
	}
}

func (s *SongHandler) BlacklistSong(c *gin.Context) {
	req, err := validateBlacklistSongReq(c)
	if err != nil {
		merrors.Validation(c, err.Error())
	}

	q := database.New(s.db)
	song, err := q.BlacklistSong(c, database.BlacklistSongParams{
		SongUri:      req.SongURI,
		PlaylistUuid: req.PlaylistUUID,
	})
	if song == 0 {
		merrors.NotFound(c, "song not found")
	} else if err != nil {
		merrors.InternalServer(c, err.Error())
	}

	c.JSON(http.StatusOK, utils.BaseResponse{
		Success: true,
		Message: "Song successfully blacklisted",
	})
}

func (s *SongHandler) GetBlacklistedSongs(c *gin.Context) {
	req, err := validateGetAllSongsReq(c)
	if err != nil {
		merrors.Validation(c, err.Error())
		return
	}

	q := database.New(s.db)
	songs, err := q.GetAllBlacklisted(c, req.PlaylistUUID)
	if errors.Is(err, pgx.ErrNoRows) {
		merrors.NotFound(c, "No Songs exist!")
		return
	} else if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BaseResponse{
		Success:    true,
		Message:    "Songs successfully retrieved",
		Data:       songs,
		StatusCode: http.StatusOK,
	})
}

func (s *SongHandler) DeleteBlacklistSong(c *gin.Context) {
	req, err := validateBlacklistSongReq(c)
	if err != nil {
		merrors.Validation(c, err.Error())
		return
	}

	q := database.New(s.db)
	song, err := q.DeleteBlacklist(c, database.DeleteBlacklistParams{
		PlaylistUuid: req.PlaylistUUID,
		SongUri:      req.SongURI,
	})
	if song == 0 {
		merrors.NotFound(c, "song not found!")
		return
	} else if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BaseResponse{
		Success:    true,
		Message:    "Song removed from blacklist",
		StatusCode: http.StatusOK,
	})
}
