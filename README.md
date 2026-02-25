<p align="center">
  <img src="https://img.shields.io/badge/platform-macOS-blue?style=flat-square&logo=apple" />
  <img src="https://img.shields.io/badge/language-Go-00ADD8?style=flat-square&logo=go" />
  <img src="https://img.shields.io/badge/TUI-Bubbletea-ff69b4?style=flat-square" />
  <img src="https://img.shields.io/badge/license-MIT-green?style=flat-square" />
</p>

<h1 align="center">🍃 freshbox</h1>

<p align="center">
  <strong>A beautiful TUI for setting up a fresh macOS machine in minutes.</strong><br/>
  <sub>一个漂亮的终端界面工具，帮你几分钟内配置好全新的 Mac。</sub>
</p>

<p align="center">
  <a href="#-quick-install">Quick Install</a> •
  <a href="#-features">Features</a> •
  <a href="#-what-gets-installed">What Gets Installed</a> •
  <a href="#-usage">Usage</a> •
  <a href="#-中文说明">中文说明</a> •
  <a href="#-credits">Credits</a>
</p>

---

## Preview

```
  ╭─────────────────────────────╮
  │   🍃 freshbox  v1.0.0      │
  │   macOS Setup Assistant     │
  ╰─────────────────────────────╯

  🔧 Development Tools

    ▸ ■ ̶H̶o̶m̶e̶b̶r̶e̶w̶ (4.4.x) — macOS package manager
      ■ ̶G̶i̶t̶ (2.x.x) — Distributed version control system
      □ Java (JDK) — Java development kit for JVM-based development
      □ uv — Ultra-fast Python package manager by Astral
      □ fnm — Fast Node.js version manager written in Rust
      □ Rust (rustup) — Systems programming language with memory safety
      ■ ̶G̶o̶ (1.26.0) — Statically typed language by Google

  📦 Applications

      ■ Zed — High-performance code editor by the Atom creators
      □ Kaku — Lightweight terminal app built on WezTerm by tw93
      □ Mole — macOS system cleaner to free up disk space by tw93
```

## 🚀 Quick Install

One command to install freshbox:

```bash
curl -fsSL https://raw.githubusercontent.com/kittors/freshbox/main/install.sh | bash
```

Or build from source:

```bash
git clone https://github.com/kittors/freshbox.git
cd freshbox
go build -o freshbox .
./freshbox
```

## ✨ Features

- **🌐 Bilingual** — English / 中文 interface, selected at startup
- **🔧 Smart Detection** — Auto-detects installed tools, shows versions with ~~strikethrough~~
- **📦 Node.js Versions** — Multi-select Node.js versions to install via fnm
- **📱 App Installer** — One-click install for essential macOS apps
- **🤖 AI Tools Config** — Full setup for Codex and Claude Code with config file generation
- **🔌 MCP Servers** — Select from 11 popular MCP servers to configure
- **🎨 Extra Setup** — Zed theme, Kaku terminal, Karabiner shortcuts, dev workspace
- **🖥 System Defaults** — Set default browser, editor, media player
- **✨ Beautiful TUI** — Rounded borders, colors, smooth multi-page navigation

## 📦 What Gets Installed

### Development Tools

| Tool | Description | Install Method |
|------|-------------|----------------|
| [Homebrew](https://brew.sh/) | macOS package manager, the foundation for everything else | Official install script |
| [Git](https://git-scm.com/) | Distributed version control system | `brew install git` |
| [Java (OpenJDK)](https://openjdk.org/) | Java development kit for JVM-based development | `brew install openjdk` |
| [Maven](https://maven.apache.org/) | Java project build and dependency management | `brew install maven` |
| [Gradle](https://gradle.org/) | Flexible build automation tool for JVM projects | `brew install gradle` |
| [Python](https://www.python.org/) | General-purpose programming language | `brew install python` |
| [uv](https://github.com/astral-sh/uv) | Ultra-fast Python package manager by Astral | `brew install uv` |
| [fnm](https://github.com/Schniz/fnm) | Fast Node.js version manager written in Rust | `brew install fnm` |
| [Rust](https://www.rust-lang.org/) | Systems programming language with memory safety | `rustup` official installer |
| [Go](https://go.dev/) | Statically typed language by Google for scalable systems | `brew install go` |

### Applications

| App | Description | Install Method |
|-----|-------------|----------------|
| [Google Chrome](https://www.google.com/chrome/) | Web browser by Google | `brew install --cask google-chrome` |
| [Zed](https://zed.dev/) | High-performance code editor by the Atom creators | `brew install --cask zed` |
| [IINA](https://iina.io/) | Modern media player for macOS | `brew install --cask iina` |
| [Kaku](https://github.com/tw93/Kaku) | Lightweight terminal app built on WezTerm by tw93 | `brew install --cask kaku` |
| [Karabiner-Elements](https://karabiner-elements.pqrs.org/) | Powerful keyboard customizer for macOS | `brew install --cask karabiner-elements` |
| [Mole](https://github.com/tw93/Mole) | macOS system cleaner to free up disk space by tw93 | `brew install --cask mole` |

### AI Tools

| Tool | Description | Install Method |
|------|-------------|----------------|
| [Codex](https://github.com/openai/codex) | OpenAI's AI coding assistant CLI | `npm install -g @openai/codex` |
| [Claude Code](https://github.com/anthropics/claude-code) | Anthropic's AI coding assistant CLI | `npm install -g @anthropic-ai/claude-code` |

### MCP Servers (Selectable)

| Server | Description |
|--------|-------------|
| [Playwright](https://github.com/anthropics/anthropic-cookbook) | Browser automation and testing |
| [Context7](https://github.com/anthropics/anthropic-cookbook) | Contextual code understanding |
| [Filesystem](https://github.com/anthropics/anthropic-cookbook) | Local file system access |
| [GitHub](https://github.com/anthropics/anthropic-cookbook) | GitHub API integration |
| [Memory](https://github.com/anthropics/anthropic-cookbook) | Persistent memory across sessions |
| [Sequential Thinking](https://github.com/anthropics/anthropic-cookbook) | Step-by-step reasoning |
| [Fetch](https://github.com/anthropics/anthropic-cookbook) | HTTP request capabilities |
| [Brave Search](https://github.com/anthropics/anthropic-cookbook) | Web search via Brave |
| [Slack](https://github.com/anthropics/anthropic-cookbook) | Slack workspace integration |
| [Google Maps](https://github.com/anthropics/anthropic-cookbook) | Location and maps API |
| [SQLite](https://github.com/anthropics/anthropic-cookbook) | Local SQLite database access |

### Extra Setup (What the Scripts Do)

#### 🎨 Zed Catppuccin Blur Theme

Clones the [catppuccin-blur](https://github.com/jenslys/zed-catppuccin-blur) theme, applies a custom icy blue tint to both light and dark variants, and configures Zed to auto-switch based on system appearance:
- **Light mode** → Catppuccin Latte with `#e8f0ff` blue tint
- **Dark mode** → Catppuccin Mocha with `#181c2e` blue tint
- Writes to `~/.config/zed/themes/catppuccin-blur.json` and `~/.config/zed/settings.json`

#### 🐚 Kaku Terminal Setup

Initializes [Kaku](https://github.com/tw93/Kaku) (a lightweight terminal by tw93) with a full config and clones 4 essential zsh plugins:
- [zsh-autosuggestions](https://github.com/zsh-users/zsh-autosuggestions) — Fish-like command auto-completion
- [zsh-completions](https://github.com/zsh-users/zsh-completions) — Additional completion definitions
- [zsh-syntax-highlighting](https://github.com/zsh-users/zsh-syntax-highlighting) — Real-time command syntax coloring
- [zsh-z](https://github.com/agkozak/zsh-z) — Fast directory jumping based on frecency
- Writes to `~/.config/kaku/kaku.lua` and `~/.config/kaku/zsh/plugins/`

#### ⌨️ Karabiner ⌃⌥⌘T → Kaku

Sets up [Karabiner-Elements](https://karabiner-elements.pqrs.org/) with a keyboard shortcut to quick-launch Kaku:
- **Shortcut**: `Control + Option + Command + T`
- **Behavior**: If Finder is active, opens Kaku in the selected folder's directory. Otherwise opens Kaku normally.
- Creates `~/.local/bin/open-kaku.sh` (the launcher script) and `~/.config/karabiner/karabiner.json`

#### 📁 Developer Workspace

Creates a standardized `~/Developer` directory structure for organizing all projects:

```
~/Developer/
├── opensource/      Personal open-source projects
├── boundless/       Company projects
├── freelance/       Freelance / contract work
├── playground/      Learning, demos, experiments
├── design/          UI designs, icons, assets
├── notes/           Technical notes, docs, blog drafts
├── scripts/         Automation scripts, CLI tools
└── archive/         Completed / archived projects
```

Also configures Finder via `defaults write`:
- Show hidden files and file extensions
- Show path bar and status bar
- Default to list view
- Search current folder only
- New Finder windows open `~/Developer`
- Restarts Finder to apply changes

### Generated Config Files

<details>
<summary>Codex — <code>~/.codex/config.toml</code> + <code>auth.json</code></summary>

```toml
# config.toml
model = "o4-mini"
thinking_level = "medium"
base_url = "https://api.openai.com/v1"
```

```json
// auth.json
{
  "api_key": "sk-..."
}
```
</details>

<details>
<summary>Claude Code — <code>~/.claude/settings.json</code> + <code>mcp_servers.json</code></summary>

```json
// settings.json
{
  "model": "claude-sonnet-4-6",
  "env": {
    "ANTHROPIC_API_KEY": "sk-ant-...",
    "ANTHROPIC_BASE_URL": "https://api.anthropic.com"
  }
}
```

```json
// settings.json (mcpServers merged in)
{
  "model": "claude-sonnet-4-6",
  "mcpServers": {
    "Playwright": {
      "command": "npx",
      "args": ["-y", "@playwright/mcp", "--headless"]
    }
  }
}
```
</details>

## 🎮 Usage

```bash
freshbox
```

| Key | Action |
|-----|--------|
| `↑` `↓` / `j` `k` | Navigate items |
| `Space` | Toggle selection |
| `a` | Select all |
| `n` | Deselect all |
| `Tab` / `Enter` | Next page |
| `Shift+Tab` | Previous page |
| `q` | Quit / Go back |

### Workflow

```
🌐 Language → 👋 Welcome → 🔧 Dev Tools → 📦 Apps → 📦 Node.js
  → 🤖 AI Tools → ⚙️ Codex Config → ⚙️ Claude Config
  → 🔌 MCP Servers → 🎨 Extra Setup → 🖥 System Defaults
  → ⏳ Installing → ✅ Done!
```

## 📁 Project Structure

```
freshbox/
├── main.go                          # Entry point
├── install.sh                       # curl-based quick installer
├── internal/
│   ├── checker/checker.go           # System detection with descriptions
│   ├── installer/installer.go       # Install logic (brew/rustup/npm/fnm)
│   ├── config/config.go             # AI tool config generation (Codex/Claude/MCP)
│   ├── setup/setup.go               # Zed theme, Kaku init, Karabiner, dev workspace
│   └── ui/
│       ├── i18n.go                  # Bilingual text (EN/ZH)
│       ├── styles.go                # Lipgloss styles
│       ├── model.go                 # Bubbletea multi-page TUI (13 pages)
│       └── install.go               # Async install queue with progress
├── go.mod
└── go.sum
```

---

## 🇨🇳 中文说明

**freshbox** 是一个 macOS 装机 TUI 工具，帮你在全新的 Mac 上一键配置开发环境。

### 快速安装

```bash
curl -fsSL https://raw.githubusercontent.com/kittors/freshbox/main/install.sh | bash
```

### 功能

- 🌐 中英文双语界面，启动时选择
- 🔧 自动检测已安装的开发工具，显示版本号（已安装的划删除线），每个工具附带一句话介绍
- 📦 通过 [fnm](https://github.com/Schniz/fnm) 安装和管理多个 Node.js 版本
- 📱 一键安装常用软件：[Chrome](https://www.google.com/chrome/)、[Zed](https://zed.dev/)、[IINA](https://iina.io/)、[Kaku](https://github.com/tw93/Kaku)、[Karabiner](https://karabiner-elements.pqrs.org/)、[Mole](https://github.com/tw93/Mole)
- 🤖 配置 AI 开发工具（[Codex](https://github.com/openai/codex)、[Claude Code](https://github.com/anthropics/claude-code)），自动生成配置文件
- 🔌 勾选配置 11 个流行的 MCP 服务
- 🎨 额外配置：
  - Zed [Catppuccin Blur](https://github.com/jenslys/zed-catppuccin-blur) 冰蓝主题，自动跟随系统明暗
  - [Kaku](https://github.com/tw93/Kaku) 终端完整初始化 + 4 个 zsh 插件
  - [Karabiner](https://karabiner-elements.pqrs.org/) 快捷键 `⌃⌥⌘T` 快速启动 Kaku
  - `~/Developer` 开发工作区目录结构 + Finder 定制化
- 🖥 设置系统默认浏览器、编辑器、播放器

### 操作方式

| 按键 | 操作 |
|------|------|
| `↑` `↓` / `j` `k` | 上下移动 |
| `空格` | 切换选中 |
| `a` | 全选 |
| `n` | 全不选 |
| `Tab` / `Enter` | 下一页 |
| `Shift+Tab` | 上一页 |
| `q` | 退出 / 返回 |

---

## 🙏 Credits

freshbox is built with and installs tools from these amazing open-source projects:

| Project | Author | Description |
|---------|--------|-------------|
| [Bubbletea](https://github.com/charmbracelet/bubbletea) | [Charm](https://github.com/charmbracelet) | TUI framework for Go |
| [Lipgloss](https://github.com/charmbracelet/lipgloss) | [Charm](https://github.com/charmbracelet) | Style definitions for terminal UIs |
| [Homebrew](https://brew.sh/) | [Homebrew](https://github.com/Homebrew) | The missing package manager for macOS |
| [fnm](https://github.com/Schniz/fnm) | [Schniz](https://github.com/Schniz) | Fast and simple Node.js version manager |
| [uv](https://github.com/astral-sh/uv) | [Astral](https://github.com/astral-sh) | Ultra-fast Python package manager |
| [Zed](https://zed.dev/) | [Zed Industries](https://github.com/zed-industries) | High-performance code editor |
| [Catppuccin Blur](https://github.com/jenslys/zed-catppuccin-blur) | [jenslys](https://github.com/jenslys) | Catppuccin theme with blur for Zed |
| [Kaku](https://github.com/tw93/Kaku) | [tw93](https://github.com/tw93) | Lightweight macOS terminal |
| [Mole](https://github.com/tw93/Mole) | [tw93](https://github.com/tw93) | macOS system cleaner |
| [IINA](https://iina.io/) | [IINA](https://github.com/iina/iina) | Modern media player for macOS |
| [Karabiner-Elements](https://karabiner-elements.pqrs.org/) | [pqrs.org](https://github.com/pqrs-org) | Keyboard customizer for macOS |
| [Codex](https://github.com/openai/codex) | [OpenAI](https://github.com/openai) | AI coding assistant CLI |
| [Claude Code](https://github.com/anthropics/claude-code) | [Anthropic](https://github.com/anthropics) | AI coding assistant CLI |
| [zsh-autosuggestions](https://github.com/zsh-users/zsh-autosuggestions) | [zsh-users](https://github.com/zsh-users) | Fish-like autosuggestions for zsh |
| [zsh-syntax-highlighting](https://github.com/zsh-users/zsh-syntax-highlighting) | [zsh-users](https://github.com/zsh-users) | Syntax highlighting for zsh |
| [zsh-z](https://github.com/agkozak/zsh-z) | [agkozak](https://github.com/agkozak) | Fast directory jumping |

---

## License

MIT

## Author

[@kittors](https://github.com/kittors)
