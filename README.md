# Orogenesis

Everybody seems to have a static site generator. This is mine.

Orogenesis builds [ironicmtn](http://www.ironicmtn.com).

## Subprograms

**oro-build**: constructs one or more pages when executed with YAML
configuration file(s) as arguments

```
    oro-build myconfig1.yaml [myconfig2.yaml...]
```

**oro-watch**: watches a directory filled with config files and rebuilds the
relevant page(s) whenever source files are modified

```
    oro-watch page_configs/
```

## YAML Keys

*template*: path to the template file (*required*)

*title-raw*: page title

*title-path*: path to file containing the page title

*header-raw*: page header in escaped HTML

*header-path*: path to file containing the page header

*nav-raw*: nav bar content in escaped HTML

*nav-path* path to file containing nav bar content

*body-raw*: page body in escaped HTML

*body-path*: path to file containing the page body

*footer-raw*: page footer in escaped HTML

*footer-path*: path to file containing the page footer

*output-html*: path for output

All paths are relative to the template path.

