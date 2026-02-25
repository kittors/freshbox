package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kittors/freshbox/internal/checker"
	"github.com/kittors/freshbox/internal/config"
)

// Page represents the current TUI page
type Page int

const (
	PageLang Page = iota
	PageWelcome
	PageDevTools
	PageApps
	PageFnmVersions
	PageAITools
	PageCodexConfig
	PageClaudeConfig
	PageMCP
	PageExtraSetup
	PageSystemDefaults
	PageInstalling
	PageDone
)

func pageNames(t T) []string {
	return []string{
		t.LangTitle,
		t.PageWelcome,
		t.PageDevTools,
		t.PageApps,
		t.PageNodeVer,
		t.PageAITools,
		t.PageCodexCfg,
		t.PageClaudeCfg,
		t.PageMCP,
		t.PageExtraSetup,
		t.PageSysDefaults,
		t.PageInstalling,
		t.PageDone,
	}
}

// InstallMsg is sent when an install completes
type InstallMsg struct {
	Name string
	Err  error
}

// FnmVersionsMsg is sent when fnm versions are fetched
type FnmVersionsMsg struct {
	Versions []string
	Err      error
}

type Model struct {
	page         Page
	lang         Lang
	t            T
	langCursor   int // 0=English, 1=Chinese
	width        int
	height       int
	devTools     []*checker.Item
	apps         []*checker.Item
	aiTools      []*checker.Item
	mcps         []config.MCPServer
	fnmVersions  []string
	cursor       int
	selected     map[string]bool
	fnmSelected  map[string]bool
	mcpSelected  map[string]bool
	sysDefaults  map[string]bool
	extraSetup   map[string]bool

	// text inputs for config
	inputs       []textinput.Model
	inputFocus   int
	inputPage    Page // which config page we're on

	// codex config values
	codexModel   string
	codexThink   string
	codexURL     string
	codexKey     string

	// claude config values
	claudeModel  string
	claudeURL    string
	claudeKey    string

	// install progress
	installLog   []installLogEntry
	installing   bool
	installDone  bool
	installQueue []installTask
	installIdx   int
	installTotal int
	currentTask  string
	spinner      spinner.Model

	// error
	err error
}

func NewModel() Model {
	devTools := checker.DevTools()
	checker.CheckAll(devTools)

	apps := checker.Apps()
	for _, a := range apps {
		checker.CheckApp(a)
	}

	aiTools := checker.AITools()
	checker.CheckAll(aiTools)

	mcps := config.AvailableMCPs()

	m := Model{
		page:        PageLang,
		lang:        LangEN,
		t:           GetText(LangEN),
		devTools:    devTools,
		apps:        apps,
		aiTools:     aiTools,
		mcps:        mcps,
		spinner:     NewSpinner(),
		selected:    make(map[string]bool),
		fnmSelected: make(map[string]bool),
		mcpSelected: make(map[string]bool),
		sysDefaults: map[string]bool{
			"browser_chrome": true,
			"editor_zed":     true,
			"player_iina":    true,
		},
		extraSetup: map[string]bool{
			"zed_theme":      true,
			"kaku_init":      true,
			"karabiner_kaku": true,
			"dev_workspace":  true,
		},
	}

	// pre-select uninstalled items
	for _, item := range devTools {
		if item.Status == checker.NotInstalled {
			m.selected[item.Name] = true
		}
	}
	for _, item := range apps {
		if item.Status == checker.NotInstalled {
			m.selected[item.Name] = true
		}
	}
	for _, item := range aiTools {
		if item.Status == checker.NotInstalled {
			m.selected[item.Name] = true
		}
	}
	// pre-select popular MCPs
	for _, mcp := range mcps[:4] {
		m.mcpSelected[mcp.Name] = true
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case InstallMsg:
		cmd := m.HandleInstallMsg(msg)
		return m, cmd

	case installDoneMsg:
		m.installing = false
		m.installDone = true
		m.page = PageDone
		return m, nil

	case spinner.TickMsg:
		if m.installing {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
		return m, nil

	case FnmVersionsMsg:
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			m.fnmVersions = msg.Versions
			// only show top 30
			if len(m.fnmVersions) > 30 {
				m.fnmVersions = m.fnmVersions[:30]
			}
		}
		return m, nil

	case tea.KeyMsg:
		// Language selection page has its own key handling
		if m.page == PageLang {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				if m.langCursor > 0 {
					m.langCursor--
				}
			case "down", "j":
				if m.langCursor < 1 {
					m.langCursor++
				}
			case "enter", " ":
				if m.langCursor == 0 {
					m.lang = LangEN
				} else {
					m.lang = LangZH
				}
				m.t = GetText(m.lang)
				m.page = PageWelcome
				m.cursor = 0
			}
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "q":
			if m.page == PageWelcome || m.page == PageDone {
				return m, tea.Quit
			}
			// on other pages, q goes back
			if m.page > PageWelcome && !m.installing {
				m.page--
				m.cursor = 0
				return m, nil
			}
			return m, tea.Quit

		case "tab", "right", "l":
			if !m.installing && m.page < PageDone {
				return m.nextPage()
			}

		case "shift+tab", "left", "h":
			if !m.installing && m.page > PageWelcome {
				m.page--
				m.cursor = 0
			}

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			max := m.currentListLen() - 1
			if m.cursor < max {
				m.cursor++
			}

		case " ":
			m.toggleCurrent()

		case "a":
			m.selectAll()

		case "n":
			m.selectNone()

		case "enter":
			if m.page == PageWelcome {
				m.page = PageDevTools
				m.cursor = 0
				return m, nil
			}
			if m.page == PageDone {
				return m, tea.Quit
			}
			return m.nextPage()
		}

		// handle text input on config pages
		if m.page == PageCodexConfig || m.page == PageClaudeConfig {
			return m.updateInputs(msg)
		}
	}

	return m, nil
}

func (m Model) nextPage() (Model, tea.Cmd) {
	switch m.page {
	case PageDevTools:
		m.page = PageApps
	case PageApps:
		// check if fnm is selected, if so go to fnm page
		if m.selected["fnm"] {
			m.page = PageFnmVersions
		} else {
			m.page = PageAITools
		}
	case PageFnmVersions:
		m.page = PageAITools
	case PageAITools:
		if m.selected["Codex"] {
			m.page = PageCodexConfig
			m.initCodexInputs()
		} else if m.selected["Claude Code"] {
			m.page = PageClaudeConfig
			m.initClaudeInputs()
		} else {
			m.page = PageMCP
		}
	case PageCodexConfig:
		m.saveCodexInputs()
		if m.selected["Claude Code"] {
			m.page = PageClaudeConfig
			m.initClaudeInputs()
		} else {
			m.page = PageMCP
		}
	case PageClaudeConfig:
		m.saveClaudeInputs()
		m.page = PageMCP
	case PageMCP:
		m.page = PageExtraSetup
	case PageExtraSetup:
		m.page = PageSystemDefaults
	case PageSystemDefaults:
		m.page = PageInstalling
		m.installing = true
		return m, m.startInstallSequence()
	default:
		if m.page < PageDone {
			m.page++
		}
	}
	m.cursor = 0
	return m, nil
}

func (m *Model) initCodexInputs() {
	m.inputs = make([]textinput.Model, 4)
	placeholders := []string{"Model (e.g. o4-mini)", "Thinking level (low/medium/high)", "Base URL", "API Key"}
	defaults := []string{"o4-mini", "medium", "https://api.openai.com/v1", ""}
	for i := range m.inputs {
		t := textinput.New()
		t.Placeholder = placeholders[i]
		t.SetValue(defaults[i])
		if i == 3 {
			t.EchoMode = textinput.EchoPassword
		}
		if i == 0 {
			t.Focus()
		}
		m.inputs[i] = t
	}
	m.inputFocus = 0
	m.inputPage = PageCodexConfig
}

func (m *Model) initClaudeInputs() {
	m.inputs = make([]textinput.Model, 3)
	placeholders := []string{"Model (e.g. claude-sonnet-4-6)", "Base URL", "API Key"}
	defaults := []string{"claude-sonnet-4-6", "https://api.anthropic.com", ""}
	for i := range m.inputs {
		t := textinput.New()
		t.Placeholder = placeholders[i]
		t.SetValue(defaults[i])
		if i == 2 {
			t.EchoMode = textinput.EchoPassword
		}
		if i == 0 {
			t.Focus()
		}
		m.inputs[i] = t
	}
	m.inputFocus = 0
	m.inputPage = PageClaudeConfig
}

func (m *Model) saveCodexInputs() {
	if len(m.inputs) >= 4 {
		m.codexModel = m.inputs[0].Value()
		m.codexThink = m.inputs[1].Value()
		m.codexURL = m.inputs[2].Value()
		m.codexKey = m.inputs[3].Value()
	}
}

func (m *Model) saveClaudeInputs() {
	if len(m.inputs) >= 3 {
		m.claudeModel = m.inputs[0].Value()
		m.claudeURL = m.inputs[1].Value()
		m.claudeKey = m.inputs[2].Value()
	}
}

func (m Model) updateInputs(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "tab", "down":
		m.inputFocus++
		if m.inputFocus >= len(m.inputs) {
			m.inputFocus = 0
		}
	case "shift+tab", "up":
		m.inputFocus--
		if m.inputFocus < 0 {
			m.inputFocus = len(m.inputs) - 1
		}
	case "enter":
		// if on last input, go next page
		if m.inputFocus == len(m.inputs)-1 {
			return m.nextPage()
		}
		m.inputFocus++
		if m.inputFocus >= len(m.inputs) {
			m.inputFocus = 0
		}
	default:
		// update the focused input
		var cmd tea.Cmd
		m.inputs[m.inputFocus], cmd = m.inputs[m.inputFocus].Update(msg)
		return m, cmd
	}

	// update focus
	for i := range m.inputs {
		if i == m.inputFocus {
			m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
	return m, nil
}

func (m *Model) toggleCurrent() {
	switch m.page {
	case PageDevTools:
		if m.cursor < len(m.devTools) {
			item := m.devTools[m.cursor]
			if item.Status == checker.NotInstalled {
				m.selected[item.Name] = !m.selected[item.Name]
			}
		}
	case PageApps:
		if m.cursor < len(m.apps) {
			item := m.apps[m.cursor]
			if item.Status == checker.NotInstalled {
				m.selected[item.Name] = !m.selected[item.Name]
			}
		}
	case PageAITools:
		if m.cursor < len(m.aiTools) {
			item := m.aiTools[m.cursor]
			if item.Status == checker.NotInstalled {
				m.selected[item.Name] = !m.selected[item.Name]
			}
		}
	case PageFnmVersions:
		if m.cursor < len(m.fnmVersions) {
			v := m.fnmVersions[m.cursor]
			m.fnmSelected[v] = !m.fnmSelected[v]
		}
	case PageMCP:
		if m.cursor < len(m.mcps) {
			name := m.mcps[m.cursor].Name
			m.mcpSelected[name] = !m.mcpSelected[name]
		}
	case PageSystemDefaults:
		keys := []string{"browser_chrome", "editor_zed", "player_iina"}
		if m.cursor < len(keys) {
			m.sysDefaults[keys[m.cursor]] = !m.sysDefaults[keys[m.cursor]]
		}
	case PageExtraSetup:
		keys := []string{"zed_theme", "kaku_init", "karabiner_kaku", "dev_workspace"}
		if m.cursor < len(keys) {
			m.extraSetup[keys[m.cursor]] = !m.extraSetup[keys[m.cursor]]
		}
	}
}

func (m *Model) selectAll() {
	switch m.page {
	case PageDevTools:
		for _, item := range m.devTools {
			if item.Status == checker.NotInstalled {
				m.selected[item.Name] = true
			}
		}
	case PageApps:
		for _, item := range m.apps {
			if item.Status == checker.NotInstalled {
				m.selected[item.Name] = true
			}
		}
	case PageMCP:
		for _, mcp := range m.mcps {
			m.mcpSelected[mcp.Name] = true
		}
	}
}

func (m *Model) selectNone() {
	switch m.page {
	case PageDevTools:
		for _, item := range m.devTools {
			m.selected[item.Name] = false
		}
	case PageApps:
		for _, item := range m.apps {
			m.selected[item.Name] = false
		}
	case PageMCP:
		for _, mcp := range m.mcps {
			m.mcpSelected[mcp.Name] = false
		}
	}
}

func (m Model) currentListLen() int {
	switch m.page {
	case PageDevTools:
		return len(m.devTools)
	case PageApps:
		return len(m.apps)
	case PageAITools:
		return len(m.aiTools)
	case PageFnmVersions:
		return len(m.fnmVersions)
	case PageMCP:
		return len(m.mcps)
	case PageSystemDefaults:
		return 3
	case PageExtraSetup:
		return 4
	default:
		return 1
	}
}

// View renders the TUI
func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var b strings.Builder

	// Header with tabs
	b.WriteString(m.renderHeader())
	b.WriteString("\n")

	// Main content
	switch m.page {
	case PageLang:
		b.WriteString(m.renderLangSelect())
	case PageWelcome:
		b.WriteString(m.renderWelcome())
	case PageDevTools:
		b.WriteString(m.renderCheckList("🔧 "+m.t.TitleDevTools, m.devTools))
	case PageApps:
		b.WriteString(m.renderCheckList("📦 "+m.t.TitleApps, m.apps))
	case PageFnmVersions:
		b.WriteString(m.renderFnmVersions())
	case PageAITools:
		b.WriteString(m.renderCheckList("🤖 "+m.t.TitleAITools, m.aiTools))
	case PageCodexConfig:
		b.WriteString(m.renderConfigForm("Codex Configuration"))
	case PageClaudeConfig:
		b.WriteString(m.renderConfigForm("Claude Code Configuration"))
	case PageMCP:
		b.WriteString(m.renderMCPList())
	case PageExtraSetup:
		b.WriteString(m.renderExtraSetup())
	case PageSystemDefaults:
		b.WriteString(m.renderSystemDefaults())
	case PageInstalling:
		b.WriteString(m.renderInstallProgress())
	case PageDone:
		b.WriteString(m.renderDone())
	}

	// Footer
	b.WriteString("\n")
	b.WriteString(m.renderFooter())

	return b.String()
}

func (m Model) renderHeader() string {
	logo := LogoStyle.Render("  ╭─────────────────────────────╮")
	logo += "\n" + LogoStyle.Render("  │   🍃 freshbox  v1.0.0      │")
	logo += "\n" + LogoStyle.Render("  │   macOS Setup Assistant     │")
	logo += "\n" + LogoStyle.Render("  ╰─────────────────────────────╯")

	// Tab bar
	names := pageNames(m.t)
	var tabs []string
	for i, name := range names {
		if Page(i) == m.page {
			tabs = append(tabs, ActiveTabStyle.Render(name))
		} else if Page(i) <= m.page {
			tabs = append(tabs, InactiveTabStyle.Render(name))
		}
	}
	tabBar := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)

	return logo + "\n\n" + tabBar + "\n"
}

func (m Model) renderWelcome() string {
	arrow := lipgloss.NewStyle().Foreground(Cyan).Render("→")
	welcome := SubtitleStyle.Render(m.t.WelcomeTitle) + "\n\n"
	welcome += "  " + m.t.WelcomeDesc + "\n\n"
	welcome += "  " + arrow + " " + m.t.WelcomeDevTools + "\n"
	welcome += "  " + arrow + " " + m.t.WelcomeFnm + "\n"
	welcome += "  " + arrow + " " + m.t.WelcomeApps + "\n"
	welcome += "  " + arrow + " " + m.t.WelcomeAI + "\n"
	welcome += "  " + arrow + " " + m.t.WelcomeMCP + "\n"
	welcome += "  " + arrow + " " + m.t.WelcomeSys + "\n"
	welcome += "\n\n  " + SelectedStyle.Render(m.t.WelcomeStart) + ", " + DimStyle.Render(m.t.WelcomeQuit)
	return BoxStyle.Render(welcome)
}

func (m Model) renderCheckList(title string, items []*checker.Item) string {
	var b strings.Builder
	b.WriteString(SubtitleStyle.Render(title) + "\n\n")

	for i, item := range items {
		cursor := "  "
		if i == m.cursor {
			cursor = CursorStyle.Render("▸ ")
		}

		desc := ""
		if item.Desc != "" {
			desc = DimStyle.Render(" — " + item.Desc)
		}

		if item.Status == checker.Installed {
			name := InstalledStyle.Render(item.Name)
			ver := VersionStyle.Render(" (" + item.Version + ")")
			b.WriteString(fmt.Sprintf("  %s %s %s%s%s\n", cursor, CheckedStyle.Render("■"), name, ver, desc))
		} else {
			check := UncheckedStyle.Render("□")
			if m.selected[item.Name] {
				check = CheckedStyle.Render("■")
			}
			name := NotInstalledStyle.Render(item.Name)
			b.WriteString(fmt.Sprintf("  %s %s %s%s\n", cursor, check, name, desc))
		}
	}

	return BoxStyle.Render(b.String())
}

func (m Model) renderFnmVersions() string {
	var b strings.Builder
	b.WriteString(SubtitleStyle.Render("📦 "+m.t.TitleFnmVer) + "\n\n")

	if len(m.fnmVersions) == 0 {
		b.WriteString("  " + DimStyle.Render(m.t.FnmHint) + "\n")
		b.WriteString("  " + DimStyle.Render(m.t.FnmLTSHint) + "\n\n")
		// show common versions as fallback
		common := []string{"v22.x (Current)", "v20.x (LTS)", "v18.x (LTS)", "v16.x (Maintenance)"}
		for i, v := range common {
			cursor := "  "
			if i == m.cursor {
				cursor = CursorStyle.Render("▸ ")
			}
			check := UncheckedStyle.Render("□")
			if m.fnmSelected[v] {
				check = CheckedStyle.Render("■")
			}
			b.WriteString(fmt.Sprintf("  %s %s %s\n", cursor, check, v))
		}
	} else {
		for i, v := range m.fnmVersions {
			cursor := "  "
			if i == m.cursor {
				cursor = CursorStyle.Render("▸ ")
			}
			check := UncheckedStyle.Render("□")
			if m.fnmSelected[v] {
				check = CheckedStyle.Render("■")
			}
			b.WriteString(fmt.Sprintf("  %s %s %s\n", cursor, check, v))
		}
	}

	return BoxStyle.Render(b.String())
}

func (m Model) renderConfigForm(title string) string {
	var b strings.Builder
	b.WriteString(SubtitleStyle.Render("⚙️  "+title) + "\n\n")

	labels := []string{}
	if m.inputPage == PageCodexConfig {
		labels = []string{m.t.CfgModel, m.t.CfgThinkLevel, m.t.CfgBaseURL, m.t.CfgAPIKey}
	} else {
		labels = []string{m.t.CfgModel, m.t.CfgBaseURL, m.t.CfgAPIKey}
	}

	for i, input := range m.inputs {
		label := LabelStyle.Render(labels[i] + ":")
		field := input.View()
		if i == m.inputFocus {
			b.WriteString(fmt.Sprintf("  %s %s  %s\n\n", CursorStyle.Render("▸"), label, field))
		} else {
			b.WriteString(fmt.Sprintf("    %s  %s\n\n", label, field))
		}
	}

	return BoxStyle.Render(b.String())
}

func (m Model) renderMCPList() string {
	var b strings.Builder
	b.WriteString(SubtitleStyle.Render("🔌 "+m.t.TitleMCP) + "\n")
	b.WriteString(DimStyle.Render("  "+m.t.TitleMCPDesc) + "\n\n")

	for i, mcp := range m.mcps {
		cursor := "  "
		if i == m.cursor {
			cursor = CursorStyle.Render("▸ ")
		}
		check := UncheckedStyle.Render("□")
		if m.mcpSelected[mcp.Name] {
			check = CheckedStyle.Render("■")
		}
		name := lipgloss.NewStyle().Foreground(White).Render(mcp.Name)
		cmd := DimStyle.Render(fmt.Sprintf(" (%s %s)", mcp.Command, strings.Join(mcp.Args, " ")))
		b.WriteString(fmt.Sprintf("  %s %s %s%s\n", cursor, check, name, cmd))
	}

	return BoxStyle.Render(b.String())
}

func (m Model) renderExtraSetup() string {
	var b strings.Builder
	b.WriteString(SubtitleStyle.Render("🎨 "+m.t.TitleExtraSetup) + "\n")
	b.WriteString(DimStyle.Render("  "+m.t.TitleExtraSetupDesc) + "\n\n")

	extras := []struct {
		key   string
		label string
		desc  string
	}{
		{"zed_theme", m.t.ExtraZedTheme, m.t.ExtraZedThemeDesc},
		{"kaku_init", m.t.ExtraKakuInit, m.t.ExtraKakuInitDesc},
		{"karabiner_kaku", m.t.ExtraKarabiner, m.t.ExtraKarabinerDesc},
		{"dev_workspace", m.t.ExtraDevWorkspace, m.t.ExtraDevWorkspaceDesc},
	}

	for i, e := range extras {
		cursor := "  "
		if i == m.cursor {
			cursor = CursorStyle.Render("▸ ")
		}
		check := UncheckedStyle.Render("□")
		if m.extraSetup[e.key] {
			check = CheckedStyle.Render("■")
		}
		name := lipgloss.NewStyle().Foreground(White).Render(e.label)
		desc := DimStyle.Render("    " + e.desc)
		b.WriteString(fmt.Sprintf("  %s %s %s\n%s\n\n", cursor, check, name, desc))
	}

	return BoxStyle.Render(b.String())
}

func (m Model) renderSystemDefaults() string {
	var b strings.Builder
	b.WriteString(SubtitleStyle.Render("🖥  "+m.t.TitleSysDefault) + "\n\n")

	defaults := []struct {
		key   string
		label string
		desc  string
	}{
		{"browser_chrome", m.t.DefBrowser, m.t.DefBrowserDesc},
		{"editor_zed", m.t.DefEditor, m.t.DefEditorDesc},
		{"player_iina", m.t.DefPlayer, m.t.DefPlayerDesc},
	}

	for i, d := range defaults {
		cursor := "  "
		if i == m.cursor {
			cursor = CursorStyle.Render("▸ ")
		}
		check := UncheckedStyle.Render("□")
		if m.sysDefaults[d.key] {
			check = CheckedStyle.Render("■")
		}
		name := lipgloss.NewStyle().Foreground(White).Render(d.label)
		desc := DimStyle.Render("  " + d.desc)
		b.WriteString(fmt.Sprintf("  %s %s %s\n%s\n", cursor, check, name, desc))
	}

	return BoxStyle.Render(b.String())
}

func (m Model) renderDone() string {
	// count errors
	errCount := 0
	for _, entry := range m.installLog {
		if !entry.success {
			errCount++
		}
	}

	var done string
	if errCount == 0 {
		done = SuccessStyle.Render("  ✓ "+m.t.DoneReady) + "\n\n"
		done += "  " + m.t.DoneMsg + "\n"
	} else {
		done = SuccessStyle.Render("  ✓ "+m.t.DoneReady) + "\n\n"
		done += "  " + m.t.DoneMsg + "\n\n"
		done += ErrorStyle.Render(fmt.Sprintf("  ⚠ %d errors occurred.", errCount)) + "\n"
		done += DimStyle.Render("  Full error log: ~/.freshbox/install.log") + "\n"
	}
	done += "\n  " + m.t.DoneExit
	return BoxStyle.Render(done)
}

func (m Model) renderFooter() string {
	help := "  " + m.t.FooterNav
	if m.page == PageCodexConfig || m.page == PageClaudeConfig {
		help = "  " + m.t.FooterForm
	}
	return HelpStyle.Render(help)
}

func (m Model) renderLangSelect() string {
	var b strings.Builder
	b.WriteString(SubtitleStyle.Render("🌐 "+m.t.LangPrompt) + "\n\n")

	langs := []struct {
		label string
		desc  string
	}{
		{"English", "Use English as the interface language"},
		{"中文", "使用中文作为界面语言"},
	}

	for i, l := range langs {
		cursor := "  "
		if i == m.langCursor {
			cursor = CursorStyle.Render("▸ ")
		}
		radio := UncheckedStyle.Render("○")
		if i == m.langCursor {
			radio = CheckedStyle.Render("●")
		}
		name := lipgloss.NewStyle().Foreground(White).Render(l.label)
		desc := DimStyle.Render("  " + l.desc)
		b.WriteString(fmt.Sprintf("  %s %s %s\n%s\n\n", cursor, radio, name, desc))
	}

	return BoxStyle.Render(b.String())
}
