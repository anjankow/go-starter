package qa

import (
	"net/http"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/auth"
	"allaboutapps.dev/aw/go-starter/internal/api/middleware"
	"allaboutapps.dev/aw/go-starter/internal/types"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/go-openapi/swag"
	"github.com/labstack/echo/v4"
)

func PostQAPushRoute(s *api.Server) *echo.Route {
	return s.Router.APIV1QA.POST("/push", postQAPushHandler(s), middleware.EnableRoute(s.Config.Push.EnableTestEndpoint))
}

func postQAPushHandler(s *api.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		log := util.LogFromContext(ctx)

		user := auth.UserFromEchoContext(c)

		var body types.PushTestPayload
		if err := util.BindAndValidateBody(c, &body); err != nil {
			return err
		}

		if body.Title == nil {
			body.Title = swag.String("Hello")
		}

		if body.Message == nil {
			body.Message = swag.String("World")
		}

		if body.Data == nil {
			body.Data = make(map[string]string)
		}

		err := s.Push.SendToUser(ctx, user, swag.StringValue(body.Title), swag.StringValue(body.Message), body.Data, false)
		if err != nil {
			log.Debug().Err(err).Msg("Error while sending push to user.")
			return err
		}

		log.Debug().Msg("Successfully sent push message.")

		return c.String(http.StatusOK, "Success")
	}
}
