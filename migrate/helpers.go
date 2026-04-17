package migrate

import (
	"strconv"
	"strings"
)

// parseFilename extracts version, name, and direction from a migration filename.
// Expected format: NNNNNN_name.up.sql or NNNNNN_name.down.sql
// Example: "20260417103000_add_users_table.up.sql"
//
//	-> version=20260417103000, name="add_users_table", direction="up"
func parseFilename(name string) (version uint, migName, direction string, ok bool) {
	// Must end with .sql
	if !strings.HasSuffix(name, ".sql") {
		return 0, "", "", false
	}
	base := strings.TrimSuffix(name, ".sql")

	// Must end with .up or .down
	var dir string
	switch {
	case strings.HasSuffix(base, ".up"):
		dir = "up"
		base = strings.TrimSuffix(base, ".up")
	case strings.HasSuffix(base, ".down"):
		dir = "down"
		base = strings.TrimSuffix(base, ".down")
	default:
		return 0, "", "", false
	}

	// Split version and name: "20260417103000_add_users_table"
	parts := strings.SplitN(base, "_", 2)
	if len(parts) != 2 {
		return 0, "", "", false
	}

	v, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return 0, "", "", false
	}

	return uint(v), parts[1], dir, true
}

// sanitizeName converts a user-provided name into a safe filename component.
// Converts to lowercase and replaces spaces/special chars with underscores.
func sanitizeName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	// Replace spaces and common separators with underscores
	replacer := strings.NewReplacer(
		" ", "_",
		"-", "_",
		".", "_",
		"/", "_",
		"\\", "_",
	)
	name = replacer.Replace(name)

	// Keep only alphanumeric and underscores
	var b strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
