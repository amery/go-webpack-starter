package main

import (
	"log"

	"github.com/creasty/defaults"

	"github.com/amery/go-webpack-starter/web/server"
)

var cfg = NewConfig()

type ServerConfig struct {
	Server      server.ServerConfig
	Development bool
}

func NewConfig() *ServerConfig {
	c := &ServerConfig{}

	if err := defaults.Set(c); err != nil {
		log.Fatal(err)
	}

	return c
}
