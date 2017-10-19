package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"

	"lanni.com/wiki"
)

/*func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}*/

var (
	addr      = flag.Bool("addr", false, "find open address and print to final-port.txt")
	templates = template.Must(template.ParseFiles("edit.html", "view.html"))
	validPath = regexp.MustCompile("^/(edit|save|view|json)/([0-9A-Za-z]+)$")
)

func getTitle(path string) string {
	p := strings.TrimRight(path, "/")
	title := path[strings.LastIndex(p, "/")+1:]
	return title
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	//path := strings.TrimRight(r.URL.Path, "/")
	//title := path[strings.LastIndex(path, "/")+1:]
	p, err := wiki.LoadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	//title := getTitle(r.URL.Path)
	p, err := wiki.LoadPage(title)
	if err != nil {
		p = &wiki.Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	//title := getTitle(r.URL.Path)
	body := r.FormValue("body")
	p := &wiki.Page{Title: title, Body: []byte(body)}
	p.Save()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func jsonHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, _ := wiki.LoadPage(title)
	encoder := json.NewEncoder(w)
	encoder.Encode(p)
}

func renderTemplate(w http.ResponseWriter, tpl string, p *wiki.Page) {
	err := templates.ExecuteTemplate(w, tpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	//t, _ := template.ParseFiles(tpl + ".html")
	//t.Execute(w, p)
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
	flag.Parse()
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/json/", makeHandler(jsonHandler))
	if *addr {
		l, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile("fina-port.txt", []byte(l.Addr().String()), 0644)
		if err != nil {
			log.Fatal(err)
		}
		s := &http.Server{}
		s.Serve(l)
		return
	}

	http.ListenAndServe(":8080", nil)
}
