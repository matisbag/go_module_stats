# Go Module Stats

## Assignment Instructions

This project is part of a programming assignment. Below are the instructions for the assignment:

### Task

Write a Go program, `index_godev`, that queries the official index of Go modules over HTTP. The program should accumulate modules and their versions by forge, generating a table organized by forge. The table format and sorting criteria are specified as follows:


```yaml
                         Forge                Modules    Versions
                  github.com                235782     1464248
                      k8s.io                   332       70345
                  gitlab.com                 4086       24155
                    gopkg.in                 2134       11976
                     . . . . . . . . . . . . . . . . . . . .
          gitlab.brurberg.no                    1           1
                    _Totals_                253732     1692307

```

### Sorting Criteria

The table is sorted based on the following criteria:

1. Versions in descending order
2. Modules in descending order
3. Forge in ascending order

### Protocol Specification

The program adheres to the protocol specified at [https://index.golang.org/](https://index.golang.org/).

### Usage

To limit the load submitted to Google, the program includes the `Disable-Module-Fetch: true` header in HTTP requests. Additionally, it handles errors gracefully.

### Implementation Details

- JSON deserialization is performed using the `encoding/json` package.
- The table is formatted using the `tabwriter` package.


## Running the Program

```bash
go run main.go
```