package api

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/i18n"
	"allaboutapps.dev/aw/go-starter/internal/mailer"
	"allaboutapps.dev/aw/go-starter/internal/mailer/transport"
	"allaboutapps.dev/aw/go-starter/internal/push"
	"allaboutapps.dev/aw/go-starter/internal/push/provider"
	"github.com/rs/zerolog/log"
)

///////////////////////////////////////////////
// PROVIDERS
// https://github.com/google/wire/blob/main/docs/guide.md#defining-providers
///////////////////////////////////////////////

func NewDB(cfg config.Server) (*sql.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := newDBConnection(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize database: %w", err)
	}

	return db, nil
}

func newDBConnection(ctx context.Context, cfg config.Server) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.Database.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("Failed to open DB connection: %w", err)
	}

	if cfg.Database.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	}
	if cfg.Database.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	}
	if cfg.Database.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("Failed to ping DB: %w", err)
	}

	return db, nil
}

func NewMailer(cfg config.Server) (m *mailer.Mailer, err error) {

	switch config.MailerTransporter(cfg.Mailer.Transporter) {
	case config.MailerTransporterMock:
		log.Warn().Msg("Initializing mock mailer")
		m = mailer.New(cfg.Mailer, transport.NewMock())
	case config.MailerTransporterSMTP:
		m = mailer.New(cfg.Mailer, transport.NewSMTP(cfg.SMTP))
	default:
		return nil, fmt.Errorf("Unsupported mail transporter: %s", cfg.Mailer.Transporter)
	}

	if err := m.ParseTemplates(); err != nil {
		return nil, fmt.Errorf("Failed to parse templates: %w", err)
	}

	return m, nil
}

func NewPush(cfg config.Server, db *sql.DB) (*push.Service, error) {
	pusher := push.New(db)

	if cfg.Push.UseFCMProvider {
		fcmProvider, err := provider.NewFCM(cfg.FCMConfig)
		if err != nil {
			return nil, err
		}
		pusher.RegisterProvider(fcmProvider)
	}

	if cfg.Push.UseMockProvider {
		log.Warn().Msg("Initializing mock push provider")
		mockProvider := provider.NewMock(push.ProviderTypeFCM)
		pusher.RegisterProvider(mockProvider)
	}

	if pusher.GetProviderCount() < 1 {
		log.Warn().Msg("No providers registered for push service")
	}

	return pusher, nil
}

func NewI18n(cfg config.Server) (*i18n.Service, error) {
	return i18n.New(cfg.I18n)
}
