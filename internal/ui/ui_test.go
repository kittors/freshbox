package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kittors/freshbox/internal/checker"
)

// --- Model Creation ---

func TestNewModel(t *testing.T) {
	m := NewModel()

	if m.page != PageLang {
		t.Errorf("initial page = %d, want PageLang (%d)", m.page, PageLang)
	}
	if m.lang != LangEN {
		t.Errorf("initial lang = %d, want LangEN", m.lang)
	}
	if len(m.devTools) == 0 {
		t.Error("devTools should not be empty")
	}
	if len(m.apps) == 0 {
		t.Error("apps should not be empty")
	}
	if len(m.aiTools) == 0 {
		t.Error("aiTools should not be empty")
	}
	if len(m.mcps) == 0 {
		t.Error("mcps should not be empty")
	}
	if m.selected == nil {
		t.Error("selected map should be initialized")
	}
	if m.fnmSelected == nil {
		t.Error("fnmSelected map should be initialized")
	}
	if m.mcpSelected == nil {
		t.Error("mcpSelected map should be initialized")
	}

	// Verify system defaults are pre-selected
	if !m.sysDefaults["browser_chrome"] {
		t.Error("browser_chrome should be pre-selected")
	}
	if !m.sysDefaults["editor_zed"] {
		t.Error("editor_zed should be pre-selected")
	}
	if !m.sysDefaults["player_iina"] {
		t.Error("player_iina should be pre-selected")
	}

	// Verify extra setup is pre-selected
	if !m.extraSetup["zed_theme"] {
		t.Error("zed_theme should be pre-selected")
	}
	if !m.extraSetup["dev_workspace"] {
		t.Error("dev_workspace should be pre-selected")
	}

	// Verify first 4 MCPs are pre-selected
	for i, mcp := range m.mcps {
		if i < 4 && !m.mcpSelected[mcp.Name] {
			t.Errorf("MCP %s should be pre-selected", mcp.Name)
		}
	}
}

func TestModelInit(t *testing.T) {
	m := NewModel()
	cmd := m.Init()
	if cmd != nil {
		t.Error("Init should return nil")
	}
}

// --- Language Selection ---

func TestLangSelection(t *testing.T) {
	m := NewModel()

	// Start on language page
	if m.page != PageLang {
		t.Fatal("should start on PageLang")
	}

	// Move down to Chinese
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(Model)
	if m.langCursor != 1 {
		t.Errorf("langCursor = %d, want 1", m.langCursor)
	}

	// Move back up to English
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(Model)
	if m.langCursor != 0 {
		t.Errorf("langCursor = %d, want 0", m.langCursor)
	}

	// Can't go above 0
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(Model)
	if m.langCursor != 0 {
		t.Errorf("langCursor should stay at 0")
	}

	// Select English
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(Model)
	if m.page != PageWelcome {
		t.Errorf("page = %d, want PageWelcome after lang select", m.page)
	}
	if m.lang != LangEN {
		t.Error("lang should be LangEN")
	}
}

func TestLangSelectionChinese(t *testing.T) {
	m := NewModel()

	// Move to Chinese
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(Model)

	// Select with space
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeySpace})
	m = updated.(Model)
	if m.lang != LangZH {
		t.Error("lang should be LangZH")
	}
	if m.page != PageWelcome {
		t.Error("should advance to welcome page")
	}
}

// --- Navigation ---

func TestWelcomeToDevTools(t *testing.T) {
	m := createModelOnPage(PageWelcome)

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(Model)
	if m.page != PageDevTools {
		t.Errorf("page = %d, want PageDevTools", m.page)
	}
}

func TestQuitFromWelcome(t *testing.T) {
	m := createModelOnPage(PageWelcome)

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Error("q on welcome page should trigger Quit")
	}
}

func TestQuitFromDone(t *testing.T) {
	m := createModelOnPage(PageDone)

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Error("enter on done page should trigger Quit")
	}
}

func TestBackNavigation(t *testing.T) {
	m := createModelOnPage(PageApps)

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	m = updated.(Model)
	if m.page != PageDevTools {
		t.Errorf("shift+tab from Apps should go to DevTools, got %d", m.page)
	}
}

func TestTabForward(t *testing.T) {
	m := createModelOnPage(PageDevTools)

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(Model)
	if m.page != PageApps {
		t.Errorf("tab from DevTools should go to Apps, got %d", m.page)
	}
}

// --- Cursor Movement ---

func TestCursorUpDown(t *testing.T) {
	m := createModelOnPage(PageDevTools)
	m.cursor = 0

	// Move down
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(Model)
	if m.cursor != 1 {
		t.Errorf("cursor = %d, want 1", m.cursor)
	}

	// Move up
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(Model)
	if m.cursor != 0 {
		t.Errorf("cursor = %d, want 0", m.cursor)
	}

	// Can't go below 0
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(Model)
	if m.cursor != 0 {
		t.Error("cursor should not go below 0")
	}
}

func TestCursorVimKeys(t *testing.T) {
	m := createModelOnPage(PageDevTools)

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(Model)
	if m.cursor != 1 {
		t.Error("j should move cursor down")
	}

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = updated.(Model)
	if m.cursor != 0 {
		t.Error("k should move cursor up")
	}
}

// --- Selection ---

func TestToggleSelection(t *testing.T) {
	m := createModelOnPage(PageDevTools)

	// Find first not-installed item
	idx := -1
	for i, item := range m.devTools {
		if item.Status == checker.NotInstalled {
			idx = i
			break
		}
	}
	if idx == -1 {
		t.Skip("all dev tools are installed, cannot test toggle")
	}

	m.cursor = idx
	item := m.devTools[idx]
	wasSel := m.selected[item.Name]

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	m = updated.(Model)
	if m.selected[item.Name] == wasSel {
		t.Error("space should toggle selection")
	}
}

func TestSelectAllDevTools(t *testing.T) {
	m := createModelOnPage(PageDevTools)

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	m = updated.(Model)

	for _, item := range m.devTools {
		if item.Status == checker.NotInstalled && !m.selected[item.Name] {
			t.Errorf("%s should be selected after 'a'", item.Name)
		}
	}
}

func TestSelectNoneDevTools(t *testing.T) {
	m := createModelOnPage(PageDevTools)

	// First select all
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	m = updated.(Model)

	// Then select none
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	m = updated.(Model)

	for _, item := range m.devTools {
		if m.selected[item.Name] {
			t.Errorf("%s should be deselected after 'n'", item.Name)
		}
	}
}

func TestToggleMCP(t *testing.T) {
	m := createModelOnPage(PageMCP)
	m.cursor = 0

	name := m.mcps[0].Name
	wasSel := m.mcpSelected[name]

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	m = updated.(Model)
	if m.mcpSelected[name] == wasSel {
		t.Error("space should toggle MCP selection")
	}
}

func TestToggleExtraSetup(t *testing.T) {
	m := createModelOnPage(PageExtraSetup)
	m.cursor = 0

	wasSet := m.extraSetup["zed_theme"]
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	m = updated.(Model)
	if m.extraSetup["zed_theme"] == wasSet {
		t.Error("space should toggle extra setup")
	}
}

func TestToggleSystemDefaults(t *testing.T) {
	m := createModelOnPage(PageSystemDefaults)
	m.cursor = 0

	wasSet := m.sysDefaults["browser_chrome"]
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	m = updated.(Model)
	if m.sysDefaults["browser_chrome"] == wasSet {
		t.Error("space should toggle system defaults")
	}
}

// --- currentListLen ---

func TestCurrentListLen(t *testing.T) {
	m := NewModel()
	m.lang = LangEN
	m.t = GetText(LangEN)

	tests := []struct {
		page Page
		min  int
	}{
		{PageDevTools, 12},
		{PageApps, 7},
		{PageAITools, 2},
		{PageMCP, 11},
		{PageSystemDefaults, 3},
		{PageExtraSetup, 4},
		{PageWelcome, 1},
	}

	for _, tt := range tests {
		m.page = tt.page
		got := m.currentListLen()
		if got < tt.min {
			t.Errorf("page %d: listLen = %d, want >= %d", tt.page, got, tt.min)
		}
	}
}

// --- Window Size ---

func TestWindowSizeMsg(t *testing.T) {
	m := NewModel()
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = updated.(Model)
	if m.width != 120 || m.height != 40 {
		t.Errorf("width=%d, height=%d, want 120x40", m.width, m.height)
	}
}

// --- View Rendering ---

func TestViewLoadingState(t *testing.T) {
	m := NewModel()
	m.width = 0 // simulates no WindowSizeMsg yet
	view := m.View()
	if view != "Loading..." {
		t.Errorf("view = %q, want 'Loading...'", view)
	}
}

func TestViewRendersWithSize(t *testing.T) {
	m := NewModel()
	m.width = 80
	m.height = 24
	m.page = PageWelcome
	m.lang = LangEN
	m.t = GetText(LangEN)

	view := m.View()
	if len(view) == 0 {
		t.Error("view should not be empty")
	}
	if !strings.Contains(view, "freshbox") {
		t.Error("view should contain 'freshbox' branding")
	}
}

func TestRenderCheckList(t *testing.T) {
	m := createModelOnPage(PageDevTools)
	m.width = 80
	m.height = 40

	view := m.View()
	if !strings.Contains(view, "freshbox") {
		t.Error("view should contain 'freshbox'")
	}
}

// --- Config Form ---

func TestInitCodexInputs(t *testing.T) {
	m := NewModel()
	m.initCodexInputs()

	if len(m.inputs) != 4 {
		t.Errorf("codex inputs = %d, want 4", len(m.inputs))
	}
	if m.inputPage != PageCodexConfig {
		t.Errorf("inputPage = %d, want PageCodexConfig", m.inputPage)
	}
	if m.inputFocus != 0 {
		t.Error("inputFocus should start at 0")
	}
	// Verify defaults
	if m.inputs[0].Value() != "o4-mini" {
		t.Errorf("default model = %q", m.inputs[0].Value())
	}
	if m.inputs[1].Value() != "medium" {
		t.Errorf("default thinking = %q", m.inputs[1].Value())
	}
}

func TestInitClaudeInputs(t *testing.T) {
	m := NewModel()
	m.initClaudeInputs()

	if len(m.inputs) != 3 {
		t.Errorf("claude inputs = %d, want 3", len(m.inputs))
	}
	if m.inputPage != PageClaudeConfig {
		t.Errorf("inputPage = %d, want PageClaudeConfig", m.inputPage)
	}
	if m.inputs[0].Value() != "claude-sonnet-4-6" {
		t.Errorf("default model = %q", m.inputs[0].Value())
	}
}

func TestSaveCodexInputs(t *testing.T) {
	m := NewModel()
	m.initCodexInputs()
	m.inputs[0].SetValue("gpt-4")
	m.inputs[1].SetValue("high")
	m.inputs[2].SetValue("https://custom.api")
	m.inputs[3].SetValue("sk-key")

	m.saveCodexInputs()
	if m.codexModel != "gpt-4" {
		t.Errorf("codexModel = %q", m.codexModel)
	}
	if m.codexThink != "high" {
		t.Errorf("codexThink = %q", m.codexThink)
	}
	if m.codexURL != "https://custom.api" {
		t.Errorf("codexURL = %q", m.codexURL)
	}
	if m.codexKey != "sk-key" {
		t.Errorf("codexKey = %q", m.codexKey)
	}
}

func TestSaveClaudeInputs(t *testing.T) {
	m := NewModel()
	m.initClaudeInputs()
	m.inputs[0].SetValue("claude-3")
	m.inputs[1].SetValue("https://custom.url")
	m.inputs[2].SetValue("sk-ant-test")

	m.saveClaudeInputs()
	if m.claudeModel != "claude-3" {
		t.Errorf("claudeModel = %q", m.claudeModel)
	}
	if m.claudeURL != "https://custom.url" {
		t.Errorf("claudeURL = %q", m.claudeURL)
	}
	if m.claudeKey != "sk-ant-test" {
		t.Errorf("claudeKey = %q", m.claudeKey)
	}
}

// --- Install Queue ---

func TestBuildInstallQueue_Empty(t *testing.T) {
	m := NewModel()
	// Deselect everything
	for k := range m.selected {
		m.selected[k] = false
	}
	for k := range m.mcpSelected {
		m.mcpSelected[k] = false
	}
	for k := range m.sysDefaults {
		m.sysDefaults[k] = false
	}
	for k := range m.extraSetup {
		m.extraSetup[k] = false
	}

	queue := m.buildInstallQueue()
	if len(queue) != 0 {
		t.Errorf("empty selection should produce empty queue, got %d items", len(queue))
	}
}

func TestBuildInstallQueue_WithSysDefaults(t *testing.T) {
	m := NewModel()
	// Deselect all except system defaults
	for k := range m.selected {
		m.selected[k] = false
	}
	for k := range m.mcpSelected {
		m.mcpSelected[k] = false
	}
	for k := range m.extraSetup {
		m.extraSetup[k] = false
	}
	m.sysDefaults = map[string]bool{
		"browser_chrome": true,
		"editor_zed":     false,
		"player_iina":    false,
	}

	queue := m.buildInstallQueue()
	found := false
	for _, task := range queue {
		if strings.Contains(task.name, "Chrome") {
			found = true
		}
	}
	if !found {
		t.Error("queue should contain Chrome default browser task")
	}
}

func TestBuildInstallQueue_WithExtraSetup(t *testing.T) {
	m := NewModel()
	for k := range m.selected {
		m.selected[k] = false
	}
	for k := range m.mcpSelected {
		m.mcpSelected[k] = false
	}
	for k := range m.sysDefaults {
		m.sysDefaults[k] = false
	}
	m.extraSetup = map[string]bool{
		"zed_theme":      true,
		"kaku_init":      false,
		"karabiner_kaku": false,
		"dev_workspace":  true,
	}

	queue := m.buildInstallQueue()
	names := []string{}
	for _, task := range queue {
		names = append(names, task.name)
	}

	foundZed := false
	foundDev := false
	for _, n := range names {
		if strings.Contains(n, "Zed") {
			foundZed = true
		}
		if strings.Contains(n, "Developer Workspace") || strings.Contains(n, "Finder") {
			foundDev = true
		}
	}
	if !foundZed {
		t.Error("queue should contain Zed theme task")
	}
	if !foundDev {
		t.Error("queue should contain dev workspace task")
	}
}

// --- i18n ---

func TestGetText(t *testing.T) {
	en := GetText(LangEN)
	zh := GetText(LangZH)

	if en.WelcomeTitle == "" {
		t.Error("EN WelcomeTitle is empty")
	}
	if zh.WelcomeTitle == "" {
		t.Error("ZH WelcomeTitle is empty")
	}
	if en.WelcomeTitle == zh.WelcomeTitle {
		t.Error("EN and ZH should have different WelcomeTitle")
	}

	// Verify all page names are non-empty
	enPages := []string{en.PageWelcome, en.PageDevTools, en.PageApps, en.PageAITools, en.PageMCP, en.PageDone}
	for _, p := range enPages {
		if p == "" {
			t.Error("empty EN page name found")
		}
	}
}

func TestPageNames(t *testing.T) {
	en := GetText(LangEN)
	names := pageNames(en)
	if len(names) != 13 {
		t.Errorf("pageNames returned %d items, want 13", len(names))
	}
	for i, name := range names {
		if name == "" {
			t.Errorf("pageNames[%d] is empty", i)
		}
	}
}

// --- Page Constants ---

func TestPageConstants(t *testing.T) {
	pages := []Page{
		PageLang, PageWelcome, PageDevTools, PageApps, PageFnmVersions,
		PageAITools, PageCodexConfig, PageClaudeConfig, PageMCP,
		PageExtraSetup, PageSystemDefaults, PageInstalling, PageDone,
	}

	// Verify they are sequential
	for i, p := range pages {
		if int(p) != i {
			t.Errorf("Page %d has value %d, want %d", i, int(p), i)
		}
	}
}

// --- Install Done Message ---

func TestInstallDoneMsg(t *testing.T) {
	m := NewModel()
	m.page = PageInstalling
	m.installing = true

	updated, _ := m.Update(installDoneMsg{})
	m = updated.(Model)
	if m.installing {
		t.Error("installing should be false after installDoneMsg")
	}
	if !m.installDone {
		t.Error("installDone should be true")
	}
	if m.page != PageDone {
		t.Errorf("page = %d, want PageDone", m.page)
	}
}

// --- Spinner ---

func TestNewSpinner(t *testing.T) {
	s := NewSpinner()
	if len(s.Spinner.Frames) == 0 {
		t.Error("spinner should have frames")
	}
}

// --- Helper ---

func createModelOnPage(p Page) Model {
	m := NewModel()
	m.page = p
	m.lang = LangEN
	m.t = GetText(LangEN)
	m.width = 80
	m.height = 40
	m.cursor = 0
	return m
}
