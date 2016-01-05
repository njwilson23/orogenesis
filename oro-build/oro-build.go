/* oro-build

	oro-build config.yaml [config1.yaml [...]]

Tool to construct static HTML pages.
*/

package main

import (
	"fmt"
	"github.com/njwilson23/orogenesis"
	"html/template"
	"os"
)

func main() {

	var args []string
	args = os.Args[1:]
	if len(args) == 0 {
		fmt.Println("At least one content file must be specified")
	}

	var config map[string]string
	var data map[string]template.HTML
	var fnmhtml string
	for _, configPath := range args {
		if _, err := os.Stat(configPath); !os.IsNotExist(err) {

			// Read configuration
			fmt.Println("parsing", configPath)
			config, data, err = orogenesis.ReadConfig(configPath)
			if err != nil {
				fmt.Println(err)
				fmt.Println("  building from", configPath, "failed")
				break
			}
			fmt.Println("  using template at", config["oro-template"])

			fnmhtml, err = orogenesis.BuildPage(config, data)
			if err != nil {
				fmt.Println(err)
				fmt.Println("  building from", configPath, "failed")
				break
			}
			fmt.Println(" ", fnmhtml, "written")

		} else {
			fmt.Println(configPath, "does not exist")
		}
	}
}
