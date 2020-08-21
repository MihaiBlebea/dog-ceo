package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MihaiBlebea/dog-ceo/api"
	"github.com/MihaiBlebea/dog-ceo/cache"
	"github.com/MihaiBlebea/dog-ceo/dog"
	"github.com/MihaiBlebea/dog-ceo/template"
	"github.com/sirupsen/logrus"
)

var c *cache.Service

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	dogService := dog.New(logger)

	c := &cache.Service{
		DogService: dogService,
		Logger:     logger,
	}

	templateService := template.New(c)

	server := api.New(
		templateService,
		logger,
	)

	httpPort := fmt.Sprintf(":%s", os.Getenv("HTTP_PORT"))
	logger.Info(fmt.Sprintf("Application started HTTP server on port %s", httpPort))

	err := http.ListenAndServe(httpPort, *server.Handler())
	if err != nil {
		logger.Fatal(err)
	}
}
