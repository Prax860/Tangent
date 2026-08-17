package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/prax860/tangent/internal/types"
)

func Show(response types.Response) {

	var b strings.Builder

	// ── Title ───────────────────────────────────
	b.WriteString(CompactTitle.Render("⚡ Tangent"))
	b.WriteString("\n\n")

	// ── Intent + Workspace (aligned 2-col) ─────
	intentRow := fmt.Sprintf(
		"  %-11s %s",
		CompactLabel.Render("Intent"),
		CompactValue.Render(string(response.Request.Intent)),
	)
	workspaceRow := fmt.Sprintf(
		"  %-11s %s",
		CompactLabel.Render("Workspace"),
		CompactValue.Render(string(response.Request.Workspace)),
	)
	b.WriteString(intentRow + "\n" + workspaceRow + "\n\n")

	// ── Command with arrow ──────────────────────
	b.WriteString(
		"  " +
			CompactArrow.Render("→") +
			" " +
			CompactCommand.Render(response.Command.Command) +
			"\n\n",
	)

	// ── Safety pill ─────────────────────────────
	if response.Command.Safe {
		b.WriteString("  " + CompactSafe.Render("✓ SAFE"))
	} else {
		b.WriteString("  " + CompactUnsafe.Render("✗ UNSAFE"))
	}

	// ── Interactive pill if applicable ──────────
	if response.Command.Interactive {
		b.WriteString("\n  " + CompactInteractive.Render("⚡ INTERACTIVE"))
	}

	b.WriteString("\n")

	fmt.Println(b.String())
}

// ─── Compact UI styles ──────────────────────────────────

var (
	CompactTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#1B1030")).
			Padding(0, 2)

	CompactLabel = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#9C6BFF"))

	CompactValue = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E8EBFF"))

	CompactArrow = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00F0FF")).
			Bold(true)

	CompactCommand = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#39FF88")).
			Bold(true)

	CompactSafe = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00110A")).
			Background(lipgloss.Color("#39FF88")).
			Bold(true).
			Padding(0, 2)

	CompactUnsafe = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#2A0008")).
			Background(lipgloss.Color("#FF4D6D")).
			Bold(true).
			Padding(0, 2)

	CompactInteractive = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#0D0F1A")).
				Background(lipgloss.Color("#FFD23F")).
				Bold(true).
				Padding(0, 2)
)
