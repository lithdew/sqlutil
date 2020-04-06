# sqlutil

Utilities for working with SQL in Go.

- Append an opinionated UTF-8 string representation of an `interface{}` to a byte slice, with byte slices being encoded into RFC-4648 Base64.
- Convert `*sql.Rows` into JSON with minimal allocations.
- Convert `*sql.Rows` into CSV with minimal allocations.

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
BenchmarkRowsToJSON-8            1217752             10560 ns/op             616 B/op         29 allocs/op
BenchmarkRowsToCSV-8             1000000             10157 ns/op             616 B/op         29 allocs/op
```