package table

import (
	"fmt"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/golden"
	"github.com/muesli/termenv"
)

// README Examples

func TestREADME(t *testing.T) {
	t.Run("table.New()", func(t *testing.T) {
		lipgloss.SetColorProfile(termenv.TrueColor)
		rows := [][]string{
			{"Chinese", "您好", "你好"},
			{"Japanese", "こんにちは", "やあ"},
			{"Arabic", "أهلين", "أهلا"},
			{"Russian", "Здравствуйте", "Привет"},
			{"Spanish", "Hola", "¿Qué tal?"},
		}
		var (
			purple    = lipgloss.Color("99")
			gray      = lipgloss.Color("245")
			lightGray = lipgloss.Color("241")

			headerStyle  = lipgloss.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
			cellStyle    = lipgloss.NewStyle().Padding(0, 1).Width(14)
			oddRowStyle  = cellStyle.Foreground(gray)
			evenRowStyle = cellStyle.Foreground(lightGray)
		)

		table := New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(purple)).
			StyleFunc(func(row, col int) lipgloss.Style {
				switch {
				case row == HeaderRow:
					return headerStyle
				case row%2 == 0:
					return evenRowStyle
				default:
					return oddRowStyle
				}
			}).
			Headers("LANGUAGE", "FORMAL", "INFORMAL").
			Rows(rows...)

		// You can also add tables row-by-row
		table.Row("English", "You look absolutely fabulous.", "How's it going?")
		t.Log(table.String())
		golden.RequireEqual(t, []byte(table.String()))
	})
}

// Other

func ExampleFilter() {
	data := NewStringData().
		Item("Chinese", "Nǐn hǎo", "Nǐ hǎo").
		Item("French", "Bonjour", "Salut").
		Item("Japanese", "こんにちは", "やあ").
		Item("Russian", "Zdravstvuyte", "Privet").
		Item("Spanish", "Hola", "¿Qué tal?")

	filter := NewFilter(data).Filter(func(row int) bool {
		return data.At(row, 0) != "French"
	})

	table := New().
		Border(lipgloss.NormalBorder()).
		StyleFunc(TableStyle).
		Headers("LANGUAGE", "FORMAL", "INFORMAL").
		Data(filter)
	fmt.Print(table.String())
	// Output:
	// ┌──────────┬──────────────┬───────────┐
	// │ LANGUAGE │    FORMAL    │ INFORMAL  │
	// ├──────────┼──────────────┼───────────┤
	// │ Chinese  │ Nǐn hǎo      │ Nǐ hǎo    │
	// │ Japanese │ こんにちは   │ やあ      │
	// │ Russian  │ Zdravstvuyte │ Privet    │
	// │ Spanish  │ Hola         │ ¿Qué tal? │
	// └──────────┴──────────────┴───────────┘
}

func ExampleStyleFunc() {
	HeaderStyle := lipgloss.NewStyle().Align(lipgloss.Center)
	EvenRowStyle := lipgloss.NewStyle().Align(lipgloss.Right)

	t := New().
		Width(16).
		Headers("Name", "Age").
		Row("Kini", "4").
		Row("Eli", "1").
		Row("Iris", "102").
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == HeaderRow:
				return HeaderStyle
			case row%2 == 0:
				return EvenRowStyle
			default:
				return lipgloss.NewStyle()
			}
		})
	fmt.Print(t.String())
	// Output:
	// ╭───────┬──────╮
	// │ Name  │ Age  │
	// ├───────┼──────┤
	// │   Kini│     4│
	// │Eli    │1     │
	// │   Iris│   102│
	// ╰───────┴──────╯
}
