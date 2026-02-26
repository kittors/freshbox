package checker

import (
	"testing"
)

func TestDevToolsReturnsExpectedItems(t *testing.T) {
	items := DevTools()
	if len(items) != 12 {
		t.Errorf("expected 12 dev tools, got %d", len(items))
	}

	expected := map[string]string{
		"Homebrew":      "brew",
		"Git":           "git",
		"Java (JDK)":    "java",
		"Maven":         "mvn",
		"Gradle":        "gradle",
		"Python":        "python3",
		"uv":            "uv",
		"fnm":           "fnm",
		"pnpm":          "pnpm",
		"Bun":           "bun",
		"Rust (rustup)": "rustup",
		"Go":            "go",
	}

	for _, item := range items {
		cmd, ok := expected[item.Name]
		if !ok {
			t.Errorf("unexpected dev tool: %s", item.Name)
			continue
		}
		if item.Cmd != cmd {
			t.Errorf("tool %s: cmd = %q, want %q", item.Name, item.Cmd, cmd)
		}
		if item.Category != "dev" {
			t.Errorf("tool %s: category = %q, want %q", item.Name, item.Category, "dev")
		}
		if item.Desc == "" {
			t.Errorf("tool %s: missing description", item.Name)
		}
	}
}

func TestAppsReturnsExpectedItems(t *testing.T) {
	items := Apps()
	if len(items) != 7 {
		t.Errorf("expected 7 apps, got %d", len(items))
	}

	for _, item := range items {
		if item.Category != "app" {
			t.Errorf("app %s: category = %q, want %q", item.Name, item.Category, "app")
		}
		if item.BrewName == "" {
			t.Errorf("app %s: missing BrewName", item.Name)
		}
		if item.Desc == "" {
			t.Errorf("app %s: missing description", item.Name)
		}
	}
}

func TestAIToolsReturnsExpectedItems(t *testing.T) {
	items := AITools()
	if len(items) != 2 {
		t.Errorf("expected 2 AI tools, got %d", len(items))
	}
	names := map[string]bool{}
	for _, item := range items {
		names[item.Name] = true
		if item.Category != "ai" {
			t.Errorf("AI tool %s: category = %q, want %q", item.Name, item.Category, "ai")
		}
	}
	if !names["Codex"] {
		t.Error("missing Codex in AI tools")
	}
	if !names["Claude Code"] {
		t.Error("missing Claude Code in AI tools")
	}
}

func TestCheckSetsStatusCorrectly(t *testing.T) {
	// Use a known-to-exist command
	item := &Item{
		Name:    "TestBash",
		Cmd:     "bash",
		VerFlag: "--version",
	}
	Check(item)
	if item.Status != Installed {
		t.Error("bash should be detected as installed")
	}
	if item.Version == "" {
		t.Error("bash should have a version string")
	}

	// Use a command that doesn't exist
	missing := &Item{
		Name:    "TestMissing",
		Cmd:     "nonexistent-command-freshbox-test",
		VerFlag: "--version",
	}
	Check(missing)
	if missing.Status != NotInstalled {
		t.Error("nonexistent command should be NotInstalled")
	}
	if missing.Version != "" {
		t.Error("nonexistent command should have empty version")
	}
}

func TestCheckAllProcessesAllItems(t *testing.T) {
	items := []*Item{
		{Name: "exists", Cmd: "echo", VerFlag: ""},
		{Name: "missing", Cmd: "freshbox-fake-cmd-does-not-exist", VerFlag: "--version"},
	}
	CheckAll(items)

	if items[0].Status != Installed {
		t.Error("echo should be detected as installed")
	}
	if items[1].Status != NotInstalled {
		t.Error("fake cmd should be not installed")
	}
}

func TestCheckWithNoVerFlag(t *testing.T) {
	item := &Item{
		Name: "TestNoVer",
		Cmd:  "echo",
	}
	Check(item)
	if item.Status != Installed {
		t.Error("echo should be installed even without VerFlag")
	}
	if item.Version != "" {
		t.Error("version should be empty when VerFlag is empty")
	}
}

func TestResolveCmdFindsPathBinaries(t *testing.T) {
	// bash should always be findable
	path := resolveCmd("bash")
	if path == "" {
		t.Error("resolveCmd should find bash")
	}

	// A nonexistent command should return ""
	path = resolveCmd("freshbox-nonexistent-binary-xyz")
	if path != "" {
		t.Errorf("resolveCmd should return empty for nonexistent cmd, got %q", path)
	}
}

func TestCheckAppForNonCask(t *testing.T) {
	item := &Item{
		Name:     "TestCLI",
		Cmd:      "echo",
		IsCask:   false,
		BrewName: "echo",
	}
	CheckApp(item)
	// Should fall through to Check() since IsCask is false
	if item.Status != Installed {
		t.Error("non-cask item with valid cmd should be installed")
	}
}

func TestStatusConstants(t *testing.T) {
	if NotInstalled != 0 {
		t.Errorf("NotInstalled should be 0, got %d", NotInstalled)
	}
	if Installed != 1 {
		t.Errorf("Installed should be 1, got %d", Installed)
	}
}
