package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kittors/freshbox/internal/checker"
	"github.com/kittors/freshbox/internal/config"
	"github.com/kittors/freshbox/internal/installer"
	"github.com/kittors/freshbox/internal/setup"
)

// installDoneMsg signals all installs are complete
type installDoneMsg struct{}

// spinnerTickMsg wraps spinner tick
type spinnerTickMsg struct {
	msg tea.Msg
}

// buildInstallQueue builds the ordered list of things to install
func (m *Model) buildInstallQueue() []installTask {
	var queue []installTask

	// 1. Dev tools
	for _, item := range m.devTools {
		if !m.selected[item.Name] || item.Status == checker.Installed {
			continue
		}
		task := installTask{name: item.Name}
		switch item.Name {
		case "Homebrew":
			task.fn = func() error { return installer.InstallHomebrew() }
		case "Rust (rustup)":
			task.fn = func() error { return installer.InstallRust() }
		default:
			brewName := item.BrewName
			isCask := item.IsCask
			if brewName != "" {
				task.fn = func() error { return installer.BrewInstall(brewName, isCask) }
			}
		}
		if task.fn != nil {
			queue = append(queue, task)
		}
	}

	// 2. Apps
	for _, item := range m.apps {
		if !m.selected[item.Name] || item.Status == checker.Installed {
			continue
		}
		brewName := item.BrewName
		isCask := item.IsCask
		task := installTask{
			name: item.Name,
			fn:   func() error { return installer.BrewInstall(brewName, isCask) },
		}
		queue = append(queue, task)
	}

	// 3. AI tools
	for _, item := range m.aiTools {
		if !m.selected[item.Name] || item.Status == checker.Installed {
			continue
		}
		switch item.Name {
		case "Codex":
			queue = append(queue, installTask{
				name: "Codex CLI",
				fn:   func() error { return installer.InstallCodex() },
			})
		case "Claude Code":
			queue = append(queue, installTask{
				name: "Claude Code",
				fn:   func() error { return installer.InstallClaudeCode() },
			})
		}
	}

	// 4. fnm Node versions
	for v, sel := range m.fnmSelected {
		if !sel {
			continue
		}
		ver := v
		queue = append(queue, installTask{
			name: "Node.js " + ver,
			fn:   func() error { return installer.FnmInstallNode(ver) },
		})
	}

	// 5. Codex config
	if m.codexKey != "" || m.codexURL != "" {
		queue = append(queue, installTask{
			name: "Codex config (config.toml + auth.json)",
			fn: func() error {
				err := config.WriteCodexConfig(config.CodexConfig{
					Model:         m.codexModel,
					ThinkingLevel: m.codexThink,
					BaseURL:       m.codexURL,
				})
				if err != nil {
					return err
				}
				return config.WriteCodexAuth(config.CodexAuth{APIKey: m.codexKey})
			},
		})
	}

	// 6. Claude config
	if m.claudeKey != "" || m.claudeURL != "" {
		queue = append(queue, installTask{
			name: "Claude Code config",
			fn: func() error {
				return config.WriteClaudeConfig(config.ClaudeConfig{
					Model:   m.claudeModel,
					BaseURL: m.claudeURL,
					APIKey:  m.claudeKey,
				})
			},
		})
	}

	// 7. MCP servers
	var selectedMCPs []config.MCPServer
	for _, mcp := range m.mcps {
		if m.mcpSelected[mcp.Name] {
			selectedMCPs = append(selectedMCPs, mcp)
		}
	}
	if len(selectedMCPs) > 0 {
		// Check if Claude Code is selected, configured, or already installed
		claudeReady := m.selected["Claude Code"] || m.claudeKey != "" || m.claudeURL != ""
		if !claudeReady {
			for _, item := range m.aiTools {
				if item.Name == "Claude Code" && item.Status == checker.Installed {
					claudeReady = true
					break
				}
			}
		}
		if claudeReady {
			queue = append(queue, installTask{
				name: "MCP servers for Claude Code",
				fn:   func() error { return config.WriteMCPConfig(selectedMCPs, "claude") },
			})
		}

		// Check if Codex is selected, configured, or already installed
		codexReady := m.selected["Codex"] || m.codexKey != "" || m.codexURL != ""
		if !codexReady {
			for _, item := range m.aiTools {
				if item.Name == "Codex" && item.Status == checker.Installed {
					codexReady = true
					break
				}
			}
		}
		if codexReady {
			queue = append(queue, installTask{
				name: "MCP servers for Codex",
				fn:   func() error { return config.WriteMCPConfig(selectedMCPs, "codex") },
			})
		}
	}

	// 8. System defaults
	if m.sysDefaults["browser_chrome"] {
		queue = append(queue, installTask{
			name: "Set default browser → Chrome",
			fn:   func() error { return installer.SetDefaultBrowser() },
		})
	}
	if m.sysDefaults["editor_zed"] {
		queue = append(queue, installTask{
			name: "Set default editor → Zed",
			fn: func() error {
				cmd := exec.Command("bash", "-c", `defaults write com.apple.LaunchServices/com.apple.launchservices.secure LSHandlers -array-add '{"LSHandlerContentType"="public.plain-text";"LSHandlerRoleAll"="dev.zed.Zed";}'`)
				return cmd.Run()
			},
		})
	}
	if m.sysDefaults["player_iina"] {
		queue = append(queue, installTask{
			name: "Set default player → IINA",
			fn: func() error {
				types := []string{"public.movie", "public.video", "public.audio"}
				for _, t := range types {
					cmd := exec.Command("bash", "-c", fmt.Sprintf(`defaults write com.apple.LaunchServices/com.apple.launchservices.secure LSHandlers -array-add '{"LSHandlerContentType"="%s";"LSHandlerRoleAll"="com.colliderli.iina";}'`, t))
					_ = cmd.Run()
				}
				return nil
			},
		})
	}

	// 9. Java JAVA_HOME
	if m.selected["Java (JDK)"] {
		queue = append(queue, installTask{
			name: "Configure JAVA_HOME",
			fn:   func() error { return installer.SetJavaHome() },
		})
	}

	// 10. Extra setup
	if m.extraSetup["zed_theme"] {
		queue = append(queue, installTask{
			name: "Zed Catppuccin Blur Theme",
			fn:   func() error { return setup.SetupZedTheme() },
		})
	}
	if m.extraSetup["kaku_init"] {
		queue = append(queue, installTask{
			name: "Kaku Terminal Setup (config + zsh plugins)",
			fn:   func() error { return setup.SetupKaku() },
		})
	}
	if m.extraSetup["karabiner_kaku"] {
		queue = append(queue, installTask{
			name: "Karabiner ⌃⌥⌘T → Kaku shortcut",
			fn:   func() error { return setup.SetupKarabiner() },
		})
	}
	if m.extraSetup["dev_workspace"] {
		queue = append(queue, installTask{
			name: "Developer Workspace + Finder config",
			fn:   func() error { return setup.SetupDevWorkspace() },
		})
	}

	return queue
}

type installTask struct {
	name string
	fn   func() error
}

// startInstallSequence kicks off the install with progress reporting
func (m *Model) startInstallSequence() tea.Cmd {
	queue := m.buildInstallQueue()
	if len(queue) == 0 {
		return func() tea.Msg {
			return installDoneMsg{}
		}
	}

	m.installQueue = queue
	m.installIdx = 0
	m.installTotal = len(queue)
	m.currentTask = queue[0].name

	// write log header
	appendLog(fmt.Sprintf("=== freshbox install started (%d tasks) ===", len(queue)))

	// start spinner + first install concurrently
	return tea.Batch(m.spinner.Tick, m.runNextInstall())
}

func (m *Model) runNextInstall() tea.Cmd {
	if m.installIdx >= len(m.installQueue) {
		return func() tea.Msg { return installDoneMsg{} }
	}

	task := m.installQueue[m.installIdx]

	return func() tea.Msg {
		time.Sleep(80 * time.Millisecond)
		err := task.fn()
		return InstallMsg{Name: task.name, Err: err}
	}
}

// HandleInstallMsg processes install results and triggers next install
func (m *Model) HandleInstallMsg(msg InstallMsg) tea.Cmd {
	if msg.Err != nil {
		fullErr := strings.TrimSpace(msg.Err.Error())
		// write full error to log file
		appendLog(fmt.Sprintf("[FAIL] %s\n       %s", msg.Name, fullErr))

		errMsg := fullErr
		if len(errMsg) > 60 {
			errMsg = errMsg[:60] + "..."
		}
		m.installLog = append(m.installLog, installLogEntry{
			name:    msg.Name,
			success: false,
			errMsg:  errMsg,
		})
	} else {
		appendLog(fmt.Sprintf("[ OK ] %s", msg.Name))
		m.installLog = append(m.installLog, installLogEntry{
			name:    msg.Name,
			success: true,
		})
	}

	m.installIdx++
	if m.installIdx >= len(m.installQueue) {
		m.currentTask = ""
		return func() tea.Msg { return installDoneMsg{} }
	}

	m.currentTask = m.installQueue[m.installIdx].name
	return m.runNextInstall()
}

type installLogEntry struct {
	name    string
	success bool
	errMsg  string
}

// logFilePath returns the path to the install error log
func logFilePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".freshbox", "install.log")
}

// appendLog writes a line to the install log file
func appendLog(line string) {
	path := logFilePath()
	os.MkdirAll(filepath.Dir(path), 0755)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	f.WriteString(time.Now().Format("2006-01-02 15:04:05") + "  " + line + "\n")
}

// NewSpinner creates a styled spinner
func NewSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Spinner{
		Frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		FPS:    80 * time.Millisecond,
	}
	return s
}

// renderInstallProgress renders the full install page with spinner, progress bar, and log
func (m Model) renderInstallProgress() string {
	var b strings.Builder

	total := m.installTotal
	done := len(m.installLog)
	pct := 0
	if total > 0 {
		pct = done * 100 / total
	}

	// Title with progress count
	title := fmt.Sprintf("⏳ %s  [%d/%d]", m.t.TitleInstalling, done, total)
	b.WriteString(SubtitleStyle.Render(title) + "\n\n")

	// Progress bar
	barWidth := 40
	filled := 0
	if total > 0 {
		filled = barWidth * done / total
	}
	if filled > barWidth {
		filled = barWidth
	}
	empty := barWidth - filled

	bar := ProgressStyle.Render(strings.Repeat("█", filled)) +
		DimStyle.Render(strings.Repeat("░", empty)) +
		DimStyle.Render(fmt.Sprintf(" %d%%", pct))
	b.WriteString("  " + bar + "\n\n")

	// Completed items log (show last 10 to avoid overflow)
	logStart := 0
	if len(m.installLog) > 10 {
		logStart = len(m.installLog) - 10
		b.WriteString(DimStyle.Render(fmt.Sprintf("  ... %d more above\n", logStart)))
	}
	for i := logStart; i < len(m.installLog); i++ {
		entry := m.installLog[i]
		if entry.success {
			b.WriteString(fmt.Sprintf("  %s %s\n",
				SuccessStyle.Render("✓"),
				DimStyle.Render(entry.name)))
		} else {
			b.WriteString(fmt.Sprintf("  %s %s  %s\n",
				ErrorStyle.Render("✗"),
				entry.name,
				ErrorStyle.Render(entry.errMsg)))
		}
	}

	// Current task with spinner
	if m.currentTask != "" {
		spinnerView := ProgressStyle.Render(m.spinner.View())
		taskName := lipgloss.NewStyle().Foreground(Cyan).Bold(true).Render(m.currentTask)
		b.WriteString(fmt.Sprintf("\n  %s %s %s\n",
			spinnerView,
			lipgloss.NewStyle().Foreground(Yellow).Render("Installing"),
			taskName))
	}

	// Upcoming tasks preview (next 3)
	upcoming := []string{}
	for i := m.installIdx + 1; i < len(m.installQueue) && len(upcoming) < 3; i++ {
		upcoming = append(upcoming, m.installQueue[i].name)
	}
	if len(upcoming) > 0 {
		b.WriteString("\n" + DimStyle.Render("  Next up:") + "\n")
		for _, name := range upcoming {
			b.WriteString(DimStyle.Render("    ○ "+name) + "\n")
		}
	}

	return BoxStyle.Render(b.String())
}
