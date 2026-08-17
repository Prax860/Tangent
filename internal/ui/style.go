package ui

import "github.com/charmbracelet/lipgloss"

// ─── Palette ──────────────────────────────────────────────

var (
	NeonBlue    = lipgloss.Color("#00F0FF")
	ElectricVio = lipgloss.Color("#B388FF")
	HotPink     = lipgloss.Color("#FF3CAC")
	AcidGreen   = lipgloss.Color("#39FF88")
	DangerRed   = lipgloss.Color("#FF4D6D")

	GlassDark = lipgloss.Color("#0D0F1A")
	GlassMid  = lipgloss.Color("#151827")

	SoftWhite = lipgloss.Color("#E8EBFF")
	MutedGrey = lipgloss.Color("#7A7F9E")
)

// ─── Title ────────────────────────────────────────────────

var TitleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(SoftWhite).
	Background(lipgloss.Color("#1B1030")).
	Padding(0, 3).
	MarginBottom(1)

// Accent dot used in headers.
var accentDot = lipgloss.NewStyle().
	Foreground(NeonBlue).
	Render("●")

// ─── Section Headers ─────────────────────────────────────

var HeaderStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(NeonBlue).
	Background(GlassMid).
	Padding(0, 1).
	MarginTop(1)

// ─── Body ─────────────────────────────────────────────────

var PromptStyle = lipgloss.NewStyle().
	Foreground(SoftWhite).
	Italic(true).
	PaddingLeft(1)

var LabelStyle = lipgloss.NewStyle().
	Foreground(ElectricVio).
	Bold(true)

var ValueStyle = lipgloss.NewStyle().
	Foreground(SoftWhite)

var CommandStyle = lipgloss.NewStyle().
	Foreground(AcidGreen).
	Background(lipgloss.Color("#0A0C14")).
	Bold(true).
	Padding(0, 1).
	MarginLeft(1)

var ExplanationStyle = lipgloss.NewStyle().
	Foreground(MutedGrey).
	Italic(true).
	PaddingLeft(1)

// ─── Safety ────────────────────────────────────────────────

var SafeStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#00110A")).
	Background(AcidGreen).
	Bold(true).
	Padding(0, 2)

var DangerStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#2A0008")).
	Background(DangerRed).
	Bold(true).
	Padding(0, 2)

// ─── Main Container ───────────────────────────────────────

var BoxStyle = lipgloss.NewStyle().
	Foreground(SoftWhite).
	Background(GlassDark).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(ElectricVio).
	Padding(1, 3).
	Margin(1, 0)

// ─── Arguments ────────────────────────────────────────────

var KeyPillStyle = lipgloss.NewStyle().
	Foreground(GlassDark).
	Background(NeonBlue).
	Bold(true).
	Padding(0, 1).
	MarginRight(1)