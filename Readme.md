# Gopular

CLI tool for fetching the most popular repositories written in Go.

## How to use it

1.  add your GitHub token as environment variable.

    ```bash
    export GITHUB_TOKEN= VALUE
    ```

2.  build the binaries for your system.

    ```bash
    go build .
    ```

3.  Usage:
    gopular popular [flags]

    Flags:

          --p string   Programming Language (default "Go")
          --d string   Date in format of yyyy-mm-dd (default "2014-01-01")
          --c uint     Count (default 10)

## Todo

- [x] Fetch the most popular repositories.
- [x] Filter for the programming language.
- [x] Get popular repositories created after specific date.
- [x] Control the number of fetched repositories.
- [x] Build the tool using cobra.
- [ ] Test.
