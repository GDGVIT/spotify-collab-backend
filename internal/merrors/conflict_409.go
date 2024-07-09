package merrors

import (
	"net/http"

	"spotify-collab/internal/utils"

	"github.com/gin-gonic/gin"
)

/* -------------------------------------------------------------------------- */
/*                                Conflict 409                                */
/* -------------------------------------------------------------------------- */

func Conflict(ctx *gin.Context, err string) {
	var res utils.BaseResponse
	var smerror utils.Error
	errorCode := http.StatusConflict
	smerror.Code = errorCode
	smerror.Type = errorType.conflict
	smerror.Message = err
	res.Error = smerror
	ctx.JSON(errorCode, res)
	ctx.Abort()
}
