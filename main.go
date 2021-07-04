package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sing3demons/api/database"
	"github.com/sing3demons/api/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var dir string
	flag.StringVar(&dir, "dir", ".", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	database.InitDatabase()

	r := mux.NewRouter().StrictSlash(true)
	r.Use(handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Hello, World"})
	})

	// Choose the folder to serve
	staticDir := "uploads"

	uploadDir := [...]string{"products", "users"}
	for _, path := range uploadDir {
		os.MkdirAll(staticDir+"/"+path, 0755)
	}

	// Create the route
	r.PathPrefix("/" + staticDir + "/").Handler(http.StripPrefix("/"+staticDir+"/", http.FileServer(http.Dir("./"+staticDir+"/"))))

	routes.Serve(r)

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	srv := &http.Server{
		Handler:      handlers.CORS(originsOk, headersOk, methodsOk)(loggedRouter),
		Addr:         ":" + os.Getenv("PORT"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Running on port : %s  \n", os.Getenv("PORT"))
	log.Fatal(srv.ListenAndServe())
}
