package main

import (
	server "github.com/bear1278/MusicWave"
	"github.com/bear1278/MusicWave/configs"
	"github.com/bear1278/MusicWave/pkg/handlers"
	"github.com/bear1278/MusicWave/pkg/repository"
	"github.com/bear1278/MusicWave/pkg/service"
	"log"
)

func main() {
	cfg, err := configs.Init()
	if err != nil {
		log.Fatalf("error occured while read config: %s", err.Error())
	}

	db, err := repository.MySqlDB(*cfg)
	if err != nil {
		log.Fatalf("failed to initialize database: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handler := handlers.NewHandler(services)

	srv := new(server.Server)
	err = srv.Run(cfg.Port, handler.InitRoutes())
	if err != nil {
		log.Fatalf("error occured while run http server: %s", err.Error())
	}
}
