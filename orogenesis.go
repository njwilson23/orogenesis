package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Page struct {
	TemplatePath string `yaml:"template"`
	TitlePath    string `yaml:"title"`
	TitleRaw     string `yaml:"title-raw"`
	HeaderPath   string `yaml:"header"`
	HeaderRaw    string `yaml:"header-raw"`
	BodyPath     string `yaml:"body"`
	BodyRaw      string `yaml:"body-raw"`
	FooterPath   string `yaml:"footer"`
	FooterRaw    string `yaml:"footer-raw"`
	Output       string `yaml:"output-html"`
}

func (p Page) String() string {
	return fmt.Sprintf("%v\n%s\n%s\n%s\n", p.Title, p.Header, p.Body, p.Footer)
}

func (page *Page) gethtml(raw *string, path *string) template.HTML {
	var html string
	if &raw == nil {
		htmlbytes, _ := ioutil.ReadFile(*path)
		html = string(htmlbytes)
	} else {
		html = *raw
	}
	return template.HTML(html)
}

func (page *Page) Title() template.HTML {
	return page.gethtml(&page.TitleRaw, &page.TitlePath)
}

func (page *Page) Header() template.HTML {
	return page.gethtml(&page.HeaderRaw, &page.HeaderPath)
}

func (page *Page) Body() template.HTML {
	return page.gethtml(&page.BodyRaw, &page.BodyPath)
}

func (page *Page) Footer() template.HTML {
	return page.gethtml(&page.FooterRaw, &page.FooterPath)
}

func readconfig(path string) (*Page, error) {
	var page Page
	var err error

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return &page, err
	}
	err = yaml.Unmarshal(data, &page)
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

	var pagePtr *Page
	var fnm_html, templatepath string
	for _, fnm := range args {
		if _, err := os.Stat(fnm); !os.IsNotExist(err) {

			fmt.Println("parsing", fnm)
			pagePtr, err = readconfig(fnm)
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println("using template at", pagePtr.TemplatePath)

			if &pagePtr.Output != nil {
				fnm_html = pagePtr.Output
			} else {
				fnm_html = filepath.Base(fnm[:len(fnm)-5]) + ".html"
			}

			fmt.Println("writing to", fnm_html)
			fout, err := os.Create(fnm_html)
			defer fout.Close()
			if err != nil {
				fmt.Println(err)
				break
			}

			templatepath = filepath.Join(filepath.Dir(fnm), pagePtr.TemplatePath)
			err = buildpage(templatepath, fout, pagePtr)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println(fnm, "does not exist")
		}
	}
}
