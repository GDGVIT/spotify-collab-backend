package merrors

import (
	"net/http"

	"spotify-collab/internal/utils"

	"github.com/gin-gonic/gin"
)

/* -------------------------------------------------------------------------- */
/*                                Conflict 409                                */
/* -------------------------------------------------------------------------- */

func NotFound(ctx *gin.Context, err string) {
	var res utils.BaseResponse
	var smerror utils.Error
	errorCode := http.StatusNotFound
	smerror.Code = errorCode
	smerror.Type = errorType.NotFound
	smerror.Message = err
	res.Error = &smerror
	ctx.JSON(errorCode, res)
	ctx.Abort()
}
