/* oro-build

	oro-build config.yaml [config1.yaml [...]]

Tool to construct static HTML pages.
*/

package main

import (
	"fmt"
	"github.com/njwilson23/orogenesis"
	"os"
)

func main() {

	var args []string
	args = os.Args[1:]
	if len(args) == 0 {
		fmt.Println("At least one content file must be specified")
	}

	var pagePtr *orogenesis.Page
	var fnmhtml string
	for _, configpath := range args {
		if _, err := os.Stat(configpath); !os.IsNotExist(err) {

			// Read configuration
			fmt.Println("parsing", configpath)
			pagePtr, err = orogenesis.ReadConfig(configpath)
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println("  using template at", pagePtr.TemplatePath)

			fnmhtml, err = orogenesis.BuildPage(configpath, pagePtr)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(" ", fnmhtml, "written")
			}

		} else {
			fmt.Println(configpath, "does not exist")
		}
	}
}
