package push

import (
	"context"
	"database/sql"
	"errors"

	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/util"
)

type ProviderType string

const (
	ProviderTypeFCM ProviderType = "fcm"
	ProviderTypeAPN ProviderType = "apn"
)

type Service struct {
	DB       *sql.DB
	provider map[ProviderType]Provider
}

type ProviderSendResponse struct {
	// token the message was sent to
	Token string

	// flag to indicate if the token is still valid
	// not every error means that the token is invalid
	Valid bool

	// ogiginal error
	Err error
}

type Provider interface {
	Send(token string, title string, message string, data map[string]string, silent bool, collapseKey ...string) ProviderSendResponse
	SendWithContext(ctx context.Context, token string, title string, message string, data map[string]string, silent bool, collapseKey ...string) ProviderSendResponse
	GetProviderType() ProviderType

	// DEPRECATED: SendMulticast
	// Allows to send same notification to multiple receivers.
	//
	// This interface function is deprecated and might be removed with future releases.
	// Please use sendMulticastWithProvider instead defined in push package.
	SendMulticast(tokens []string, title string, message string, data map[string]string, silent bool, collapseKey ...string) []ProviderSendResponse
}

func New(db *sql.DB) *Service {
	return &Service{
		DB:       db,
		provider: make(map[ProviderType]Provider),
	}
}

func (s *Service) RegisterProvider(p Provider) {
	s.provider[p.GetProviderType()] = p
}

func (s *Service) ResetProviders() {
	s.provider = make(map[ProviderType]Provider)
}

func (s *Service) GetProviderCount() int {
	return len(s.provider)
}

func (s *Service) SendToUser(ctx context.Context, user *models.User, title string, message string, data map[string]string, silent bool, collapseKey ...string) error {
	if s.GetProviderCount() < 1 {
		return errors.New("No provider found")
	}
	log := util.LogFromContext(ctx)

	for k, p := range s.provider {
		// get all registered tokens for provider
		pushTokens, err := user.PushTokens(models.PushTokenWhere.Provider.EQ(string(k))).All(ctx, s.DB)
		if err != nil {
			return err
		}

		var tokens []string
		for _, token := range pushTokens {
			tokens = append(tokens, token.Token)
		}

		responseSlice := s.sendMulticastWithProvider(ctx, p, tokens, title, message, data, silent, collapseKey...)
		tokenToDelete := make([]string, 0)
		for _, res := range responseSlice {
			if res.Err != nil && res.Valid {
				log.Debug().Err(res.Err).Str("token", res.Token).Str("provider", string(p.GetProviderType())).Msgf("Error while sending push message to provider with valid token.")
			}

			if !res.Valid {
				tokenToDelete = append(tokenToDelete, res.Token)
			}
		}
		// delete invalid tokens
		_, err = user.PushTokens(models.PushTokenWhere.Token.IN(tokenToDelete)).DeleteAll(ctx, s.DB)
		if err != nil {
			log.Debug().Err(err).Str("provider", string(p.GetProviderType())).Msg("Could not delete invalid tokens for provider")
			return err
		}
	}

	return nil
}

// sendMulticastWithProvider allows to send same notification to multiple receivers via one provider.
func (s *Service) sendMulticastWithProvider(ctx context.Context, p Provider, tokens []string, title string, message string, data map[string]string, silent bool, collapseKey ...string) []ProviderSendResponse {
	responseSlice := make([]ProviderSendResponse, 0)

	for _, token := range tokens {
		responseSlice = append(responseSlice, p.SendWithContext(ctx, token, title, message, data, silent, collapseKey...))
	}

	return responseSlice
}
