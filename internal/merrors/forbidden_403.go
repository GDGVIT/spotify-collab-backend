package merrors

import (
	"net/http"

	"spotify-collab/internal/utils"

	"github.com/gin-gonic/gin"
)

/* -------------------------------------------------------------------------- */
/*                                FORBIDDEN 403                               */
/* -------------------------------------------------------------------------- */
func Forbidden(ctx *gin.Context, err string) {
	var res utils.BaseResponse
	var smerror utils.Error
	errorCode := http.StatusForbidden
	smerror.Code = errorCode
	smerror.Type = errorType.Forbidden
	smerror.Message = err
	res.Error = smerror
	ctx.JSON(errorCode, res)
	ctx.Abort()
}
