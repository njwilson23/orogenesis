package main

import (
	"fmt"
	"testing"
)

func TestRawSources(t *testing.T) {
	page := Page{"tests/example-template.html",
		"", "<b>Title</b>", // title
		"", "<link href=\"example.css\" rel=\"stylesheet\">\n<!-- Javascript examples -->", // header
		"", "<p>Test with raw sources</p>", // body
		"", "<div class=\"centered\">This content is in the Public Domain</div>", // footer
		"tests/nofile.html", // output path
	}

	title := page.Title()
	header := page.Header()
	body := page.Body()
	footer := page.Footer()

	fmt.Println(title)
	fmt.Println(header)
	fmt.Println(body)
	fmt.Println(footer)

}

func TestPathSources(t *testing.T) {
	page := Page{"tests/example-template.html",
		"tests/example/title.html", "", // title
		"tests/example/header.html", "", // header
		"tests/example/body.html", "", // body
		"tests/example/footer.html", "", // footer
		"tests/nofile.html", // output path
	}

	title := page.Title()
	header := page.Header()
	body := page.Body()
	footer := page.Footer()

	fmt.Println(title)
	fmt.Println(header)
	fmt.Println(body)
	fmt.Println(footer)

}
