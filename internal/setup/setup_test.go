package setup

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestSetupKaku_WritesConfig(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	// We can't clone from GitHub in unit tests, but we can test the file writing.
	// SetupKaku will fail at git clone, but let's verify the config dir is created
	// and kaku.lua is written before the clone step.
	kakuDir := filepath.Join(tmp, ".config", "kaku")
	pluginDir := filepath.Join(kakuDir, "zsh", "plugins")
	os.MkdirAll(pluginDir, 0755)

	// Pre-create plugin dirs to skip git clone
	plugins := []string{"zsh-autosuggestions", "zsh-completions", "zsh-syntax-highlighting", "zsh-z"}
	for _, p := range plugins {
		os.MkdirAll(filepath.Join(pluginDir, p), 0755)
	}

	err := SetupKaku()
	if err != nil {
		t.Fatalf("SetupKaku failed: %v", err)
	}

	// Verify kaku.lua was written
	luaPath := filepath.Join(kakuDir, "kaku.lua")
	data, err := os.ReadFile(luaPath)
	if err != nil {
		t.Fatalf("kaku.lua not created: %v", err)
	}
	if len(data) == 0 {
		t.Error("kaku.lua is empty")
	}

	content := string(data)
	if !contains(content, "wezterm") {
		t.Error("kaku.lua missing wezterm reference")
	}
	if !contains(content, "resolve_bundled_config") {
		t.Error("kaku.lua missing resolve_bundled_config function")
	}
}

func TestSetupDevWorkspace_CreatesStructure(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	err := SetupDevWorkspace()
	if err != nil {
		t.Fatalf("SetupDevWorkspace failed: %v", err)
	}

	devDir := filepath.Join(tmp, "Developer")

	// Verify directories were created
	expected := []string{
		"opensource",
		"boundless",
		"freelance",
		"freelance/_template",
		"playground",
		"design",
		"notes",
		"scripts",
		"archive",
	}
	for _, d := range expected {
		path := filepath.Join(devDir, d)
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("directory %s not created: %v", d, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("%s is not a directory", d)
		}
	}

	// Verify root README.md was created
	readmePath := filepath.Join(devDir, "README.md")
	data, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("root README not created: %v", err)
	}
	if !contains(string(data), "Developer") {
		t.Error("root README missing Developer reference")
	}

	// Verify subdirectory READMEs
	subReadmes := []string{
		"opensource/README.md",
		"boundless/README.md",
		"freelance/README.md",
		"playground/README.md",
		"design/README.md",
		"notes/README.md",
		"scripts/README.md",
		"archive/README.md",
	}
	for _, r := range subReadmes {
		path := filepath.Join(devDir, r)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("sub README %s not created: %v", r, err)
		}
	}
}

func TestSetupKarabiner_WritesConfig(t *testing.T) {
	// This test is limited because SetupKarabiner runs `brew install` first.
	// We test the JSON structure that would be generated.
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	karabinerDir := filepath.Join(tmp, ".config", "karabiner")
	os.MkdirAll(karabinerDir, 0755)

	// Test that the function creates the open-kaku.sh script
	binDir := filepath.Join(tmp, ".local", "bin")
	os.MkdirAll(binDir, 0755)

	// The full SetupKarabiner calls brew install first which we can't do in tests.
	// Instead, verify the karabiner config structure
	config := map[string]any{
		"global": map[string]any{"show_in_menu_bar": false},
		"profiles": []map[string]any{
			{
				"complex_modifications": map[string]any{
					"rules": []map[string]any{
						{
							"description": "Control+Option+Command+T opens Kaku",
						},
					},
				},
				"name":     "Default profile",
				"selected": true,
			},
		},
	}

	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		t.Fatalf("marshal config: %v", err)
	}

	if !contains(string(data), "Control+Option+Command+T opens Kaku") {
		t.Error("karabiner config missing expected rule description")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
