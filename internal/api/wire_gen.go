// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package api

import (
	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/i18n"
	"allaboutapps.dev/aw/go-starter/internal/mailer"
	"allaboutapps.dev/aw/go-starter/internal/persistence"
	"database/sql"
)

import (
	_ "github.com/lib/pq"
)

// Injectors from wire.go:

// InitNewServer returns a new Server instance.
// All the components are initialized via go wire according to the configuration.
// WARNING! Exceptions are Echo and Router, which are not initialized.
// After this call make sure that router.Init(s) is invoked.
func InitNewServer(cfg config.Server) (*Server, error) {
	db, err := persistence.NewDB(cfg)
	if err != nil {
		return nil, err
	}
	configMailer := config.GetMailerConfig(cfg)
	mailerMailer, err := mailer.NewWithConfig(configMailer)
	if err != nil {
		return nil, err
	}
	service, err := NewPush(cfg, db)
	if err != nil {
		return nil, err
	}
	configI18n := config.GetI18nConfig(cfg)
	i18nService, err := i18n.New(configI18n)
	if err != nil {
		return nil, err
	}
	server := newServerWithComponents(cfg, db, mailerMailer, service, i18nService)
	return server, nil
}

// InitNewServerWithDB returns a new Server instance with the given DB instance
// All the other components are initialized via go wire according to the configuration.
// WARNING! Exceptions are Echo and Router, which are not initialized.
// After this call make sure that router.Init(s) is invoked.
func InitNewServerWithDB(cfg config.Server, db *sql.DB) (*Server, error) {
	configMailer := config.GetMailerConfig(cfg)
	mailerMailer, err := mailer.NewWithConfig(configMailer)
	if err != nil {
		return nil, err
	}
	service, err := NewPush(cfg, db)
	if err != nil {
		return nil, err
	}
	configI18n := config.GetI18nConfig(cfg)
	i18nService, err := i18n.New(configI18n)
	if err != nil {
		return nil, err
	}
	server := newServerWithComponents(cfg, db, mailerMailer, service, i18nService)
	return server, nil
}
