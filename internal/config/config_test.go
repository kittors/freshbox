package config

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteClaudeMCP(t *testing.T) {
	servers := []MCPServer{
		{Name: "test-freshbox-pw", Command: "npx", Args: []string{"-y", "@playwright/mcp@latest", "--headless"}},
	}

	err := WriteClaudeMCP(servers)
	if err != nil {
		t.Fatalf("WriteClaudeMCP failed: %v", err)
	}

	// Verify via claude mcp get
	out, err := exec.Command("claude", "mcp", "get", "test-freshbox-pw").CombinedOutput()
	if err != nil {
		t.Fatalf("claude mcp get failed: %v (%s)", err, string(out))
	}
	outStr := string(out)
	if !strings.Contains(outStr, "test-freshbox-pw") {
		t.Errorf("MCP server not found in claude mcp get output: %s", outStr)
	}
	if !strings.Contains(outStr, "--headless") {
		t.Error("--headless flag not found in claude mcp config")
	}
	t.Logf("claude mcp get output:\n%s", outStr)

	// Clean up
	exec.Command("claude", "mcp", "remove", "-s", "user", "test-freshbox-pw").Run()
}

func TestWriteCodexMCP(t *testing.T) {
	servers := []MCPServer{
		{Name: "test-freshbox-ctx", Command: "npx", Args: []string{"-y", "@upstash/context7-mcp@latest"}},
	}

	err := WriteCodexMCP(servers)
	if err != nil {
		t.Fatalf("WriteCodexMCP failed: %v", err)
	}

	// Verify via codex mcp list
	out, _ := exec.Command("codex", "mcp", "list").CombinedOutput()
	if !strings.Contains(string(out), "test-freshbox-ctx") {
		t.Errorf("MCP server not found in codex mcp list: %s", string(out))
	}

	// Verify startup_timeout_sec = 60 was added
	home, _ := os.UserHomeDir()
	configData, _ := os.ReadFile(filepath.Join(home, ".codex", "config.toml"))
	configStr := string(configData)
	if !strings.Contains(configStr, "startup_timeout_sec = 60") {
		t.Error("startup_timeout_sec = 60 not found in config.toml")
		t.Logf("config.toml:\n%s", configStr)
	} else {
		t.Log("startup_timeout_sec = 60 found in config.toml")
	}

	// Clean up
	exec.Command("codex", "mcp", "remove", "test-freshbox-ctx").Run()
}

func TestAvailableMCPs(t *testing.T) {
	mcps := AvailableMCPs()
	if len(mcps) != 11 {
		t.Errorf("expected 11 MCP servers, got %d", len(mcps))
	}

	// Verify all packages use @latest
	for _, mcp := range mcps {
		for _, arg := range mcp.Args {
			if strings.HasPrefix(arg, "@") && !strings.HasPrefix(arg, "-") && arg != "-y" {
				if !strings.HasSuffix(arg, "@latest") {
					t.Errorf("MCP %s package %s missing @latest suffix", mcp.Name, arg)
				}
			}
		}
	}

	// Verify Playwright has --headless
	for _, mcp := range mcps {
		if mcp.Name == "Playwright" {
			hasHeadless := false
			for _, arg := range mcp.Args {
				if arg == "--headless" {
					hasHeadless = true
				}
			}
			if !hasHeadless {
				t.Error("Playwright MCP missing --headless flag")
			}
		}
	}

	// Verify Filesystem has home dir path
	home, _ := os.UserHomeDir()
	for _, mcp := range mcps {
		if mcp.Name == "Filesystem" {
			hasHome := false
			for _, arg := range mcp.Args {
				if arg == home {
					hasHome = true
				}
			}
			if !hasHome {
				t.Errorf("Filesystem MCP missing home dir path, args: %v", mcp.Args)
			}
		}
	}

	// Verify Puppeteer is removed
	for _, mcp := range mcps {
		if mcp.Name == "Puppeteer" {
			t.Error("Puppeteer should be removed from default MCP list")
		}
	}
}

func TestWriteClaudeConfig(t *testing.T) {
	home, _ := os.UserHomeDir()
	settingsPath := filepath.Join(home, ".claude", "settings.json")

	// Backup original
	original, _ := os.ReadFile(settingsPath)
	defer os.WriteFile(settingsPath, original, 0600)

	err := WriteClaudeConfig(ClaudeConfig{
		Model:   "claude-sonnet-4-6",
		BaseURL: "https://api.anthropic.com",
		APIKey:  "sk-ant-test123",
	})
	if err != nil {
		t.Fatalf("WriteClaudeConfig failed: %v", err)
	}

	data, _ := os.ReadFile(settingsPath)
	var result map[string]any
	json.Unmarshal(data, &result)

	if result["model"] != "claude-sonnet-4-6" {
		t.Errorf("model = %v, want claude-sonnet-4-6", result["model"])
	}

	env, ok := result["env"].(map[string]any)
	if !ok {
		t.Fatal("env not found or not a map")
	}
	if env["ANTHROPIC_API_KEY"] != "sk-ant-test123" {
		t.Errorf("ANTHROPIC_API_KEY = %v, want sk-ant-test123", env["ANTHROPIC_API_KEY"])
	}
	if env["ANTHROPIC_BASE_URL"] != "https://api.anthropic.com" {
		t.Errorf("ANTHROPIC_BASE_URL = %v", env["ANTHROPIC_BASE_URL"])
	}

	// Verify original keys preserved
	if _, ok := result["skipDangerousModePermissionPrompt"]; !ok {
		t.Error("skipDangerousModePermissionPrompt was lost")
	}

	// Verify no mcpServers leaked into settings.json
	if _, ok := result["mcpServers"]; ok {
		t.Error("mcpServers should NOT be in settings.json")
	}

	t.Logf("settings.json after config write:\n%s", string(data))
}

func TestPreDownloadMCPPackages(t *testing.T) {
	servers := []MCPServer{
		{Name: "Context7", Command: "npx", Args: []string{"-y", "@upstash/context7-mcp@latest"}},
	}
	err := PreDownloadMCPPackages(servers)
	if err != nil {
		t.Logf("PreDownload warning (non-fatal): %v", err)
	} else {
		t.Log("PreDownload succeeded")
	}
}
