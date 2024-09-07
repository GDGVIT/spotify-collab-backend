package playlists

import (
	"errors"
	"fmt"
	"net/http"
	"spotify-collab/internal/controllers/v1/auth"
	"spotify-collab/internal/database"
	"spotify-collab/internal/merrors"
	"spotify-collab/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

type PlaylistHandler struct {
	db          *pgxpool.Pool
	spotifyauth *spotifyauth.Authenticator
}

func Handler(db *pgxpool.Pool, spotifyAuth *spotifyauth.Authenticator) *PlaylistHandler {
	return &PlaylistHandler{
		db:          db,
		spotifyauth: spotifyAuth,
	}
}

func (p *PlaylistHandler) CreatePlaylist(c *gin.Context) {
	req, err := validateCreatePlaylist(c)
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

	tx, err := p.db.Begin(c)
	if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}
	defer tx.Rollback(c)
	qtx := database.New(p.db).WithTx(tx)

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

	if !oauthToken.Valid() {
		oauthToken, err = p.spotifyauth.RefreshToken(c, oauthToken)
		if err != nil {
			merrors.InternalServer(c, fmt.Sprintf("Couldn't get access token %s", err))
			return
		}

		_, err := qtx.UpdateToken(c, database.UpdateTokenParams{
			Refresh:  []byte(oauthToken.RefreshToken),
			Access:   []byte(oauthToken.AccessToken),
			UserUuid: user.UserUUID,
		})
		if err != nil {
			merrors.InternalServer(c, err.Error())
			return
		}
	}

	// Generate Playlist from spotify
	client := spotify.New(p.spotifyauth.Client(c, oauthToken))
	spotifyPlaylist, err := client.CreatePlaylistForUser(c, token.SpotifyID, req.Name, "", true, false)
	if err != nil {
		merrors.InternalServer(c, fmt.Sprintf("Error while creating spotify playlist: %s", err.Error()))
		return
	}

	playlist, err := qtx.CreatePlaylist(c, database.CreatePlaylistParams{
		PlaylistID:   spotifyPlaylist.ID.String(),
		UserUuid:     user.UserUUID,
		Name:         req.Name,
		PlaylistCode: GeneratePlaylistCode(6),
	})
	// Check name already exists
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
		Message:    "playlist successfully created",
		Data:       playlist,
		StatusCode: http.StatusOK,
	})
}

func (p *PlaylistHandler) ListPlaylists(c *gin.Context) {
	u, ok := c.Get("user")
	if !ok {
		panic(" user failed to set in context ")
	}
	user := u.(*auth.ContextUser)
	if user == auth.AnonymousUser {
		merrors.Unauthorized(c, "This action is forbidden.")
		return
	}

	q := database.New(p.db)
	playlists, err := q.ListPlaylists(c, user.UserUUID)
	if len(playlists) == 0 {
		merrors.NotFound(c, "No Playlists exist!")
		return
	} else if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BaseResponse{
		Success:    true,
		Message:    "Playlists successfully retrieved",
		Data:       playlists,
		StatusCode: http.StatusOK,
	})
}

func (p *PlaylistHandler) GetPlaylist(c *gin.Context) {
	req, err := validateGetPlaylistReq(c)
	if err != nil {
		merrors.Validation(c, err.Error())
		return
	}
	// No need for check err since binding checks uuid
	uuid, _ := uuid.Parse(req.PlaylistUUID)

	q := database.New(p.db)
	playlist, err := q.GetPlaylist(c, uuid)
	if errors.Is(err, pgx.ErrNoRows) {
		merrors.NotFound(c, "Playlist not found")
		return
	} else if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BaseResponse{
		Success:    true,
		Message:    "Playlist successfully retrieved",
		Data:       playlist,
		StatusCode: http.StatusOK,
	})
}

func (p *PlaylistHandler) UpdatePlaylist(c *gin.Context) {
	req, err := validateUpdatePlaylistReq(c)
	if err != nil {
		merrors.Validation(c, err.Error())
		return
	}
	// No need for check err since binding checks uuid
	uuid, _ := uuid.Parse(req.PlaylistUUID)

	q := database.New(p.db)
	playlist, err := q.UpdatePlaylistName(c, database.UpdatePlaylistNameParams{
		Name:         req.Name,
		PlaylistUuid: uuid,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		merrors.NotFound(c, "Playlist not found")
		return
	} else if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BaseResponse{
		Success:    true,
		Message:    "Playlist successfully updated",
		Data:       playlist,
		StatusCode: http.StatusOK,
	})
}

func (p *PlaylistHandler) DeletePlaylist(c *gin.Context) {
	req, err := validateDeletePlaylistReq(c)
	if err != nil {
		merrors.Validation(c, err.Error())
		return
	}
	// No need for check err since binding checks uuid
	uuid, _ := uuid.Parse(req.PlaylistUUID)

	q := database.New(p.db)
	rows, err := q.DeletePlaylist(c, uuid)
	if rows == 0 {
		merrors.NotFound(c, "Playlist not found")
		return
	} else if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BaseResponse{
		Success:    true,
		Message:    "Playlist successfully deleted",
		StatusCode: http.StatusOK,
	})
}
