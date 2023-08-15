package main

import (
	"log"
	"scp-copy/config"
	"scp-copy/scp"
)

func main() {
	log.Println("Begin backup script")

	// Not sure if this is the "go" way of doing things
	sqlConfig := config.GetSqlConfig()
	fsConfig := config.GetFsConfig()
	bastionConfig := config.GetBastionConfig()

	scp.BackupDB(sqlConfig, fsConfig, bastionConfig)
}
