package respone

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neathhatake/the_Scan/pkg/apperrors"
	"github.com/neathhatake/the_Scan/pkg/logger"
)


type Envelop struct {
	Success bool `json:"success"`
	Message string	`json:"message,omitempty"` 
	Data any		`json:"data,omitempty"`
	Meta map[string]any `json:"meta,omitempty"`
	Error string `json:"error,omitempty"`
}


func OK(c *gin.Context , data any){
	c.JSON(http.StatusOK,Envelop{Success: true, Data: data})
}

func Created(c *gin.Context, data any){
	c.JSON(http.StatusCreated, Envelop{
		Success: true,
		Data: data,
	})
}

func NoContent(c *gin.Context){
	c.Status(http.StatusNoContent)
}

func Message(c *gin.Context, msg string){
	c.JSON(http.StatusOK, Envelop{
		Success:  true,
		Message: msg,
	})
}


func Error(c *gin.Context, err error){
	var appErr *apperrors.AppError
	if errors.As(err , &appErr){
		if appErr.Code >= 500 {
			logger.Log.Errorw("internal errror", "path",c.FullPath(), "method",c.Request.Method,"cause",appErr.Err,)
			c.JSON(appErr.Code, Envelop{Success: false, Error: appErr.Message})
			return 	
		}
		c.JSON(appErr.Code, Envelop{Success: false, Error: appErr.Message})
		return 
	}
	
	
	//Untype Error -> 500
	logger.Log.Errorw("unhandle error","path",c.FullPath(), "method",c.Request.Method,"cause",err,)
	c.JSON(http.StatusInternalServerError, Envelop{Success: false, Error: http.StatusText(http.StatusInternalServerError)}) //http.StatusText(http.StatusInternalServerErrron)
	

}


// ValidationError converts a binding validation error to an AppError.
func ValidationError(err error) *apperrors.AppError {
	return apperrors.BadRequest(err.Error())
}

// OKPaginated sends a 200 response with paginated data.
func OKPaginated(c *gin.Context, data any, total int64, page, limit int) {
	c.JSON(http.StatusOK, Envelop{Success: true, Data: data, Meta: map[string]any{
		"total": total,
		"page":  page,
		"limit": limit,
	}})
}
