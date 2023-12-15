package push_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/push"
	"allaboutapps.dev/aw/go-starter/internal/push/provider"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func TestSendMessageSuccess(t *testing.T) {
	test.WithTestPusher(t, func(p *push.Service, db *sql.DB) {
		ctx := context.Background()
		fixtures := test.Fixtures()

		user1 := fixtures.User1

		err := p.SendToUser(ctx, user1, "Hello", "World")
		assert.NoError(t, err)

		tokenCount, err2 := user1.PushTokens().Count(ctx, db)
		require.NoError(t, err2)
		assert.Equal(t, int64(2), tokenCount)
	})
}

func TestSendMessageSuccessWithGenericError(t *testing.T) {
	test.WithTestPusher(t, func(p *push.Service, db *sql.DB) {
		ctx := context.Background()
		fixtures := test.Fixtures()

		user1 := fixtures.User1

		// provoke error from mock provider
		err := p.SendToUser(ctx, user1, "other error", "World")
		assert.NoError(t, err)

		tokenCount, err2 := user1.PushTokens().Count(ctx, db)
		require.NoError(t, err2)
		assert.Equal(t, int64(2), tokenCount)
	})
}

func TestSendMessageWithInvalidToken(t *testing.T) {
	test.WithTestPusher(t, func(p *push.Service, db *sql.DB) {
		ctx := context.Background()
		fixtures := test.Fixtures()
		user1 := fixtures.User1

		user1InvalidPushToken := models.PushToken{
			ID:       "55c37bc8-f245-40b3-bdef-14dee35b10bd",
			Token:    "d5ded380-3285-4243-8a9c-72cc3f063fee",
			UserID:   user1.ID,
			Provider: models.ProviderTypeFCM,
		}
		err := user1InvalidPushToken.Insert(ctx, db, boil.Infer())
		require.NoError(t, err)

		tokenCount, err2 := user1.PushTokens().Count(ctx, db)
		require.NoError(t, err2)
		require.Equal(t, int64(3), tokenCount)

		err = p.SendToUser(ctx, user1, "Hello", "World")
		assert.NoError(t, err)

		tokenCount, err2 = user1.PushTokens().Count(ctx, db)
		require.NoError(t, err2)
		assert.Equal(t, int64(2), tokenCount)
	})
}

func TestSendMessageWithNoProvider(t *testing.T) {
	test.WithTestPusher(t, func(p *push.Service, db *sql.DB) {

		ctx := context.Background()
		fixtures := test.Fixtures()

		p.ResetProviders()
		require.Equal(t, 0, p.GetProviderCount())

		user1 := fixtures.User1

		err := p.SendToUser(ctx, user1, "Hello", "World")
		assert.Error(t, err)

		tokenCount, err2 := user1.PushTokens().Count(ctx, db)
		require.NoError(t, err2)
		assert.Equal(t, int64(2), tokenCount)
	})
}

func TestSendMessageWithMultipleProvider(t *testing.T) {
	test.WithTestPusherGoMock(t, []push.ProviderType{push.ProviderTypeFCM, push.ProviderTypeAPN}, func(p *push.Service, db *sql.DB, mockProvider map[push.ProviderType]*provider.GomockProvider) {
		ctx := context.Background()
		fixtures := test.Fixtures()

		user1 := fixtures.User1

		// make FCM provider fail but with a valid token
		mockProvider[push.ProviderTypeFCM].EXPECT().SendMulticast(gomock.Any(), "Hello", "World").
			Return([]push.ProviderSendResponse{
				{
					Err:   errors.New("Failed to send"),
					Valid: true,
				},
			})

		// make APN provider return no error but invalidate the token
		mockProvider[push.ProviderTypeAPN].EXPECT().SendMulticast(gomock.Any(), "Hello", "World").
			Return([]push.ProviderSendResponse{
				{
					Token: fixtures.User1PushTokenAPN.Token,
					Valid: false,
				},
			})

		err := p.SendToUser(ctx, user1, "Hello", "World")
		assert.NoError(t, err)

		// FCM token should still exist
		require.NoError(t, fixtures.User1PushToken.Reload(ctx, db))
		// APN token should be gone
		require.ErrorIs(t, fixtures.User1PushTokenAPN.Reload(ctx, db), sql.ErrNoRows)
	})
}

func TestSendMessageSuccessWithGomock(t *testing.T) {
	test.WithTestPusherGoMock(t, []push.ProviderType{push.ProviderTypeFCM}, func(p *push.Service, db *sql.DB, mockProvider map[push.ProviderType]*provider.GomockProvider) {
		ctx := context.Background()
		fixtures := test.Fixtures()

		user1 := fixtures.User1
		existingToken := fixtures.User1PushToken

		mockProvider[push.ProviderTypeFCM].EXPECT().SendMulticast(gomock.Any(), "Hello", "World").
			Return([]push.ProviderSendResponse{
				{
					Token: existingToken.Token,
					Valid: true,
				},
			})

		err := p.SendToUser(ctx, user1, "Hello", "World")
		assert.NoError(t, err)

		tokenCount, err2 := user1.PushTokens().Count(ctx, db)
		require.NoError(t, err2)
		assert.Equal(t, int64(2), tokenCount)
	})
}
