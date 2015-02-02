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

// Returns a template.HTML type with the HTML content from the raw string or
// path in a Page pointer. Used to construct specific getter methods.
func (page *Page) gethtml(raw *string, path *string) template.HTML {
	var html string
	if len(*raw) == 0 {
		if len(*path) == 0 {
			html = ""
		} else {
			htmlbytes, err := ioutil.ReadFile(*path)
			if err != nil {
				fmt.Println(err)
			}
			html = string(htmlbytes)
		}
	} else {
		html = *raw
	}
	return template.HTML(html)
}

// Return HTML types for page content
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

// Read a configuration file and return a pointer to a Page struct
func ReadConfig(path string) (*Page, error) {
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

func BuildPage(rootpath string, fout *os.File, page *Page) error {
	templatepath := filepath.Join(rootpath, page.TemplatePath)
	templatebytes, err := ioutil.ReadFile(templatepath)
	if err != nil {
		return err
	}
	t := template.Must(template.New("unnamed").Parse(string(templatebytes)))
	err = t.Execute(fout, page)
	return err
}

func OutputPath(page *Page, configpath, rootpath string) string {
	var fnmhtml string
	if len(page.Output) == 0 {
		fnmhtml = filepath.Base(configpath[:len(configpath)-5]) + ".html"
	} else {
		fnmhtml = filepath.Join(rootpath, page.Output)
	}
	return fnmhtml
}

func main() {

	var args []string
	args = os.Args[1:]
	if len(args) == 0 {
		fmt.Println("At least one content file must be specified")
	}

	var pagePtr *Page
	var fnmhtml, rootpath string
	for _, configpath := range args {
		if _, err := os.Stat(configpath); !os.IsNotExist(err) {

			// Read configuration
			fmt.Println("parsing", configpath)
			rootpath = filepath.Dir(configpath)
			pagePtr, err = ReadConfig(configpath)
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println("using template at", pagePtr.TemplatePath)
			fnmhtml = OutputPath(pagePtr, configpath, rootpath)
			fmt.Println("writing to", fnmhtml)

			// Write output
			fout, err := os.Create(fnmhtml)
			if err != nil {
				fmt.Println(err)
				break
			}

			err = BuildPage(rootpath, fout, pagePtr)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println(configpath, "does not exist")
		}
	}
}
