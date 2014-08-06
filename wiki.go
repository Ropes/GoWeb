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

type D3Page struct {
	Title string
	Data  interface{}
}

func pageName(name string) string {
	return "pages/" + name + ".txt"
}

func (p *Page) save() error {
	filename := pageName(p.Title)
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func staticName(name string) string {
	return "static/" + name
}

func loadPage(title string) (*Page, error) {
	filename := pageName(title)
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func loadStatic(fpath string) ([]byte, error) {
	body, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func base_path(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, I love %s!", r.URL.Path[1:])
}

var validPath = regexp.MustCompile("^/(view|edit|save)/([a-zA-Z0-9]+)$")
var templates map[string]*template.Template

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid path")
	}
	return m[2], nil
}

func templateInit() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	templates["view"] = template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/view.html"))
	templates["edit"] = template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/edit.html"))
	templates["d3"] = template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/ddd.html"))
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) error {
	t, ok := templates[tmpl]
	if !ok {
		return errors.New("Template not found!")
	}

	err := t.Execute(w, p)
	if err != nil {
		fmt.Println(err)
	}
	return nil
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

func d3Handler(w http.ResponseWriter, r *http.Request) {
	t, ok := templates["d3"]
	if !ok {
		fmt.Println(errors.New("Template not found!"))
	}
	dp := &D3Page{Title: "D3 Demo", Data: nil}
	err := t.Execute(w, dp)
	if err != nil {
		fmt.Println(err)
	}
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	fpath := r.URL.Path[1:]
	static, err := loadStatic(fpath)
	//fmt.Fprintf(w, "Hello, I love %s!", r.URL.Path[1:])
	if err != nil {
		fmt.Fprintf(w, "Unable to return static asset: %s!", fpath)
	}
	fmt.Fprintf(w, "%s", static)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
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
	templateInit()
	flag.Parse()
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editFile))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/", base_path)

	//D3 Example handlers
	http.HandleFunc("/d3", d3Handler)
	http.HandleFunc("/static/", staticHandler)

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
