package orogenesis

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
	NavPath      string `yaml:"nav"`
	NavRaw       string `yaml:"nav-raw"`
	BodyPath     string `yaml:"body"`
	BodyRaw      string `yaml:"body-raw"`
	FooterPath   string `yaml:"footer"`
	FooterRaw    string `yaml:"footer-raw"`
	Output       string `yaml:"output-html,omitempty"`
}

func (p Page) String() string {
	return fmt.Sprintf("%v\n%s\n%s\n%s\n%s\n", p.Title, p.Header, p.Nav, p.Body, p.Footer)
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

func (page *Page) Nav() template.HTML {
	return page.gethtml(&page.NavRaw, &page.NavPath)
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
	if len(page.TemplatePath) != 0 {
		page.TemplatePath = filepath.Join(basepath, page.TemplatePath)
	}
	if len(page.TitlePath) != 0 {
		page.TitlePath = filepath.Join(basepath, page.TitlePath)
	}
	if len(page.HeaderPath) != 0 {
		page.HeaderPath = filepath.Join(basepath, page.HeaderPath)
	}
	if len(page.NavPath) != 0 {
		page.NavPath = filepath.Join(basepath, page.NavPath)
	}
	if len(page.BodyPath) != 0 {
		page.BodyPath = filepath.Join(basepath, page.BodyPath)
	}
	if len(page.FooterPath) != 0 {
		page.FooterPath = filepath.Join(basepath, page.FooterPath)
	}

	return &page, err
}

func BuildPage(configpath string, page *Page) (string, error) {

	fnmhtml := OutputPath(page, configpath)

	fout, err := os.Create(fnmhtml)
	if err != nil {
		fmt.Println(err)
	}

	//templatepath := filepath.Join(filepath.Dir(configpath), page.TemplatePath)
	templatebytes, err := ioutil.ReadFile(page.TemplatePath)
	if err != nil {
		return "", err
	}
	t := template.Must(template.New("unnamed").Parse(string(templatebytes)))

	err = t.Execute(fout, page)
	return fnmhtml, err
}

func OutputPath(page *Page, configpath string) string {
	var fnmhtml string

	if len(page.Output) == 0 {
		fnmhtml = filepath.Base(configpath[:len(configpath)-5]) + ".html"
	} else {
		rootpath := filepath.Dir(configpath)
		fnmhtml = filepath.Join(rootpath, page.Output)
	}

	return fnmhtml
}
