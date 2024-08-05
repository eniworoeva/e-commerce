package server

import (
	"context"
	"e-commerce/internal/api"
	"e-commerce/internal/repository"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// Run injects all dependencies needed to run the app
func Run(db *gorm.DB, port string) {
	//Create a new instance of our repository
	newRepo := repository.NewDB(db)

	//Create a new instance of our handler
	Handler := api.NewHTTPHandler(newRepo)
	//Create a new router
	router := SetupRouter(Handler, newRepo)

	//Create a new server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	fmt.Printf("Listening and serving HTTP on : %v\n", port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	log.Println("Receive terminate and shutdown gracefully", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

// Params is a data model of the data in our environment variable
type Params struct {
	Port  string
	DbUrl string
}

// InitDBParams gets environment variables needed to run the app
func InitDBParams() Params {

	//for debugging
	/* pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	errEnv := godotenv.Load(filepath.Join(pwd, "../.env")) */

	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Println("Error loading .env file")
	}

	port := os.Getenv("PORT")
	dbURL := os.Getenv("DATABASE_URL")
	if port == "" {
		port = "8080"
	}

	return Params{
		Port:  port,
		DbUrl: dbURL,
	}
}
