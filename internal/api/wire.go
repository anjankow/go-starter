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

// NewServer returns NewServer instance with all the components initialized.
// Exceptions are Echo and Router, which are not initialized.
// After this call make sure that router.Init(s) is invoked.
func NewServer(
	cfg config.Server,
	db *sql.DB,
	mailer *mailer.Mailer,
	push *push.Service,
	i18n *i18n.Service,
) *Server {
	return &Server{
		Config: cfg,
		DB:     db,
		Mailer: mailer,
		Push:   push,
		I18n:   i18n,
	}
}

func InitNewServer(
	cfg config.Server,
) (*Server, error) {
	wire.Build(NewServer, InitDB, InitMailer, InitPush, InitI18n)
	return new(Server), nil
}

func InitNewServerWithDB(
	cfg config.Server,
	db *sql.DB,
) (*Server, error) {
	wire.Build(NewServer, InitMailer, InitPush, InitI18n)
	return new(Server), nil
}
