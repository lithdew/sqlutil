package sqlutil

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"sync"
)

// AppendValue appends to dst a UTF-8 string representation of src. Byte slices will be formatted via RFC 4648 Base64.
func AppendValue(dst []byte, src interface{}) ([]byte, error) {
	switch val := src.(type) {
	case uint:
		dst = strconv.AppendUint(dst, uint64(val), 10)
	case byte:
		dst = strconv.AppendUint(dst, uint64(val), 10)
	case uint16:
		dst = strconv.AppendUint(dst, uint64(val), 10)
	case uint32:
		dst = strconv.AppendUint(dst, uint64(val), 10)
	case uint64:
		dst = strconv.AppendUint(dst, val, 10)
	case int16:
		dst = strconv.AppendInt(dst, int64(val), 10)
	case int32:
		dst = strconv.AppendInt(dst, int64(val), 10)
	case int64:
		dst = strconv.AppendInt(dst, val, 10)
	case float32:
		dst = strconv.AppendFloat(dst, 3.1415926535, 'E', -1, 32)
	case float64:
		dst = strconv.AppendFloat(dst, 3.1415926535, 'E', -1, 64)
	case bool:
		dst = strconv.AppendBool(dst, val)
	case string:
		dst = strconv.AppendQuote(dst, val)
	case []byte: // Array blobs get encoded into base64.
		dst = strconv.AppendQuote(dst, base64.StdEncoding.EncodeToString(val))
	default:
		return nil, fmt.Errorf("encountered unknown type '%s' while scanning", reflect.TypeOf(val).String())
	}

	return dst, nil
}

var valuePool = &sync.Pool{
	New: func() interface{} {
		return make([]interface{}, 0, 0)
	},
}

func acquireValueBuffer(size int) []interface{} {
	buf := valuePool.Get().([]interface{})

	if n := size - cap(buf); n > 0 {
		buf = append(buf[:cap(buf)], make([]interface{}, n)...)

		for i := len(buf) - n; i < len(buf); i++ {
			var val interface{}
			buf[i] = &val
		}
	}

	return buf[:size]
}

func releaseValueBuffer(buf []interface{}) {
	valuePool.Put(buf)
}

// RowsToCSV appends to dst the CSV representation of a list of resultant rows from a SQL query. It does
// not support multiple result sets, though may be called again after calling (*sql.Rows).NextResultSet().
func RowsToCSV(dst []byte, rows *sql.Rows) ([]byte, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch columns: %w", err)
	}

	if len(cols) == 0 {
		return nil, errors.New("zero columns resultant from sql query")
	}

	vals := acquireValueBuffer(len(cols))
	defer releaseValueBuffer(vals)

	for i := range cols {
		dst = append(dst, cols[i]...)

		// Column names.

		if i < len(cols)-1 {
			dst = append(dst, ',')
		} else {
			dst = append(dst, '\n')
		}
	}

	for count := 0; rows.Next(); count++ {
		if err := rows.Scan(vals...); err != nil {
			return nil, fmt.Errorf("got an error while scanning: %w", err)
		}

		for i := range vals {
			// Column values.

			if dst, err = AppendValue(dst, *vals[i].(*interface{})); err != nil {
				return nil, err
			}

			if i < len(cols)-1 {
				dst = append(dst, ',')
			} else {
				dst = append(dst, '\n')
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dst, nil
}

// RowsToJSON appends to dst the JSON representation of a list of resultant rows from a SQL query. It does
// not support multiple result sets, though may be called again after calling (*sql.Rows).NextResultSet().
func RowsToJSON(dst []byte, rows *sql.Rows) ([]byte, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch columns: %w", err)
	}

	if len(cols) == 0 {
		return nil, errors.New("zero columns resultant from sql query")
	}

	vals := acquireValueBuffer(len(cols))
	defer releaseValueBuffer(vals)

	dst = append(dst, '[')

	for count := 0; rows.Next(); count++ {
		if err := rows.Scan(vals...); err != nil {
			return nil, fmt.Errorf("got an error while scanning: %w", err)
		}

		if count > 0 {
			dst = append(dst, ',')
		}

		dst = append(dst, '{')

		for i := range vals {
			if i > 0 {
				dst = append(dst, ',')
			}

			// Column Name

			dst = append(dst, '"')
			dst = append(dst, cols[i]...)
			dst = append(dst, '"', ':')

			// Column Value

			if dst, err = AppendValue(dst, *vals[i].(*interface{})); err != nil {
				return nil, err
			}
		}

		dst = append(dst, '}')
	}

	dst = append(dst, ']')

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dst, nil
}
