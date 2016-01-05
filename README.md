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

Keys prefixed by `oro-` are used to declare templates and output paths. The
`html-` prefix indicates an HTML source file. The `raw-` prefix is for direct
code inclusions.

All paths are relative to the template path.

