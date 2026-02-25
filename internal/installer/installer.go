package installer

import (
	"fmt"
	"os/exec"
	"strings"
)

// BrewInstall installs a formula or cask via Homebrew
func BrewInstall(name string, isCask bool) error {
	args := []string{"install"}
	if isCask {
		args = append(args, "--cask")
	}
	args = append(args, name)
	cmd := exec.Command("brew", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, string(out))
	}
	return nil
}

// InstallHomebrew installs Homebrew itself
func InstallHomebrew() error {
	script := `/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`
	cmd := exec.Command("bash", "-c", script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, string(out))
	}
	return nil
}

// InstallCodex installs OpenAI Codex CLI via npm
func InstallCodex() error {
	cmd := exec.Command("npm", "install", "-g", "@openai/codex")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, string(out))
	}
	return nil
}

// InstallClaudeCode installs Claude Code via npm
func InstallClaudeCode() error {
	cmd := exec.Command("npm", "install", "-g", "@anthropic-ai/claude-code")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, string(out))
	}
	return nil
}

// InstallRust installs Rust via rustup
func InstallRust() error {
	cmd := exec.Command("bash", "-c", "curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, string(out))
	}
	return nil
}

// FnmInstallNode installs a specific Node.js version via fnm
func FnmInstallNode(version string) error {
	cmd := exec.Command("fnm", "install", version)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, string(out))
	}
	return nil
}

// FnmListRemote lists available Node.js versions
func FnmListRemote() ([]string, error) {
	cmd := exec.Command("fnm", "list-remote")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%s: %s", err, string(out))
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	// filter to only LTS / major versions for cleaner display
	var versions []string
	for _, line := range lines {
		v := strings.TrimSpace(line)
		if v != "" {
			versions = append(versions, v)
		}
	}
	// reverse so newest first
	for i, j := 0, len(versions)-1; i < j; i, j = i+1, j-1 {
		versions[i], versions[j] = versions[j], versions[i]
	}
	return versions, nil
}

// SetDefaultBrowser sets Chrome as default browser via LSHandlers
func SetDefaultBrowser() error {
	// Use defaults write to set Chrome as default HTTP/HTTPS handler
	// This avoids opening a browser window
	types := []struct{ scheme, role string }{
		{"http", "LSHandlerURLScheme"},
		{"https", "LSHandlerURLScheme"},
	}
	for _, t := range types {
		cmd := exec.Command("bash", "-c", fmt.Sprintf(
			`defaults write com.apple.LaunchServices/com.apple.launchservices.secure LSHandlers -array-add '{"LSHandlerURLScheme"="%s";"LSHandlerRoleAll"="com.google.chrome";}'`,
			t.scheme))
		_ = cmd.Run()
	}
	return nil
}

// SetJavaHome configures JAVA_HOME in shell profile
func SetJavaHome() error {
	cmd := exec.Command("bash", "-c", `echo 'export JAVA_HOME=$(/usr/libexec/java_home)' >> ~/.zshrc`)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, string(out))
	}
	return nil
}
