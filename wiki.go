package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
	"text/template"
)

var (
	addr = flag.Bool("addr", false, "find open address and print to final-port.txt")
)

type Page struct {
	Title string
	Body  []byte
}

func pageName(name string) string {
	return "pages/" + name + ".txt"
}

func (p *Page) save() error {
	filename := pageName(p.Title)
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := pageName(title)
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func base_path(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, I love %s!", r.URL.Path[1:])
}

var validPath = regexp.MustCompile("^/(view|edit|save)/([a-zA-Z0-9]+)$")

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid path")
	}
	return m[2], nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles("templates/" + tmpl + ".html")
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

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editFile(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func makeWikiHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	flag.Parse()
	//http.HandleFunc("/bad/", badHandler)
	http.HandleFunc("/view/", makeWikiHandler(viewHandler))
	http.HandleFunc("/edit/", makeWikiHandler(editFile))
	http.HandleFunc("/save/", makeWikiHandler(saveHandler))
	http.HandleFunc("/", base_path)

	if *addr {
		l, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile("final-port.txt", []byte(l.Addr().String()), 0644)
		if err != nil {
			log.Fatal(err)
		}
		s := &http.Server{}
		s.Serve(l)
		return
	}

	http.ListenAndServe(":8080", nil)
}
