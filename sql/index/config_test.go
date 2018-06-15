package index

import (
	"crypto/sha1"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/src-d/go-mysql-server.v0/sql"
)

func TestConfig(t *testing.T) {
	require := require.New(t)

	db, table, id := "db_name", "table_name", "index_id"
	path := filepath.Join(os.TempDir(), db, table, id)
	err := os.MkdirAll(path, 0750)

	require.NoError(err)
	defer os.RemoveAll(path)

	h1 := sha1.Sum([]byte("h1"))
	h2 := sha1.Sum([]byte("h2"))
	exh1 := sql.ExpressionHash(h1[:])
	exh2 := sql.ExpressionHash(h2[:])

	cfg1 := NewConfig(
		db,
		table,
		id,
		[]sql.ExpressionHash{exh1, exh2},
		"DriverID",
		map[string]string{
			"port": "10101",
			"host": "localhost",
		},
	)

	err = WriteConfigFile(path, cfg1)
	require.NoError(err)

	cfg2, err := ReadConfigFile(path)
	require.NoError(err)
	require.Equal(cfg1, cfg2)

	require.NoError(SetConfigFileReady(path))

	cfg3, err := ReadConfigFile(path)
	require.NoError(err)
	require.True(cfg3.Ready)
}
