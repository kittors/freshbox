package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CodexConfig represents Codex CLI configuration
type CodexConfig struct {
	Model         string
	ThinkingLevel string
	BaseURL       string
}

// CodexAuth represents Codex auth configuration
type CodexAuth struct {
	APIKey string `json:"api_key"`
}

// ClaudeConfig represents Claude Code configuration
type ClaudeConfig struct {
	Model   string
	BaseURL string
	APIKey  string
}

// MCPServer represents an MCP server configuration
type MCPServer struct {
	Name    string   `json:"name"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

// AvailableMCPs returns all popular MCP servers with correct npm package names
func AvailableMCPs() []MCPServer {
	home, _ := os.UserHomeDir()
	return []MCPServer{
		{Name: "Playwright", Command: "npx", Args: []string{"-y", "@playwright/mcp@latest", "--headless"}},
		{Name: "Context7", Command: "npx", Args: []string{"-y", "@upstash/context7-mcp@latest"}},
		{Name: "Filesystem", Command: "npx", Args: []string{"-y", "@modelcontextprotocol/server-filesystem@latest", home}},
		{Name: "GitHub", Command: "npx", Args: []string{"-y", "@modelcontextprotocol/server-github@latest"}},
		{Name: "Memory", Command: "npx", Args: []string{"-y", "@modelcontextprotocol/server-memory@latest"}},
		{Name: "Sequential Thinking", Command: "npx", Args: []string{"-y", "@modelcontextprotocol/server-sequential-thinking@latest"}},
		{Name: "Fetch", Command: "npx", Args: []string{"-y", "@modelcontextprotocol/server-fetch@latest"}},
		{Name: "Brave Search", Command: "npx", Args: []string{"-y", "@modelcontextprotocol/server-brave-search@latest"}},
		{Name: "Slack", Command: "npx", Args: []string{"-y", "@modelcontextprotocol/server-slack@latest"}},
		{Name: "Google Maps", Command: "npx", Args: []string{"-y", "@modelcontextprotocol/server-google-maps@latest"}},
		{Name: "SQLite", Command: "npx", Args: []string{"-y", "@modelcontextprotocol/server-sqlite@latest"}},
	}
}

// WriteCodexConfig merges model/thinking/baseURL into existing ~/.codex/config.toml
func WriteCodexConfig(cfg CodexConfig) error {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".codex")
	os.MkdirAll(dir, 0755)
	configPath := filepath.Join(dir, "config.toml")

	// Read existing config
	existing, _ := os.ReadFile(configPath)
	lines := strings.Split(string(existing), "\n")

	// Build a map of top-level keys to update
	updates := map[string]string{}
	if cfg.Model != "" {
		updates["model"] = fmt.Sprintf(`model = "%s"`, cfg.Model)
	}
	if cfg.ThinkingLevel != "" {
		updates["model_reasoning_effort"] = fmt.Sprintf(`model_reasoning_effort = "%s"`, cfg.ThinkingLevel)
	}

	// Track which keys we've updated in-place
	updated := map[string]bool{}
	var result []string
	inSection := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Track if we're inside a [section]
		if strings.HasPrefix(trimmed, "[") {
			inSection = true
		} else if !inSection {
			// Only replace top-level keys (before any [section])
			for key, newLine := range updates {
				if strings.HasPrefix(trimmed, key+" ") || strings.HasPrefix(trimmed, key+"=") {
					line = newLine
					updated[key] = true
					break
				}
			}
		}

		result = append(result, line)
	}

	// If we need to handle base_url via model_providers section
	if cfg.BaseURL != "" {
		providerName := "freshbox"
		updates["model_provider"] = fmt.Sprintf(`model_provider = "%s"`, providerName)

		if !updated["model_provider"] {
			result = append([]string{updates["model_provider"]}, result...)
		} else {
			for i, line := range result {
				trimmed := strings.TrimSpace(line)
				if strings.HasPrefix(trimmed, "model_provider ") || strings.HasPrefix(trimmed, "model_provider=") {
					result[i] = updates["model_provider"]
					break
				}
			}
		}

		providerSection := fmt.Sprintf(`
[model_providers.%s]
name = "openai"
base_url = "%s"
wire_api = "responses"
requires_openai_auth = true`, providerName, cfg.BaseURL)

		var cleaned []string
		skipSection := false
		for _, line := range result {
			trimmed := strings.TrimSpace(line)
			if trimmed == fmt.Sprintf("[model_providers.%s]", providerName) {
				skipSection = true
				continue
			}
			if skipSection && strings.HasPrefix(trimmed, "[") {
				skipSection = false
			}
			if !skipSection {
				cleaned = append(cleaned, line)
			}
		}
		result = cleaned
		result = append(result, providerSection)
	}

	// Append any top-level keys that weren't found in existing file
	var prepend []string
	for key, newLine := range updates {
		if !updated[key] && key != "model_provider" {
			prepend = append(prepend, newLine)
		}
	}
	if len(prepend) > 0 {
		result = append(prepend, result...)
	}

	content := strings.Join(result, "\n")
	for strings.Contains(content, "\n\n\n") {
		content = strings.ReplaceAll(content, "\n\n\n", "\n\n")
	}

	return os.WriteFile(configPath, []byte(strings.TrimSpace(content)+"\n"), 0644)
}

// WriteCodexAuth writes auth.json for Codex
func WriteCodexAuth(auth CodexAuth) error {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".codex")
	os.MkdirAll(dir, 0755)

	data, _ := json.MarshalIndent(auth, "", "  ")
	return os.WriteFile(filepath.Join(dir, "auth.json"), data, 0600)
}

// WriteClaudeConfig merges settings into existing ~/.claude/settings.json
func WriteClaudeConfig(cfg ClaudeConfig) error {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".claude")
	os.MkdirAll(dir, 0755)
	settingsPath := filepath.Join(dir, "settings.json")

	// Read existing settings
	existing := map[string]any{}
	if data, err := os.ReadFile(settingsPath); err == nil {
		json.Unmarshal(data, &existing)
	}

	// Merge new values (only non-empty)
	if cfg.Model != "" {
		existing["model"] = cfg.Model
	}

	// Build env map, preserving existing env entries
	envMap := map[string]string{}
	if existingEnv, ok := existing["env"].(map[string]any); ok {
		for k, v := range existingEnv {
			if s, ok := v.(string); ok {
				envMap[k] = s
			}
		}
	}
	if cfg.APIKey != "" {
		envMap["ANTHROPIC_API_KEY"] = cfg.APIKey
	}
	if cfg.BaseURL != "" {
		envMap["ANTHROPIC_BASE_URL"] = cfg.BaseURL
	}
	if len(envMap) > 0 {
		existing["env"] = envMap
	}

	data, _ := json.MarshalIndent(existing, "", "  ")
	return os.WriteFile(settingsPath, data, 0600)
}

// WriteClaudeMCP adds MCP servers to Claude Code via `claude mcp add -s user`
func WriteClaudeMCP(servers []MCPServer) error {
	var errs []string
	for _, s := range servers {
		// Remove existing first (ignore errors if not found)
		exec.Command("claude", "mcp", "remove", "-s", "user", s.Name).Run()

		// claude mcp add -s user <name> -- <command> <args...>
		args := []string{"mcp", "add", "-s", "user", s.Name, "--"}
		args = append(args, s.Command)
		args = append(args, s.Args...)

		cmd := exec.Command("claude", args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			outStr := strings.TrimSpace(string(out))
			errs = append(errs, fmt.Sprintf("%s: %s (%s)", s.Name, err.Error(), outStr))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("failed to add MCP servers: %s", strings.Join(errs, "; "))
	}
	return nil
}

// WriteCodexMCP adds MCP servers to Codex via `codex mcp add` with startup_timeout_sec
func WriteCodexMCP(servers []MCPServer) error {
	var errs []string
	for _, s := range servers {
		// Remove existing first (ignore errors if not found)
		exec.Command("codex", "mcp", "remove", s.Name).Run()

		// codex mcp add <name> -- <command> <args...>
		args := []string{"mcp", "add", s.Name, "--"}
		args = append(args, s.Command)
		args = append(args, s.Args...)

		cmd := exec.Command("codex", args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			outStr := strings.TrimSpace(string(out))
			errs = append(errs, fmt.Sprintf("%s: %s (%s)", s.Name, err.Error(), outStr))
		}
	}

	// Add startup_timeout_sec to each MCP server in config.toml
	if err := addCodexMCPTimeout(servers, 60); err != nil {
		errs = append(errs, fmt.Sprintf("timeout config: %s", err.Error()))
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to add MCP servers: %s", strings.Join(errs, "; "))
	}
	return nil
}

// addCodexMCPTimeout adds startup_timeout_sec to each [mcp_servers.*] section in config.toml
func addCodexMCPTimeout(servers []MCPServer, timeout int) error {
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".codex", "config.toml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	var result []string
	timeoutLine := fmt.Sprintf("startup_timeout_sec = %d", timeout)

	for i := 0; i < len(lines); i++ {
		result = append(result, lines[i])
		trimmed := strings.TrimSpace(lines[i])

		// Check if this is an [mcp_servers.*] section header
		if strings.HasPrefix(trimmed, "[mcp_servers.") && strings.HasSuffix(trimmed, "]") {
			// Check if startup_timeout_sec already exists in this section
			hasTimeout := false
			for j := i + 1; j < len(lines); j++ {
				nextTrimmed := strings.TrimSpace(lines[j])
				if strings.HasPrefix(nextTrimmed, "[") {
					break // next section
				}
				if strings.HasPrefix(nextTrimmed, "startup_timeout_sec") {
					hasTimeout = true
					break
				}
			}
			if !hasTimeout {
				result = append(result, timeoutLine)
			}
		}
	}

	content := strings.Join(result, "\n")
	return os.WriteFile(configPath, []byte(content), 0644)
}

// PreDownloadMCPPackages pre-downloads all MCP npm packages so they're cached
// and ready when the MCP client tries to connect (avoids startup timeouts)
func PreDownloadMCPPackages(servers []MCPServer) error {
	var errs []string
	for _, s := range servers {
		if s.Command != "npx" {
			continue
		}
		// Find the package name (the arg after "-y")
		var pkg string
		for i, arg := range s.Args {
			if arg == "-y" && i+1 < len(s.Args) {
				pkg = s.Args[i+1]
				break
			}
		}
		if pkg == "" {
			continue
		}

		// Use npm cache add to pre-download without executing
		cmd := exec.Command("npm", "cache", "add", pkg)
		out, err := cmd.CombinedOutput()
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s: %s", pkg, strings.TrimSpace(string(out))))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("pre-download warnings: %s", strings.Join(errs, "; "))
	}
	return nil
}

// WriteMCPConfig pre-downloads packages then writes MCP server configuration for the specified target
func WriteMCPConfig(servers []MCPServer, target string) error {
	// Pre-download all npm packages first to avoid startup timeouts
	_ = PreDownloadMCPPackages(servers)

	switch target {
	case "claude":
		return WriteClaudeMCP(servers)
	case "codex":
		return WriteCodexMCP(servers)
	default:
		return fmt.Errorf("unknown MCP target: %s", target)
	}
}
