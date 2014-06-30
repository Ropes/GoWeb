package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"text/template"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func base_path(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, I love %s!", r.URL.Path[1:])
}

var validPath = regexp.MustCompile("^(view|edit|save|bad)/[a-zA-Z0-9]+$")

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid path")
	}
	fmt.Println(m)
	return m[2], nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(w, p)
}

func badHandler(w http.ResponseWriter, r *http.Request) {
	_, err := getTitle(w, r)
	if err != nil {
		return
	}
	p := &Page{Title: "Bad page!", Body: []byte("Bad page path requested!")}
	renderTemplate(w, "badpath", p)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		http.Redirect(w, r, "/bad/", http.StatusFound)
		return
	}
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editFile(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	http.HandleFunc("/bad/", badHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editFile)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/", base_path)
	http.ListenAndServe(":8080", nil)
}
