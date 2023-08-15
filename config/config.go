package config

import (
	"fmt"
	"github.com/gofor-little/env"
	"time"
)

type FsConfig struct {
	DumpDir   string
	LocalDir  string
	PublicKey string
}

type SqlConfig struct {
	Username       string
	Password       string
	DBName         string
	Address        string
	Port           string
	DBDumpFileName string
}

type BastionConfig struct {
	Address string
	Port    string
	User    string
}

func GetSqlConfig() *SqlConfig {
	return &SqlConfig{
		Username: env.Get("DB_USERNAME", ""),
		Password: env.Get("DB_PASSWORD", ""),
		DBName:   env.Get("DB_NAME", ""),
		Address:  env.Get("DB_ADDRESS", ""),
		Port:     env.Get("DB_PORT", ""),
		DBDumpFileName: fmt.Sprintf(
			"%s-%s",
			env.Get("DB_DUMP_FILE_NAME",
				"",
			), time.Now().Format("20060102T150405")),
	}
}

func GetFsConfig() *FsConfig {
	return &FsConfig{
		DumpDir:   env.Get("DUMP_DIR", ""),
		PublicKey: env.Get("PUBLIC_KEY", ""),
		LocalDir:  env.Get("LOCAL_DIR", ""),
	}
}

func GetBastionConfig() *BastionConfig {
	return &BastionConfig{
		Address: env.Get("BASTION_ADDRESS", ""),
		Port:    env.Get("BASTION_PORT", ""),
		User:    env.Get("BASTION_USER", ""),
	}
}

func init() {
	if err := env.Load(".env"); err != nil {
		panic(err)
	}
}
