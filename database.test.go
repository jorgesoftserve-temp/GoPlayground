package main

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresWithDocker(t *testing.T) {
	pool, err := dockertest.NewPool("")
	require.NoError(t, err)

	resource, err := pool.Run(
		"postgres",
		"16",
		[]string{
			"POSTGRES_USER=test",
			"POSTGRES_PASSWORD=test",
			"POSTGRES_DB=veterinaria",
		},
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := pool.Purge(resource)
		require.NoError(t, err)
	})

	var db *sql.DB

	err = pool.Retry(func() error {
		connectionString := fmt.Sprintf(
			"postgres://test:test@localhost:%s/veterinaria?sslmode=disable",
			resource.GetPort("5432/tcp"),
		)

		db, err = sql.Open(
			"postgres",
			connectionString,
		)

		if err != nil {
			return err
		}

		return db.Ping()
	})

	require.NoError(t, err)

	t.Cleanup(func() {
		db.Close()
	})

	_, err = db.Exec(`
		CREATE TABLE pets (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			owner_name TEXT NOT NULL
		)
	`)

	require.NoError(t, err)

	_, err = db.Exec(
		`
			INSERT INTO pets (
				name,
				owner_name
			)
			VALUES ($1, $2)
		`,
		"Max",
		"Jorge",
	)

	require.NoError(t, err)

	var pet Pet

	err = db.QueryRow(
		`
			SELECT
				id,
				name,
				owner_name
			FROM pets
			LIMIT 1
		`,
	).Scan(
		&pet.ID,
		&pet.Name,
		&pet.OwnerName,
	)

	require.NoError(t, err)

	assert.Equal(
		t,
		"Max",
		pet.Name,
	)

	assert.Equal(
		t,
		"Jorge",
		pet.OwnerName,
	)
}
