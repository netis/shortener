package main

import (
	"net/http"
	"fmt"
	"shortener/store"
)

var st = store.NewURLStore()

func main() {
	http.HandleFunc("/", Redirect)
	http.HandleFunc("/add", Add)
	http.ListenAndServe(":8080", nil)
}

func Add(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	fmt.Printf("url = %s\n", url)
	if url == "" {
		fmt.Fprintf(w, AddForm)
		return
	}
	key := st.Put(url)
	fmt.Fprintf(w, "http://localhost:8080/%s", key)
}

const AddForm = `
	<html>
	<form method="POST" action="/add">
	URL: <input type="text" name="url">
	<input type="submit" value="Add">
	</form>
	</html>`

func Redirect(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]
	fmt.Printf("Path is %s\n", r.URL.Path)
	url := st.Get(key)
	fmt.Printf("url is %s\n", url)
	if url == "" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "http://"+url, http.StatusFound)
}

