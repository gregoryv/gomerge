# Changelog
All notable changes to this project will be documented in this file.

The format is based on http://keepachangelog.com/en/1.0.0/
and this project adheres to http://semver.org/spec/v2.0.0.html.

## [0.6.0] - 2022-02-19

- Fix Split for files without imports
- Command moved to github.com/gregoryv/gomerge/cmd/gomerge

## [0.5.0] - 2022-02-19

- Add flag -i to include src filename in merged result
- Include src header if any

## [0.4.0] - 2022-02-19

- Rewrite parsing, simplify command to only merge two files
- Fix issue of missing declarations

## [0.3.0] - 2022-02-15

- Don't panic when trying to merge non go files
- Fix multifile merge

## [0.2.1] - 2022-02-14

- Filter out duplicate imports

## [0.2.0] - 2022-02-14

- Implement flag -w for writing to destination

## [0.1.0] - 2022-02-14

- Add gomerge for merging go files and removing duplicate imports
