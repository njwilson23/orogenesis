# Orogenesis

Orogenesis is the tool used to build [ironicmtn](http://www.ironicmtn.com).

## Subprograms

**oro-build**: constructs one or more pages when executed with YAML
configuration file(s) as arguments

## YAML Keys

**template**: path to the template file (*required*)

**title-raw**: page title

**title-path**: path to file containing the page title

**header-raw**: page header in escaped HTML

**header-path**: path to file containing the page header

**body-raw**: page body in escaped HTML

**body-path**: path to file containing the page body

**footer-raw**: page footer in escaped HTML

**footer-path**: path to file containing the page footer

**output-html**: path for output

All paths are relative to the template path.

