package config

type PushService struct {
	UseFCMProvider     bool
	UseMockProvider    bool
	UseAPNSProvider    bool
	PushPayloadDebug   bool
	EnableTestEndpoint bool
}
