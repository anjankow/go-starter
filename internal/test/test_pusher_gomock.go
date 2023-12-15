package test

import (
	"database/sql"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/push"
	"allaboutapps.dev/aw/go-starter/internal/push/provider"
	"github.com/golang/mock/gomock"
)

func WithTestPusherGoMock(t *testing.T, providerTypes []push.ProviderType, closure func(p *push.Service, db *sql.DB, mockProvider map[push.ProviderType]*provider.GomockProvider)) {
	t.Helper()

	// create mock controller - to verify the expectations after the test run
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// create requested mock providers
	gomockProviders := make(map[push.ProviderType]*provider.GomockProvider, len(providerTypes))
	for _, providerType := range providerTypes {
		provider := provider.NewGomockProvider(mockCtrl)
		// make this mock always identify itself as a certain provider type
		provider.EXPECT().GetProviderType().AnyTimes().Return(providerType)

		gomockProviders[providerType] = provider
	}

	WithTestDatabase(t, func(db *sql.DB) {
		t.Helper()
		closure(NewTestPusherGomock(t, gomockProviders, db), db, gomockProviders)
	})
}

func NewTestPusherGomock(t *testing.T, gomockProviders map[push.ProviderType]*provider.GomockProvider, db *sql.DB) *push.Service {
	t.Helper()

	pushService := push.New(db)
	for _, provider := range gomockProviders {
		pushService.RegisterProvider(provider)
	}

	return pushService
}
