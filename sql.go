package sqlutil

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/valyala/fastjson"
	"reflect"
	"strconv"
)

// RowsToJSON appends to dst the JSON representation of a list of row result sets from a SQL query.
func RowsToJSON(dst []byte, rows *sql.Rows) ([]byte, error) {
	var a fastjson.Arena

	results := a.NewArray()

	cols, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch columns: %w", err)
	}

	vals := make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		var val interface{}
		vals[i] = &val
	}

	// TODO(kenta): take into account that there may be multiple result sets in one query

	for count := 0; rows.Next(); count++ {
		if err := rows.Scan(vals...); err != nil {
			return nil, fmt.Errorf("got an error while scanning: %w", err)
		}

		row := a.NewObject()

		for i := range vals {
			switch val := (*(vals[i].(*interface{}))).(type) {
			case uint:
				row.Set(cols[i], a.NewNumberString(strconv.FormatUint(uint64(val), 10)))
			case byte:
				row.Set(cols[i], a.NewNumberString(strconv.FormatUint(uint64(val), 10)))
			case uint16:
				row.Set(cols[i], a.NewNumberString(strconv.FormatUint(uint64(val), 10)))
			case uint32:
				row.Set(cols[i], a.NewNumberString(strconv.FormatUint(uint64(val), 10)))
			case uint64:
				row.Set(cols[i], a.NewNumberString(strconv.FormatUint(val, 10)))
			case int16:
				row.Set(cols[i], a.NewNumberString(strconv.FormatInt(int64(val), 10)))
			case int32:
				row.Set(cols[i], a.NewNumberString(strconv.FormatInt(int64(val), 10)))
			case int64:
				row.Set(cols[i], a.NewNumberString(strconv.FormatInt(val, 10)))
			case bool:
				if val {
					row.Set(cols[i], a.NewTrue())
				} else {
					row.Set(cols[i], a.NewFalse())
				}
			case string:
				row.Set(cols[i], a.NewString(val))
			case []byte: // Array blobs get encoded into base64.
				row.Set(cols[i], a.NewString(base64.StdEncoding.EncodeToString(val)))
			default:
				return nil, fmt.Errorf("got unknown type '%s' while scanning", reflect.TypeOf(val).String())
			}
		}

		results.SetArrayItem(count, row)
	}

	return results.MarshalTo(dst), nil
}
