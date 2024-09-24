package main

import (
	"context"
	"go-microservices/product-api/handler"
	"go-microservices/product-api/migrate"
	"go-microservices/product-api/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/nicholasjackson/env"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":3000", "Bind address for the server")

func main() {

	env.Parse()
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	// helloHandler := handler.NewHello(l)
	// goodByeHandler := handler.NewGoodBye(l)

	/*Load the config from .env file*/
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	/*Load the database*/
	database, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Could not load database", err.Error())
	}

	/*Load the tables*/
	err = migrate.MigrateAllTables(database)
	if err != nil {
		log.Fatal("Could not migrate books database")
	}

	productHandler := handler.NewProducts(l, database)

	// sm := http.NewServeMux()
	// sm.Handle("/", productHandler)

	/*Using mux of gureilla framework which
	removes the need of ServeHTTP*/
	sm := mux.NewRouter()

	//router for getting all products
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", productHandler.GetProducts)

	//router for adding a product
	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/products/{id:[0-9]+}", productHandler.UpdateProducts)
	putRouter.Use(productHandler.MiddlewareValidateProduct)

	//router for updating a product
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/products", productHandler.AddProducts)
	postRouter.Use(productHandler.MiddlewareValidateProduct)

	//router for adding an image to a product
	imagePutRouter := sm.Methods(http.MethodPut).Subrouter()
	imagePutRouter.HandleFunc("/products/{id:[0-9]+}/image", productHandler.UploadFile)

	//router for adding an image to a product
	imageGetRouter := sm.Methods(http.MethodGet).Subrouter()
	imageGetRouter.HandleFunc("/products/{id:[0-9]+}/image", productHandler.DownloadFile)

	//Set up swagger
	ops := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(ops, nil)
	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	//Set up CORS issue
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	//create a new server

	s := http.Server{
		Addr:         *bindAddress,
		Handler:      ch(sm),
		ErrorLog:     l,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	//start the server
	go func() {
		l.Println("starting the server")

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	//// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	s.Shutdown(ctx)

}
