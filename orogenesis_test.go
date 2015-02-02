package orogenesis

import (
	"fmt"
	"os"
	"testing"
)

func TestRawSources(t *testing.T) {
	page := Page{"tests/example-template.html",
		"", "<b>Title fragment</b>", // title
		"", "<link href=\"example.css\" rel=\"stylesheet\">\n<!-- Javascript examples -->", // header
		"", "<p>Test with raw sources</p>", // body
		"", "<div class=\"centered\">This content is in the Public Domain</div>", // footer
		"tests/nofile.html", // output path
	}

	if string(page.Title()) != "<b>Title fragment</b>" {
		fmt.Println(page.Title())
		t.Fail()
	}

	if string(page.Header()) != "<link href=\"example.css\" rel=\"stylesheet\">\n<!-- Javascript examples -->" {
		fmt.Println(page.Header())
		t.Fail()
	}

	if string(page.Body()) != "<p>Test with raw sources</p>" {
		fmt.Println(page.Body())
		t.Fail()
	}

	if string(page.Footer()) != "<div class=\"centered\">This content is in the Public Domain</div>" {
		fmt.Println(page.Footer())
		t.Fail()
	}
}

func TestPathSources(t *testing.T) {
	page := Page{"tests/example-template.html",
		"tests/example/title.html", "", // title
		"tests/example/header.html", "", // header
		"tests/example/body.html", "", // body
		"tests/example/footer.html", "", // footer
		"tests/nofile.html", // output path
	}

	if string(page.Title()) != "<b>Title fragment</b>\n" {
		t.Fail()
	}

	if string(page.Header()) != "<link href=\"example.css\" rel=\"stylesheet\">\n<!-- Javascript examples -->\n" {
		t.Fail()
	}

	if string(page.Body()) != "<p>Test with path sources</p>\n" {
		t.Fail()
	}

	if string(page.Footer()) != "<div class=\"centered\">This content is in the Public Domain</div>\n" {
		t.Fail()
	}
}

func TestPageConstruction(t *testing.T) {
	page, err := ReadConfig("tests/example.yaml")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	// Write output
	fout, err := os.Create("test.html")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	err = BuildPage("tests", fout, page)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	if page == nil {
		t.Fail()
	}
}
