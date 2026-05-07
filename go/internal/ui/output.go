package ui

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/Kivayan/azctx/go/internal/service"
)

type Field struct {
	Name  string
	Value string
}

func PrintMessagePanel(title string, message string) {
	PrintPanel(title, splitLines(message))
}

func PrintFieldPanel(title string, fields []Field, extraLines ...string) {
	lines := make([]string, 0, len(fields)+len(extraLines)+1)
	for _, field := range fields {
		lines = append(lines, fmt.Sprintf("%s: %s", field.Name, field.Value))
	}
	lines = append(lines, extraLines...)
	PrintPanel(title, lines)
}

func PrintContextTable(title string, summaries []service.ContextSummary) {
	fmt.Println(title)
	fmt.Println()

	writer := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	_, _ = fmt.Fprintln(writer, "ID\tName")
	_, _ = fmt.Fprintln(writer, "--\t----")
	for _, summary := range summaries {
		_, _ = fmt.Fprintf(writer, "%s\t%s\n", summary.ContextID, summary.ContextName)
	}
	_ = writer.Flush()
	fmt.Println()
	fmt.Println("Use 'azctx list --verbose' for detailed information")
}

func PrintPanel(title string, lines []string) {
	cleaned := normalizeLines(lines)
	width := len(title)
	for _, line := range cleaned {
		if len(line) > width {
			width = len(line)
		}
	}

	border := "+" + strings.Repeat("-", width+2) + "+"
	fmt.Println(title)
	fmt.Println(border)
	for _, line := range cleaned {
		fmt.Printf("| %-*s |\n", width, line)
	}
	fmt.Println(border)
}

func normalizeLines(lines []string) []string {
	if len(lines) == 0 {
		return []string{""}
	}

	normalized := make([]string, 0, len(lines))
	for _, line := range lines {
		parts := splitLines(line)
		normalized = append(normalized, parts...)
	}
	return normalized
}

func splitLines(text string) []string {
	return strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
}