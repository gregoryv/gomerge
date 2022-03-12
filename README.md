[gomerge](https://pkg.go.dev/github.com/gregoryv/gomerge) - constructs for merging Go files

When dealing with over structured(to many directory and files)
repositories, one route to tidying them is to merge files with related
concepts. This often simplifies additional refactoring.

The provided cmd/gomerge removes the manual steps of concatenating go
files and removing duplicate package imports.

## Quick start

    $ go install github.com/gregoryv/gomerge/cmd/gomerge@latest
	$ gomerge -h
	Usage: gomerge [OPTIONS] DST SRC
    Options
      -i    include src filename in merged as comment
      -r    removes source after merge(only with -w)
      -w    writes result to destination file

