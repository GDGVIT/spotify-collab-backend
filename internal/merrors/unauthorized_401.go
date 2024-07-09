package merrors

import (
	"net/http"

	"spotify-collab/internal/utils"

	"github.com/gin-gonic/gin"
)

/* -------------------------------------------------------------------------- */
/*                           Unauthorized Error 401                           */
/* -------------------------------------------------------------------------- */
func Unauthorized(ctx *gin.Context, err string) {
	var res utils.BaseResponse
	var smerror utils.Error
	errorCode := http.StatusUnauthorized
	smerror.Code = errorCode
	smerror.Type = errorType.Unauthorized
	smerror.Message = err

	res.Error = smerror

	ctx.JSON(errorCode, res)
	ctx.Abort()
}
