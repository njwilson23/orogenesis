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
	Output       string `yaml:"output-html,omitempty"`
}

func (p Page) String() string {
	return fmt.Sprintf("%v\n%s\n%s\n%s\n", p.Title, p.Header, p.Body, p.Footer)
}

func (page *Page) gethtml(raw *string, path *string) template.HTML {
	var html string
	if len(*raw) == 0 {
		htmlbytes, err := ioutil.ReadFile(*path)
		if err != nil {
			fmt.Println(err)
		}
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

	basepath := filepath.Dir(path)
	if len(page.TitlePath) != 0 {
		page.TitlePath = filepath.Join(basepath, page.TitlePath)
	}
	if len(page.HeaderPath) != 0 {
		page.HeaderPath = filepath.Join(basepath, page.HeaderPath)
	}
	if len(page.BodyPath) != 0 {
		page.BodyPath = filepath.Join(basepath, page.BodyPath)
	}
	if len(page.FooterPath) != 0 {
		page.FooterPath = filepath.Join(basepath, page.FooterPath)
	}

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

			if len(pagePtr.Output) == 0 {
				fnm_html = filepath.Base(fnm[:len(fnm)-5]) + ".html"
			} else {
				fnm_html = pagePtr.Output
			}

			fmt.Println("writing to", fnm_html)
			fout, err := os.Create(fnm_html)
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
