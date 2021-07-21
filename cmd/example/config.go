package main

import (
	"log"

	"go.sancus.dev/config/flags"

	"github.com/amery/go-webpack-starter/web/server"
)

var cfg = NewConfig()

type ServerConfig struct {
	Server      server.ServerConfig
	Development bool
}

func NewConfig() *ServerConfig {
	c := &ServerConfig{}

	if err := flags.SetDefaults(c); err != nil {
		log.Fatal(err)
	}

	return c
}
