package main

import (
	"eWallet/config"
	"eWallet/internal/handler"
	"eWallet/internal/service"
	"eWallet/internal/storage"
	"fmt"
	"log"
)

func main() {
	cnfg := config.New()
	stM, err := storage.NewDBStorage(cnfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	sr := service.NewService(stM)
	if err != nil {
		log.Fatalf("Error creating service: %s", err)
		return
	}
	h := handler.NewHandler(sr, cnfg)
	h.Start()
	stM.Close()

}
