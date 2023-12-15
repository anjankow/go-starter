package qa_test

import (
	"net/http"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestPostQAPushDefault(t *testing.T) {
	cfg := config.DefaultServiceConfigFromEnv()
	cfg.Push.EnableTestEndpoint = true
	test.WithTestServerConfigurable(t, cfg, func(s *api.Server) {
		fixtures := test.Fixtures()

		res := test.PerformRequest(t, s, "POST", "/api/v1/qa/push", nil, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))

		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
	})
}

func TestPostQAPush(t *testing.T) {
	cfg := config.DefaultServiceConfigFromEnv()
	cfg.Push.EnableTestEndpoint = true
	cfg.Push.PushPayloadDebug = true
	cfg.Push.UseMockProvider = false
	cfg.Push.UseAPNSProvider = true
	cfg.Push.UseFCMProvider = true
	test.WithTestServerConfigurable(t, cfg, func(s *api.Server) {
		fixtures := test.Fixtures()

		payload := test.GenericPayload{
			"title":   "Baking tips",
			"message": "How to avoid setting up your house in fire",
			"data": test.GenericPayload{
				"author": "me",
			},
		}

		res := test.PerformRequest(t, s, "POST", "/api/v1/qa/push", payload, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))

		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
	})
}

func TestPostQAPushUnauthorized(t *testing.T) {
	cfg := config.DefaultServiceConfigFromEnv()
	cfg.Push.EnableTestEndpoint = true
	test.WithTestServerConfigurable(t, cfg, func(s *api.Server) {
		res := test.PerformRequest(t, s, "POST", "/api/v1/qa/push", nil, nil)

		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode)
	})
}

func TestPostQAPushDisabled(t *testing.T) {
	cfg := config.DefaultServiceConfigFromEnv()
	cfg.Push.EnableTestEndpoint = false
	test.WithTestServerConfigurable(t, cfg, func(s *api.Server) {
		fixtures := test.Fixtures()

		res := test.PerformRequest(t, s, "POST", "/api/v1/qa/push", nil, test.HeadersWithAuth(t, fixtures.User1AccessToken1.Token))

		assert.Equal(t, http.StatusNotFound, res.Result().StatusCode)
	})
}
