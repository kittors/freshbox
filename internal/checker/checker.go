package checker

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Status int

const (
	NotInstalled Status = iota
	Installed
)

type Item struct {
	Name      string
	Desc      string // one-line description
	Cmd       string // command to check
	VerFlag   string // flag to get version, e.g. "--version"
	Status    Status
	Version   string
	Category  string
	InstallFn func() error // custom install function, nil = use default brew
	BrewName  string       // brew formula/cask name
	IsCask    bool
}

// resolveCmd finds the command binary, checking extra paths for known tools
func resolveCmd(cmd string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = ""
	}

	// Check tool-specific paths first (before PATH, which may have stubs)
	// Rust: ~/.cargo/bin
	if cmd == "rustup" || cmd == "rustc" || cmd == "cargo" {
		cargoPath := filepath.Join(home, ".cargo", "bin", cmd)
		if _, err := os.Stat(cargoPath); err == nil {
			return cargoPath
		}
	}
	// Java: brew's openjdk (macOS /usr/bin/java is a stub that fails without JDK)
	if cmd == "java" {
		brewJava := "/opt/homebrew/opt/openjdk/bin/java"
		if _, err := os.Stat(brewJava); err == nil {
			return brewJava
		}
	}

	// Try PATH
	if p, err := exec.LookPath(cmd); err == nil {
		return p
	}

	// Fallback: ~/.cargo/bin for any command
	cargoPath := filepath.Join(home, ".cargo", "bin", cmd)
	if _, err := os.Stat(cargoPath); err == nil {
		return cargoPath
	}

	return ""
}

func Check(item *Item) {
	cmdPath := resolveCmd(item.Cmd)
	if cmdPath == "" {
		item.Status = NotInstalled
		item.Version = ""
		return
	}

	if item.VerFlag != "" {
		out, err := exec.Command(cmdPath, item.VerFlag).CombinedOutput()
		if err != nil {
			// Command exists but --version fails (e.g. macOS /usr/bin/java stub)
			item.Status = NotInstalled
			item.Version = ""
			return
		}
		ver := strings.TrimSpace(string(out))
		if idx := strings.Index(ver, "\n"); idx > 0 {
			ver = ver[:idx]
		}
		item.Version = ver
	}

	item.Status = Installed
}

func CheckAll(items []*Item) {
	for _, item := range items {
		Check(item)
	}
}

func DevTools() []*Item {
	return []*Item{
		{Name: "Homebrew", Desc: "macOS package manager", Cmd: "brew", VerFlag: "--version", Category: "dev", BrewName: ""},
		{Name: "Git", Desc: "Distributed version control system", Cmd: "git", VerFlag: "--version", Category: "dev", BrewName: "git"},
		{Name: "Java (JDK)", Desc: "Java development kit for JVM-based development", Cmd: "java", VerFlag: "--version", Category: "dev", BrewName: "openjdk"},
		{Name: "Maven", Desc: "Java project build and dependency management", Cmd: "mvn", VerFlag: "--version", Category: "dev", BrewName: "maven"},
		{Name: "Gradle", Desc: "Flexible build automation tool for JVM projects", Cmd: "gradle", VerFlag: "--version", Category: "dev", BrewName: "gradle"},
		{Name: "Python", Desc: "General-purpose programming language", Cmd: "python3", VerFlag: "--version", Category: "dev", BrewName: "python"},
		{Name: "uv", Desc: "Ultra-fast Python package manager by Astral", Cmd: "uv", VerFlag: "--version", Category: "dev", BrewName: "uv"},
		{Name: "fnm", Desc: "Fast Node.js version manager written in Rust", Cmd: "fnm", VerFlag: "--version", Category: "dev", BrewName: "fnm"},
		{Name: "pnpm", Desc: "Fast, disk-efficient package manager for Node.js", Cmd: "pnpm", VerFlag: "--version", Category: "dev", BrewName: "pnpm"},
		{Name: "Bun", Desc: "All-in-one JavaScript runtime, bundler, and package manager", Cmd: "bun", VerFlag: "--version", Category: "dev", BrewName: "bun"},
		{Name: "Rust (rustup)", Desc: "Systems programming language with memory safety", Cmd: "rustup", VerFlag: "--version", Category: "dev", BrewName: "rustup"},
		{Name: "Go", Desc: "Statically typed language by Google for scalable systems", Cmd: "go", VerFlag: "version", Category: "dev", BrewName: "go"},
	}
}

func Apps() []*Item {
	return []*Item{
		{Name: "Google Chrome", Desc: "Web browser by Google", Cmd: "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", VerFlag: "--version", Category: "app", BrewName: "google-chrome", IsCask: true},
		{Name: "Zed", Desc: "High-performance code editor by the Atom creators", Cmd: "/Applications/Zed.app/Contents/MacOS/cli", VerFlag: "--version", Category: "app", BrewName: "zed", IsCask: true},
		{Name: "IINA", Desc: "Modern media player for macOS", Cmd: "/Applications/IINA.app/Contents/MacOS/IINA", VerFlag: "", Category: "app", BrewName: "iina", IsCask: true},
		{Name: "Kaku", Desc: "Lightweight terminal app built on WezTerm by tw93", Cmd: "kaku", VerFlag: "--version", Category: "app", BrewName: "tw93/tap/kakuku", IsCask: true},
		{Name: "Karabiner-Elements", Desc: "Powerful keyboard customizer for macOS", Cmd: "/Applications/Karabiner-Elements.app/Contents/MacOS/Karabiner-Elements", VerFlag: "", Category: "app", BrewName: "karabiner-elements", IsCask: true},
		{Name: "Mole", Desc: "macOS system cleaner to free up disk space by tw93", Cmd: "mo", VerFlag: "", Category: "app", BrewName: "tw93/tap/mole", IsCask: false},
		{Name: "Tabby", Desc: "Modern open-source terminal with SSH and serial support", Cmd: "/Applications/Tabby.app/Contents/MacOS/Tabby", VerFlag: "", Category: "app", BrewName: "tabby", IsCask: true},
	}
}

func AITools() []*Item {
	return []*Item{
		{Name: "Codex", Desc: "OpenAI's AI coding assistant CLI", Cmd: "codex", VerFlag: "--version", Category: "ai", BrewName: ""},
		{Name: "Claude Code", Desc: "Anthropic's AI coding assistant CLI", Cmd: "claude", VerFlag: "--version", Category: "ai", BrewName: ""},
	}
}

// CheckApp checks if a macOS .app exists
func CheckApp(item *Item) {
	if item.IsCask {
		appPaths := map[string]string{
			"google-chrome":      "/Applications/Google Chrome.app",
			"zed":                "/Applications/Zed.app",
			"iina":               "/Applications/IINA.app",
			"tw93/tap/kakuku":    "/Applications/Kaku.app",
			"karabiner-elements": "/Applications/Karabiner-Elements.app",
			"tabby":              "/Applications/Tabby.app",
		}
		if p, ok := appPaths[item.BrewName]; ok {
			if _, err := os.Stat(p); err == nil {
				item.Status = Installed
				ver, verErr := exec.Command("defaults", "read", p+"/Contents/Info.plist", "CFBundleShortVersionString").CombinedOutput()
				if verErr == nil {
					item.Version = strings.TrimSpace(string(ver))
				}
				return
			}
		}
	}
	Check(item)
}
