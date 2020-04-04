package sqlutil

import (
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"
)

import _ "github.com/mattn/go-sqlite3"

func createDatabase(t testing.TB) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})

	db.SetMaxOpenConns(1)

	_, err = db.Exec("CREATE TABLE test (id integer primary key autoincrement, name varchar)")
	require.NoError(t, err)

	test := []string{"a", "b", "c", "d"}

	for _, name := range test {
		_, err = db.Exec("INSERT INTO test (name) VALUES (?)", name)
		require.NoError(t, err)
	}

	return db
}

func createQuery(t testing.TB, db *sql.DB) *sql.Stmt {
	t.Helper()

	stmt, err := db.Prepare("SELECT * FROM test")
	require.NoError(t, err)
	return stmt
}

func TestRowsToJSON(t *testing.T) {
	stmt := createQuery(t, createDatabase(t))
	rows, err := stmt.Query()
	require.NoError(t, err)

	json, err := RowsToJSON(nil, rows)
	require.NoError(t, err)

	require.EqualValues(
		t, `[{"id":1,"name":"a"},{"id":2,"name":"b"},{"id":3,"name":"c"},{"id":4,"name":"d"}]`, string(json),
	)
}

func BenchmarkRowsToJSON(b *testing.B) {
	stmt := createQuery(b, createDatabase(b))

	b.ReportAllocs()
	b.ResetTimer()

	var buf [1024]byte

	for i := 0; i < b.N; i++ {
		rows, err := stmt.Query()
		if err != nil {
			b.Fatal(err)
		}

		if _, err = RowsToJSON(buf[:0], rows); err != nil {
			b.Fatal(err)
		}
	}
}
