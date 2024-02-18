# Custom Log Hunter
Custom Log Searcher is a command-line tool for searching through log files based on custom criteria using Gojq queries.

## Prerequisites
Before using Custom Log Searcher, ensure that you have the following installed:

- Go programming language (https://golang.org/doc/install)

## Installation

1. Clone the repository:

  ```bash
  git clone https://github.com/rahulmishra24/Custom-Log-Hunter
  ```

2. Change into the porject directoy:
   
 ```bash
 cd Custom-Log-Hunter
 ```
3. Build the executable
   
  ```bash
  go build
  ```
## Usage

Custom Log Searcher allows you to search through log files using Gojq queries.

```bash
./logger -query "<your-query>" -dir "<directory-path>"
```

## Examples

```bash
./logger -query ".message | test(\"ERROR\")" -dir "logs"
```



