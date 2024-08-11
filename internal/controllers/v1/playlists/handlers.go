package playlists

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"spotify-collab/internal/database"
	"spotify-collab/internal/merrors"
	"spotify-collab/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
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
		PlaylistID: "f",
		UserUuid:   req.UserUUID,
		Name:       req.Name,
	})
	// Check name already exists
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
	// No need for check err since binding checks uuid
	uuid, _ := uuid.Parse(req.PlaylistUUID)

	q := database.New(p.db)
	playlist, err := q.GetPlaylist(c, uuid)
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
	// No need for check err since binding checks uuid
	uuid, _ := uuid.Parse(req.PlaylistUUID)

	q := database.New(p.db)
	playlist, err := q.UpdatePlaylistName(c, database.UpdatePlaylistNameParams{
		Name:         req.Name,
		PlaylistUuid: uuid,
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

// currently running get id everytime can remove it later if when creating account we associate user directly with spotify id in db

func GetUserId(AccessToken string) (string, error) {

	url := "https://api.spotify.com/v1/me"

	spotifyReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	spotifyReq.Header.Set("Authorization", "Bearer "+AccessToken)
	spotifyReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(spotifyReq)
	if err != nil {
		fmt.Println("Error making request:", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-OK HTTP status: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("error parsing JSON: %v", err)
	}

	userId, ok := result["id"].(string)
	if !ok {
		return "", fmt.Errorf("user ID not found in response")
	}

	return userId, nil
}

// Create playlist

func (p *PlaylistHandler) CreatePlaylistSpotify(c *gin.Context) {
	req, err := validateCreatePlaylistSpotifyReq(c)
	if err != nil {
		merrors.Validation(c, err.Error())
		return
	}

	if req.AccessToken == "" {
		merrors.Validation(c, "Access token is required")
		return
	}

	if req.PlaylistName == "" {
		merrors.Validation(c, "Playlist Name is required")
	}

	uid, err := GetUserId(req.AccessToken)
	if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	requestBody := CreatePlaylistSpotifyReqBody{
		PlaylistName:    req.PlaylistName,
		IsPublic:        req.IsPublic || true,
		IsCollaborative: req.IsCollaborative || false,
		Description:     req.Description,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error marshaling request body:", err)
		return
	}

	fmt.Println("Request Body:", string(body))

	url := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", uid)

	fmt.Println("URL:", url)

	spotifyReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set headers
	spotifyReq.Header.Set("Authorization", "Bearer "+req.AccessToken)
	spotifyReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(spotifyReq)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(responseBody))

	if resp.StatusCode == http.StatusOK {
		var responseBodyMap map[string]interface{}
		err := json.NewDecoder(bytes.NewBuffer(responseBody)).Decode(&responseBodyMap)
		if err != nil {
			fmt.Println("Error decoding response body:", err)
			return
		}

		snapshotID := responseBodyMap["snapshot_id"].(string)
		fmt.Println("Added!!\n", snapshotID)
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Error: Unauthorized - Invalid or expired access token")
	} else {
		fmt.Println("Error:", resp.Status, string(responseBody))
	}

}
