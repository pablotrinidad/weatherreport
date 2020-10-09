# Weather Report

Weather report is a CLI used to obtain weather information from a given dataset of flight tickets.

## Folder structure

All source contained here is laid out as a valid [**Go module**](https://blog.golang.org/using-go-modules)
such that it can be easily installed and built with commands like `go get github.com/pablotrinidad/weatherreport`
and documentation gets automatically generated and published by tools like [pkg.go.dev](https://pkg.go.dev/).

Folders:
* `cli`: Is the application entrypoint that handles the dataset parsing and communication with the storage layer.
* `store`: Is the storage layer of the application, it abstracts away the communication with the third-party
services and cache system.

* `docs`: Project description written in LaTeX. `.tex` and `.pdf` output file can be found there. 

## Running the application

#### Environment variables
Set the following environment variables:

```build
OPENWEATHER_API_KEY
```

_You can use the command `export VAR_NAME=VAR_VALUE`_

#### Building and running

1. Place the source under `$GOPATH/src/github.com/pablotrinidad/weathereport`.
2. cd into `cli`, i.e: `cd cli/`
2. Run the app with `go run . -d DATASET_FILE -f DATASET_FORMAT`.