//go:build wireinject
// +build wireinject

package api

import (
	"database/sql"

	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/i18n"
	"allaboutapps.dev/aw/go-starter/internal/mailer"
	"allaboutapps.dev/aw/go-starter/internal/push"
	"github.com/google/wire"
)

// newServerWithComponents is used by wire to initialize the server components.
// Components not listed here won't be handled by wire and should be initialized separately.
// Components which shouldn't be handled must be labeled `wire:"-"` in Server struct.
func newServerWithComponents(
	cfg config.Server,
	db *sql.DB,
	mail *mailer.Mailer,
	pusher *push.Service,
	i18n *i18n.Service,
) *Server {
	return &Server{
		Config: cfg,
		DB:     db,
		Mailer: mail,
		Push:   pusher,
		I18n:   i18n,
	}
}

// InitNewServer returns a new Server instance.
// All the components are initialized via go wire according to the configuration.
// WARNING! Exceptions are Echo and Router, which are not initialized.
// After this call make sure that router.Init(s) is invoked.
func InitNewServer(
	cfg config.Server,
) (*Server, error) {
	wire.Build(newServerWithComponents, InitDB, InitMailer, InitPush, InitI18n)
	return new(Server), nil
}

// InitNewServerWithDB returns a new Server instance with the given DB instance
// All the other components are initialized via go wire according to the configuration.
// WARNING! Exceptions are Echo and Router, which are not initialized.
// After this call make sure that router.Init(s) is invoked.
func InitNewServerWithDB(
	cfg config.Server,
	db *sql.DB,
) (*Server, error) {
	wire.Build(newServerWithComponents, InitMailer, InitPush, InitI18n)
	return new(Server), nil
}
