// Package migrations carries the SQL migrations that create and evolve the
// data stores used by the Epileptic application component.
package migrations

import "embed"

// FS holds the embedded SQL migration files.
//
//go:embed *.sql
var FS embed.FS
