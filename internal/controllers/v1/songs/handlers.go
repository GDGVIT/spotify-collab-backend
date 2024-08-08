package songs

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"spotify-collab/internal/database"
	"spotify-collab/internal/merrors"
	"spotify-collab/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SongHandler struct {
	db *pgxpool.Pool
}

func Handler(db *pgxpool.Pool) *SongHandler {
	return &SongHandler{
		db: db,
	}
}

func (s *SongHandler) AddSongToEvent(c *gin.Context) {
	req, err := validateAddSongToEventReq(c)
	if err != nil {
		merrors.Validation(c, err.Error())
	}

	q := database.New(s.db)
	event, err := q.GetEventUUIDByCode(c, req.EventCode)
	if errors.Is(sql.ErrNoRows, err) {
		merrors.NotFound(c, "Event not found")
	} else if err != nil {
		merrors.InternalServer(c, err.Error())
	}

	playlist, err := q.GetPlaylistUUIDByEventUUID(c, event)
	if errors.Is(sql.ErrNoRows, err) {
		merrors.NotFound(c, "no playlist found")
	} else if err != nil {
		merrors.InternalServer(c, err.Error())
	}

	// TODO: Check if valid song, passes config -> not greater than count, not blacklisted, other configs
	_, err = q.AddSong(c, database.AddSongParams{
		SongUri:      req.SongURI,
		PlaylistUuid: playlist,
	})
	if err != nil {
		merrors.InternalServer(c, err.Error())
	}

	c.JSON(http.StatusOK, utils.BaseResponse{
		Success:    true,
		Message:    "Song successfully added",
		StatusCode: http.StatusOK,
	})
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

func (s *SongHandler) GetAllSongs(c *gin.Context) {
	req, err := validateGetAllSongsReq(c)
	if err != nil {
		merrors.Validation(c, err.Error())
		return
	}

	q := database.New(s.db)
	songs, err := q.GetAllSongs(c, req.PlaylistUUID)
	if errors.Is(err, sql.ErrNoRows) {
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
func (s *SongHandler) GetBlacklistedSongs(c *gin.Context) {
	req, err := validateGetAllSongsReq(c)
	if err != nil {
		merrors.Validation(c, err.Error())
		return
	}

	q := database.New(s.db)
	songs, err := q.GetAllBlacklisted(c, req.PlaylistUUID)
	if errors.Is(err, sql.ErrNoRows) {
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

func (s *SongHandler) AddSongToPlaylist(c *gin.Context) {
	req, err := validateAddSongToPlaylist(c)
	if err != nil {
		merrors.Validation(c, err.Error())
		return
	}

	if req.AccessToken == "" {
		merrors.Validation(c, "Access token is required")
		return
	}

	if req.PlaylistID == "" {
		merrors.Validation(c, "Playlist ID is required")
		return
	}

	if len(req.SongURIList) == 0 {
		merrors.Validation(c, "At least one song URI is required")
		return
	}

	requestBody := RequestBody{Uris: req.SongURIList}
	body, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error marshaling request body:", err)
		return
	}

	fmt.Println("Request Body:", string(body))

	url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", req.PlaylistID)

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

	// Debug: Print the response status and body for further analysis
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
