//go:build wireinject
// +build wireinject

package api

import (
	"database/sql"

	"allaboutapps.dev/aw/go-starter/internal/config"
	"github.com/google/wire"
)

///////////////////////////////////////////////
// INJECTORS
// https://github.com/google/wire/blob/main/docs/guide.md#injectors
///////////////////////////////////////////////

// InitNewServer returns a new Server instance.
// All the components are initialized via go wire according to the configuration.
// WARNING! Exceptions are Echo and Router, which are not initialized.
// After this call make sure that router.Init(s) is invoked.
func InitNewServer(
	cfg config.Server,
) (*Server, error) {
	wire.Build(newServerWithComponents, NewDB, NewMailer, NewPush, NewI18n)
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
	wire.Build(newServerWithComponents, NewMailer, NewPush, NewI18n)
	return new(Server), nil
}
