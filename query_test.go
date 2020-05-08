package sqlutil

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseNamedQuery(t *testing.T) {
	result, err := ParseNamedQuery("select * from items where id = :id and test = :test")
	require.NoError(t, err)

	require.EqualValues(t, "select * from items where id = ? and test = ?", result.Parsed)
	require.EqualValues(t, []string{"id", "test"}, result.Names)

	result, err = ParseNamedQuery("select * from items")
	require.NoError(t, err)
	require.EqualValues(t, result.Query, result.Parsed)
	require.Empty(t, result.Names)
}
