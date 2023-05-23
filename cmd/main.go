package main

import (
	"WB_Intern/internal/broker"
	"WB_Intern/internal/repository"
	"WB_Intern/internal/repository/cache"
	"WB_Intern/internal/service"
	"github.com/sirupsen/logrus"
)


func main() {
	logger := logrus.New()

	// db lvl
	repo := repository.NewRepo(logger)
	// cache lvl
	newCache := cache.NewCache(logger, repo)
	// broker lvl
	newBroker := broker.NewBroker(repo, newCache, logger)
	// service lvl
	service := services.NewService(logger, repo, newCache, newBroker)

	// execute all
	service.Run()
}