package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	_ "net/http/pprof"

	"github.com/andersjanmyr/jobs/controllers"
	"github.com/andersjanmyr/jobs/models"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/panicwrap"
	log "github.com/sirupsen/logrus"
)

func main() {
	_, err := panicwrap.BasicWrap(panicHandler)
	if err != nil {
		// Something went wrong setting up the panic wrapper. Unlikely,
		// but possible.
		panic(err)
	}
	port := 5555

	log.SetOutput(os.Stdout)
	var router = mux.NewRouter().StrictSlash(true)
	loggedRouter := handlers.LoggingHandler(os.Stdout, slowMiddleware(router))
	db, err := gorm.Open("postgres", "host=localhost user=jobs dbname=jobs sslmode=disable password=jobs")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	jobRepo := models.NewPgJobRepo(db)
	_, _ = jobRepo.Add(models.NewJob("One"))
	_, _ = jobRepo.Add(models.NewJob("Two"))
	setupRouter(router.PathPrefix("/jobs"), controllers.NewJobController(jobRepo))

	go func() {
		log.Print("Profile server started on port 6060")
		log.Fatal(http.ListenAndServe("127.0.0.1:6060", nil))
	}()
	log.Print("Server started on port ", strconv.Itoa(port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), loggedRouter))
}

func panicHandler(output string) {
	// output contains the full output (including stack traces) of the
	// panic. Put it in a file or something.
	log.Panic(output)
	os.Exit(1)
}

func slowMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for i := 0; i < 10000; i++ {
			fmt.Fprint(os.Stderr, i)
		}
		next.ServeHTTP(w, r)
	})
}
