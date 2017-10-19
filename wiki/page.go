package wiki

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) Save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func LoadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	fp, _ := filepath.Abs(filename)
	fmt.Println(fp)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}
