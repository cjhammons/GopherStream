# Gopher Stream

A self-hosted music streaming solution. It scans all the music in a specified library directory and serves them via REST API for browsing and streaming.


# Dev Guide

## Configuration for Local Running

To run locally you first need to select a directory to use as your Library. This directory can have as many subdirectories as you like, it will not affect functionality. Make note of the filepath of the directory.

Clone the repositoty.

Create a file called ``config.json`` in the root directory of the Gopher Stream project and populate it like so:

```json
{
  "libraryDirectory": "<your library directory>"
}
```

Then build and run:
```bash
go build
./gopherstream
```
