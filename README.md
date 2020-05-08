# sqlutil

[![MIT License](https://img.shields.io/apm/l/atomic-design-ui.svg?)](LICENSE)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/lithdew/sqlutil)
[![Discord Chat](https://img.shields.io/discord/697002823123992617)](https://discord.gg/HZEbkeQ)

Utilities for working with SQL in Go.

- Append an opinionated UTF-8 string representation of an `interface{}` to a byte slice, with byte slices being encoded into RFC-4648 Base64.
- Convert `*sql.Rows` into JSON with minimal allocations.
- Convert `*sql.Rows` into CSV with minimal allocations.
- Parse/evaluate SQL statements with named parameters.

## Example

Given the SQL query:

```sqlite
CREATE TABLE test (id integer primary key autoincrement, name varchar);

INSERT INTO test (name) VALUES ('a');
INSERT INTO test (name) VALUES ('b');
INSERT INTO test (name) VALUES ('c');
INSERT INTO test (name) VALUES ('d');

SELECT * FROM test;
```

`sqlutil.RowsToJSON` would yield:

```json
[
  {
    "id": 1,
    "name": "a"
  },
  {
    "id": 2,
    "name": "b"
  },
  {
    "id": 3,
    "name": "c"
  },
  {
    "id": 4,
    "name": "d"
  }
]
```

`sqlutil.RowsToCSV` would yield:

```csv
id,name
1,"a"
2,"b"
3,"c"
4,"d"
```


## Benchmarks

```
go test -bench=. -benchmem -benchtime=10s

goos: linux
goarch: amd64
pkg: github.com/lithdew/sqlutil
BenchmarkRowsToJSON-8            1213502              9531 ns/op             584 B/op         27 allocs/op
BenchmarkRowsToCSV-8             1268085              9991 ns/op             584 B/op         27 allocs/op
```