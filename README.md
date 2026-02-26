<p align="center">
  <img src="https://img.shields.io/badge/platform-macOS-blue?style=for-the-badge&logo=apple&logoColor=white" />
  <img src="https://img.shields.io/badge/language-Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/TUI-Bubbletea-ff69b4?style=for-the-badge" />
  <img src="https://img.shields.io/badge/license-MIT-green?style=for-the-badge" />
</p>

<h1 align="center">🍃 freshbox</h1>

<p align="center">
  <strong>Set up a fresh Mac in minutes — not hours.</strong><br/>
  <sub>一个漂亮的终端界面工具，帮你几分钟内配置好全新的 Mac。</sub>
</p>

<p align="center">
  <a href="#-quick-install">Install</a> •
  <a href="#-features">Features</a> •
  <a href="#️-what-gets-installed">Catalog</a> •
  <a href="#-usage">Usage</a> •
  <a href="#-中文说明">中文</a> •
  <a href="#-contributing">Contributing</a>
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

> Already-installed tools are shown with **strikethrough** and their version — you only install what's missing.

---

## 🚀 Quick Install

**One command:**

```bash
curl -fsSL https://raw.githubusercontent.com/kittors/freshbox/main/install.sh | bash
```

**Or build from source:**

```bash
git clone https://github.com/kittors/freshbox.git
cd freshbox
go build -o freshbox .
./freshbox
```

The installer auto-detects your architecture (Apple Silicon / Intel) and either builds from source (if Go is available) or downloads a pre-built binary from GitHub Releases.

---

## ✨ Features

| | Feature | Description |
|---|---|---|
| 🌐 | **Bilingual Interface** | Full English / 中文 interface — choose at startup |
| 🔧 | **Smart Detection** | Auto-detects installed tools, shows versions, greys out what's already there |
| 📦 | **Node.js Manager** | Multi-select Node.js versions to install via [fnm](https://github.com/Schniz/fnm) |
| 📱 | **App Installer** | One-click install for curated macOS apps via Homebrew Cask |
| 🤖 | **AI Tool Config** | Full setup for Codex & Claude Code — model, API key, base URL |
| 🔌 | **MCP Servers** | Select from 11 popular MCP servers to configure for Claude Code & Codex |
| 🎨 | **Theme & Terminal** | Zed Catppuccin Blur theme, Kaku terminal + 4 zsh plugins |
| ⌨️ | **Keyboard Shortcuts** | Karabiner `⌃⌥⌘T` → opens Kaku in Finder's current folder |
| 📁 | **Dev Workspace** | Create organized `~/Developer` directory + Finder customization |
| 🖥 | **System Defaults** | Set default browser, editor, and media player |
| ✨ | **Beautiful TUI** | Rounded borders, spinner progress, smooth multi-page navigation |
| 📝 | **Install Logging** | Full install log at `~/.freshbox/install.log` for troubleshooting |

---

## 🗂️ What Gets Installed

<details open>
<summary><strong>🔧 Development Tools</strong></summary>

| Tool | Description | Method |
|------|-------------|--------|
| [Homebrew](https://brew.sh/) | The missing package manager for macOS | Official script |
| [Git](https://git-scm.com/) | Distributed version control | `brew install` |
| [Java (OpenJDK)](https://openjdk.org/) | JDK for JVM-based development | `brew install openjdk` |
| [Maven](https://maven.apache.org/) | Java build & dependency manager | `brew install` |
| [Gradle](https://gradle.org/) | Flexible build automation for JVM | `brew install` |
| [Python](https://www.python.org/) | General-purpose language | `brew install` |
| [uv](https://github.com/astral-sh/uv) | Ultra-fast Python package manager | `brew install` |
| [fnm](https://github.com/Schniz/fnm) | Fast Node.js version manager (Rust) | `brew install` |
| [Rust](https://www.rust-lang.org/) | Systems language with memory safety | `rustup` installer |
| [Go](https://go.dev/) | Statically typed language by Google | `brew install` |

</details>

<details>
<summary><strong>📦 Applications</strong></summary>

| App | Description | Method |
|-----|-------------|--------|
| [Google Chrome](https://www.google.com/chrome/) | Web browser by Google | `brew --cask` |
| [Zed](https://zed.dev/) | High-performance code editor | `brew --cask` |
| [IINA](https://iina.io/) | Modern media player for macOS | `brew --cask` |
| [Kaku](https://github.com/tw93/Kaku) | Lightweight terminal by tw93 | `brew --cask` |
| [Karabiner-Elements](https://karabiner-elements.pqrs.org/) | Keyboard customizer | `brew --cask` |
| [Mole](https://github.com/tw93/Mole) | macOS system cleaner by tw93 | `brew install` |
| [Tabby](https://tabby.sh/) | Modern terminal with SSH support | `brew --cask` |

</details>

<details>
<summary><strong>🤖 AI Tools</strong></summary>

| Tool | Description | Method |
|------|-------------|--------|
| [Codex](https://github.com/openai/codex) | OpenAI's AI coding CLI | `npm install -g` |
| [Claude Code](https://github.com/anthropics/claude-code) | Anthropic's AI coding CLI | `npm install -g` |

</details>

<details>
<summary><strong>🔌 MCP Servers</strong> (11 available)</summary>

| Server | Package |
|--------|---------|
| Playwright | `@playwright/mcp` |
| Context7 | `@upstash/context7-mcp` |
| Filesystem | `@modelcontextprotocol/server-filesystem` |
| GitHub | `@modelcontextprotocol/server-github` |
| Memory | `@modelcontextprotocol/server-memory` |
| Sequential Thinking | `@modelcontextprotocol/server-sequential-thinking` |
| Fetch | `@modelcontextprotocol/server-fetch` |
| Brave Search | `@modelcontextprotocol/server-brave-search` |
| Slack | `@modelcontextprotocol/server-slack` |
| Google Maps | `@modelcontextprotocol/server-google-maps` |
| SQLite | `@modelcontextprotocol/server-sqlite` |

</details>

<details>
<summary><strong>🎨 Extra Setup</strong></summary>

#### Zed Catppuccin Blur Theme

Clones [catppuccin-blur](https://github.com/jenslys/zed-catppuccin-blur), applies a custom icy blue tint, and auto-switches based on system appearance:

- **Light** → Catppuccin Latte with `#e8f0ff` tint
- **Dark** → Catppuccin Mocha with `#181c2e` tint

#### Kaku Terminal Setup

Initializes [Kaku](https://github.com/tw93/Kaku) with a full config and 4 essential zsh plugins:

- [zsh-autosuggestions](https://github.com/zsh-users/zsh-autosuggestions) — Fish-like auto-completion
- [zsh-completions](https://github.com/zsh-users/zsh-completions) — Additional completions
- [zsh-syntax-highlighting](https://github.com/zsh-users/zsh-syntax-highlighting) — Real-time syntax coloring
- [zsh-z](https://github.com/agkozak/zsh-z) — Fast directory jumping

#### Karabiner ⌃⌥⌘T → Kaku

Sets up `Ctrl+Option+Cmd+T` to quick-launch Kaku — opens in Finder's current directory if Finder is active.

#### Developer Workspace

Creates `~/Developer` with an organized structure + configures Finder (hidden files, path bar, list view, default to `~/Developer`):

```
~/Developer/
├── opensource/     Personal open-source projects
├── boundless/      Company projects
├── freelance/      Freelance / contract work
├── playground/     Learning & experiments
├── design/         UI designs, icons, assets
├── notes/          Technical notes & blog drafts
├── scripts/        Automation scripts & CLI tools
└── archive/        Completed / archived projects
```

</details>

<details>
<summary><strong>⚙️ Generated Config Files</strong></summary>

**Codex** — `~/.codex/config.toml` + `auth.json`

```toml
model = "o4-mini"
model_reasoning_effort = "medium"

[model_providers.freshbox]
name = "openai"
base_url = "https://api.openai.com/v1"
```

**Claude Code** — `~/.claude/settings.json`

```json
{
  "model": "claude-sonnet-4-6",
  "env": {
    "ANTHROPIC_API_KEY": "sk-ant-...",
    "ANTHROPIC_BASE_URL": "https://api.anthropic.com"
  }
}
```

</details>

---

## 🎮 Usage

```bash
freshbox
```

### Keyboard Shortcuts

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
🌐 Language  →  👋 Welcome  →  🔧 Dev Tools  →  📦 Apps  →  📦 Node.js
  →  🤖 AI Tools  →  ⚙️ Codex Config  →  ⚙️ Claude Config
  →  🔌 MCP Servers  →  🎨 Extra Setup  →  🖥 System Defaults
  →  ⏳ Installing...  →  ✅ Done!
```

---

## 📁 Project Structure

```
freshbox/
├── main.go                           # Entry point
├── install.sh                        # curl-based quick installer
├── internal/
│   ├── checker/
│   │   ├── checker.go                # System detection & version checking
│   │   └── checker_test.go           # 9 tests
│   ├── config/
│   │   ├── config.go                 # AI tool config generation (Codex/Claude/MCP)
│   │   └── config_test.go            # 14 tests
│   ├── installer/
│   │   ├── installer.go              # Install logic (brew/rustup/npm/fnm)
│   │   └── installer_test.go         # 7 tests
│   ├── setup/
│   │   ├── setup.go                  # Zed theme, Kaku init, Karabiner, workspace
│   │   └── setup_test.go             # 3 tests
│   └── ui/
│       ├── model.go                  # Bubbletea multi-page TUI (13 pages)
│       ├── install.go                # Async install queue with progress
│       ├── i18n.go                   # Bilingual text (EN/ZH)
│       ├── styles.go                 # Lipgloss styles
│       └── ui_test.go                # 33 tests
├── go.mod
└── go.sum
```

---

## 🧪 Testing

Run the full test suite:

```bash
go test ./... -v
```

The project has **55+ unit tests** covering:

- **checker** — tool detection, version parsing, registry completeness
- **config** — config merge logic, MCP timeout injection, JSON round-trips
- **installer** — brew args, fnm operations, system defaults
- **setup** — directory creation, config file generation
- **ui** — model lifecycle, navigation, selection, i18n, install queue

All tests use `t.TempDir()` and `t.Setenv("HOME", ...)` for complete isolation — no side effects on your real config files.

---

## 🇨🇳 中文说明

**freshbox** 是一个 macOS 装机 TUI 工具，帮你在全新的 Mac 上一键配置开发环境。

### 快速安装

```bash
curl -fsSL https://raw.githubusercontent.com/kittors/freshbox/main/install.sh | bash
```

### 功能亮点

- 🌐 中英文双语界面，启动时选择
- 🔧 自动检测已安装工具并显示版本号（已安装的划删除线）
- 📦 通过 fnm 安装和管理多个 Node.js 版本
- 📱 一键安装常用软件：Chrome、Zed、IINA、Kaku、Karabiner、Mole、Tabby
- 🤖 配置 AI 开发工具（Codex、Claude Code），自动生成配置文件
- 🔌 勾选配置 11 个流行的 MCP 服务
- 🎨 额外配置：Zed 冰蓝主题 / Kaku 终端初始化 / Karabiner 快捷键 / 开发工作区
- 🖥 设置系统默认浏览器、编辑器、播放器
- 📝 完整安装日志保存在 `~/.freshbox/install.log`

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

## 🤝 Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

```bash
# Clone and run locally
git clone https://github.com/kittors/freshbox.git
cd freshbox
go run .

# Run tests
go test ./... -v
```

---

## 🙏 Credits

Built with these amazing open-source projects:

| Project | Author | Role |
|---------|--------|------|
| [Bubbletea](https://github.com/charmbracelet/bubbletea) | [Charm](https://github.com/charmbracelet) | TUI framework |
| [Lipgloss](https://github.com/charmbracelet/lipgloss) | [Charm](https://github.com/charmbracelet) | Terminal styling |
| [Homebrew](https://brew.sh/) | [Homebrew](https://github.com/Homebrew) | Package manager |
| [fnm](https://github.com/Schniz/fnm) | [Schniz](https://github.com/Schniz) | Node.js version manager |
| [uv](https://github.com/astral-sh/uv) | [Astral](https://github.com/astral-sh) | Python package manager |
| [Zed](https://zed.dev/) | [Zed Industries](https://github.com/zed-industries) | Code editor |
| [Catppuccin Blur](https://github.com/jenslys/zed-catppuccin-blur) | [jenslys](https://github.com/jenslys) | Zed theme |
| [Kaku](https://github.com/tw93/Kaku) | [tw93](https://github.com/tw93) | Terminal |
| [Mole](https://github.com/tw93/Mole) | [tw93](https://github.com/tw93) | System cleaner |
| [IINA](https://iina.io/) | [IINA](https://github.com/iina/iina) | Media player |
| [Karabiner](https://karabiner-elements.pqrs.org/) | [pqrs.org](https://github.com/pqrs-org) | Keyboard customizer |

---

<p align="center">
  MIT License • Made with 💚 by <a href="https://github.com/kittors">@kittors</a>
</p>
