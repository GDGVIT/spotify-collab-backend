package merrors

import (
	"net/http"

	"spotify-collab/internal/utils"

	"github.com/gin-gonic/gin"
)

/* -------------------------------------------------------------------------- */
/*                            VALIDATION ERROR 422                            */
/* -------------------------------------------------------------------------- */
func Validation(ctx *gin.Context, err string) {
	var res utils.BaseResponse
	var smerror utils.Error
	errorCode := http.StatusUnprocessableEntity

	smerror.Code = errorCode
	smerror.Type = errorType.validation
	smerror.Message = err

	res.Error = smerror

	ctx.JSON(errorCode, res)
	ctx.Abort()
}
