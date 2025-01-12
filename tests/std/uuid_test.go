package std

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestStdUUID(t *testing.T) {
	if conn, err := sql.Open("clickhouse", "clickhouse://127.0.0.1:9000"); assert.NoError(t, err) {
		const ddl = `
			CREATE TABLE test_uuid (
				  Col1 UUID
				, Col2 UUID
			) Engine Memory
		`
		if _, err := conn.Exec("DROP TABLE IF EXISTS test_uuid"); assert.NoError(t, err) {
			if _, err := conn.Exec(ddl); assert.NoError(t, err) {
				scope, err := conn.Begin()
				if !assert.NoError(t, err) {
					return
				}

				if batch, err := scope.Prepare("INSERT INTO test_uuid"); assert.NoError(t, err) {
					var (
						col1Data = uuid.New()
						col2Data = uuid.New()
					)
					if _, err := batch.Exec(col1Data, col2Data); assert.NoError(t, err) {
						if assert.NoError(t, scope.Commit()) {
							var (
								col1 uuid.UUID
								col2 uuid.UUID
							)
							if err := conn.QueryRow("SELECT * FROM test_uuid").Scan(&col1, &col2); assert.NoError(t, err) {
								assert.Equal(t, col1Data, col1)
								assert.Equal(t, col2Data, col2)
							}
						}
					}
				}
			}
		}
	}
}

func TestStdNullableUUID(t *testing.T) {
	if conn, err := sql.Open("clickhouse", "clickhouse://127.0.0.1:9000"); assert.NoError(t, err) {
		const ddl = `
			CREATE TABLE test_uuid (
				  Col1 Nullable(UUID)
				, Col2 Nullable(UUID)
			) Engine Memory
		`
		if _, err := conn.Exec("DROP TABLE IF EXISTS test_uuid"); assert.NoError(t, err) {
			if _, err := conn.Exec(ddl); assert.NoError(t, err) {
				scope, err := conn.Begin()
				if assert.NoError(t, err) {
					return
				}
				if batch, err := conn.Prepare("INSERT INTO test_uuid"); assert.NoError(t, err) {
					var (
						col1Data = uuid.New()
						col2Data = uuid.New()
					)
					if _, err := batch.Exec(col1Data, col2Data); assert.NoError(t, err) {
						if assert.NoError(t, scope.Commit()) {
							var (
								col1 *uuid.UUID
								col2 *uuid.UUID
							)
							if err := conn.QueryRow("SELECT * FROM test_uuid").Scan(&col1, &col2); assert.NoError(t, err) {
								assert.Equal(t, col1Data, *col1)
								assert.Equal(t, col2Data, *col2)
							}
						}
					}
				}
			}
		}
		if _, err := conn.Exec("TRUNCATE TABLE test_uuid"); !assert.NoError(t, err) {
			return
		}
		scope, err := conn.Begin()
		if assert.NoError(t, err) {
			return
		}
		if batch, err := scope.Prepare("INSERT INTO test_uuid"); assert.NoError(t, err) {
			var col1Data = uuid.New()
			if _, err := batch.Exec(col1Data, nil); assert.NoError(t, err) {
				if assert.NoError(t, scope.Commit()) {
					var (
						col1 *uuid.UUID
						col2 *uuid.UUID
					)
					if err := conn.QueryRow("SELECT * FROM test_uuid").Scan(&col1, &col2); assert.NoError(t, err) {
						if assert.Nil(t, col2) {
							assert.Equal(t, col1Data, *col1)
						}
					}
				}
			}
		}
	}
}
