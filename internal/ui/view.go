package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kittors/freshbox/internal/checker"
	"github.com/kittors/freshbox/internal/version"
)

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

	// Constrain output to terminal size to prevent scroll overflow
	content := b.String()
	return lipgloss.Place(m.width, m.height, lipgloss.Left, lipgloss.Top, content)
}

func (m Model) renderHeader() string {
	logo := LogoStyle.Render("  ╭─────────────────────────────╮")
	logo += "\n" + LogoStyle.Render(fmt.Sprintf("  │   🍃 freshbox  %-12s │", version.Version))
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
