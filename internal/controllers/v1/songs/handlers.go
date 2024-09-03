package songs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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

func (s *SongHandler) AcceptSongToPlaylist(c *gin.Context) {
	// _, err = qtx.AddSong(c, database.AddSongParams{
	// 	SongUri:      req.SongURI,
	// 	PlaylistUuid: playlist,
	// })
	// if err != nil {
	// 	merrors.InternalServer(c, err.Error())
	// }
}

// Adding a song through spotify api to the playlist
func (s *SongHandler) AddSongToPlaylist(c *gin.Context) {
	req, err := validateAddSongToPlaylistReq(c)
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

	playlist, err := qtx.GetPlaylistIDByCode(c, req.PlaylistCode)
	if errors.Is(err, pgx.ErrNoRows) {
		merrors.NotFound(c, "no playlist found")
		return
	} else if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}
	// TODO: Get UserUUID from the context instead of request body
	userUUID, _ := uuid.Parse("f76f7a84-6a5a-4f49-892b-b40864ce7165")

	token, err := qtx.GetOAuthToken(c, userUUID)
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
		oauthToken, err = s.spotifyauth.RefreshToken(c, oauthToken)
		if err != nil {
			merrors.InternalServer(c, fmt.Sprintf("Couldn't get access token %s", err))
			return
		}

		_, err := qtx.UpdateToken(c, database.UpdateTokenParams{
			Refresh:  []byte(oauthToken.RefreshToken),
			Access:   []byte(oauthToken.AccessToken),
			UserUuid: userUUID,
		})
		if err != nil {
			merrors.InternalServer(c, err.Error())
			return
		}
	}

	client := spotify.New(s.spotifyauth.Client(c, oauthToken))
	_, err = client.AddTracksToPlaylist(c, spotify.ID(playlist), spotify.ID(req.SongURI))
	if err != nil {
		merrors.InternalServer(c, fmt.Sprintf("Error while adding to playlist: %s", err.Error()))
		return
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

func (s *SongHandler) KaranAddSongToPlaylist(c *gin.Context) {
	req, err := validateKaranAddSongToPlaylist(c)
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
