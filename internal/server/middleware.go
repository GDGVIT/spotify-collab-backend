package server

import (
	"errors"
	"spotify-collab/internal/controllers/v1/auth"
	"spotify-collab/internal/database"
	"spotify-collab/internal/merrors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (s *Server) authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Add("Vary", "Authorization")

		authorizationHeader := ctx.Request.Header.Get("Authorization")

		if authorizationHeader == "" {
			ctx.Set("user", auth.AnonymousUser)
			ctx.Next()
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			ctx.Writer.Header().Set("WWW-Authenticate", "Bearer")
			merrors.Unauthorized(ctx, "invalid or missing authentication token.")
			return
		}

		token := headerParts[1]

		q := database.New(s.db)

		user, err := q.GetUserByToken(ctx, []byte(token))
		if err != nil {
			switch {
			case errors.Is(err, pgx.ErrNoRows):
				ctx.Writer.Header().Set("WWW-Authenticate", "Bearer")
				merrors.Unauthorized(ctx, "invalid or missing authentication token. Please log in again.")
			default:
				merrors.InternalServer(ctx, err.Error())
			}
			return
		}
		contextUser := auth.ContextUser{
			UserUUID:  user.UserUuid,
			SpotifyID: user.SpotifyID,
		}

		ctx.Set("user", &contextUser)

		ctx.Next()
	}
}
