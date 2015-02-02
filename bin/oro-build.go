package main

import (
	"fmt"
	"github.com/njwilson23/orogenesis"
	"os"
	"path/filepath"
)

func main() {

	var args []string
	args = os.Args[1:]
	if len(args) == 0 {
		fmt.Println("At least one content file must be specified")
	}

	var pagePtr *orogenesis.Page
	var fnmhtml, rootpath string
	for _, configpath := range args {
		if _, err := os.Stat(configpath); !os.IsNotExist(err) {

			// Read configuration
			fmt.Println("parsing", configpath)
			rootpath = filepath.Dir(configpath)
			pagePtr, err = orogenesis.ReadConfig(configpath)
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println("using template at", pagePtr.TemplatePath)
			fnmhtml = orogenesis.OutputPath(pagePtr, configpath, rootpath)
			fmt.Println("writing to", fnmhtml)

			// Write output
			fout, err := os.Create(fnmhtml)
			if err != nil {
				fmt.Println(err)
				break
			}

			err = orogenesis.BuildPage(rootpath, fout, pagePtr)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println(configpath, "does not exist")
		}
	}
}
