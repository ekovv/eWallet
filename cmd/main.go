package main

import (
	"eWallet/config"
	"eWallet/internal/handler"
	"eWallet/internal/service"
	"eWallet/internal/storage"
	"log"
)

func main() {
	cnfg := config.New()
	stM, err := storage.NewDBStorage(cnfg)
	if err != nil {
		log.Fatalf("Error creating storage: %s", err)
		return
	}
	sr := service.NewService(stM, cnfg)
	if err != nil {
		log.Fatalf("Error creating service: %s", err)
		return
	}
	h := handler.NewHandler(sr, cnfg)
	h.Start()
	stM.Close()

}
