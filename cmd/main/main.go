package main

import (
	server "github.com/bear1278/MusicWave"
	"github.com/bear1278/MusicWave/pkg/handlers"
	"log"
)

/*func hundler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hi, %s", r.URL.Path[1:])
	t, err := template.ParseFiles("./public/index.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, "")
}*/

func main() {
	/*	http.HandleFunc("/", hundler)
		log.Fatal(http.ListenAndServe(":8080", nil))*/
	handler := new(handlers.Handler)
	srv := new(server.Server)
	err := srv.Run("8080", handler.InitRoutes())
	if err != nil {
		log.Fatalf("error occured while run http server: %s", err.Error())
	}
}
