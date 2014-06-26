package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
)

func main() {
	log.Println("Listening...")

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	r := mux.NewRouter()

	r.HandleFunc("/{page}", serveHandler)
	r.HandleFunc("/", rootHandler)
	//r.NotFoundHandler(redirectToRoot)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":3000", nil))

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func redirectToRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func serveHandler(w http.ResponseWriter, r *http.Request) {
	lp := path.Join("templates", "layout.html")
	fp := path.Join("templates", r.URL.Path)

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	templates, err := template.ParseFiles(lp, fp)
	if err != nil {
		log.Print(err)
		http.Error(w, "500 Internal Server Error", 500)
		return
	}
	templates.ExecuteTemplate(w, "layout", nil)
}
