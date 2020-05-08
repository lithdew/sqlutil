package sqlutil

import (
	"bytes"
	"errors"
	"github.com/lithdew/bytesutil"
	"strings"
)

var ErrNamedQueryMalformed = errors.New("named query is malformed")

// NamedQuery is the parsed result of a SQL statement with named parameters.
type NamedQuery struct {
	Query  string   `json:"query"`
	Parsed string   `json:"parsed"`
	Names  []string `json:"names"`
}

// ParseNamedQuery parses a SQL statement with named parameters, replaces those named parameters with ?, and returns
// the parameter names sorted in the order from the start to the end of the SQL statement.
func ParseNamedQuery(query string) (res NamedQuery, err error) {
	res.Query = query

	buf := strings.Builder{}
	ptr := bytesutil.Slice(query)

	for {
		i := bytes.IndexByte(ptr, ':')
		if i == -1 {
			buf.Write(ptr)
			break
		}

		buf.Write(ptr[:i])
		buf.WriteByte('?')

		ptr = ptr[i+1:]

		end := bytes.IndexByte(ptr, ' ')
		if end == -1 {
			name := ptr
			if len(name) == 0 {
				return res, ErrNamedQueryMalformed
			}
			res.Names = append(res.Names, bytesutil.String(name))
			break
		}

		name := ptr[:end]
		if len(name) == 0 {
			return res, ErrNamedQueryMalformed
		}
		res.Names, ptr = append(res.Names, bytesutil.String(name)), ptr[end:]
	}

	res.Parsed = buf.String()

	return res, err
}
