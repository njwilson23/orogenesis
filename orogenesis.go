package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

type Page struct {
	Template string
	Title    string
	Header   string
	Body     string
	Footer   string
}

func (p Page) String() string {
	return fmt.Sprintf("%v\n%s\n%s\n%s\n", p.Title, p.Header, p.Body, p.Footer)
}

func readcontent(path string) (*Page, error) {
	var page Page
	var err error
	pagedata, err := ioutil.ReadFile(path)
	if err != nil {
		return &page, err
	}
	err = yaml.Unmarshal(pagedata, &page)
	return &page, err
}

func buildpage(templatefile string, fout *os.File, page *Page) error {
	templatebytes, err := ioutil.ReadFile(templatefile)
	if err != nil {
		return err
	}
	templatename := templatefile[:len(templatefile)-5]
	t := template.Must(template.New(templatename).Parse(string(templatebytes)))
	err = t.Execute(fout, page)
	return err
}

func main() {

	var args []string
	args = os.Args[1:]
	if len(args) == 0 {
		fmt.Println("At least one content file must be specified")
	}

	var page *Page
	var fnm_html, templatepath string
	for _, fnm := range args {
		if _, err := os.Stat(fnm); !os.IsNotExist(err) {

			fmt.Println("parsing", fnm)
			page, err = readcontent(fnm)
			if err != nil {
				fmt.Println(err)
				break
			}

			fnm_html = filepath.Base(fnm[:len(fnm)-5]) + ".html"
			fmt.Println("writing to", fnm_html)
			fout, err := os.Create(fnm_html)
			defer fout.Close()
			if err != nil {
				fmt.Println(err)
				break
			}

			templatepath = filepath.Join(filepath.Dir(fnm), page.Template)
			err = buildpage(templatepath, fout, page)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println(fnm, "does not exist")
		}
	}
}
