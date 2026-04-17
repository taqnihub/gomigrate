package migrate

import (
	"testing"
)

func TestParseFilename(t *testing.T) {
	tests := []struct {
		filename      string
		wantVersion   uint
		wantName      string
		wantDirection string
		wantOK        bool
	}{
		{
			filename:      "000001_add_users_table.up.sql",
			wantVersion:   1,
			wantName:      "add_users_table",
			wantDirection: "up",
			wantOK:        true,
		},
		{
			filename:      "000002_add_products.down.sql",
			wantVersion:   2,
			wantName:      "add_products",
			wantDirection: "down",
			wantOK:        true,
		},
		{
			filename:      "20260417103000_timestamp_version.up.sql",
			wantVersion:   20260417103000,
			wantName:      "timestamp_version",
			wantDirection: "up",
			wantOK:        true,
		},
		{
			filename: "README.md",
			wantOK:   false,
		},
		{
			filename: "not_a_migration.txt",
			wantOK:   false,
		},
		{
			filename: "missing_direction.sql",
			wantOK:   false,
		},
		{
			filename: "abc_invalid_version.up.sql",
			wantOK:   false,
		},
		{
			filename: "001.up.sql", // no name part
			wantOK:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.filename, func(t *testing.T) {
			version, name, direction, ok := parseFilename(tc.filename)

			if ok != tc.wantOK {
				t.Errorf("ok: got %v, want %v", ok, tc.wantOK)
				return
			}

			if !tc.wantOK {
				return // don't check other fields when we expect failure
			}

			if version != tc.wantVersion {
				t.Errorf("version: got %d, want %d", version, tc.wantVersion)
			}
			if name != tc.wantName {
				t.Errorf("name: got %q, want %q", name, tc.wantName)
			}
			if direction != tc.wantDirection {
				t.Errorf("direction: got %q, want %q", direction, tc.wantDirection)
			}
		})
	}
}

func TestSanitizeName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"add_users_table", "add_users_table"},
		{"Add Users Table", "add_users_table"},
		{"add-users-table", "add_users_table"},
		{"Add Users Table!", "add_users_table"},
		{"  trim  spaces  ", "trim__spaces"},
		{"MixED CaSe", "mixed_case"},
		{"user.email.index", "user_email_index"},
		{"path/with/slashes", "path_with_slashes"},
		{"unicode_é", "unicode_"}, // non-ASCII stripped
		{"", ""},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := sanitizeName(tc.input)
			if got != tc.want {
				t.Errorf("sanitizeName(%q):\n  got:  %q\n  want: %q", tc.input, got, tc.want)
			}
		})
	}
}
