package middleware

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/danielmesquitta/tasks-api/internal/app/restapi/dto"
	"github.com/danielmesquitta/tasks-api/internal/domain/entity"
	"github.com/labstack/echo/v4"
)

var mapErrTypeToStatusCode = map[entity.ErrType]int{
	entity.ErrTypeForbidden:    http.StatusForbidden,
	entity.ErrTypeUnauthorized: http.StatusUnauthorized,
	entity.ErrTypeValidation:   http.StatusBadRequest,
	entity.ErrTypeUnknown:      http.StatusInternalServerError,
	entity.ErrTypeNotFound:     http.StatusNotFound,
}

func (m *Middleware) ErrorHandler(
	defaultErrorHandler echo.HTTPErrorHandler,
) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		if appErr, ok := err.(*entity.Err); ok {
			statusCode := mapErrTypeToStatusCode[appErr.Type]
			if internalServerError := statusCode >= 500 || statusCode == 0; internalServerError {
				req := c.Request()

				requestData := map[string]any{}
				_ = c.Bind(&requestData)

				slog.Error(
					appErr.Error(),
					"url",
					req.URL.Path,
					"body",
					requestData,
					"query",
					c.QueryParams(),
					"params",
					c.ParamValues(),
					"stacktrace",
					appErr.StackTrace,
				)

				err = c.JSON(
					statusCode,
					dto.ErrorResponseDTO{Message: "internal server error"},
				)
				if err != nil {
					log.Println(err)
				}
				return
			}

			err = c.JSON(
				statusCode,
				dto.ErrorResponseDTO{Message: appErr.Error()},
			)
			if err != nil {
				log.Println(err)
			}
			return
		}

		defaultErrorHandler(err, c)
	}
}
