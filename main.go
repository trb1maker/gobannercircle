package main

import (
	"embed"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/trb1maker/gobannercircle/cmd"
)

//go:embed migrations/postgres
var fs embed.FS

func main() {
	cmd.SendMigrations(fs)
	cmd.Execute()
}
