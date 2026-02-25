package ui

// Lang represents the selected language
type Lang int

const (
	LangEN Lang = iota
	LangZH
)

// T holds all translatable strings
type T struct {
	// Page names
	PageWelcome     string
	PageDevTools    string
	PageApps        string
	PageNodeVer     string
	PageAITools     string
	PageCodexCfg    string
	PageClaudeCfg   string
	PageMCP         string
	PageExtraSetup  string
	PageSysDefaults string
	PageInstalling  string
	PageDone        string

	// Welcome
	WelcomeTitle    string
	WelcomeDesc     string
	WelcomeDevTools string
	WelcomeFnm      string
	WelcomeApps     string
	WelcomeAI       string
	WelcomeMCP      string
	WelcomeSys      string
	WelcomeStart    string
	WelcomeQuit     string

	// Language selection
	LangTitle       string
	LangPrompt      string

	// Section titles
	TitleDevTools   string
	TitleApps       string
	TitleAITools    string
	TitleFnmVer     string
	TitleMCP        string
	TitleMCPDesc    string
	TitleSysDefault string
	TitleInstalling string
	TitleDone       string

	// Extra setup
	TitleExtraSetup     string
	TitleExtraSetupDesc string
	ExtraZedTheme       string
	ExtraZedThemeDesc   string
	ExtraKakuInit       string
	ExtraKakuInitDesc   string
	ExtraKarabiner      string
	ExtraKarabinerDesc  string
	ExtraDevWorkspace   string
	ExtraDevWorkspaceDesc string

	// Config form
	CfgModel        string
	CfgThinkLevel   string
	CfgBaseURL      string
	CfgAPIKey       string

	// System defaults
	DefBrowser      string
	DefBrowserDesc  string
	DefEditor       string
	DefEditorDesc   string
	DefPlayer       string
	DefPlayerDesc   string

	// Fnm
	FnmHint         string
	FnmLTSHint      string

	// Install
	InstallPrepare  string

	// Done
	DoneMsg         string
	DoneReady       string
	DoneExit        string

	// Footer
	FooterNav       string
	FooterForm      string
}

var texts = map[Lang]T{
	LangEN: {
		PageWelcome:     "Welcome",
		PageDevTools:    "Dev Tools",
		PageApps:        "Apps",
		PageNodeVer:     "Node.js Versions",
		PageAITools:     "AI Tools",
		PageCodexCfg:    "Codex Config",
		PageClaudeCfg:   "Claude Config",
		PageMCP:         "MCP Servers",
		PageSysDefaults: "System Defaults",
		PageInstalling:  "Installing...",
		PageDone:        "Done!",

		WelcomeTitle:    "Welcome to freshbox!",
		WelcomeDesc:     "This tool will help you set up your new Mac with:",
		WelcomeDevTools: "Development tools (brew, git, java, python, rust, go...)",
		WelcomeFnm:      "Node.js version management via fnm",
		WelcomeApps:     "Applications (Chrome, Zed, IINA, Kaku, Karabiner)",
		WelcomeAI:       "AI tools (Codex, Claude Code) with full config",
		WelcomeMCP:      "MCP servers (Playwright, Context7, and more)",
		WelcomeSys:      "System defaults + Zed theme, Kaku setup, dev workspace",
		WelcomeStart:    "Press Enter to get started",
		WelcomeQuit:     "q to quit",

		LangTitle:       "Language / 语言",
		LangPrompt:      "Select your language / 选择语言",

		TitleDevTools:   "Development Tools",
		TitleApps:       "Applications",
		TitleAITools:    "AI Tools",
		TitleFnmVer:     "Select Node.js Versions to Install",
		TitleMCP:        "MCP Servers",
		TitleMCPDesc:    "Select MCP servers to configure for your AI tools",
		TitleSysDefault: "System Defaults",
		TitleInstalling: "Installing...",
		TitleDone:       "All done!",

		TitleExtraSetup:       "Extra Setup",
		TitleExtraSetupDesc:   "Optional configurations to enhance your workflow",
		ExtraZedTheme:         "Zed Catppuccin Blur Theme",
		ExtraZedThemeDesc:     "Install catppuccin-blur theme with icy blue tint, auto light/dark mode",
		ExtraKakuInit:         "Kaku Terminal Setup",
		ExtraKakuInitDesc:     "Initialize Kaku config + zsh plugins (autosuggestions, completions, syntax-highlighting, z)",
		ExtraKarabiner:        "Karabiner ⌃⌥⌘T → Kaku",
		ExtraKarabinerDesc:    "Set up Ctrl+Option+Cmd+T to quick-launch Kaku (opens Finder's current folder)",
		ExtraDevWorkspace:     "Developer Workspace",
		ExtraDevWorkspaceDesc: "Create ~/Developer directory structure + configure Finder (hidden files, path bar, list view)",

		PageExtraSetup: "Extra Setup",

		CfgModel:        "Model",
		CfgThinkLevel:   "Thinking Level",
		CfgBaseURL:      "Base URL",
		CfgAPIKey:       "API Key",

		DefBrowser:      "Default Browser → Google Chrome",
		DefBrowserDesc:  "Set Chrome as system default browser",
		DefEditor:       "Default Editor → Zed",
		DefEditorDesc:   "Set Zed as global code editor",
		DefPlayer:       "Default Player → IINA",
		DefPlayerDesc:   "Set IINA as default media player",

		FnmHint:         "fnm will be installed first, then you can select Node versions.",
		FnmLTSHint:      "Showing common LTS versions:",

		InstallPrepare:  "Preparing installation...",

		DoneMsg:         "Your Mac is set up and ready to go.",
		DoneReady:       "All done!",
		DoneExit:        "Press Enter or q to exit.",

		FooterNav:       "↑/↓ navigate • space toggle • a all • n none • tab next • shift+tab back • q quit",
		FooterForm:      "↑/↓ navigate fields • tab next field • enter confirm • shift+tab back",
	},
	LangZH: {
		PageWelcome:     "欢迎",
		PageDevTools:    "开发工具",
		PageApps:        "应用程序",
		PageNodeVer:     "Node.js 版本",
		PageAITools:     "AI 工具",
		PageCodexCfg:    "Codex 配置",
		PageClaudeCfg:   "Claude 配置",
		PageMCP:         "MCP 服务",
		PageSysDefaults: "系统默认",
		PageInstalling:  "安装中...",
		PageDone:        "完成！",

		WelcomeTitle:    "欢迎使用 freshbox！",
		WelcomeDesc:     "这个工具将帮助你配置新 Mac：",
		WelcomeDevTools: "开发工具（brew、git、java、python、rust、go...）",
		WelcomeFnm:      "通过 fnm 管理 Node.js 多版本",
		WelcomeApps:     "常用应用（Chrome、Zed、IINA、Kaku、Karabiner）",
		WelcomeAI:       "AI 工具（Codex、Claude Code）完整配置",
		WelcomeMCP:      "MCP 服务（Playwright、Context7 等）",
		WelcomeSys:      "系统默认设置 + Zed 主题、Kaku 配置、开发工作区",
		WelcomeStart:    "按 Enter 开始",
		WelcomeQuit:     "q 退出",

		LangTitle:       "Language / 语言",
		LangPrompt:      "Select your language / 选择语言",

		TitleDevTools:   "开发工具",
		TitleApps:       "应用程序",
		TitleAITools:    "AI 工具",
		TitleFnmVer:     "选择要安装的 Node.js 版本",
		TitleMCP:        "MCP 服务",
		TitleMCPDesc:    "选择要为 AI 工具配置的 MCP 服务",
		TitleSysDefault: "系统默认设置",
		TitleInstalling: "安装中...",
		TitleDone:       "全部完成！",

		TitleExtraSetup:       "额外配置",
		TitleExtraSetupDesc:   "可选的工作流增强配置",
		ExtraZedTheme:         "Zed Catppuccin Blur 主题",
		ExtraZedThemeDesc:     "安装 catppuccin-blur 主题，冰蓝色调，自动跟随系统明暗模式",
		ExtraKakuInit:         "Kaku 终端初始化",
		ExtraKakuInitDesc:     "初始化 Kaku 配置 + zsh 插件（自动补全、语法高亮、目录跳转等）",
		ExtraKarabiner:        "Karabiner ⌃⌥⌘T → Kaku",
		ExtraKarabinerDesc:    "设置 Ctrl+Option+Cmd+T 快速启动 Kaku（自动打开 Finder 当前目录）",
		ExtraDevWorkspace:     "开发工作区",
		ExtraDevWorkspaceDesc: "创建 ~/Developer 目录结构 + 配置 Finder（显示隐藏文件、路径栏、列表视图）",

		PageExtraSetup: "额外配置",

		CfgModel:        "模型",
		CfgThinkLevel:   "思考级别",
		CfgBaseURL:      "接口地址",
		CfgAPIKey:       "API 密钥",

		DefBrowser:      "默认浏览器 → Google Chrome",
		DefBrowserDesc:  "将 Chrome 设为系统默认浏览器",
		DefEditor:       "默认编辑器 → Zed",
		DefEditorDesc:   "将 Zed 设为全局代码编辑器",
		DefPlayer:       "默认播放器 → IINA",
		DefPlayerDesc:   "将 IINA 设为默认媒体播放器",

		FnmHint:         "fnm 将先被安装，之后你可以选择 Node 版本。",
		FnmLTSHint:      "显示常用 LTS 版本：",

		InstallPrepare:  "正在准备安装...",

		DoneMsg:         "你的 Mac 已配置完成，准备就绪。",
		DoneReady:       "全部完成！",
		DoneExit:        "按 Enter 或 q 退出。",

		FooterNav:       "↑/↓ 导航 • 空格 切换 • a 全选 • n 全不选 • tab 下一步 • shift+tab 上一步 • q 退出",
		FooterForm:      "↑/↓ 切换字段 • tab 下一字段 • enter 确认 • shift+tab 返回",
	},
}

// GetText returns the translation struct for the given language
func GetText(lang Lang) T {
	return texts[lang]
}
