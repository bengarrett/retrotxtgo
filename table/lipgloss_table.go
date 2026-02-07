// Package table provides lipgloss-based table formatting.
package table

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// LipglossTable renders a table using lipgloss styling.
func LipglossTable(wr io.Writer, rows []Row) error {
	if wr == nil {
		wr = io.Discard
	}

	// Create lipgloss styles
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("231")).
		Padding(0, 1)

	cellStyle := lipgloss.NewStyle().
		Padding(0, 1)

	specialCellStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")).
		Padding(0, 1)

	nonTableCellStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("244")).
		Padding(0, 1)

	tableOnlyCellStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("141")).
		Padding(0, 1)

	// Calculate column widths
	colWidths := calculateColumnWidths(rows)

	// Create header
	header := createHeader(headerStyle, colWidths)

	// Create rows
	var rowStrings []string
	for _, row := range rows {
		rowString := createRow(&row, cellStyle, specialCellStyle, nonTableCellStyle, tableOnlyCellStyle, colWidths)
		rowStrings = append(rowStrings, rowString)
	}

	// Build the table
	table := borderStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			header,
			lipgloss.JoinVertical(lipgloss.Left, rowStrings...),
		),
	)

	// Write the table
	fmt.Fprintln(wr, table)
	return nil
}

func calculateColumnWidths(rows []Row) [4]int {
	var widths [4]int

	// Header widths
	headers := []string{"Formal name", "Named value", "Numeric value", "Alias value"}
	for i, header := range headers {
		if len(header) > widths[i] {
			widths[i] = len(header)
		}
	}

	// Data widths
	for _, row := range rows {
		if len(row.Name) > widths[0] {
			widths[0] = len(row.Name)
		}
		if len(row.Value) > widths[1] {
			widths[1] = len(row.Value)
		}
		if len(row.Numeric) > widths[2] {
			widths[2] = len(row.Numeric)
		}
		if len(row.Alias) > widths[3] {
			widths[3] = len(row.Alias)
		}
	}

	// Add some padding
	for i := range widths {
		widths[i] += 2
	}

	return widths
}

func createHeader(style lipgloss.Style, widths [4]int) string {
	headerCells := []string{
		style.Render(fitString("Formal name", widths[0])),
		style.Render(fitString("Named value", widths[1])),
		style.Render(fitString("Numeric value", widths[2])),
		style.Render(fitString("Alias value", widths[3])),
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, headerCells...)
}

func createRow(row *Row, cellStyle, specialCellStyle, nonTableCellStyle, tableOnlyCellStyle lipgloss.Style, widths [4]int) string {
	var cells []string

	// Determine the appropriate style based on the row content
	if strings.HasPrefix(row.Name, "* ") {
		row.Name = strings.TrimPrefix(row.Name, "* ")
		cells = append(cells,
			specialCellStyle.Render(fitString(row.Name, widths[0])),
			specialCellStyle.Render(fitString(row.Value, widths[1])),
			specialCellStyle.Render(fitString(row.Numeric, widths[2])),
			specialCellStyle.Render(fitString(row.Alias, widths[3])),
		)
	} else if strings.HasPrefix(row.Name, "† ") {
		row.Name = strings.TrimPrefix(row.Name, "† ")
		cells = append(cells,
			nonTableCellStyle.Render(fitString(row.Name, widths[0])),
			nonTableCellStyle.Render(fitString(row.Value, widths[1])),
			nonTableCellStyle.Render(fitString(row.Numeric, widths[2])),
			nonTableCellStyle.Render(fitString(row.Alias, widths[3])),
		)
	} else if strings.HasPrefix(row.Name, "⁑ ") {
		row.Name = strings.TrimPrefix(row.Name, "⁑ ")
		cells = append(cells,
			tableOnlyCellStyle.Render(fitString(row.Name, widths[0])),
			tableOnlyCellStyle.Render(fitString(row.Value, widths[1])),
			tableOnlyCellStyle.Render(fitString(row.Numeric, widths[2])),
			tableOnlyCellStyle.Render(fitString(row.Alias, widths[3])),
		)
	} else {
		cells = append(cells,
			cellStyle.Render(fitString(row.Name, widths[0])),
			cellStyle.Render(fitString(row.Value, widths[1])),
			cellStyle.Render(fitString(row.Numeric, widths[2])),
			cellStyle.Render(fitString(row.Alias, widths[3])),
		)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, cells...)
}

func fitString(s string, width int) string {
	if len(s) < width {
		return s + strings.Repeat(" ", width-len(s))
	}
	return s[:width]
}
