package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xorcare/blockchain"

	"Cryptogo/config"
	"Cryptogo/internal/db/postgres"
	"Cryptogo/internal/server"
	"Cryptogo/pkg/logger"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	log := logger.GetLogger(cfg)

	db, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	bcClient := blockchain.New()

	g := gin.New()
	s := server.NewServer(g, cfg, db, log, bcClient)

	s.Run()
}
