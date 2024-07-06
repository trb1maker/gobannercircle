package main

import (
	"embed"

	"github.com/trb1maker/gobannercircle/cmd"
)

//go:embed migrations/postgres/*
var migrations embed.FS

func main() {
	cmd.MigrationsPipe(migrations)
	cmd.Execute()
}
