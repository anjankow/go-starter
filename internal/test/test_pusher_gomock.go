package test

import (
	"database/sql"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/push"
	"allaboutapps.dev/aw/go-starter/internal/push/provider"
	"github.com/golang/mock/gomock"
)

func WithTestPusherGoMock(t *testing.T, providerType push.ProviderType, closure func(p *push.Service, db *sql.DB, mockProvider *provider.GomockProvider)) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	gomockProvider := provider.NewGomockProvider(mockCtrl)

	gomockProvider.EXPECT().GetProviderType().AnyTimes().Return(providerType)

	WithTestDatabase(t, func(db *sql.DB) {
		t.Helper()
		closure(NewTestPusherGomock(t, gomockProvider, db), db, gomockProvider)
	})

	mockCtrl.Finish()
}

func NewTestPusherGomock(t *testing.T, gomockProvider *provider.GomockProvider, db *sql.DB) *push.Service {
	t.Helper()

	pushService := push.New(db)
	pushService.RegisterProvider(gomockProvider)

	return pushService
}
