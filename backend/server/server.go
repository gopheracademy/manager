package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gopheracademy/manager/backend/generated"
	"github.com/pacedotdev/oto/otohttp"
)

// Run starts the http server and blocks
func Run() {
	var conferenceServer ConferenceService
	server := otohttp.NewServer()
	generated.RegisterConferenceService(server, conferenceServer)
	http.Handle("/oto/", server)
	http.Handle("/", http.FileServer(http.Dir(".")))
	fmt.Println("listening at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
