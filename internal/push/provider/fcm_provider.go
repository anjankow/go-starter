package provider

import (
	"context"
	"errors"
	"net/http"

	"allaboutapps.dev/aw/go-starter/internal/push"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"google.golang.org/api/fcm/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

type FCM struct {
	Config  FCMConfig
	service *fcm.Service
}

type FCMConfig struct {
	GoogleApplicationCredentials string `json:"-"` // sensitive
	ProjectID                    string
	ValidateOnly                 bool
	DebugPayload                 bool
}

func NewFCM(config FCMConfig, opts ...option.ClientOption) (*FCM, error) {
	ctx := context.Background()
	fcmService, err := fcm.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &FCM{
		Config:  config,
		service: fcmService,
	}, nil
}

func (p *FCM) GetProviderType() push.ProviderType {
	return push.ProviderTypeFCM
}

func (p *FCM) Send(token string, title string, message string, data map[string]string, silent bool, collapseKey ...string) push.ProviderSendResponse {
	ctx := context.Background()
	return p.SendWithContext(ctx, token, title, message, data, silent, collapseKey...)
}

func (p *FCM) SendWithContext(ctx context.Context, token string, title string, message string, data map[string]string, silent bool, collapseKey ...string) push.ProviderSendResponse {
	log := util.LogFromContext(ctx)

	// https: //godoc.org/google.golang.org/api/fcm/v1#SendMessageRequest
	// https://firebase.google.com/docs/cloud-messaging/send-message#rest
	data["title"] = title
	data["message"] = message

	messageRequest := &fcm.SendMessageRequest{
		ValidateOnly: p.Config.ValidateOnly,
		Message: &fcm.Message{
			Token: token,
			Data:  data,
		},
	}

	if silent {
		messageRequest.Message.Android = &fcm.AndroidConfig{
			Priority: "NORMAL",
		}
		if len(collapseKey) == 1 {
			messageRequest.Message.Android.CollapseKey = collapseKey[0]
		}
	}

	res, err := p.service.Projects.Messages.Send("projects/"+p.Config.ProjectID, messageRequest).Do()
	if p.Config.DebugPayload {
		log.Debug().Str("token", token).Interface("message", messageRequest.Message).Msg("FCM notification")
		log.Debug().Str("token", token).Interface("response", res).Msg("FCM response")
	}

	valid := true
	if err != nil {

		// convert to original error and determine if the token was at fault
		var gErr *googleapi.Error
		if errors.As(err, &gErr) {
			valid = !(gErr.Code == http.StatusNotFound || gErr.Code == http.StatusBadRequest)
		}
	}

	return push.ProviderSendResponse{
		Token: token,
		Valid: valid,
		Err:   err,
	}
}

func (p *FCM) SendMulticast(tokens []string, title string, message string, data map[string]string, silent bool, collapseKey ...string) []push.ProviderSendResponse {
	return sendMulticastWithProvider(p, tokens, title, message, data, silent, collapseKey...)
}
