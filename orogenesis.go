package orogenesis

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// First tries to convert `raw-[name]` to HTML. If the key is missing, tries instead
// to read a filename from `html-[name]`.
func getRawHTML(config map[string]string, key *string) (template.HTML, error) {
	var result template.HTML

	ss := []string{"raw-", *key}
	tempkey := strings.Join(ss, "")
	html, ok := config[tempkey]
	if !ok {
		return result, errors.New(fmt.Sprintf("key raw-%s not found", key))
	}
	result = template.HTML(html)
	return result, nil
}

func getExternalHTML(config map[string]string, key *string) (template.HTML, error) {
	var result template.HTML

	ss := []string{"html-", *key}
	tempkey := strings.Join(ss, "")
	html, ok := config[tempkey]
	if !ok {
		return result, errors.New(fmt.Sprintf("key html-%s not found", key))
	}
	htmlbytes, err := ioutil.ReadFile(html)
	if err != nil {
		return result, err
	}
	html = string(htmlbytes)
	result = template.HTML(html)
	return result, nil
}

// Read a configuration file and return a pointer to a Page struct
func ReadConfig(path string) (map[string]string, map[string]template.HTML, error) {

	var config map[string]string
	var htmlMap map[string]template.HTML
	var err error

	htmlMap = make(map[string]template.HTML)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return config, htmlMap, err
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, htmlMap, err
	}

	basepath := filepath.Dir(path)

	var ok bool
	var templatePath string
	templatePath, ok = config["oro-template"]
	if !ok {
		fmt.Println("'oro-template' key not defined")
		return config, htmlMap, errors.New("required 'oro-template' key not defined")
	}
	config["oro-template"] = filepath.Join(basepath, templatePath)

	var outputPath string
	outputPath, ok = config["oro-output"]
	if !ok {
		fmt.Println("'oro-output' key not defined")
		outputPath = filepath.Base(path[:len(path)-5]) + ".html"
	}
	config["oro-output"] = filepath.Join(basepath, outputPath)

	var keystem string
	var html template.HTML
	for key := range config {
		if strings.HasPrefix(key, "raw-") {
			keystem = strings.Replace(key, "raw-", "", 1)
			html, err = getRawHTML(config, &keystem)
			if err != nil {
				fmt.Println(err)
			}
			htmlMap[keystem] = html
			delete(config, key)
		} else if strings.HasPrefix(key, "html-") {
			config[key] = filepath.Join(basepath, config[key])
			keystem = strings.Replace(key, "html-", "", 1)
			html, err = getExternalHTML(config, &keystem)
			if err != nil {
				fmt.Println(err)
			}
			htmlMap[keystem] = html
			delete(config, key)
		}
	}
	return config, htmlMap, nil
}

func BuildPage(config map[string]string, data map[string]template.HTML) (string, error) {

	fnmHTML := config["oro-output"]
	templatePath := config["oro-template"]
	templateBytes, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return "", err
	}

	fout, err := os.Create(fnmHTML)
	if err != nil {
		return "", err
	}
	tmpl := template.Must(template.New("toplevel").Parse(string(templateBytes)))

	err = tmpl.Execute(fout, data)
	return fnmHTML, err
}
