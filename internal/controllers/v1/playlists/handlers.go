package playlists

import (
	"errors"
	"net/http"
	"spotify-collab/internal/database"
	"spotify-collab/internal/merrors"
	"spotify-collab/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlaylistHandler struct {
	db *pgxpool.Pool
}

func Handler(db *pgxpool.Pool) *PlaylistHandler {
	return &PlaylistHandler{
		db: db,
	}
}

func (p *PlaylistHandler) CreatePlaylist(c *gin.Context) {
	req, err := validateCreatePlaylist(c)
	if err != nil {
		merrors.Validation(c, err.Error())
		return
	}

	tx, err := p.db.Begin(c)
	if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}
	defer tx.Rollback(c)
	qtx := database.New(p.db).WithTx(tx)

	// Generate Playlist from spotify
	playlist, err := qtx.CreatePlaylist(c, database.CreatePlaylistParams{
		PlaylistID: "",
		UserUuid:   req.UserUUID,
		Name:       req.Name,
	})
	if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	config, err := qtx.CreateDefaultConfiguration(c, playlist.PlaylistUuid)
	if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	err = tx.Commit(c)
	if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"playlist": playlist,
		"config":   config,
	})
}

func (p *PlaylistHandler) ListPlaylists(c *gin.Context) {
	req, err := validateListPlaylistsReq(c)
	if err != nil {
		merrors.Validation(c, err.Error())
		return
	}

	q := database.New(p.db)
	playlists, err := q.ListPlaylists(c, req.UserUUID)
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

	q := database.New(p.db)
	playlist, err := q.GetPlaylist(c, req.PlaylistUUID)
	if errors.Is(pgx.ErrNoRows, err) {
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

	q := database.New(p.db)
	playlist, err := q.UpdatePlaylistName(c, database.UpdatePlaylistNameParams{
		Name:         req.Name,
		PlaylistUuid: req.PlaylistUUID,
	})
	if errors.Is(pgx.ErrNoRows, err) {
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

	q := database.New(p.db)
	rows, err := q.DeletePlaylist(c, req.PlaylistUUID)
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

func (p *PlaylistHandler) UpdateConfiguration(c *gin.Context) {
	req, err := validateUpdateConfigurationReq(c)
	if err != nil {
		binding.EnableDecoderDisallowUnknownFields = true
		merrors.Validation(c, err.Error())
		return
	}

	q := database.New(p.db)
	config, err := q.GetConfiguration(c, req.PlaylistUUID)
	if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	if req.Explicit == nil {
		req.Explicit = &config.Explicit
	}
	if req.RequireApproval == nil {
		req.RequireApproval = &config.RequireApproval
	}
	if req.MaxSong == 0 {
		req.MaxSong = config.MaxSong
	}

	params := database.UpdateConfigurationParams{
		PlaylistUuid:    req.PlaylistUUID,
		Explicit:        *req.Explicit,
		RequireApproval: *req.RequireApproval,
		MaxSong:         req.MaxSong,
	}

	config, err = q.UpdateConfiguration(c, params)
	if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BaseResponse{
		Success:    true,
		Message:    "Configuration successfully updated",
		Data:       config,
		StatusCode: http.StatusOK,
	})

}
