# Grep in Go
A simple grep implementation in Go. Inspired by Go Bootcamp here https://one2n.io/go-bootcamp/go-projects/grep-in-go

## Build
To build the application, run the following command:
```sh
go build -o grep main.go
```

## Usage
### Search within a single file
```sh
grep <pattern> <file_path>
```
Example:
```sh
grep "error" logs.txt
```

### Search within all files in a directory (including subdirectories)
```sh
grep <pattern> <directory_path>
```
Example:
```sh
grep "error" ./logs
```

## Assumptions
1. If a directory is provided as input, all files within its subdirectories are also considered.
2. All files in the given directory are assumed to be text files.
