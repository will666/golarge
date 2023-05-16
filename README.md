# golarge

Look for large files recursively from given directory path

## Features

- Recursively analyse files
- export results to text file
- export results to JSON file
- logging to stdout or file
- concurrent file processing

## Manual build from source (local)

Install Go/Golang

Follow [Installation instructions](https://go.dev/doc/install).

Clone the repository

```shell
git clone https://github.com/will666/golarge.git
cd golarge
make prod
```

## Install binay (remote)

Note: golang must be installed locally

```shell
go install https://github.com/will666/golarge@latest
```

## Run binary (remote)

Note: golang must be installed locally

```shell
go run https://github.com/will666/golarge@latest /foo
```

## Usage

```
Usage: golarge [OPTIONS] <PATH>

Look for files bigger than 1GiB from given directory path

Options:
  -o <string>     Output path (default: list.txt)
  -j      	      Enable export to JSON file
  -v      	      Display warnings & error instead of logging to file
  -t      	      Enable concurrent file processing

Examples:
  golarge /foo/bar
  golarge -v /bar
  golarge -o list.txt -j /foo/bar
  golarge -t /foo
```

### See help

```shell
golarge -h
```

### Search files at **/foo/bar**, save result to file **large_files.txt**.

```shell
golarge -o large_files.txt /foo/bar
```

### Search files at **/foo**, save result to json format (default: **list.json**).

```shell
golarge -j /foo
```

### Concurrent file processing (goroutines)

Enabling concurrent processing can greatly decrease time to result on most machines with fast storage, it could be slower on old hardware and mechanical storage. When enabled, each file is processed in an independant thread (Goroutine).

```shell
golarge -t /bar
```

## Licence

[MIT](LICENSE)

## Note

[clarge](https://github.com/will666/clarge): Another implementation in C