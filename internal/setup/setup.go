package setup

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// --- Zed Catppuccin Blur Theme ---

// SetupZedTheme clones catppuccin-blur, applies blue tint, configures Zed settings
func SetupZedTheme() error {
	home, _ := os.UserHomeDir()
	themeDir := filepath.Join(home, ".config", "zed", "themes")
	themeFile := filepath.Join(themeDir, "catppuccin-blur.json")
	settingsFile := filepath.Join(home, ".config", "zed", "settings.json")

	os.MkdirAll(themeDir, 0755)

	// Clone theme repo to temp dir
	tmpDir, err := os.MkdirTemp("", "zed-theme-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	cmd := exec.Command("git", "clone", "--depth", "1", "--quiet",
		"https://github.com/jenslys/zed-catppuccin-blur.git", filepath.Join(tmpDir, "repo"))
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("clone theme: %s %w", string(out), err)
	}

	// Copy theme file
	src := filepath.Join(tmpDir, "repo", "themes", "catppuccin-blur.json")
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("read theme: %w", err)
	}
	if err := os.WriteFile(themeFile, data, 0644); err != nil {
		return fmt.Errorf("write theme: %w", err)
	}

	// Apply blue tint customizations + higher opacity for readability via python3
	pyScript := `
import json, os

theme_file = os.environ["THEME_FILE"]
with open(theme_file, "r") as f:
    data = json.load(f)

for theme in data["themes"]:
    name = theme["name"]
    style = theme["style"]
    if name == "Catppuccin Latte (Blur) [Light]":
        for key, val in {
            "elevated_surface.background": "#e8f0ff",
            "surface.background": "#e8f0ffc8",
            "background": "#e8f0ffd0",
            "status_bar.background": "#e8f0ffd0",
            "title_bar.background": "#e8f0ffd0",
            "tab.active_background": "#e8f0ffc0",
            "ghost_element.background": "#e8f0ff90",
            "ghost_element.hover": "#e8f0ffc0",
            "panel.overlay_background": "#e8f0ff",
        }.items():
            if key in style:
                style[key] = val
    elif name == "Catppuccin Mocha (Blur) [Light]":
        for key, val in {
            "elevated_surface.background": "#161a28",
            "surface.background": "#181c2ec8",
            "background": "#181c2ed0",
            "status_bar.background": "#181c2ed0",
            "title_bar.background": "#181c2ed0",
            "title_bar.inactive_background": "#151928",
            "tab.active_background": "#161a28c0",
            "ghost_element.background": "#161a2890",
            "ghost_element.hover": "#161a28c0",
            "panel.overlay_background": "#181c2e",
        }.items():
            if key in style:
                style[key] = val

# Also increase opacity on all other [Light] variants for consistency
opacity_map = {"99": "d0", "8c": "c8", "90": "c0", "60": "90"}
bg_keys = ["background", "surface.background", "status_bar.background",
           "title_bar.background", "tab.active_background",
           "ghost_element.background", "ghost_element.hover"]
for theme in data["themes"]:
    if not theme["name"].endswith("[Light]"):
        continue
    if theme["name"] in ("Catppuccin Latte (Blur) [Light]", "Catppuccin Mocha (Blur) [Light]"):
        continue
    style = theme["style"]
    for key in bg_keys:
        if key not in style:
            continue
        val = style[key]
        if len(val) == 9 and val.startswith("#"):
            old_alpha = val[7:9]
            if old_alpha in opacity_map:
                style[key] = val[:7] + opacity_map[old_alpha]

with open(theme_file, "w") as f:
    json.dump(data, f, indent=2, ensure_ascii=False)
`
	pyCmd := exec.Command("python3", "-c", pyScript)
	pyCmd.Env = append(os.Environ(), "THEME_FILE="+themeFile)
	if out, err := pyCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("apply tint: %s %w", string(out), err)
	}

	// Configure Zed settings
	themeConfig := map[string]any{
		"mode":  "system",
		"light": "Catppuccin Latte (Blur) [Light]",
		"dark":  "Catppuccin Mocha (Blur) [Light]",
	}

	if _, err := os.Stat(settingsFile); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(settingsFile), 0755)
		settings := map[string]any{"theme": themeConfig}
		d, _ := json.MarshalIndent(settings, "", "  ")
		return os.WriteFile(settingsFile, append(d, '\n'), 0644)
	}

	// Update existing settings
	pyUpdate := `
import json, re, os

settings_path = os.environ["SETTINGS_FILE"]
with open(settings_path, "r") as f:
    content = f.read()

content = re.sub(r"//.*$", "", content, flags=re.MULTILINE)
content = re.sub(r",\s*([}\]])", r"\1", content)

data = json.loads(content)
data["theme"] = {
    "mode": "system",
    "light": "Catppuccin Latte (Blur) [Light]",
    "dark": "Catppuccin Mocha (Blur) [Light]"
}

with open(settings_path, "w") as f:
    json.dump(data, f, indent=2, ensure_ascii=False)
    f.write("\n")
`
	pyCmd2 := exec.Command("python3", "-c", pyUpdate)
	pyCmd2.Env = append(os.Environ(), "SETTINGS_FILE="+settingsFile)
	if out, err := pyCmd2.CombinedOutput(); err != nil {
		return fmt.Errorf("update settings: %s %w", string(out), err)
	}

	return nil
}

// --- Kaku Terminal ---

// SetupKaku initializes Kaku config and installs zsh plugins
// Note: Kaku app installation is handled separately via brew (tw93/tap/kakuku)
func SetupKaku() error {
	home, _ := os.UserHomeDir()
	kakuDir := filepath.Join(home, ".config", "kaku")
	pluginDir := filepath.Join(kakuDir, "zsh", "plugins")
	os.MkdirAll(pluginDir, 0755)

	// Write kaku.lua config
	kakuLua := `local wezterm = require 'wezterm'

local function resolve_bundled_config()
  local resource_dir = wezterm.executable_dir:gsub('MacOS/?$', 'Resources')
  local bundled = resource_dir .. '/kaku.lua'
  local f = io.open(bundled, 'r')
  if f then
    f:close()
    return bundled
  end

  local dev_bundled = wezterm.executable_dir .. '/../../assets/macos/Kaku.app/Contents/Resources/kaku.lua'
  f = io.open(dev_bundled, 'r')
  if f then
    f:close()
    return dev_bundled
  end

  local app_bundled = '/Applications/Kaku.app/Contents/Resources/kaku.lua'
  f = io.open(app_bundled, 'r')
  if f then
    f:close()
    return app_bundled
  end

  local home = os.getenv('HOME') or ''
  local home_bundled = home .. '/Applications/Kaku.app/Contents/Resources/kaku.lua'
  f = io.open(home_bundled, 'r')
  if f then
    f:close()
    return home_bundled
  end

  return nil
end

local config = {}
local bundled = resolve_bundled_config()

if bundled then
  local ok, loaded = pcall(dofile, bundled)
  if ok and type(loaded) == 'table' then
    config = loaded
  else
    wezterm.log_error('Kaku: failed to load bundled defaults from ' .. bundled)
  end
else
  wezterm.log_error('Kaku: bundled defaults not found')
end

return config
`
	if err := os.WriteFile(filepath.Join(kakuDir, "kaku.lua"), []byte(kakuLua), 0644); err != nil {
		return fmt.Errorf("write kaku.lua: %w", err)
	}

	// Install zsh plugins
	plugins := map[string]string{
		"zsh-autosuggestions":     "https://github.com/zsh-users/zsh-autosuggestions.git",
		"zsh-completions":         "https://github.com/zsh-users/zsh-completions.git",
		"zsh-syntax-highlighting": "https://github.com/zsh-users/zsh-syntax-highlighting.git",
		"zsh-z":                   "https://github.com/agkozak/zsh-z.git",
	}

	for name, repo := range plugins {
		dest := filepath.Join(pluginDir, name)
		if _, err := os.Stat(dest); err == nil {
			continue // already exists
		}
		cmd := exec.Command("git", "clone", "--depth", "1", "--quiet", repo, dest)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("clone %s: %s %w", name, string(out), err)
		}
	}

	return nil
}

// --- Karabiner Elements ---

// SetupKarabiner installs Karabiner-Elements and configures Ctrl+Opt+Cmd+T to open Kaku
func SetupKarabiner() error {
	// Install via brew cask
	cmd := exec.Command("brew", "install", "--cask", "karabiner-elements")
	if out, err := cmd.CombinedOutput(); err != nil {
		if !strings.Contains(string(out), "already installed") {
			return fmt.Errorf("install karabiner: %s %w", string(out), err)
		}
	}

	home, _ := os.UserHomeDir()

	// Create open-kaku.sh script
	binDir := filepath.Join(home, ".local", "bin")
	os.MkdirAll(binDir, 0755)

	openKakuScript := `#!/bin/bash
# Get the selected folder in Finder, or Finder's current directory
DIR=$(osascript -e '
tell application "System Events"
    set frontApp to name of first application process whose frontmost is true
end tell
if frontApp is "Finder" then
    tell application "Finder"
        try
            set sel to selection
            if (count of sel) > 0 then
                set theItem to item 1 of sel
                if class of theItem is folder or class of theItem is disk then
                    return POSIX path of (theItem as alias)
                else
                    return POSIX path of (container of theItem as alias)
                end if
            else
                return POSIX path of (target of front window as alias)
            end if
        on error
            return POSIX path of (path to home folder)
        end try
    end tell
else
    return ""
end if
' 2>/dev/null)

if [ -n "$DIR" ] && [ -d "$DIR" ]; then
    /opt/homebrew/bin/kaku start --cwd "$DIR"
else
    open /Applications/Kaku.app
fi
`
	scriptPath := filepath.Join(binDir, "open-kaku.sh")
	if err := os.WriteFile(scriptPath, []byte(openKakuScript), 0755); err != nil {
		return fmt.Errorf("write open-kaku.sh: %w", err)
	}

	// Write karabiner config
	karabinerDir := filepath.Join(home, ".config", "karabiner")
	os.MkdirAll(karabinerDir, 0755)

	karabinerConfig := map[string]any{
		"global": map[string]any{"show_in_menu_bar": false},
		"profiles": []map[string]any{
			{
				"complex_modifications": map[string]any{
					"rules": []map[string]any{
						{
							"description": "Control+Option+Command+T opens Kaku",
							"manipulators": []map[string]any{
								{
									"from": map[string]any{
										"key_code": "t",
										"modifiers": map[string]any{
											"mandatory": []string{"control", "option", "command"},
										},
									},
									"to":   []map[string]any{{"shell_command": scriptPath}},
									"type": "basic",
								},
							},
						},
					},
				},
				"name":     "Default profile",
				"selected": true,
				"virtual_hid_keyboard": map[string]any{
					"keyboard_type_v2": "ansi",
				},
			},
		},
	}

	// If existing config, merge our rule into it
	existingPath := filepath.Join(karabinerDir, "karabiner.json")
	if existData, err := os.ReadFile(existingPath); err == nil {
		var existing map[string]any
		if json.Unmarshal(existData, &existing) == nil {
			// Use existing config, just ensure our rule is present
			karabinerConfig = existing
		}
	}

	data, _ := json.MarshalIndent(karabinerConfig, "", "    ")
	return os.WriteFile(filepath.Join(karabinerDir, "karabiner.json"), append(data, '\n'), 0644)
}

// --- macOS Dev Workspace ---

// SetupDevWorkspace creates the developer directory structure and configures Finder
func SetupDevWorkspace() error {
	home, _ := os.UserHomeDir()
	devDir := filepath.Join(home, "Developer")

	dirs := []string{
		"opensource",
		"boundless",
		"freelance/_template",
		"playground",
		"design",
		"notes",
		"scripts",
		"archive",
	}

	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(devDir, d), 0755); err != nil {
			return fmt.Errorf("create dir %s: %w", d, err)
		}
	}

	// Write root README
	rootReadme := `# ~/Developer 目录体系

> kittors 的 macOS 开发工作区，统一管理所有项目、资源和工具。

## 目录结构

` + "```" + `
~/Developer/
├── opensource/      个人开源项目
├── boundless/       无境科技公司项目
├── freelance/       自由职业 / 外包项目
├── playground/      学习、Demo、实验性代码
├── design/          UI 设计稿、图标、素材
├── notes/           技术笔记、文档、博客草稿
├── scripts/         自动化脚本、CLI 工具
└── archive/         已完成或不再维护的项目归档
` + "```" + `

## 使用规范

1. 新项目根据类型放入对应目录，不要堆在根目录
2. 项目目录名使用小写 + 连字符（kebab-case），如 ` + "`my-awesome-project`" + `
3. 项目不再活跃时移入 ` + "`archive/`" + `，保持工作区整洁
4. 外包项目统一放 ` + "`freelance/客户名/项目名`" + `
5. 临时实验代码放 ` + "`playground/`" + `，避免污染正式项目目录
`
	if err := os.WriteFile(filepath.Join(devDir, "README.md"), []byte(rootReadme), 0644); err != nil {
		return fmt.Errorf("write root README: %w", err)
	}

	// Write sub-directory READMEs
	readmes := map[string]string{
		"opensource":          "# opensource\n\n个人开源项目目录。存放在 GitHub 上公开维护的项目。\n",
		"boundless":           "# boundless\n\n无境科技（Boundless Tech）公司项目目录。\n",
		"freelance":           "# freelance\n\n自由职业 / 外包项目目录。按客户名称组织。\n",
		"freelance/_template": "# freelance/_template\n\n外包项目模板目录。新建客户项目时可从此复制初始结构。\n",
		"playground":          "# playground\n\n学习、Demo、教程跟练、实验性代码。\n",
		"design":              "# design\n\nUI 设计稿、图标、素材资源。\n",
		"notes":               "# notes\n\n技术笔记、文档、博客草稿。\n",
		"scripts":             "# scripts\n\n自动化脚本、CLI 工具、dotfiles 配置。\n",
		"archive":             "# archive\n\n已完成或不再维护的项目归档。\n",
	}

	for subdir, content := range readmes {
		if err := os.WriteFile(filepath.Join(devDir, subdir, "README.md"), []byte(content), 0644); err != nil {
			return fmt.Errorf("write %s README: %w", subdir, err)
		}
	}

	// Configure Finder
	finderCmds := [][]string{
		{"defaults", "write", "com.apple.finder", "AppleShowAllFiles", "-bool", "true"},
		{"defaults", "write", "NSGlobalDomain", "AppleShowAllExtensions", "-bool", "true"},
		{"defaults", "write", "com.apple.finder", "ShowPathbar", "-bool", "true"},
		{"defaults", "write", "com.apple.finder", "ShowStatusBar", "-bool", "true"},
		{"defaults", "write", "com.apple.finder", "FXPreferredViewStyle", "-string", "Nlsv"},
		{"defaults", "write", "com.apple.finder", "FXDefaultSearchScope", "-string", "SCcf"},
		{"defaults", "write", "com.apple.finder", "FXEnableExtensionChangeWarning", "-bool", "false"},
		{"defaults", "write", "com.apple.finder", "NewWindowTarget", "-string", "PfLo"},
		{"defaults", "write", "com.apple.finder", "NewWindowTargetPath", "-string", "file://" + devDir + "/"},
	}

	for _, args := range finderCmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Run()
	}

	// Restart Finder
	exec.Command("killall", "Finder").Run()

	return nil
}
