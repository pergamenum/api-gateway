package utilities

// TODO: Move to go-utils-gin.
import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ErrorResponse(ctx *gin.Context, status int, message ...any) {

	var str string
	if len(message) == 1 {
		m := message[0]
		if v, ok := m.(string); ok {
			str = v
		}
		if v, ok := m.(error); ok {
			var valErrs validator.ValidationErrors
			if errors.As(v, &valErrs) {
				var sb strings.Builder
				sb.WriteString("required: ")
				for _, e := range valErrs {
					sb.WriteString(fmt.Sprintf("(%s) ", e.Field()))
				}
				str = strings.TrimSpace(sb.String())
			} else {
				str = v.Error()
			}
		}
	}

	jo := struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}{
		Status:  status,
		Message: str,
	}
	ctx.AbortWithStatusJSON(status, jo)
}
