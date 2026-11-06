package config

import (
	"os"
	"testing"
)

func TestSetUser(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		wantErr  bool
		wantUser string
	}{
		{name: "valid name", input: "Marcel", wantErr: false, wantUser: "Marcel"},
		{name: "blank name", input: "", wantErr: true},
		{name: "whitespace only", input: "   ", wantErr: true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Redirect the home dir to a throwaway folder so we don't write the actual config file
			home := t.TempDir()
			t.Setenv("HOME", home)        // Linux / macOS
			t.Setenv("USERPROFILE", home) // Windows

			cfg := &Config{}
			err := cfg.SetUser(c.input)

			if c.wantErr {
				if err == nil {
					t.Fatalf("SetUser(%q): expected an error, got nil", c.input)
				}
				if cfg.CurrentUserName != "" {
					t.Errorf("SetUser(%q): CurrentUserName changed to %q on error",
						c.input, cfg.CurrentUserName)
				}
				entries, err := os.ReadDir(home)
				if err != nil {
					t.Fatalf("reading temp dir: %v", err)
				}
				if len(entries) != 0 {
					t.Errorf("SetUser(%q): expected no file written, found %d", c.input, len(entries))
				}
				return
			}

			if err != nil {
				t.Fatalf("SetUser(%q): unexpected error: %v", c.input, err)
			}
			if cfg.CurrentUserName != c.wantUser {
				t.Errorf("CurrentUserName = %q, want %q", cfg.CurrentUserName, c.wantUser)
			}

			// check if persisted not just saved in memory
			saved, err := Read()
			if err != nil {
				t.Fatalf("Read() after SetUser: %v", err)
			}
			if saved.CurrentUserName != c.wantUser {
				t.Errorf("saved CurrentUserName = %q, want %q", saved.CurrentUserName, c.wantUser)
			}
		})
	}
}
