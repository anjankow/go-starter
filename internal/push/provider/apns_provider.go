package provider

import (
	"context"
	"fmt"

	"allaboutapps.dev/aw/go-starter/internal/push"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/payload"
	"github.com/sideshow/apns2/token"
)

type APNS struct {
	config APNSConfiguration
	client *apns2.Client
}

type APNSConfiguration struct {
	AuthKeyPath  string
	Topic        string
	DebugPayload bool
}

func NewAPNS(cfg APNSConfiguration) (*APNS, error) {

	authKey, err := token.AuthKeyFromFile(cfg.AuthKeyPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to get client auth key from file: %w", err)
	}

	token := &token.Token{
		AuthKey: authKey,
	}

	client := apns2.NewTokenClient(token).Production()

	return &APNS{
		config: cfg,
		client: client,
	}, nil
}

func (a *APNS) Send(token string, title string, message string, data map[string]string, silent bool, collapseKey ...string) push.ProviderSendResponse {
	ctx := context.Background()
	return a.SendWithContext(ctx, token, title, message, data, silent, collapseKey...)
}

func (a *APNS) SendWithContext(ctx context.Context, token string, title string, message string, data map[string]string, silent bool, collapseKey ...string) push.ProviderSendResponse {
	log := util.LogFromContext(ctx)

	notification := &apns2.Notification{
		DeviceToken: token,
	}

	var p *payload.Payload

	if silent {
		notification.PushType = apns2.PushTypeBackground
		notification.Priority = apns2.PriorityLow
		if len(collapseKey) == 1 {
			notification.CollapseID = collapseKey[0]
		}
		p = payload.NewPayload().ContentAvailable()

	} else {
		p = payload.NewPayload().AlertTitle(title).Alert(message)
	}

	for key, value := range data {
		p = p.Custom(key, value)
	}

	notification.Payload = p

	var res *apns2.Response
	var err error

	notification.Topic = a.config.Topic

	res, err = a.client.PushWithContext(ctx, notification)

	if a.config.DebugPayload {
		log.Debug().Str("token", token).Interface("notification", notification).Interface("payload", notification.Payload).Msg("APNS notification")
		log.Debug().Str("token", token).Interface("response", res).Msg("APNS response")
	}

	if err != nil {
		err = fmt.Errorf(`Push via APNS failed with error: "%w" and reason: "%s"`, err, res.Reason)
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		err = fmt.Errorf(`Push via APNS failed with status code: "%d" and reason: "%s"`, res.StatusCode, res.Reason)
	}

	return push.ProviderSendResponse{
		Token: token,
		Err:   err,
		Valid: res.Reason != apns2.ReasonBadDeviceToken,
	}
}

func (a *APNS) GetProviderType() push.ProviderType {
	return push.ProviderTypeAPN
}
