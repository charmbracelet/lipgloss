package table

import (
	"fmt"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/x/exp/golden"
)

var TableStyle = func(row, col int) lipgloss.Style {
	switch {
	case row == HeaderRow:
		return lipgloss.NewStyle().Padding(0, 1).Align(lipgloss.Center)
	case row%2 == 0:
		return lipgloss.NewStyle().Padding(0, 1)
	default:
		return lipgloss.NewStyle().Padding(0, 1)
	}
}

func TestTable(t *testing.T) {
	table := New().
		Border(lipgloss.NormalBorder()).
		StyleFunc(TableStyle).
		Headers("LANGUAGE", "FORMAL", "INFORMAL").
		Row("Chinese", "Nǐn hǎo", "Nǐ hǎo").
		Row("French", "Bonjour", "Salut").
		Row("Japanese", "こんにちは", "やあ").
		Row("Russian", "Zdravstvuyte", "Privet").
		Row("Spanish", "Hola", "¿Qué tal?")

	golden.RequireEqual(t, []byte(table.String()))
}

func TestTableExample(t *testing.T) {
	HeaderStyle := lipgloss.NewStyle().Padding(0, 1).Align(lipgloss.Center)
	EvenRowStyle := lipgloss.NewStyle().Padding(0, 1)
	OddRowStyle := lipgloss.NewStyle().Padding(0, 1)

	rows := [][]string{
		{"Chinese", "您好", "你好"},
		{"Japanese", "こんにちは", "やあ"},
		{"Russian", "Здравствуйте", "Привет"},
		{"Spanish", "Hola", "¿Qué tal?"},
	}

	table := New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == HeaderRow:
				return HeaderStyle
			case row%2 == 0:
				return EvenRowStyle
			default:
				return OddRowStyle
			}
		}).
		Headers("LANGUAGE", "FORMAL", "INFORMAL").
		Rows(rows...)

	// You can also add tables row-by-row
	table.Row("English", "You look absolutely fabulous.", "How's it going?")

	golden.RequireEqual(t, []byte(table.String()))
}

func TestTableEmpty(t *testing.T) {
	table := New().
		Border(lipgloss.NormalBorder()).
		StyleFunc(TableStyle).
		Headers("LANGUAGE", "FORMAL", "INFORMAL")

	golden.RequireEqual(t, []byte(table.String()))
}

func TestTableOffset(t *testing.T) {
	table := New().
		Border(lipgloss.NormalBorder()).
		StyleFunc(TableStyle).
		Headers("LANGUAGE", "FORMAL", "INFORMAL").
		Row("Chinese", "Nǐn hǎo", "Nǐ hǎo").
		Row("French", "Bonjour", "Salut").
		Row("Japanese", "こんにちは", "やあ").
		Row("Russian", "Zdravstvuyte", "Privet").
		Row("Spanish", "Hola", "¿Qué tal?").
		Offset(1)

	golden.RequireEqual(t, []byte(table.String()))
}

func TestTableBorder(t *testing.T) {
	rows := [][]string{
		{"Chinese", "Nǐn hǎo", "Nǐ hǎo"},
		{"French", "Bonjour", "Salut"},
		{"Japanese", "こんにちは", "やあ"},
		{"Russian", "Zdravstvuyte", "Privet"},
		{"Spanish", "Hola", "¿Qué tal?"},
	}

	table := New().
		Border(lipgloss.DoubleBorder()).
		StyleFunc(TableStyle).
		Headers("LANGUAGE", "FORMAL", "INFORMAL").
		Rows(rows...)

	golden.RequireEqual(t, []byte(table.String()))
}

func TestTableSetRows(t *testing.T) {
	rows := [][]string{
		{"Chinese", "Nǐn hǎo", "Nǐ hǎo"},
		{"French", "Bonjour", "Salut"},
		{"Japanese", "こんにちは", "やあ"},
		{"Russian", "Zdravstvuyte", "Privet"},
		{"Spanish", "Hola", "¿Qué tal?"},
	}
	table := New().
		Border(lipgloss.NormalBorder()).
		StyleFunc(TableStyle).
		Headers("LANGUAGE", "FORMAL", "INFORMAL").
		Rows(rows...)

	golden.RequireEqual(t, []byte(table.String()))
}

func TestMoreCellsThanHeaders(t *testing.T) {
	rows := [][]string{
		{"Chinese", "Nǐn hǎo", "Nǐ hǎo"},
		{"French", "Bonjour", "Salut"},
		{"Japanese", "こんにちは", "やあ"},
		{"Russian", "Zdravstvuyte", "Privet"},
		{"Spanish", "Hola", "¿Qué tal?"},
	}
	table := New().
		Border(lipgloss.NormalBorder()).
		StyleFunc(TableStyle).
		Headers("LANGUAGE", "FORMAL").
		Rows(rows...)

	golden.RequireEqual(t, []byte(table.String()))
}

func TestMoreCellsThanHeadersExtra(t *testing.T) {
	rows := [][]string{
		{"Chinese", "Nǐn hǎo", "Nǐ hǎo"},
		{"French", "Bonjour", "Salut", "Salut"},
		{"Japanese", "こんにちは", "やあ"},
		{"Russian", "Zdravstvuyte", "Privet", "Privet", "Privet"},
		{"Spanish", "Hola", "¿Qué tal?"},
	}

	table := New().
		Border(lipgloss.NormalBorder()).
		StyleFunc(TableStyle).
		Headers("LANGUAGE", "FORMAL").
		Rows(rows...)

	golden.RequireEqual(t, []byte(table.String()))
}

func TestTableNoHeaders(t *testing.T) {
	table := New().
		Border(lipgloss.NormalBorder()).
		StyleFunc(TableStyle).
		Row("Chinese", "Nǐn hǎo", "Nǐ hǎo").
		Row("French", "Bonjour", "Salut").
		Row("Japanese", "こんにちは", "やあ").
		Row("Russian", "Zdravstvuyte", "Privet").
		Row("Spanish", "Hola", "¿Qué tal?")

	golden.RequireEqual(t, []byte(table.String()))
}

func TestTableNoColumnSeparators(t *testing.T) {
	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(false).
		StyleFunc(TableStyle).
		Row("Chinese", "Nǐn hǎo", "Nǐ hǎo").
		Row("French", "Bonjour", "Salut").
		Row("Japanese", "こんにちは", "やあ").
		Row("Russian", "Zdravstvuyte", "Privet").
		Row("Spanish", "Hola", "¿Qué tal?")

	golden.RequireEqual(t, []byte(table.String()))
}

func TestTableNoColumnSeparatorsWithHeaders(t *testing.T) {
	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(false).
		StyleFunc(TableStyle).
		Headers("LANGUAGE", "FORMAL", "INFORMAL").
		Row("Chinese", "Nǐn hǎo", "Nǐ hǎo").
		Row("French", "Bonjour", "Salut").
		Row("Japanese", "こんにちは", "やあ").
		Row("Russian", "Zdravstvuyte", "Privet").
		Row("Spanish", "Hola", "¿Qué tal?")

	golden.RequireEqual(t, []byte(table.String()))
}

func TestBorderColumnsWithExtraRows(t *testing.T) {
	rows := [][]string{
		{"Chinese", "Nǐn hǎo", "Nǐ hǎo"},
		{"French", "Bonjour", "Salut", "Salut"},
		{"Japanese", "こんにちは", "やあ"},
		{"Russian", "Zdravstvuyte", "Privet", "Privet", "Privet"},
		{"Spanish", "Hola", "¿Qué tal?"},
	}

	table := New().
		Border(lipgloss.NormalBorder()).
		BorderColumn(false).
		StyleFunc(TableStyle).
		Headers("LANGUAGE", "FORMAL").
		Rows(rows...)

	golden.RequireEqual(t, []byte(table.String()))
}

func TestNew(t *testing.T) {
	table := New()
	expected := ""
	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableUnsetBorders(t *testing.T) {
	rows := [][]string{
		{"Chinese", "Nǐn hǎo", "Nǐ hǎo"},
		{"French", "Bonjour", "Salut"},
		{"Japanese", "こんにちは", "やあ"},
		{"Russian", "Zdravstvuyte", "Privet"},
		{"Spanish", "Hola", "¿Qué tal?"},
	}

	table := New().
		Border(lipgloss.NormalBorder()).
		StyleFunc(TableStyle).
		Headers("LANGUAGE", "FORMAL", "INFORMAL").
		Rows(rows...).
		BorderTop(false).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false)

	golden.RequireEqual(t, []byte(table.String()))
}

func TestTableUnsetHeaderSeparator(t *testing.T) {
	rows := [][]string{
		{"Chinese", "Nǐn hǎo", "Nǐ hǎo"},
		{"French", "Bonjour", "Salut"},
		{"Japanese", "こんにちは", "やあ"},
		{"Russian", "Zdravstvuyte", "Privet"},
		{"Spanish", "Hola", "¿Qué tal?"},
	}

	table := New().
		Border(lipgloss.NormalBorder()).
		StyleFunc(TableStyle).
		Headers("LANGUAGE", "FORMAL", "INFORMAL").
		Rows(rows...).
		BorderHeader(false).
		BorderTop(false).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false)

	golden.RequireEqual(t, []byte(table.String()))
}

func TestTableUnsetHeaderSeparatorWithBorder(t *testing.T) {
	rows := [][]string{
		{"Chinese", "Nǐn hǎo", "Nǐ hǎo"},
		{"French", "Bonjour", "Salut"},
		{"Japanese", "こんにちは", "やあ"},
		{"Russian", "Zdravstvuyte", "Privet"},
		{"Spanish", "Hola", "¿Qué tal?"},
	}

	table := New().
		Border(lipgloss.NormalBorder()).
		StyleFunc(TableStyle).
		Headers("LANGUAGE", "FORMAL", "INFORMAL").
		Rows(rows...).
		BorderHeader(false)

	golden.RequireEqual(t, []byte(table.String()))
}

func TestTableRowSeparators(t *testing.T) {
	rows := [][]string{
		{"Chinese", "Nǐn hǎo", "Nǐ hǎo"},
		{"French", "Bonjour", "Salut"},
		{"Japanese", "こんにちは", "やあ"},
		{"Russian", "Zdravstvuyte", "Privet"},
		{"Spanish", "Hola", "¿Qué tal?"},
	}

	table := New().
		Border(lipgloss.NormalBorder()).
		StyleFunc(TableStyle).
		BorderRow(true).
		Headers("LANGUAGE", "FORMAL", "INFORMAL").
		Rows(rows...)

	golden.RequireEqual(t, []byte(table.String()))
}

func TestTableHeights(t *testing.T) {
	styleFunc := func(row, col int) lipgloss.Style {
		if row == HeaderRow {
			return lipgloss.NewStyle().Padding(0, 1)
		}
		if col == 0 {
			return lipgloss.NewStyle().Width(18).Padding(1)
		}
		return lipgloss.NewStyle().Width(25).Padding(1, 2)
	}

	rows := [][]string{
		{"Chutar o balde", `Literally translates to "kick the bucket." It's used when someone gives up or loses patience.`},
		{"Engolir sapos", `Literally means "to swallow frogs." It's used to describe someone who has to tolerate or endure unpleasant situations.`},
		{"Arroz de festa", `Literally means "party rice." It´s used to refer to someone who shows up everywhere.`},
	}

	table := New().
		Border(lipgloss.NormalBorder()).
		StyleFunc(styleFunc).
		Headers("EXPRESSION", "MEANING").
		Rows(rows...)

	golden.RequireEqual(t, []byte(table.String()))
}

func TestTableMultiLineRowSeparator(t *testing.T) {
	styleFunc := func(row, col int) lipgloss.Style {
		if row == HeaderRow {
			return lipgloss.NewStyle().Padding(0, 1)
		}
		if col == 0 {
			return lipgloss.NewStyle().Width(18).Padding(1)
		}
		return lipgloss.NewStyle().Width(25).Padding(1, 2)
	}

	table := New().
		Border(lipgloss.NormalBorder()).
		StyleFunc(styleFunc).
		Headers("EXPRESSION", "MEANING").
		BorderRow(true).
		Row("Chutar o balde", `Literally translates to "kick the bucket." It's used when someone gives up or loses patience.`).
		Row("Engolir sapos", `Literally means "to swallow frogs." It's used to describe someone who has to tolerate or endure unpleasant situations.`).
		Row("Arroz de festa", `Literally means "party rice." It´s used to refer to someone who shows up everywhere.`)

	golden.RequireEqual(t, []byte(table.String()))
}

func TestTableWidthExpand(t *testing.T) {
	rows := [][]string{
		{"Chinese", "Nǐn hǎo", "Nǐ hǎo"},
		{"French", "Bonjour", "Salut"},
		{"Japanese", "こんにちは", "やあ"},
		{"Russian", "Zdravstvuyte", "Privet"},
		{"Spanish", "Hola", "¿Qué tal?"},
	}

	table := New().
		Width(80).
		StyleFunc(TableStyle).
		Border(lipgloss.NormalBorder()).
		Headers("LANGUAGE", "FORMAL", "INFORMAL").
		Rows(rows...)

	if lipgloss.Width(table.String()) != 80 {
		t.Fatalf("expected table width to be 80, got %d", lipgloss.Width(table.String()))
	}

	golden.RequireEqual(t, []byte(table.String()))
}

func TestTableWidthShrink(t *testing.T) {
	rows := [][]string{
		{"Chinese", "Nǐn hǎo", "Nǐ hǎo"},
		{"French", "Bonjour", "Salut"},
		{"Japanese", "こんにちは", "やあ"},
		{"Russian", "Zdravstvuyte", "Privet"},
		{"Spanish", "Hola", "¿Qué tal?"},
	}

	table := New().
		Width(30).
		StyleFunc(TableStyle).
		Border(lipgloss.NormalBorder()).
		Headers("LANGUAGE", "FORMAL", "INFORMAL").
		Rows(rows...)

	golden.RequireEqual(t, []byte(table.String()))
}

func TestTableWidthSmartCrop(t *testing.T) {
	rows := [][]string{
		{"Kini", "40", "New York"},
		{"Eli", "30", "London"},
		{"Iris", "20", "Paris"},
	}

	table := New().
		Width(25).
		StyleFunc(TableStyle).
		Border(lipgloss.NormalBorder()).
		Headers("Name", "Age of Person", "Location").
		Rows(rows...)

	golden.RequireEqual(t, []byte(table.String()))
}

func TestTableWidthSmartCropExtensive(t *testing.T) {
	rows := [][]string{
		{"Chinese", "您好", "你好"},
		{"Japanese", "こんにちは", "やあ"},
		{"Arabic", "أهلين", "أهلا"},
		{"Russian", "Здравствуйте", "Привет"},
		{"Spanish", "Hola", "¿Qué tal?"},
		{"English", "You look absolutely fabulous.", "How's it going?"},
	}

	table := New().
		Width(18).
		StyleFunc(TableStyle).
		Border(lipgloss.ThickBorder()).
		Headers("LANGUAGE", "FORMAL", "INFORMAL").
		Wrap(false).
		Rows(rows...)

	golden.RequireEqual(t, []byte(table.String()))
}

func TestTableWidthSmartCropTiny(t *testing.T) {
	rows := [][]string{
		{"Chinese", "您好", "你好"},
		{"Japanese", "こんにちは", "やあ"},
		{"Russian", "Здравствуйте", "Привет"},
		{"Spanish", "Hola", "¿Qué tal?"},
		{"English", "You look absolutely fabulous.", "How's it going?"},
	}

	table := New().
		Width(1).
		StyleFunc(TableStyle).
		Border(lipgloss.NormalBorder()).
		Headers("LANGUAGE", "FORMAL", "INFORMAL").
		Rows(rows...)

	golden.RequireEqual(t, []byte(table.String()))
}

func TestTableWidths(t *testing.T) {
	rows := [][]string{
		{"Chinese", "Nǐn hǎo", "Nǐ hǎo"},
		{"French", "Bonjour", "Salut"},
		{"Japanese", "こんにちは", "やあ"},
		{"Russian", "Zdravstvuyte", "Privet"},
		{"Spanish", "Hola", "¿Qué tal?"},
	}

	table := New().
		Width(30).
		StyleFunc(TableStyle).
		BorderLeft(false).
		BorderRight(false).
		Border(lipgloss.NormalBorder()).
		BorderColumn(false).
		Headers("LANGUAGE", "FORMAL", "INFORMAL").
		Rows(rows...)

	golden.RequireEqual(t, []byte(table.String()))
}

func TestTableWidthShrinkNoBorders(t *testing.T) {
	rows := [][]string{
		{"Chinese", "Nǐn hǎo", "Nǐ hǎo"},
		{"French", "Bonjour", "Salut"},
		{"Japanese", "こんにちは", "やあ"},
		{"Russian", "Zdravstvuyte", "Privet"},
		{"Spanish", "Hola", "¿Qué tal?"},
	}

	table := New().
		Width(30).
		StyleFunc(TableStyle).
		BorderLeft(false).
		BorderRight(false).
		Border(lipgloss.NormalBorder()).
		BorderColumn(false).
		Headers("LANGUAGE", "FORMAL", "INFORMAL").
		Rows(rows...)

	golden.RequireEqual(t, []byte(table.String()))
}

func TestFilter(t *testing.T) {
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

	golden.RequireEqual(t, []byte(table.String()))
}

func TestFilterInverse(t *testing.T) {
	data := NewStringData().
		Item("Chinese", "Nǐn hǎo", "Nǐ hǎo").
		Item("French", "Bonjour", "Salut").
		Item("Japanese", "こんにちは", "やあ").
		Item("Russian", "Zdravstvuyte", "Privet").
		Item("Spanish", "Hola", "¿Qué tal?")

	filter := NewFilter(data).Filter(func(row int) bool {
		return data.At(row, 0) == "French"
	})

	table := New().
		Border(lipgloss.NormalBorder()).
		StyleFunc(TableStyle).
		Headers("LANGUAGE", "FORMAL", "INFORMAL").
		Data(filter)

	golden.RequireEqual(t, []byte(table.String()))
}

func TestTableANSI(t *testing.T) {
	const code = "\x1b[31mC\x1b[0m\x1b[32mo\x1b[0m\x1b[34md\x1b[0m\x1b[33me\x1b[0m"

	rows := [][]string{
		{"Apple", "Red", "\x1b[31m31\x1b[0m"},
		{"Lime", "Green", "\x1b[32m32\x1b[0m"},
		{"Banana", "Yellow", "\x1b[33m33\x1b[0m"},
		{"Blueberry", "Blue", "\x1b[34m34\x1b[0m"},
	}

	table := New().
		Width(29).
		StyleFunc(TableStyle).
		Border(lipgloss.NormalBorder()).
		Headers("Fruit", "Color", code).
		Rows(rows...)

	golden.RequireEqual(t, []byte(table.String()))
}

func TestTableHeightExact(t *testing.T) {
	table := New().
		Height(9).
		Border(lipgloss.NormalBorder()).
		StyleFunc(TableStyle).
		Headers("LANGUAGE", "FORMAL", "INFORMAL").
		Row("Chinese", "Nǐn hǎo", "Nǐ hǎo").
		Row("French", "Bonjour", "Salut").
		Row("Japanese", "こんにちは", "やあ").
		Row("Russian", "Zdravstvuyte", "Privet").
		Row("Spanish", "Hola", "¿Qué tal?")

	golden.RequireEqual(t, []byte(table.String()))
}

func TestTableHeightExtra(t *testing.T) {
	table := New().
		Height(100).
		Border(lipgloss.NormalBorder()).
		StyleFunc(TableStyle).
		Headers("LANGUAGE", "FORMAL", "INFORMAL").
		Row("Chinese", "Nǐn hǎo", "Nǐ hǎo").
		Row("French", "Bonjour", "Salut").
		Row("Japanese", "こんにちは", "やあ").
		Row("Russian", "Zdravstvuyte", "Privet").
		Row("Spanish", "Hola", "¿Qué tal?")

	golden.RequireEqual(t, []byte(table.String()))
}

func TestTableHeightShrink(t *testing.T) {
	table := New().
		Height(8).
		Border(lipgloss.NormalBorder()).
		StyleFunc(TableStyle).
		Headers("LANGUAGE", "FORMAL", "INFORMAL").
		Row("Chinese", "Nǐn hǎo", "Nǐ hǎo").
		Row("French", "Bonjour", "Salut").
		Row("Japanese", "こんにちは", "やあ").
		Row("Russian", "Zdravstvuyte", "Privet").
		Row("Spanish", "Hola", "¿Qué tal?")

	golden.RequireEqual(t, []byte(table.String()))
}

func TestTableHeightMinimum(t *testing.T) {
	table := New().
		Height(0).
		Border(lipgloss.NormalBorder()).
		StyleFunc(TableStyle).
		Headers("ID", "LANGUAGE", "FORMAL", "INFORMAL").
		Row("1", "Chinese", "Nǐn hǎo", "Nǐ hǎo").
		Row("2", "French", "Bonjour", "Salut").
		Row("3", "Japanese", "こんにちは", "やあ").
		Row("4", "Russian", "Zdravstvuyte", "Privet").
		Row("5", "Spanish", "Hola", "¿Qué tal?")

	golden.RequireEqual(t, []byte(table.String()))
}

func TestTableHeightMinimumShowData(t *testing.T) {
	table := New().
		Height(0).
		Border(lipgloss.NormalBorder()).
		StyleFunc(TableStyle).
		Headers("LANGUAGE", "FORMAL", "INFORMAL").
		Row("Chinese", "Nǐn hǎo", "Nǐ hǎo")

	golden.RequireEqual(t, []byte(table.String()))
}

func TestTableHeightWithOffset(t *testing.T) {
	// This test exists to check for a bug/edge case when the table has an
	// offset and the height is set.

	table := New().
		Height(8).
		Border(lipgloss.NormalBorder()).
		StyleFunc(TableStyle).
		Headers("LANGUAGE", "FORMAL", "INFORMAL").
		Row("Chinese", "Nǐn hǎo", "Nǐ hǎo").
		Row("French", "Bonjour", "Salut").
		Row("Japanese", "こんにちは", "やあ").
		Row("Russian", "Zdravstvuyte", "Privet").
		Row("Spanish", "Hola", "¿Qué tal?").
		Offset(1)

	golden.RequireEqual(t, []byte(table.String()))
}

func TestStyleFunc(t *testing.T) {
	tests := []struct {
		name  string
		style StyleFunc
	}{
		{
			"RightAlignedTextWithMargins",
			func(row, col int) lipgloss.Style {
				switch {
				case row == HeaderRow:
					return lipgloss.NewStyle().Align(lipgloss.Center)
				default:
					return lipgloss.NewStyle().Margin(0, 1).Align(lipgloss.Right)
				}
			},
		},
		{
			"MarginAndPaddingSet",
			// this test case uses background colors to differentiate margins
			// and padding.
			func(row, col int) lipgloss.Style {
				switch {
				case row == HeaderRow:
					return lipgloss.NewStyle().Align(lipgloss.Center)
				default:
					return lipgloss.NewStyle().
						Padding(1).
						Margin(1).
						// keeping right-aligned text as it's the most likely to
						// be broken when truncated.
						Align(lipgloss.Right).
						Background(lipgloss.Color("#874bfc"))
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			table := New().
				Border(lipgloss.NormalBorder()).
				StyleFunc(tc.style).
				Headers("LANGUAGE", "FORMAL", "INFORMAL").
				Row("Chinese", "Nǐn hǎo", "Nǐ hǎo").
				Row("French", "Bonjour", "Salut").
				Row("Japanese", "こんにちは", "やあ").
				Row("Russian", "Zdravstvuyte", "Privet").
				Row("Spanish", "Hola", "¿Qué tal?")

			golden.RequireEqual(t, []byte(table.String()))
		})
	}
}

func TestClearRows(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("had to recover: %v", r)
		}
	}()

	table := New().
		Border(lipgloss.NormalBorder()).
		Headers("LANGUAGE", "FORMAL", "INFORMAL").
		Row("Chinese", "Nǐn hǎo", "Nǐ hǎo")
	table.ClearRows()
	table.Row("French", "Bonjour", "Salut")

	// String() will try to get the rows from table.data
	table.String()
}

func TestContentWrapping(t *testing.T) {
	tests := []struct {
		name      string
		headers   []string
		data      [][]string
		wrap      bool
		styleFunc StyleFunc
	}{
		{
			"LongRowContent",
			[]string{"Name", "Description", "Type", "Required", "Default"},
			[][]string{{"command", "A command to be executed inside the container to assess its health. Each space delimited token of the command is a separate array element. Commands exiting 0 are considered to be successful probes, whilst all other exit codes are considered failures.", "yes", "hello", "yep"}},
			true,
			TableStyle,
		},
		{
			"LongRowContentNoWrap",
			[]string{"Name", "Description", "Type", "Required", "Default"},
			[][]string{{"command", "A command to be executed inside the container to assess its health. Each space delimited token of the command is a separate array element. Commands exiting 0 are considered to be successful probes, whilst all other exit codes are considered failures.", "yes", "hello", "yep"}},
			false,
			TableStyle,
		},
		{
			"LongRowContentNoWrapNoMargins",
			[]string{"Name", "Description", "Type", "Required", "Default"},
			[][]string{{"command", "A command to be executed inside the container to assess its health. Each space delimited token of the command is a separate array element. Commands exiting 0 are considered to be successful probes, whilst all other exit codes are considered failures.", "yes", "hello", "yep"}},
			false,
			func(row, col int) lipgloss.Style {
				switch {
				case row == HeaderRow:
					return lipgloss.NewStyle().Padding(0).Align(lipgloss.Center)
				default:
					return lipgloss.NewStyle().Padding(0)
				}
			},
		},
		{
			"LongRowContentNoWrapCustomMargins",
			[]string{"Name", "Description", "Type", "Required", "Default"},
			[][]string{{"command", "A command to be executed inside the container to assess its health. Each space delimited token of the command is a separate array element. Commands exiting 0 are considered to be successful probes, whilst all other exit codes are considered failures.", "yes", "hello", "yep"}},
			false,
			func(row, col int) lipgloss.Style {
				switch {
				case row == HeaderRow:
					return lipgloss.NewStyle().Padding(0, 2).Align(lipgloss.Center)
				default:
					return lipgloss.NewStyle().Padding(0, 2)
				}
			},
		},
		{
			"MissingRowContent",
			[]string{"Name", "Description", "Type", "Required", "Default"},
			[][]string{{"command", "A command to be executed inside the container to assess its health. Each space delimited token of the command is a separate array element. Commands exiting 0 are considered to be successful probes, whilst all other exit codes are considered failures.", "yes", "", ""}},
			true,
			TableStyle,
		},
		{
			"LongHeaderContentLongAndShortRows",
			[]string{"Destination", "Why are you going on this trip? Is it a hot or cold climate?", "Affordability"},
			[][]string{
				{"Mexico", "I want to go somewhere hot, dry, and affordable. Mexico has really good food, just don't drink tap water!", "$"},
				{"New York", "I'm thinking about going during the Christmas season to check out Rockefeller center. Might be cold though...", "$$$"},
				{"California", "", "$$$"},
			},
			true,
			TableStyle,
		},
		{
			"LongTextDifferentLanguages",
			[]string{"Hello", "你好", "مرحبًا", "안녕하세요"},
			[][]string{
				{
					"Lorem ipsum dolor sit amet, regione detracto eos an. Has ei quidam hendrerit intellegebat, id tamquam iudicabit necessitatibus ius, at errem officiis hendrerit mei. Exerci noster at has, sit id tota convenire, vel ex rebum inciderint liberavisse. Quaeque delectus corrumpit cu cum.",
					`耐許ヱヨカハ調出あゆ監件び理別よン國給災レホチ権輝モエフ会割もフ響3現エツ文時しだびほ経機ムイメフ敗文ヨク現義なさド請情ゆじょて憶主管州けでふく。排ゃわつげ美刊ヱミ出見ツ南者オ抜豆ハトロネ論索モネニイ任償スヲ話破リヤヨ秒止口イセソス止央のさ食周健でてつだ官送ト読聴遊容ひるべ。際ぐドらづ市居ネムヤ研校35岩6繹ごわク報拐イ革深52球ゃレスご究東スラ衝3間ラ録占たス。
禁にンご忘康ざほぎル騰般ねど事超スんいう真表何カモ自浩ヲシミ図客線るふ静王ぱーま写村月掛焼詐面ぞゃ。昇強ごントほ価保キ族85岡モテ恋困ひりこな刊並せご出来ぼぎむう点目ヲウ止環公ニレ事応タス必書タメムノ当84無信升ちひょ。価ーぐ中客テサ告覧ヨトハ極整ラ得95稿はかラせ江利ス宏丸霊ミ考整ス静将ず業巨職ノラホ収嗅ざな。`,
					"شيء قد للحكومة والكوري الأوروبيّون, بوابة تعديل واعتلاء ضرب بـ. إذ أسر اتّجة اعلان, ٣٠ اكتوبر العصبة استمرار ومن. أفاق للسيطرة التاريخ، مع بحث, كلّ اتّجة القوى مع. فبعد ايطاليا، تم حتى, لكل تم جسيمة الإحتفاظ وباستثناء, عل فرنسا وانتهاءً الإقتصادية عرض. ونتج دأبوا إحكام بال إذ. لغات عملية وتم مع, وصل بداية وبغطاء البرية بل, أي قررت بلاده فكانت حدى",
					"版応道潟部中幕爆営報門案名見壌府。博健必権次覧編仕断青場内凄新東深簿代供供。守聞書神秀同浜東波恋闘秀。未格打好作器来利阪持西焦朝三女。権幽問季負娘購合旧資健載員式活陸。未倍校朝遺続術吉迎暮広知角亡志不説空住。法省当死年勝絡聞方北投健。室分性山天態意画詳知浅方裁。変激伝阜中野品省載嗅闘額端反。中必台際造事寄民経能前作臓",
					"각급 선거관리위원회의 조직·직무범위 기타 필요한 사항은 법률로 정한다. 임시회의 회기는 30일을 초과할 수 없다. 국가는 여자의 복지와 권익의 향상을 위하여 노력하여야 한다. 국군의 조직과 편성은 법률로 정한다.",
				},
			},
			true,
			TableStyle,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			table := New().
				Headers(tc.headers...).
				Rows(tc.data...).
				Width(80).
				StyleFunc(tc.styleFunc).
				Wrap(tc.wrap)

			golden.RequireEqual(t, []byte(table.String()))
		})
	}
}

func TestContentWrapping_WithPadding(t *testing.T) {
	tests := []struct {
		name    string
		headers []string
		data    [][]string
	}{
		{
			"LongRowContent",
			[]string{"Name", "Description", "Type", "Required", "Default"},
			[][]string{{"command", "A command to be executed inside the container to assess its health. Each space delimited token of the command is a separate array element. Commands exiting 0 are considered to be successful probes, whilst all other exit codes are considered failures.", "yes", "hello", "yep"}},
		},
		{
			"MissingRowContent",
			[]string{"Name", "Description", "Type", "Required", "Default"},
			[][]string{{"command", "A command to be executed inside the container to assess its health. Each space delimited token of the command is a separate array element. Commands exiting 0 are considered to be successful probes, whilst all other exit codes are considered failures.", "yes", "", ""}},
		},
		{
			"LongHeaderContentLongAndShortRows",
			[]string{"Destination", "Why are you going on this trip? Is it a hot or cold climate?", "Affordability"},
			[][]string{
				{"Mexico", "I want to go somewhere hot, dry, and affordable. Mexico has really good food, just don't drink tap water!", "$"},
				{"New York", "I'm thinking about going during the Christmas season to check out Rockefeller center. Might be cold though...", "$$$"},
				{"California", "", "$$$"},
			},
		},
		{
			"LongTextDifferentLanguages",
			[]string{"Hello", "你好", "مرحبًا", "안녕하세요"},
			[][]string{
				{
					"",
					`耐許ヱヨカハ調出あゆ監件び理別よン國給災レホチ権輝モエフ会割もフ響3現エツ文時しだびほ経機ムイメフ敗文ヨク現義なさド請情ゆじょて憶主管州けでふく。排ゃわつげ美刊ヱミ出見ツ南者オ抜豆ハトロネ論索モネニイ任償スヲ話破リヤヨ秒止口イセソス止央のさ食周健でてつだ官送ト読聴遊容ひるべ。際ぐドらづ市居ネムヤ研校35岩6繹ごわク報拐イ革深52球ゃレスご究東スラ衝3間ラ録占たス。
禁にンご忘康ざほぎル騰般ねど事超スんいう真表何カモ自浩ヲシミ図客線るふ静王ぱーま写村月掛焼詐面ぞゃ。昇強ごントほ価保キ族85岡モテ恋困ひりこな刊並せご出来ぼぎむう点目ヲウ止環公ニレ事応タス必書タメムノ当84無信升ちひょ。価ーぐ中客テサ告覧ヨトハ極整ラ得95稿はかラせ江利ス宏丸霊ミ考整ス静将ず業巨職ノラホ収嗅ざな。`,
					"شيء قد للحكومة والكوري الأوروبيّون, بوابة تعديل واعتلاء ضرب بـ. إذ أسر اتّجة اعلان, ٣٠ اكتوبر العصبة استمرار ومن. أفاق للسيطرة التاريخ، مع بحث, كلّ اتّجة القوى مع. فبعد ايطاليا، تم حتى, لكل تم جسيمة الإحتفاظ وباستثناء, عل فرنسا وانتهاءً الإقتصادية عرض. ونتج دأبوا إحكام بال إذ. لغات عملية وتم مع, وصل بداية وبغطاء البرية بل, أي قررت بلاده فكانت حدى",
					"版応道潟部中幕爆営報門案名見壌府。博健必権次覧編仕断青場内凄新東深簿代供供。守聞書神秀同浜東波恋闘秀。未格打好作器来利阪持西焦朝三女。権幽問季負娘購合旧資健載員式活陸。未倍校朝遺続術吉迎暮広知角亡志不説空住。法省当死年勝絡聞方北投健。室分性山天態意画詳知浅方裁。変激伝阜中野品省載嗅闘額端反。中必台際造事寄民経能前作臓",
					"각급 선거관리위원회의 조직·직무범위 기타 필요한 사항은 법률로 정한다. 임시회의 회기는 30일을 초과할 수 없다. 국가는 여자의 복지와 권익의 향상을 위하여 노력하여야 한다. 국군의 조직과 편성은 법률로 정한다.",
				},
			},
		},
	}

	defaultWidth := 80
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			table := New().
				Headers(tc.headers...).
				Rows(tc.data...).
				StyleFunc(func(_, col int) lipgloss.Style {
					return lipgloss.NewStyle().Padding(0, 1)
				})
			table.Width(defaultWidth)

			// check total width.
			if got := lipgloss.Width(table.String()); got != defaultWidth {
				t.Fatalf("Table is not the correct width. got %d, want %d", got, defaultWidth)
				t.Log(table.String())
			}
			golden.RequireEqual(t, []byte(table.String()))
		})
	}
}

func TestContentWrapping_WithMargins(t *testing.T) {
	tests := []struct {
		name    string
		headers []string
		data    [][]string
	}{
		{
			"LongRowContent",
			[]string{"Name", "Description", "Type", "Required", "Default"},
			[][]string{{"command", "A command to be executed inside the container to assess its health. Each space delimited token of the command is a separate array element. Commands exiting 0 are considered to be successful probes, whilst all other exit codes are considered failures.", "yes", "hello", "yep"}},
		},
		{
			"MissingRowContent",
			[]string{"Name", "Description", "Type", "Required", "Default"},
			[][]string{{"command", "A command to be executed inside the container to assess its health. Each space delimited token of the command is a separate array element. Commands exiting 0 are considered to be successful probes, whilst all other exit codes are considered failures.", "yes", "", ""}},
		},
		{
			"LongHeaderContentLongAndShortRows",
			[]string{"Destination", "Why are you going on this trip? Is it a hot or cold climate?", "Affordability"},
			[][]string{
				{"Mexico", "I want to go somewhere hot, dry, and affordable. Mexico has really good food, just don't drink tap water!", "$"},
				{"New York", "I'm thinking about going during the Christmas season to check out Rockefeller center. Might be cold though...", "$$$"},
				{"California", "", "$$$"},
			},
		},
		{
			"LongTextDifferentLanguages",
			[]string{"Hello", "你好", "مرحبًا", "안녕하세요"},
			[][]string{
				{
					"Lorem ipsum dolor sit amet, regione detracto eos an. Has ei quidam hendrerit intellegebat, id tamquam iudicabit necessitatibus ius, at errem officiis hendrerit mei. Exerci noster at has, sit id tota convenire, vel ex rebum inciderint liberavisse. Quaeque delectus corrumpit cu cum.",
					`耐許ヱヨカハ調出あゆ監件び理別よン國給災レホチ権輝モエフ会割もフ響3現エツ文時しだびほ経機ムイメフ敗文ヨク現義なさド請情ゆじょて憶主管州けでふく。排ゃわつげ美刊ヱミ出見ツ南者オ抜豆ハトロネ論索モネニイ任償スヲ話破リヤヨ秒止口イセソス止央のさ食周健でてつだ官送ト読聴遊容ひるべ。際ぐドらづ市居ネムヤ研校35岩6繹ごわク報拐イ革深52球ゃレスご究東スラ衝3間ラ録占たス。
禁にンご忘康ざほぎル騰般ねど事超スんいう真表何カモ自浩ヲシミ図客線るふ静王ぱーま写村月掛焼詐面ぞゃ。昇強ごントほ価保キ族85岡モテ恋困ひりこな刊並せご出来ぼぎむう点目ヲウ止環公ニレ事応タス必書タメムノ当84無信升ちひょ。価ーぐ中客テサ告覧ヨトハ極整ラ得95稿はかラせ江利ス宏丸霊ミ考整ス静将ず業巨職ノラホ収嗅ざな。`,
					"شيء قد للحكومة والكوري الأوروبيّون, بوابة تعديل واعتلاء ضرب بـ. إذ أسر اتّجة اعلان, ٣٠ اكتوبر العصبة استمرار ومن. أفاق للسيطرة التاريخ، مع بحث, كلّ اتّجة القوى مع. فبعد ايطاليا، تم حتى, لكل تم جسيمة الإحتفاظ وباستثناء, عل فرنسا وانتهاءً الإقتصادية عرض. ونتج دأبوا إحكام بال إذ. لغات عملية وتم مع, وصل بداية وبغطاء البرية بل, أي قررت بلاده فكانت حدى",
					"版応道潟部中幕爆営報門案名見壌府。博健必権次覧編仕断青場内凄新東深簿代供供。守聞書神秀同浜東波恋闘秀。未格打好作器来利阪持西焦朝三女。権幽問季負娘購合旧資健載員式活陸。未倍校朝遺続術吉迎暮広知角亡志不説空住。法省当死年勝絡聞方北投健。室分性山天態意画詳知浅方裁。変激伝阜中野品省載嗅闘額端反。中必台際造事寄民経能前作臓",
					"각급 선거관리위원회의 조직·직무범위 기타 필요한 사항은 법률로 정한다. 임시회의 회기는 30일을 초과할 수 없다. 국가는 여자의 복지와 권익의 향상을 위하여 노력하여야 한다. 국군의 조직과 편성은 법률로 정한다.",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			table := New().
				Headers(tc.headers...).
				Rows(tc.data...).
				StyleFunc(func(row, col int) lipgloss.Style {
					return lipgloss.NewStyle().Margin(0, 4)
				})
			table.Width(80)
			golden.RequireEqual(t, []byte(table.String()))
		})
	}
}

func TestContentWrapping_ColumnWidth(t *testing.T) {
	tests := []struct {
		name    string
		headers []string
		data    [][]string
	}{
		{
			"LongRowContent",
			[]string{"Name", "Description", "Type", "Required", "Default"},
			[][]string{{"command", "A command to be executed inside the container to assess its health. Each space delimited token of the command is a separate array element. Commands exiting 0 are considered to be successful probes, whilst all other exit codes are considered failures.", "yes", "hello", "yep"}},
		},
		{
			"MissingRowContent",
			[]string{"Name", "Description", "Type", "Required", "Default"},
			[][]string{{"command", "A command to be executed inside the container to assess its health. Each space delimited token of the command is a separate array element. Commands exiting 0 are considered to be successful probes, whilst all other exit codes are considered failures.", "yes", "", ""}},
		},
		{
			"LongHeaderContentLongAndShortRows",
			[]string{"Destination", "Why are you going on this trip? Is it a hot or cold climate?", "Affordability"},
			[][]string{
				{"Mexico", "I want to go somewhere hot, dry, and affordable. Mexico has really good food, just don't drink tap water!", "$"},
				{"New York", "I'm thinking about going during the Christmas season to check out Rockefeller center. Might be cold though...", "$$$"},
				{"California", "", "$$$"},
			},
		},
		{
			"LongTextDifferentLanguages",
			[]string{"Hello", "你好", "مرحبًا", "안녕하세요"},
			[][]string{
				{
					"Lorem ipsum dolor sit amet, regione detracto eos an. Has ei quidam hendrerit intellegebat, id tamquam iudicabit necessitatibus ius, at errem officiis hendrerit mei. Exerci noster at has, sit id tota convenire, vel ex rebum inciderint liberavisse. Quaeque delectus corrumpit cu cum.",
					`耐許ヱヨカハ調出あゆ監件び理別よン國給災レホチ権輝モエフ会割もフ響3現エツ文時しだびほ経機ムイメフ敗文ヨク現義なさド請情ゆじょて憶主管州けでふく。排ゃわつげ美刊ヱミ出見ツ南者オ抜豆ハトロネ論索モネニイ任償スヲ話破リヤヨ秒止口イセソス止央のさ食周健でてつだ官送ト読聴遊容ひるべ。際ぐドらづ市居ネムヤ研校35岩6繹ごわク報拐イ革深52球ゃレスご究東スラ衝3間ラ録占たス。
禁にンご忘康ざほぎル騰般ねど事超スんいう真表何カモ自浩ヲシミ図客線るふ静王ぱーま写村月掛焼詐面ぞゃ。昇強ごントほ価保キ族85岡モテ恋困ひりこな刊並せご出来ぼぎむう点目ヲウ止環公ニレ事応タス必書タメムノ当84無信升ちひょ。価ーぐ中客テサ告覧ヨトハ極整ラ得95稿はかラせ江利ス宏丸霊ミ考整ス静将ず業巨職ノラホ収嗅ざな。`,
					"شيء قد للحكومة والكوري الأوروبيّون, بوابة تعديل واعتلاء ضرب بـ. إذ أسر اتّجة اعلان, ٣٠ اكتوبر العصبة استمرار ومن. أفاق للسيطرة التاريخ، مع بحث, كلّ اتّجة القوى مع. فبعد ايطاليا، تم حتى, لكل تم جسيمة الإحتفاظ وباستثناء, عل فرنسا وانتهاءً الإقتصادية عرض. ونتج دأبوا إحكام بال إذ. لغات عملية وتم مع, وصل بداية وبغطاء البرية بل, أي قررت بلاده فكانت حدى",
					"版応道潟部中幕爆営報門案名見壌府。博健必権次覧編仕断青場内凄新東深簿代供供。守聞書神秀同浜東波恋闘秀。未格打好作器来利阪持西焦朝三女。権幽問季負娘購合旧資健載員式活陸。未倍校朝遺続術吉迎暮広知角亡志不説空住。法省当死年勝絡聞方北投健。室分性山天態意画詳知浅方裁。変激伝阜中野品省載嗅闘額端反。中必台際造事寄民経能前作臓",
					"각급 선거관리위원회의 조직·직무범위 기타 필요한 사항은 법률로 정한다. 임시회의 회기는 30일을 초과할 수 없다. 국가는 여자의 복지와 권익의 향상을 위하여 노력하여야 한다. 국군의 조직과 편성은 법률로 정한다.",
				},
			},
		},
	}
	defaultWidth := 80
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			table := New().
				Headers(tc.headers...).
				Rows(tc.data...).
				StyleFunc(func(row, col int) lipgloss.Style {
					// If we set a specific cell width, it should count for all rows
					// in that column.
					if row == 0 && col == 1 {
						return lipgloss.NewStyle().Width(30)
					}
					// Set a column's width directly.
					if col == 2 {
						return lipgloss.NewStyle().Width(5)
					}
					return lipgloss.NewStyle()
				})
			table.Width(defaultWidth)
			// check total width.
			if got := lipgloss.Width(table.String()); got != defaultWidth {
				t.Log(table.String())
				t.Fatalf("Table is not the correct width. got %d, want %d", got, defaultWidth)
			}

			// check that width is overridden with a small value.
			if table.widths[2] != 5 {
				t.Log(table.String())
				t.Fatalf("Did not set correct width value at column at index %d.\ngot %d, want %d", 2, table.widths[2], 5)
			}

			// check that width is overridden with a wide value.
			if table.widths[1] != 30 {
				t.Log(table.String())
				t.Fatalf("Did not set correct width value at column at index %d.\ngot %d, want %d", 1, table.widths[1], 30)
			}

			t.Log(table.widths[2])
			golden.RequireEqual(t, []byte(table.String()))
		})
	}
}

// Test truncation for overflow and no wrap when combined.
func TestTableOverFlowNoWrap(t *testing.T) {
	// LongTextDifferentLanguages
	headers := []string{"Hello", "你好", "مرحبًا", "안녕하세요"}
	data := [][]string{
		{
			"Lorem ipsum dolor sit amet, regione detracto eos an. Has ei quidam hendrerit intellegebat, id tamquam iudicabit necessitatibus ius, at errem officiis hendrerit mei. Exerci noster at has, sit id tota convenire, vel ex rebum inciderint liberavisse. Quaeque delectus corrumpit cu cum.",
			`耐許ヱヨカハ調出あゆ監件び理別よン國給災レホチ権輝モエフ会割もフ響3現エツ文時しだびほ経機ムイメフ敗文ヨク現義なさド請情ゆじょて憶主管州けでふく。排ゃわつげ美刊ヱミ出見ツ南者オ抜豆ハトロネ論索モネニイ任償スヲ話破リヤヨ秒止口イセソス止央のさ食周健でてつだ官送ト読聴遊容ひるべ。際ぐドらづ市居ネムヤ研校35岩6繹ごわク報拐イ革深52球ゃレスご究東スラ衝3間ラ録占たス。
禁にンご忘康ざほぎル騰般ねど事超スんいう真表何カモ自浩ヲシミ図客線るふ静王ぱーま写村月掛焼詐面ぞゃ。昇強ごントほ価保キ族85岡モテ恋困ひりこな刊並せご出来ぼぎむう点目ヲウ止環公ニレ事応タス必書タメムノ当84無信升ちひょ。価ーぐ中客テサ告覧ヨトハ極整ラ得95稿はかラせ江利ス宏丸霊ミ考整ス静将ず業巨職ノラホ収嗅ざな。`,
			"شيء قد للحكومة والكوري الأوروبيّون, بوابة تعديل واعتلاء ضرب بـ. إذ أسر اتّجة اعلان, ٣٠ اكتوبر العصبة استمرار ومن. أفاق للسيطرة التاريخ، مع بحث, كلّ اتّجة القوى مع. فبعد ايطاليا، تم حتى, لكل تم جسيمة الإحتفاظ وباستثناء, عل فرنسا وانتهاءً الإقتصادية عرض. ونتج دأبوا إحكام بال إذ. لغات عملية وتم مع, وصل بداية وبغطاء البرية بل, أي قررت بلاده فكانت حدى",
			"각급 선거관리위원회의 조직·직무범위 기타 필요한 사항은 법률로 정한다. 임시회의 회기는 30일을 초과할 수 없다. 국가는 여자의 복지와 권익의 향상을 위하여 노력하여야 한다. 국군의 조직과 편성은 법률로 정한다.",
			"版応道潟部中幕爆営報門案名見壌府。博健必権次覧編仕断青場内凄新東深簿代供供。守聞書神秀同浜東波恋闘秀。未格打好作器来利阪持西焦朝三女。権幽問季負娘購合旧資健載員式活陸。未倍校朝遺続術吉迎暮広知角亡志不説空住。法省当死年勝絡聞方北投健。室分性山天態意画詳知浅方裁。変激伝阜中野品省載嗅闘額端反。中必台際造事寄民経能前作臓",
		},
		{"Welcome", "いらっしゃいませ", "مرحباً", "환영", "欢迎"},
		{"Goodbye", "さようなら", "مع السلامة", "안녕히 가세요", "再见"},
	}
	tableHeight := 6
	table := New().
		Headers(headers...).
		Rows(data...).
		StyleFunc(TableStyle).
		Height(tableHeight).
		Width(80).
		Wrap(false)
	if got := lipgloss.Height(table.String()); got != tableHeight {
		t.Fatalf("got the wrong height. got %d, want %d", got, tableHeight)
	}
	golden.RequireEqual(t, []byte(table.String()))
}

func TestCarriageReturn(t *testing.T) {
	data := [][]string{
		{"a0", "b0", "c0", "d0"},
		{"a1", "b1.0\r\nb1.1\r\nb1.2\r\nb1.3\r\nb1.4\r\nb1.5\r\nb1.6", "c1", "d1"},
		{"a2", "b2", "c2", "d2"},
		{"a3", "b3", "c3", "d3"},
	}
	table := New().Rows(data...).Border(lipgloss.NormalBorder())

	golden.RequireEqual(t, []byte(table.String()))
}

func TestTableShrinkWithOffset(t *testing.T) {
	rows := [][]string{
		{"1", "Tokyo", "Japan", "37,274,000"},
		{"2", "Delhi", "India", "32,065,760"},
		{"3", "Shanghai", "China", "28,516,904"},
		{"4", "Dhaka", "Bangladesh", "22,478,116"},
		{"5", "São Paulo", "Brazil", "22,429,800"},
		{"6", "Mexico City", "Mexico", "22,085,140"},
		{"7", "Cairo", "Egypt", "21,750,020"},
		{"8", "Beijing", "China", "21,333,332"},
		{"9", "Mumbai", "India", "20,961,472"},
		{"10", "Osaka", "Japan", "19,059,856"},
		{"11", "Chongqing", "China", "16,874,740"},
		{"12", "Karachi", "Pakistan", "16,839,950"},
		{"13", "Istanbul", "Turkey", "15,636,243"},
		{"14", "Kinshasa", "DR Congo", "15,628,085"},
		{"15", "Lagos", "Nigeria", "15,387,639"},
		{"16", "Buenos Aires", "Argentina", "15,369,919"},
		{"17", "Kolkata", "India", "15,133,888"},
		{"18", "Manila", "Philippines", "14,406,059"},
		{"19", "Tianjin", "China", "14,011,828"},
		{"20", "Guangzhou", "China", "13,964,637"},
		{"21", "Rio De Janeiro", "Brazil", "13,634,274"},
		{"22", "Lahore", "Pakistan", "13,541,764"},
		{"23", "Bangalore", "India", "13,193,035"},
		{"24", "Shenzhen", "China", "12,831,330"},
		{"25", "Moscow", "Russia", "12,640,818"},
		{"26", "Chennai", "India", "11,503,293"},
		{"27", "Bogota", "Colombia", "11,344,312"},
		{"28", "Paris", "France", "11,142,303"},
		{"29", "Jakarta", "Indonesia", "11,074,811"},
		{"30", "Lima", "Peru", "11,044,607"},
		{"31", "Bangkok", "Thailand", "10,899,698"},
		{"32", "Hyderabad", "India", "10,534,418"},
		{"33", "Seoul", "South Korea", "9,975,709"},
		{"34", "Nagoya", "Japan", "9,571,596"},
		{"35", "London", "United Kingdom", "9,540,576"},
		{"36", "Chengdu", "China", "9,478,521"},
		{"37", "Nanjing", "China", "9,429,381"},
		{"38", "Tehran", "Iran", "9,381,546"},
		{"39", "Ho Chi Minh City", "Vietnam", "9,077,158"},
		{"40", "Luanda", "Angola", "8,952,496"},
		{"41", "Wuhan", "China", "8,591,611"},
		{"42", "Xi An Shaanxi", "China", "8,537,646"},
		{"43", "Ahmedabad", "India", "8,450,228"},
		{"44", "Kuala Lumpur", "Malaysia", "8,419,566"},
		{"45", "New York City", "United States", "8,177,020"},
		{"46", "Hangzhou", "China", "8,044,878"},
		{"47", "Surat", "India", "7,784,276"},
		{"48", "Suzhou", "China", "7,764,499"},
		{"49", "Hong Kong", "Hong Kong", "7,643,256"},
		{"50", "Riyadh", "Saudi Arabia", "7,538,200"},
		{"51", "Shenyang", "China", "7,527,975"},
		{"52", "Baghdad", "Iraq", "7,511,920"},
		{"53", "Dongguan", "China", "7,511,851"},
		{"54", "Foshan", "China", "7,497,263"},
		{"55", "Dar Es Salaam", "Tanzania", "7,404,689"},
		{"56", "Pune", "India", "6,987,077"},
		{"57", "Santiago", "Chile", "6,856,939"},
		{"58", "Madrid", "Spain", "6,713,557"},
		{"59", "Haerbin", "China", "6,665,951"},
		{"60", "Toronto", "Canada", "6,312,974"},
		{"61", "Belo Horizonte", "Brazil", "6,194,292"},
		{"62", "Khartoum", "Sudan", "6,160,327"},
		{"63", "Johannesburg", "South Africa", "6,065,354"},
		{"64", "Singapore", "Singapore", "6,039,577"},
		{"65", "Dalian", "China", "5,930,140"},
		{"66", "Qingdao", "China", "5,865,232"},
		{"67", "Zhengzhou", "China", "5,690,312"},
		{"68", "Ji Nan Shandong", "China", "5,663,015"},
		{"69", "Barcelona", "Spain", "5,658,472"},
		{"70", "Saint Petersburg", "Russia", "5,535,556"},
		{"71", "Abidjan", "Ivory Coast", "5,515,790"},
		{"72", "Yangon", "Myanmar", "5,514,454"},
		{"73", "Fukuoka", "Japan", "5,502,591"},
		{"74", "Alexandria", "Egypt", "5,483,605"},
		{"75", "Guadalajara", "Mexico", "5,339,583"},
		{"76", "Ankara", "Turkey", "5,309,690"},
		{"77", "Chittagong", "Bangladesh", "5,252,842"},
		{"78", "Addis Ababa", "Ethiopia", "5,227,794"},
		{"79", "Melbourne", "Australia", "5,150,766"},
		{"80", "Nairobi", "Kenya", "5,118,844"},
		{"81", "Hanoi", "Vietnam", "5,067,352"},
		{"82", "Sydney", "Australia", "5,056,571"},
		{"83", "Monterrey", "Mexico", "5,036,535"},
		{"84", "Changsha", "China", "4,809,887"},
		{"85", "Brasilia", "Brazil", "4,803,877"},
		{"86", "Cape Town", "South Africa", "4,800,954"},
		{"87", "Jiddah", "Saudi Arabia", "4,780,740"},
		{"88", "Urumqi", "China", "4,710,203"},
		{"89", "Kunming", "China", "4,657,381"},
		{"90", "Changchun", "China", "4,616,002"},
		{"91", "Hefei", "China", "4,496,456"},
		{"92", "Shantou", "China", "4,490,411"},
		{"93", "Xinbei", "Taiwan", "4,470,672"},
		{"94", "Kabul", "Afghanistan", "4,457,882"},
		{"95", "Ningbo", "China", "4,405,292"},
		{"96", "Tel Aviv", "Israel", "4,343,584"},
		{"97", "Yaounde", "Cameroon", "4,336,670"},
		{"98", "Rome", "Italy", "4,297,877"},
		{"99", "Shijiazhuang", "China", "4,285,135"},
		{"100", "Montreal", "Canada", "4,276,526"},
	}
	table := New().
		Rows(rows...).
		Offset(80).
		Height(45)

	got := lipgloss.Height(table.String())
	if got != table.height {
		t.Fatalf("expected the height to be %d with an offset of %d. got: table with height %d\n%s", table.height, table.offset, got, table.String())
	}
}

func TestBorderStyles(t *testing.T) {
	rows := [][]string{
		{"Chinese", "Nǐn hǎo", "Nǐ hǎo"},
		{"French", "Bonjour", "Salut"},
		{"Japanese", "こんにちは", "やあ"},
		{"Russian", "Zdravstvuyte", "Privet"},
		{"Spanish", "Hola", "¿Qué tal?"},
	}

	tests := []struct {
		name             string
		borderFunc       func() lipgloss.Border
		topBottomBorders bool
	}{
		{"NormalBorder", lipgloss.NormalBorder, true},
		{"RoundedBorder", lipgloss.RoundedBorder, true},
		{"BlockBorder", lipgloss.BlockBorder, true},
		{"ThickBorder", lipgloss.ThickBorder, true},
		{"HiddenBorder", lipgloss.HiddenBorder, true},
		{"MarkdownBorder", lipgloss.MarkdownBorder, false},
		{"ASCIIBorder", lipgloss.ASCIIBorder, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			table := New().
				StyleFunc(TableStyle).
				Border(test.borderFunc()).
				Headers("LANGUAGE", "FORMAL", "INFORMAL").
				Rows(rows...).
				BorderTop(test.topBottomBorders).
				BorderBottom(test.topBottomBorders)

			golden.RequireEqual(t, []byte(table.String()))
		})
	}
}

// Examples

func ExampleTable_Wrap() {
	// LongTextDifferentLanguages
	headers := []string{"Hello", "你好", "مرحبًا", "안녕하세요"}
	data := [][]string{
		{
			"Lorem ipsum dolor sit amet, regione detracto eos an. Has ei quidam hendrerit intellegebat, id tamquam iudicabit necessitatibus ius, at errem officiis hendrerit mei. Exerci noster at has, sit id tota convenire, vel ex rebum inciderint liberavisse. Quaeque delectus corrumpit cu cum.",
			`耐許ヱヨカハ調出あゆ監件び理別よン國給災レホチ権輝モエフ会割もフ響3現エツ文時しだびほ経機ムイメフ敗文ヨク現義なさド請情ゆじょて憶主管州けでふく。排ゃわつげ美刊ヱミ出見ツ南者オ抜豆ハトロネ論索モネニイ任償スヲ話破リヤヨ秒止口イセソス止央のさ食周健でてつだ官送ト読聴遊容ひるべ。際ぐドらづ市居ネムヤ研校35岩6繹ごわク報拐イ革深52球ゃレスご究東スラ衝3間ラ録占たス。
禁にンご忘康ざほぎル騰般ねど事超スんいう真表何カモ自浩ヲシミ図客線るふ静王ぱーま写村月掛焼詐面ぞゃ。昇強ごントほ価保キ族85岡モテ恋困ひりこな刊並せご出来ぼぎむう点目ヲウ止環公ニレ事応タス必書タメムノ当84無信升ちひょ。価ーぐ中客テサ告覧ヨトハ極整ラ得95稿はかラせ江利ス宏丸霊ミ考整ス静将ず業巨職ノラホ収嗅ざな。`,
			"شيء قد للحكومة والكوري الأوروبيّون, بوابة تعديل واعتلاء ضرب بـ. إذ أسر اتّجة اعلان, ٣٠ اكتوبر العصبة استمرار ومن. أفاق للسيطرة التاريخ، مع بحث, كلّ اتّجة القوى مع. فبعد ايطاليا، تم حتى, لكل تم جسيمة الإحتفاظ وباستثناء, عل فرنسا وانتهاءً الإقتصادية عرض. ونتج دأبوا إحكام بال إذ. لغات عملية وتم مع, وصل بداية وبغطاء البرية بل, أي قررت بلاده فكانت حدى",
			"版応道潟部中幕爆営報門案名見壌府。博健必権次覧編仕断青場内凄新東深簿代供供。守聞書神秀同浜東波恋闘秀。未格打好作器来利阪持西焦朝三女。権幽問季負娘購合旧資健載員式活陸。未倍校朝遺続術吉迎暮広知角亡志不説空住。法省当死年勝絡聞方北投健。室分性山天態意画詳知浅方裁。変激伝阜中野品省載嗅闘額端反。中必台際造事寄民経能前作臓",
			"각급 선거관리위원회의 조직·직무범위 기타 필요한 사항은 법률로 정한다. 임시회의 회기는 30일을 초과할 수 없다. 국가는 여자의 복지와 권익의 향상을 위하여 노력하여야 한다. 국군의 조직과 편성은 법률로 정한다.",
		},
	}

	table := New().
		Headers(headers...).
		Rows(data...).
		StyleFunc(TableStyle).
		Width(80).
		Wrap(false)
	fmt.Println(table.String())

	table.Wrap(true)
	fmt.Println(table.String())

	// Output:
	// ╭──────────────┬───────────────┬───────────────┬───────────────┬───────────────╮
	// │    Hello     │     你好      │     مرحبًا     │  안녕하세요   │               │
	// ├──────────────┼───────────────┼───────────────┼───────────────┼───────────────┤
	// │ Lorem ipsum… │ 耐許ヱヨカハ… │ شيء قد للحكو… │ 版応道潟部中… │ 각급 선거관…  │
	// ╰──────────────┴───────────────┴───────────────┴───────────────┴───────────────╯
	// ╭──────────────┬───────────────┬───────────────┬───────────────┬───────────────╮
	// │    Hello     │     你好      │     مرحبًا     │  안녕하세요   │               │
	// ├──────────────┼───────────────┼───────────────┼───────────────┼───────────────┤
	// │ Lorem ipsum  │ 耐許ヱヨカハ  │ شيء قد        │ 版応道潟部中  │ 각급          │
	// │ dolor sit    │ 調出あゆ監件  │ للحكومة       │ 幕爆営報門案  │ 선거관리위원  │
	// │ amet,        │ び理別よン國  │ والكوري       │ 名見壌府。博  │ 회의          │
	// │ regione      │ 給災レホチ権  │ الأوروبيّون,   │ 健必権次覧編  │ 조직·직무범위 │
	// │ detracto eos │ 輝モエフ会割  │ بوابة تعديل   │ 仕断青場内凄  │ 기타 필요한   │
	// │ an. Has ei   │ もフ響3現エツ │ واعتلاء ضرب   │ 新東深簿代供  │ 사항은 법률로 │
	// │ quidam       │ 文時しだびほ  │ بـ. إذ أسر    │ 供。守聞書神  │ 정한다.       │
	// │ hendrerit    │ 経機ムイメフ  │ اتّجة اعلان,   │ 秀同浜東波恋  │ 임시회의      │
	// │ intellegebat │ 敗文ヨク現義  │ ٣٠ اكتوبر     │ 闘秀。未格打  │ 회기는 30일을 │
	// │ , id tamquam │ なさド請情ゆ  │ العصبة        │ 好作器来利阪  │ 초과할 수     │
	// │ iudicabit    │ じょて憶主管  │ استمرار ومن.  │ 持西焦朝三女  │ 없다. 국가는  │
	// │ necessitatib │ 州けでふく。  │ أفاق للسيطرة  │ 。権幽問季負  │ 여자의 복지와 │
	// │ us ius, at   │ 排ゃわつげ美  │ التاريخ، مع   │ 娘購合旧資健  │ 권익의 향상을 │
	// │ errem        │ 刊ヱミ出見ツ  │ بحث, كلّ اتّجة  │ 載員式活陸。  │ 위하여        │
	// │ officiis     │ 南者オ抜豆ハ  │ القوى مع.     │ 未倍校朝遺続  │ 노력하여야    │
	// │ hendrerit    │ トロネ論索モ  │ فبعد ايطاليا، │ 術吉迎暮広知  │ 한다. 국군의  │
	// │ mei. Exerci  │ ネニイ任償ス  │ تم حتى, لكل   │ 角亡志不説空  │ 조직과 편성은 │
	// │ noster at    │ ヲ話破リヤヨ  │ تم جسيمة      │ 住。法省当死  │ 법률로        │
	// │ has, sit id  │ 秒止口イセソ  │ الإحتفاظ      │ 年勝絡聞方北  │ 정한다.       │
	// │ tota         │ ス止央のさ食  │ وباستثناء, عل │ 投健。室分性  │               │
	// │ convenire,   │ 周健でてつだ  │ فرنسا وانتهاءً │ 山天態意画詳  │               │
	// │ vel ex rebum │ 官送ト読聴遊  │ الإقتصادية    │ 知浅方裁。変  │               │
	// │ inciderint   │ 容ひるべ。際  │ عرض. ونتج     │ 激伝阜中野品  │               │
	// │ liberavisse. │ ぐドらづ市居  │ دأبوا إحكام   │ 省載嗅闘額端  │               │
	// │ Quaeque      │ ネムヤ研校35  │ بال إذ. لغات  │ 反。中必台際  │               │
	// │ delectus     │ 岩6繹ごわク報 │ عملية وتم مع, │ 造事寄民経能  │               │
	// │ corrumpit cu │ 拐イ革深52球  │ وصل بداية     │ 前作臓        │               │
	// │ cum.         │ ゃレスご究東  │ وبغطاء البرية │               │               │
	// │              │ スラ衝3間ラ録 │ بل, أي قررت   │               │               │
	// │              │ 占たス。      │ بلاده فكانت   │               │               │
	// │              │ 禁にンご忘康  │ حدى           │               │               │
	// │              │ ざほぎル騰般  │               │               │               │
	// │              │ ねど事超スん  │               │               │               │
	// │              │ いう真表何カ  │               │               │               │
	// │              │ モ自浩ヲシミ  │               │               │               │
	// │              │ 図客線るふ静  │               │               │               │
	// │              │ 王ぱーま写村  │               │               │               │
	// │              │ 月掛焼詐面ぞ  │               │               │               │
	// │              │ ゃ。昇強ごン  │               │               │               │
	// │              │ トほ価保キ族8 │               │               │               │
	// │              │ 5岡モテ恋困ひ │               │               │               │
	// │              │ りこな刊並せ  │               │               │               │
	// │              │ ご出来ぼぎむ  │               │               │               │
	// │              │ う点目ヲウ止  │               │               │               │
	// │              │ 環公ニレ事応  │               │               │               │
	// │              │ タス必書タメ  │               │               │               │
	// │              │ ムノ当84無信  │               │               │               │
	// │              │ 升ちひょ。価  │               │               │               │
	// │              │ ーぐ中客テサ  │               │               │               │
	// │              │ 告覧ヨトハ極  │               │               │               │
	// │              │ 整ラ得95稿は  │               │               │               │
	// │              │ かラせ江利ス  │               │               │               │
	// │              │ 宏丸霊ミ考整  │               │               │               │
	// │              │ ス静将ず業巨  │               │               │               │
	// │              │ 職ノラホ収嗅  │               │               │               │
	// │              │ ざな。        │               │               │               │
	// ╰──────────────┴───────────────┴───────────────┴───────────────┴───────────────╯
}

// Check that stylized wrapped content does not go beyond its cell.
func TestWrapPreStyledContent(t *testing.T) {
	headers := []string{"Package", "Version", "Link"}
	data := [][]string{
		{"sourcegit", "0.19", lipgloss.JoinHorizontal(lipgloss.Left, lipgloss.NewStyle().Foreground(lipgloss.Color("#31BB71")).Render("https://aur.archlinux.org/packages/sourcegit-bin"))},
		{},
		{"Welcome", "いらっしゃいませ", "مرحباً", "환영", "欢迎"},
		{"Goodbye", "さようなら", "مع السلامة", "안녕히 가세요", "再见"},
	}
	table := New().
		Headers(headers...).
		Rows(data...).
		Width(80).
		Wrap(true)
	golden.RequireEqual(t, []byte(table.String()))
}

// Check that stylized wrapped content does not go beyond its cell.
func TestWrapStyleFuncContent(t *testing.T) {
	headers := []string{"Package", "Version", "Link"}
	data := [][]string{
		{"sourcegit", "0.19", "https://aur.archlinux.org/packages/sourcegit-bin"},
		{"Welcome", "いらっしゃいませ", "مرحباً"},
		{"Goodbye", "さようなら", "مع السلامة"},
	}
	table := New().
		Headers(headers...).
		Rows(data...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == HeaderRow {
				return lipgloss.NewStyle()
			}
			if strings.Contains(data[row][col], "https://") {
				return lipgloss.NewStyle().Foreground(lipgloss.Color("#31BB71"))
			}
			return lipgloss.NewStyle()
		}).
		Width(60).
		Wrap(true)
	golden.RequireEqual(t, []byte(table.String()))
}
