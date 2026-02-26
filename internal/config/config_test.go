package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// --- AvailableMCPs ---

func TestAvailableMCPs(t *testing.T) {
	mcps := AvailableMCPs()
	if len(mcps) != 11 {
		t.Errorf("expected 11 MCP servers, got %d", len(mcps))
	}

	names := make(map[string]bool)
	for _, mcp := range mcps {
		names[mcp.Name] = true

		if mcp.Command == "" {
			t.Errorf("MCP %s has empty command", mcp.Name)
		}
		if len(mcp.Args) == 0 {
			t.Errorf("MCP %s has no args", mcp.Name)
		}
	}

	// Verify key servers exist
	required := []string{"Playwright", "Context7", "Filesystem", "GitHub", "Memory"}
	for _, name := range required {
		if !names[name] {
			t.Errorf("missing required MCP server: %s", name)
		}
	}

	// Verify all npm packages use @latest
	for _, mcp := range mcps {
		for _, arg := range mcp.Args {
			if strings.Contains(arg, "@") && strings.Contains(arg, "/") {
				if !strings.HasSuffix(arg, "@latest") {
					t.Errorf("MCP %s package %s missing @latest suffix", mcp.Name, arg)
				}
			}
		}
	}

	// Verify Playwright has --headless
	for _, mcp := range mcps {
		if mcp.Name == "Playwright" {
			found := false
			for _, arg := range mcp.Args {
				if arg == "--headless" {
					found = true
				}
			}
			if !found {
				t.Error("Playwright MCP missing --headless flag")
			}
		}
	}

	// Verify Filesystem has home dir path
	home, _ := os.UserHomeDir()
	for _, mcp := range mcps {
		if mcp.Name == "Filesystem" {
			found := false
			for _, arg := range mcp.Args {
				if arg == home {
					found = true
				}
			}
			if !found {
				t.Errorf("Filesystem MCP missing home dir, args: %v", mcp.Args)
			}
		}
	}
}

// --- WriteCodexConfig ---

func TestWriteCodexConfig_NewFile(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	err := WriteCodexConfig(CodexConfig{
		Model:         "o4-mini",
		ThinkingLevel: "medium",
		BaseURL:       "",
	})
	if err != nil {
		t.Fatalf("WriteCodexConfig failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmp, ".codex", "config.toml"))
	if err != nil {
		t.Fatalf("read config.toml: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, `model = "o4-mini"`) {
		t.Errorf("config missing model, got:\n%s", content)
	}
	if !strings.Contains(content, `model_reasoning_effort = "medium"`) {
		t.Errorf("config missing thinking level, got:\n%s", content)
	}
}

func TestWriteCodexConfig_WithBaseURL(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	err := WriteCodexConfig(CodexConfig{
		Model:         "gpt-4",
		ThinkingLevel: "high",
		BaseURL:       "https://custom.api.com/v1",
	})
	if err != nil {
		t.Fatalf("WriteCodexConfig failed: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(tmp, ".codex", "config.toml"))
	content := string(data)
	if !strings.Contains(content, `model_provider = "freshbox"`) {
		t.Errorf("config missing model_provider, got:\n%s", content)
	}
	if !strings.Contains(content, `base_url = "https://custom.api.com/v1"`) {
		t.Errorf("config missing base_url, got:\n%s", content)
	}
}

func TestWriteCodexConfig_MergeExisting(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	dir := filepath.Join(tmp, ".codex")
	os.MkdirAll(dir, 0755)

	existing := `model = "old-model"
model_reasoning_effort = "low"

[some_section]
key = "value"
`
	os.WriteFile(filepath.Join(dir, "config.toml"), []byte(existing), 0644)

	err := WriteCodexConfig(CodexConfig{
		Model:         "new-model",
		ThinkingLevel: "high",
	})
	if err != nil {
		t.Fatalf("WriteCodexConfig failed: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(dir, "config.toml"))
	content := string(data)

	if !strings.Contains(content, `model = "new-model"`) {
		t.Errorf("model not updated, got:\n%s", content)
	}
	if !strings.Contains(content, `model_reasoning_effort = "high"`) {
		t.Errorf("thinking level not updated, got:\n%s", content)
	}
	if !strings.Contains(content, `[some_section]`) {
		t.Errorf("existing section was lost, got:\n%s", content)
	}
}

// --- WriteCodexAuth ---

func TestWriteCodexAuth(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	err := WriteCodexAuth(CodexAuth{APIKey: "sk-test-key-123"})
	if err != nil {
		t.Fatalf("WriteCodexAuth failed: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(tmp, ".codex", "auth.json"))
	var result map[string]string
	json.Unmarshal(data, &result)

	if result["api_key"] != "sk-test-key-123" {
		t.Errorf("api_key = %q, want %q", result["api_key"], "sk-test-key-123")
	}

	// Verify file permissions
	info, _ := os.Stat(filepath.Join(tmp, ".codex", "auth.json"))
	perm := info.Mode().Perm()
	if perm != 0600 {
		t.Errorf("auth.json permissions = %o, want 0600", perm)
	}
}

// --- WriteClaudeConfig ---

func TestWriteClaudeConfig_NewFile(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	err := WriteClaudeConfig(ClaudeConfig{
		Model:   "claude-sonnet-4-6",
		BaseURL: "https://api.anthropic.com",
		APIKey:  "sk-ant-test",
	})
	if err != nil {
		t.Fatalf("WriteClaudeConfig failed: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(tmp, ".claude", "settings.json"))
	var result map[string]any
	json.Unmarshal(data, &result)

	if result["model"] != "claude-sonnet-4-6" {
		t.Errorf("model = %v, want claude-sonnet-4-6", result["model"])
	}

	env, ok := result["env"].(map[string]any)
	if !ok {
		t.Fatal("env not found or not a map")
	}
	if env["ANTHROPIC_API_KEY"] != "sk-ant-test" {
		t.Errorf("ANTHROPIC_API_KEY = %v", env["ANTHROPIC_API_KEY"])
	}
	if env["ANTHROPIC_BASE_URL"] != "https://api.anthropic.com" {
		t.Errorf("ANTHROPIC_BASE_URL = %v", env["ANTHROPIC_BASE_URL"])
	}
}

func TestWriteClaudeConfig_MergeExisting(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	dir := filepath.Join(tmp, ".claude")
	os.MkdirAll(dir, 0755)

	existing := map[string]any{
		"existingKey": "should-survive",
		"env":         map[string]any{"EXISTING_VAR": "keep-me"},
	}
	data, _ := json.MarshalIndent(existing, "", "  ")
	os.WriteFile(filepath.Join(dir, "settings.json"), data, 0600)

	err := WriteClaudeConfig(ClaudeConfig{
		Model:  "new-model",
		APIKey: "new-key",
	})
	if err != nil {
		t.Fatalf("WriteClaudeConfig failed: %v", err)
	}

	resultData, _ := os.ReadFile(filepath.Join(dir, "settings.json"))
	var result map[string]any
	json.Unmarshal(resultData, &result)

	if result["existingKey"] != "should-survive" {
		t.Error("existing key was lost")
	}
	env := result["env"].(map[string]any)
	if env["EXISTING_VAR"] != "keep-me" {
		t.Error("existing env var was lost")
	}
	if env["ANTHROPIC_API_KEY"] != "new-key" {
		t.Errorf("new API key not set: %v", env["ANTHROPIC_API_KEY"])
	}
}

func TestWriteClaudeConfig_EmptyFields(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	// Only set model, leave base URL and API key empty
	err := WriteClaudeConfig(ClaudeConfig{Model: "claude-sonnet-4-6"})
	if err != nil {
		t.Fatalf("WriteClaudeConfig failed: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(tmp, ".claude", "settings.json"))
	var result map[string]any
	json.Unmarshal(data, &result)

	if result["model"] != "claude-sonnet-4-6" {
		t.Errorf("model = %v", result["model"])
	}
	// env should not exist since no key/url was set
	if _, ok := result["env"]; ok {
		t.Error("env should not be present when no key/url set")
	}
}

// --- PreDownloadMCPPackages ---

func TestPreDownloadMCPPackages_SkipsNonNpx(t *testing.T) {
	servers := []MCPServer{
		{Name: "custom", Command: "node", Args: []string{"server.js"}},
	}
	// should be a no-op, no error
	err := PreDownloadMCPPackages(servers)
	if err != nil {
		t.Errorf("expected no error for non-npx server, got: %v", err)
	}
}

func TestPreDownloadMCPPackages_SkipsNoYFlag(t *testing.T) {
	servers := []MCPServer{
		{Name: "no-y", Command: "npx", Args: []string{"some-pkg"}},
	}
	err := PreDownloadMCPPackages(servers)
	if err != nil {
		t.Errorf("expected no error for server without -y flag, got: %v", err)
	}
}

// --- WriteMCPConfig ---

func TestWriteMCPConfig_InvalidTarget(t *testing.T) {
	err := WriteMCPConfig([]MCPServer{{Name: "test"}}, "invalid-target")
	if err == nil {
		t.Error("expected error for invalid MCP target")
	}
	if !strings.Contains(err.Error(), "unknown MCP target") {
		t.Errorf("error should mention unknown target, got: %v", err)
	}
}

// --- addCodexMCPTimeout ---

func TestAddCodexMCPTimeout(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	dir := filepath.Join(tmp, ".codex")
	os.MkdirAll(dir, 0755)

	config := `model = "o4-mini"

[mcp_servers.test-server]
command = "npx"
args = ["-y", "@test/pkg"]
`
	os.WriteFile(filepath.Join(dir, "config.toml"), []byte(config), 0644)

	servers := []MCPServer{{Name: "test-server"}}
	err := addCodexMCPTimeout(servers, 60)
	if err != nil {
		t.Fatalf("addCodexMCPTimeout failed: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(dir, "config.toml"))
	if !strings.Contains(string(data), "startup_timeout_sec = 60") {
		t.Errorf("timeout not added, config:\n%s", string(data))
	}
}

func TestAddCodexMCPTimeout_AlreadyExists(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	dir := filepath.Join(tmp, ".codex")
	os.MkdirAll(dir, 0755)

	config := `[mcp_servers.test]
startup_timeout_sec = 30
command = "npx"
`
	os.WriteFile(filepath.Join(dir, "config.toml"), []byte(config), 0644)

	err := addCodexMCPTimeout([]MCPServer{{Name: "test"}}, 60)
	if err != nil {
		t.Fatalf("addCodexMCPTimeout failed: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(dir, "config.toml"))
	content := string(data)
	// Should NOT add a second timeout
	count := strings.Count(content, "startup_timeout_sec")
	if count > 1 {
		t.Errorf("duplicate timeout entries, count=%d, config:\n%s", count, content)
	}
}

// --- MCPServer struct ---

func TestMCPServerJSON(t *testing.T) {
	server := MCPServer{
		Name:    "test",
		Command: "npx",
		Args:    []string{"-y", "@test/pkg@latest"},
	}
	data, err := json.Marshal(server)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var result map[string]any
	json.Unmarshal(data, &result)

	if result["name"] != "test" {
		t.Errorf("name = %v", result["name"])
	}
	if result["command"] != "npx" {
		t.Errorf("command = %v", result["command"])
	}
}
