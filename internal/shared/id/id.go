// Package id is the single place primary keys get generated.
package id

import "github.com/google/uuid"

func New() uuid.UUID {
	return uuid.Must(uuid.NewV7())
}
