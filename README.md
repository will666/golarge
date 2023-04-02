# golarge

CLI that list files of about 1GB in size, recursively from a given directory path.

## Features

Output result to file via **-o** or **--output** flag (default destination: list.txt).

Output result as JSON file via **-j** or **--json**. The destination file will be as specified in **-o** flag (default destination: list.json).

## Build from source

Install Go

Follow [Installation instructions](https://go.dev/doc/install).

Clone the repository

```shell
git clone https://github.com/will666/golarge.git
cd golarge
go build .
```

## Basic usage

See help

```shell
golarge help
```

Search large files at **/tmp**, save result to file **large_files.txt**.

```shell
golarge -o large_files.txt /tmp
```

Search large files at **/tmp**, save result to json format (default: **list.json**).

```shell
golarge -j /tmp
```

## Licence

[GNU General Public License v3.0](LICENSE)