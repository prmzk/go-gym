package response

import (
	"net/http"

	"github.com/go-chi/render"
)

type JSONResponse struct {
	HTTPStatusCode int          `json:"-"`
	StatusText     string       `json:"status"`
	Data           *interface{} `json:"data,omitempty"`
	Message        string       `json:"message,omitempty"`
}

func (response *JSONResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, response.HTTPStatusCode)
	return nil
}

func SuccessResponseOK(data interface{}) render.Renderer {
	return &JSONResponse{
		HTTPStatusCode: http.StatusOK,
		StatusText:     "success",
		Data:           &data,
	}
}

func SuccessResponseCreated(data interface{}) render.Renderer {
	return &JSONResponse{
		HTTPStatusCode: http.StatusCreated,
		StatusText:     "success",
		Data:           &data,
	}
}

func ErrorResponseBadRequest(err error) render.Renderer {
	errorText := http.StatusText(http.StatusBadRequest)
	if err != nil {
		errorText = err.Error()
	}

	return &JSONResponse{
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "error",
		Message:        errorText,
	}
}

func ErrorResponseUnauthorized(err error) render.Renderer {
	errorText := http.StatusText(http.StatusUnauthorized)
	if err != nil {
		errorText = err.Error()
	}

	return &JSONResponse{
		HTTPStatusCode: http.StatusUnauthorized,
		StatusText:     "error",
		Message:        errorText,
	}
}

func ErrorResponseInternalServerError() render.Renderer {
	return &JSONResponse{
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "error",
		Message:        http.StatusText(http.StatusInternalServerError),
	}
}

func ErrorResponseNotFound() render.Renderer {
	return &JSONResponse{
		HTTPStatusCode: http.StatusNotFound,
		StatusText:     "error",
		Message:        http.StatusText(http.StatusNotFound),
	}
}
