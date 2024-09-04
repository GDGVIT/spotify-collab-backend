package auth

import (
	"errors"
	"fmt"
	"net/http"
	"spotify-collab/internal/database"
	"spotify-collab/internal/merrors"
	"spotify-collab/internal/utils"

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
	db          *pgxpool.Pool
	spotifyauth *spotifyauth.Authenticator
}

func Handler(db *pgxpool.Pool, spotifyAuth *spotifyauth.Authenticator) *AuthHandler {
	return &AuthHandler{
		db:          db,
		spotifyauth: spotifyAuth,
	}
}

func (a *AuthHandler) SpotifyLogin(c *gin.Context) {
	url := a.spotifyauth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (a *AuthHandler) SpotifyCallback(c *gin.Context) {
	tok, err := a.spotifyauth.Token(c, state, c.Request)
	if err != nil {
		merrors.Forbidden(c, fmt.Sprintf("Couldn't get token: %s", err.Error()))
		return
	}

	if st := c.Request.FormValue("state"); st != state {
		merrors.NotFound(c, "State Mismatch!")
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
			Email:        user.Email,
			SpotifyID:    spotifyID,
			Name:         user.DisplayName,
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

	// queryVal := url.Values{
	// 	"token": {tok.Plaintext},
	// }
	// reqUrl := a.frontendCallback + "?" + queryVal.Encode()
	// c.Redirect(http.StatusTemporaryRedirect, reqUrl)

	c.JSON(http.StatusOK, utils.BaseResponse{
		Success:    true,
		Message:    "Spotify user successfully authenticated",
		Data:       gin.H{"token": tok, "user": userUUID},
		StatusCode: http.StatusOK,
	})
}
