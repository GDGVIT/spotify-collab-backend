package auth

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"spotify-collab/internal/database"
	"spotify-collab/internal/merrors"
	"spotify-collab/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const state = "abc123"

type AuthHandler struct {
	db               *pgxpool.Pool
	spotifyauth      *spotifyauth.Authenticator
	frontendCallback string
	appCallback      string
}

func Handler(db *pgxpool.Pool, spotifyAuth *spotifyauth.Authenticator) *AuthHandler {
	return &AuthHandler{
		db:               db,
		spotifyauth:      spotifyAuth,
		frontendCallback: os.Getenv("FRONTEND_CALLBACK_URL"),
		appCallback:      os.Getenv("APP_CALLBACK_URL"),
	}
}

func (a *AuthHandler) SpotifyLogin(c *gin.Context) {
	var request struct {
		Platform string `uri:"platform"`
	}

	err := c.ShouldBindUri(&request)
	if err != nil {
		merrors.Validation(c, err.Error())
		return
	}
	if request.Platform != "app" && request.Platform != "web" {
		merrors.Validation(c, "either app or web (case sensitive)")
		return
	}

	state2 := state + "-" + request.Platform
	url := a.spotifyauth.AuthURL(state2)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (a *AuthHandler) SpotifyCallback(c *gin.Context) {
	st := c.Request.FormValue("state")
	stSplit := strings.Split(st, "-")
	if stSplit[0] != state || (stSplit[1] != "app" && stSplit[1] != "web") {
		merrors.NotFound(c, "State Mismatch!")
		return
	}

	state2 := state + "-" + stSplit[1]
	tok, err := a.spotifyauth.Token(c, state2, c.Request)
	if err != nil {
		merrors.Forbidden(c, fmt.Sprintf("Couldn't get token: %s", err.Error()))
		return
	}

	client := spotify.New(a.spotifyauth.Client(c, tok))
	user, err := client.CurrentUser(c)
	if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}
	spotifyID := user.ID

	var userUUID uuid.UUID

	tx, err := a.db.Begin(c)
	if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}
	defer tx.Rollback(c)
	qtx := database.New(a.db).WithTx(tx)

	userUUID, err = qtx.GetUserBySpotifyID(c, spotifyID)
	if errors.Is(err, pgx.ErrNoRows) {
		// If not, register a new user
		usr, err := qtx.CreateUser(c, database.CreateUserParams{
			Email:     user.Email,
			SpotifyID: spotifyID,
			Name:      user.DisplayName,
		})
		userUUID = usr.UserUuid
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			merrors.Validation(c, "User already exists with this spotify ID!")
			return
		} else if err != nil {
			merrors.InternalServer(c, err.Error())
			return
		}
	} else if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	_, err = qtx.NewToken(c, database.NewTokenParams{
		Refresh:  []byte(tok.RefreshToken),
		Access:   []byte(tok.AccessToken),
		Expiry:   pgtype.Timestamptz{Time: tok.Expiry, Valid: true},
		UserUuid: userUUID,
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

	queryVal := url.Values{
		"token": {tok.AccessToken},
		"user":  {userUUID.String()},
	}

	var reqURL string
	if stSplit[1] == "web" {
		reqURL = a.frontendCallback + "?" + queryVal.Encode()
	} else if stSplit[1] == "app" {
		reqURL = a.appCallback + "?" + queryVal.Encode()
	}

	// queryVal := url.Values{
	// 	"token": {tok.Plaintext},
	// }
	// reqUrl := a.frontendCallback + "?" + queryVal.Encode()
	c.Redirect(http.StatusTemporaryRedirect, reqURL)

	c.JSON(http.StatusOK, utils.BaseResponse{
		Success:    true,
		Message:    "Spotify user successfully authenticated",
		Data:       gin.H{"token": tok, "user": userUUID},
		StatusCode: http.StatusOK,
	})
}
