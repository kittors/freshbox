package checker

import (
	"os/exec"
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

func Check(item *Item) {
	path, err := exec.LookPath(item.Cmd)
	if err != nil {
		item.Status = NotInstalled
		item.Version = ""
		return
	}
	_ = path
	item.Status = Installed

	if item.VerFlag != "" {
		out, err := exec.Command(item.Cmd, item.VerFlag).CombinedOutput()
		if err == nil {
			ver := strings.TrimSpace(string(out))
			if idx := strings.Index(ver, "\n"); idx > 0 {
				ver = ver[:idx]
			}
			item.Version = ver
		}
	}
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
		}
		if p, ok := appPaths[item.BrewName]; ok {
			out, err := exec.Command("test", "-d", p).CombinedOutput()
			_ = out
			if err == nil {
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
