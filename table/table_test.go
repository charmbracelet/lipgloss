package table

import (
	"fmt"
	"strings"
	"testing"
	"unicode"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/exp/golden"
	"github.com/muesli/termenv"
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

	expected := strings.TrimSpace(`
┌──────────┬──────────────┬───────────┐
│ LANGUAGE │    FORMAL    │ INFORMAL  │
├──────────┼──────────────┼───────────┤
│ Chinese  │ Nǐn hǎo      │ Nǐ hǎo    │
│ French   │ Bonjour      │ Salut     │
│ Japanese │ こんにちは   │ やあ      │
│ Russian  │ Zdravstvuyte │ Privet    │
│ Spanish  │ Hola         │ ¿Qué tal? │
└──────────┴──────────────┴───────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
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

	expected := strings.TrimSpace(`
┌──────────┬───────────────────────────────┬─────────────────┐
│ LANGUAGE │            FORMAL             │    INFORMAL     │
├──────────┼───────────────────────────────┼─────────────────┤
│ Chinese  │ 您好                          │ 你好            │
│ Japanese │ こんにちは                    │ やあ            │
│ Russian  │ Здравствуйте                  │ Привет          │
│ Spanish  │ Hola                          │ ¿Qué tal?       │
│ English  │ You look absolutely fabulous. │ How's it going? │
└──────────┴───────────────────────────────┴─────────────────┘
`)

	if got := ansi.Strip(table.String()); got != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, got)
	}
}

func TestTableEmpty(t *testing.T) {
	table := New().
		Border(lipgloss.NormalBorder()).
		StyleFunc(TableStyle).
		Headers("LANGUAGE", "FORMAL", "INFORMAL")

	expected := strings.TrimSpace(`
┌──────────┬────────┬──────────┐
│ LANGUAGE │ FORMAL │ INFORMAL │
├──────────┼────────┼──────────┤
└──────────┴────────┴──────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
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

	expected := strings.TrimSpace(`
┌──────────┬──────────────┬───────────┐
│ LANGUAGE │    FORMAL    │ INFORMAL  │
├──────────┼──────────────┼───────────┤
│ French   │ Bonjour      │ Salut     │
│ Japanese │ こんにちは   │ やあ      │
│ Russian  │ Zdravstvuyte │ Privet    │
│ Spanish  │ Hola         │ ¿Qué tal? │
└──────────┴──────────────┴───────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
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

	expected := strings.TrimSpace(`
╔══════════╦══════════════╦═══════════╗
║ LANGUAGE ║    FORMAL    ║ INFORMAL  ║
╠══════════╬══════════════╬═══════════╣
║ Chinese  ║ Nǐn hǎo      ║ Nǐ hǎo    ║
║ French   ║ Bonjour      ║ Salut     ║
║ Japanese ║ こんにちは   ║ やあ      ║
║ Russian  ║ Zdravstvuyte ║ Privet    ║
║ Spanish  ║ Hola         ║ ¿Qué tal? ║
╚══════════╩══════════════╩═══════════╝
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
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

	expected := strings.TrimSpace(`
┌──────────┬──────────────┬───────────┐
│ LANGUAGE │    FORMAL    │ INFORMAL  │
├──────────┼──────────────┼───────────┤
│ Chinese  │ Nǐn hǎo      │ Nǐ hǎo    │
│ French   │ Bonjour      │ Salut     │
│ Japanese │ こんにちは   │ やあ      │
│ Russian  │ Zdravstvuyte │ Privet    │
│ Spanish  │ Hola         │ ¿Qué tal? │
└──────────┴──────────────┴───────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
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

	expected := strings.TrimSpace(`
┌──────────┬──────────────┬───────────┐
│ LANGUAGE │    FORMAL    │           │
├──────────┼──────────────┼───────────┤
│ Chinese  │ Nǐn hǎo      │ Nǐ hǎo    │
│ French   │ Bonjour      │ Salut     │
│ Japanese │ こんにちは   │ やあ      │
│ Russian  │ Zdravstvuyte │ Privet    │
│ Spanish  │ Hola         │ ¿Qué tal? │
└──────────┴──────────────┴───────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
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

	expected := strings.TrimSpace(`
┌──────────┬──────────────┬───────────┬────────┬────────┐
│ LANGUAGE │    FORMAL    │           │        │        │
├──────────┼──────────────┼───────────┼────────┼────────┤
│ Chinese  │ Nǐn hǎo      │ Nǐ hǎo    │        │        │
│ French   │ Bonjour      │ Salut     │ Salut  │        │
│ Japanese │ こんにちは   │ やあ      │        │        │
│ Russian  │ Zdravstvuyte │ Privet    │ Privet │ Privet │
│ Spanish  │ Hola         │ ¿Qué tal? │        │        │
└──────────┴──────────────┴───────────┴────────┴────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
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

	expected := strings.TrimSpace(`
┌──────────┬──────────────┬───────────┐
│ Chinese  │ Nǐn hǎo      │ Nǐ hǎo    │
│ French   │ Bonjour      │ Salut     │
│ Japanese │ こんにちは   │ やあ      │
│ Russian  │ Zdravstvuyte │ Privet    │
│ Spanish  │ Hola         │ ¿Qué tal? │
└──────────┴──────────────┴───────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
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

	expected := strings.TrimSpace(`
┌───────────────────────────────────┐
│ Chinese   Nǐn hǎo       Nǐ hǎo    │
│ French    Bonjour       Salut     │
│ Japanese  こんにちは    やあ      │
│ Russian   Zdravstvuyte  Privet    │
│ Spanish   Hola          ¿Qué tal? │
└───────────────────────────────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
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

	expected := strings.TrimSpace(`
┌───────────────────────────────────┐
│ LANGUAGE     FORMAL     INFORMAL  │
├───────────────────────────────────┤
│ Chinese   Nǐn hǎo       Nǐ hǎo    │
│ French    Bonjour       Salut     │
│ Japanese  こんにちは    やあ      │
│ Russian   Zdravstvuyte  Privet    │
│ Spanish   Hola          ¿Qué tal? │
└───────────────────────────────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
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

	expected := strings.TrimSpace(`
┌───────────────────────────────────────────────────┐
│ LANGUAGE     FORMAL                               │
├───────────────────────────────────────────────────┤
│ Chinese   Nǐn hǎo       Nǐ hǎo                    │
│ French    Bonjour       Salut      Salut          │
│ Japanese  こんにちは    やあ                      │
│ Russian   Zdravstvuyte  Privet     Privet  Privet │
│ Spanish   Hola          ¿Qué tal?                 │
└───────────────────────────────────────────────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
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

	expected := strings.TrimPrefix(`
 LANGUAGE │    FORMAL    │ INFORMAL  
──────────┼──────────────┼───────────
 Chinese  │ Nǐn hǎo      │ Nǐ hǎo    
 French   │ Bonjour      │ Salut     
 Japanese │ こんにちは   │ やあ      
 Russian  │ Zdravstvuyte │ Privet    
 Spanish  │ Hola         │ ¿Qué tal? `, "\n")

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", debug(expected), debug(table.String()))
	}
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

	expected := strings.TrimPrefix(`
 LANGUAGE │    FORMAL    │ INFORMAL  
 Chinese  │ Nǐn hǎo      │ Nǐ hǎo    
 French   │ Bonjour      │ Salut     
 Japanese │ こんにちは   │ やあ      
 Russian  │ Zdravstvuyte │ Privet    
 Spanish  │ Hola         │ ¿Qué tal? `, "\n")

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", debug(expected), debug(table.String()))
	}
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

	expected := strings.TrimSpace(`
┌──────────┬──────────────┬───────────┐
│ LANGUAGE │    FORMAL    │ INFORMAL  │
│ Chinese  │ Nǐn hǎo      │ Nǐ hǎo    │
│ French   │ Bonjour      │ Salut     │
│ Japanese │ こんにちは   │ やあ      │
│ Russian  │ Zdravstvuyte │ Privet    │
│ Spanish  │ Hola         │ ¿Qué tal? │
└──────────┴──────────────┴───────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
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

	expected := strings.TrimSpace(`
┌──────────┬──────────────┬───────────┐
│ LANGUAGE │    FORMAL    │ INFORMAL  │
├──────────┼──────────────┼───────────┤
│ Chinese  │ Nǐn hǎo      │ Nǐ hǎo    │
├──────────┼──────────────┼───────────┤
│ French   │ Bonjour      │ Salut     │
├──────────┼──────────────┼───────────┤
│ Japanese │ こんにちは   │ やあ      │
├──────────┼──────────────┼───────────┤
│ Russian  │ Zdravstvuyte │ Privet    │
├──────────┼──────────────┼───────────┤
│ Spanish  │ Hola         │ ¿Qué tal? │
└──────────┴──────────────┴───────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
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

	expected := strings.TrimSpace(`
┌──────────────────┬─────────────────────────┐
│ EXPRESSION       │ MEANING                 │
├──────────────────┼─────────────────────────┤
│                  │                         │
│ Chutar o balde   │  Literally translates   │
│                  │  to "kick the bucket."  │
│                  │  It's used when         │
│                  │  someone gives up or    │
│                  │  loses patience.        │
│                  │                         │
│                  │                         │
│ Engolir sapos    │  Literally means "to    │
│                  │  swallow frogs." It's   │
│                  │  used to describe       │
│                  │  someone who has to     │
│                  │  tolerate or endure     │
│                  │  unpleasant             │
│                  │  situations.            │
│                  │                         │
│                  │                         │
│ Arroz de festa   │  Literally means        │
│                  │  "party rice." It´s     │
│                  │  used to refer to       │
│                  │  someone who shows up   │
│                  │  everywhere.            │
│                  │                         │
└──────────────────┴─────────────────────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
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

	expected := strings.TrimSpace(`
┌──────────────────┬─────────────────────────┐
│ EXPRESSION       │ MEANING                 │
├──────────────────┼─────────────────────────┤
│                  │                         │
│ Chutar o balde   │  Literally translates   │
│                  │  to "kick the bucket."  │
│                  │  It's used when         │
│                  │  someone gives up or    │
│                  │  loses patience.        │
│                  │                         │
├──────────────────┼─────────────────────────┤
│                  │                         │
│ Engolir sapos    │  Literally means "to    │
│                  │  swallow frogs." It's   │
│                  │  used to describe       │
│                  │  someone who has to     │
│                  │  tolerate or endure     │
│                  │  unpleasant             │
│                  │  situations.            │
│                  │                         │
├──────────────────┼─────────────────────────┤
│                  │                         │
│ Arroz de festa   │  Literally means        │
│                  │  "party rice." It´s     │
│                  │  used to refer to       │
│                  │  someone who shows up   │
│                  │  everywhere.            │
│                  │                         │
└──────────────────┴─────────────────────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
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

	expected := strings.TrimSpace(`
┌────────────────────────┬────────────────────────────┬────────────────────────┐
│        LANGUAGE        │           FORMAL           │        INFORMAL        │
├────────────────────────┼────────────────────────────┼────────────────────────┤
│ Chinese                │ Nǐn hǎo                    │ Nǐ hǎo                 │
│ French                 │ Bonjour                    │ Salut                  │
│ Japanese               │ こんにちは                 │ やあ                   │
│ Russian                │ Zdravstvuyte               │ Privet                 │
│ Spanish                │ Hola                       │ ¿Qué tal?              │
└────────────────────────┴────────────────────────────┴────────────────────────┘
`)

	if lipgloss.Width(table.String()) != 80 {
		t.Fatalf("expected table width to be 80, got %d", lipgloss.Width(table.String()))
	}

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
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

	expected := strings.TrimSpace(`
┌─────────┬─────────┬────────┐
│ LANGUAG │ FORMAL  │ INFORM │
├─────────┼─────────┼────────┤
│ Chinese │ Nǐn hǎo │ Nǐ hǎo │
│ French  │ Bonjour │ Salut  │
│ Japanes │ こんに  │ やあ   │
│ Russian │ Zdravst │ Privet │
│ Spanish │ Hola    │ ¿Qué   │
└─────────┴─────────┴────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
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

	expected := strings.TrimSpace(`
┌──────┬─────┬──────────┐
│ Name │ Age │ Location │
├──────┼─────┼──────────┤
│ Kini │ 40  │ New York │
│ Eli  │ 30  │ London   │
│ Iris │ 20  │ Paris    │
└──────┴─────┴──────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
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
		Rows(rows...)

	expected := strings.TrimSpace(`
┏━━━━┳━━━━━┳━━━━━┓
┃ LA ┃ FOR ┃ INF ┃
┣━━━━╋━━━━━╋━━━━━┫
┃ Ch ┃ 您  ┃ 你  ┃
┃ Ja ┃ こ  ┃ や  ┃
┃ Ar ┃ أهل ┃ أهل ┃
┃ Ru ┃ Здр ┃ При ┃
┃ Sp ┃ Hol ┃ ¿Qu ┃
┃ En ┃ You ┃ How ┃
┗━━━━┻━━━━━┻━━━━━┛
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
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

	expected := strings.TrimSpace(`
┌
│
├
│
│
│
│
│
└
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
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

	expected := strings.TrimSpace(`
──────────────────────────────
 LANGUAGE  FORMAL   INFORMAL  
──────────────────────────────
 Chinese   Nǐn hǎo  Nǐ hǎo    
 French    Bonjour  Salut     
 Japanese  こんに   やあ      
 Russian   Zdravst  Privet    
 Spanish   Hola     ¿Qué tal? 
──────────────────────────────
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
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

	expected := strings.TrimSpace(`
──────────────────────────────
 LANGUAGE  FORMAL   INFORMAL  
──────────────────────────────
 Chinese   Nǐn hǎo  Nǐ hǎo    
 French    Bonjour  Salut     
 Japanese  こんに   やあ      
 Russian   Zdravst  Privet    
 Spanish   Hola     ¿Qué tal? 
──────────────────────────────
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
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

	expected := strings.TrimSpace(`
┌──────────┬──────────────┬───────────┐
│ LANGUAGE │    FORMAL    │ INFORMAL  │
├──────────┼──────────────┼───────────┤
│ Chinese  │ Nǐn hǎo      │ Nǐ hǎo    │
│ Japanese │ こんにちは   │ やあ      │
│ Russian  │ Zdravstvuyte │ Privet    │
│ Spanish  │ Hola         │ ¿Qué tal? │
└──────────┴──────────────┴───────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
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

	expected := strings.TrimSpace(`
┌──────────┬─────────┬──────────┐
│ LANGUAGE │ FORMAL  │ INFORMAL │
├──────────┼─────────┼──────────┤
│ French   │ Bonjour │ Salut    │
└──────────┴─────────┴──────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
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

	expected := strings.TrimSpace(`
┌───────────┬────────┬──────┐
│   Fruit   │ Color  │ Code │
├───────────┼────────┼──────┤
│ Apple     │ Red    │ 31   │
│ Lime      │ Green  │ 32   │
│ Banana    │ Yellow │ 33   │
│ Blueberry │ Blue   │ 34   │
└───────────┴────────┴──────┘
`)

	if stripString(table.String()) != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, stripString(table.String()))
	}
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

	expected := strings.TrimSpace(`
┌──────────┬──────────────┬───────────┐
│ LANGUAGE │    FORMAL    │ INFORMAL  │
├──────────┼──────────────┼───────────┤
│ Chinese  │ Nǐn hǎo      │ Nǐ hǎo    │
│ French   │ Bonjour      │ Salut     │
│ Japanese │ こんにちは   │ やあ      │
│ Russian  │ Zdravstvuyte │ Privet    │
│ Spanish  │ Hola         │ ¿Qué tal? │
└──────────┴──────────────┴───────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
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

	expected := strings.TrimSpace(`
┌──────────┬──────────────┬───────────┐
│ LANGUAGE │    FORMAL    │ INFORMAL  │
├──────────┼──────────────┼───────────┤
│ Chinese  │ Nǐn hǎo      │ Nǐ hǎo    │
│ French   │ Bonjour      │ Salut     │
│ Japanese │ こんにちは   │ やあ      │
│ Russian  │ Zdravstvuyte │ Privet    │
│ Spanish  │ Hola         │ ¿Qué tal? │
└──────────┴──────────────┴───────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
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

	expected := strings.TrimSpace(`
┌──────────┬──────────────┬───────────┐
│ LANGUAGE │    FORMAL    │ INFORMAL  │
├──────────┼──────────────┼───────────┤
│ Chinese  │ Nǐn hǎo      │ Nǐ hǎo    │
│ French   │ Bonjour      │ Salut     │
│ Japanese │ こんにちは   │ やあ      │
│ …        │ …            │ …         │
└──────────┴──────────────┴───────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
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

	expected := strings.TrimSpace(`
┌────┬──────────┬──────────────┬───────────┐
│ ID │ LANGUAGE │    FORMAL    │ INFORMAL  │
├────┼──────────┼──────────────┼───────────┤
│ …  │ …        │ …            │ …         │
└────┴──────────┴──────────────┴───────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestTableHeightMinimumShowData(t *testing.T) {
	table := New().
		Height(0).
		Border(lipgloss.NormalBorder()).
		StyleFunc(TableStyle).
		Headers("LANGUAGE", "FORMAL", "INFORMAL").
		Row("Chinese", "Nǐn hǎo", "Nǐ hǎo")

	expected := strings.TrimSpace(`
┌──────────┬─────────┬──────────┐
│ LANGUAGE │ FORMAL  │ INFORMAL │
├──────────┼─────────┼──────────┤
│ Chinese  │ Nǐn hǎo │ Nǐ hǎo   │
└──────────┴─────────┴──────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
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

	expected := strings.TrimSpace(`
┌──────────┬──────────────┬───────────┐
│ LANGUAGE │    FORMAL    │ INFORMAL  │
├──────────┼──────────────┼───────────┤
│ French   │ Bonjour      │ Salut     │
│ Japanese │ こんにちは   │ やあ      │
│ Russian  │ Zdravstvuyte │ Privet    │
│ Spanish  │ Hola         │ ¿Qué tal? │
└──────────┴──────────────┴───────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
	}
}

func TestStyleFunc(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	tests := []struct {
		name  string
		style StyleFunc
	}{
		{
			"right-aligned text with margins",
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
			"margin and padding set",
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

// TODO add tests to check calculations for
// - long text, small header
// - long header, small text
// - padding and wrapping

func TestContentWrapping(t *testing.T) {
	tests := []struct {
		name    string
		headers []string
		data    [][]string
	}{
		//		{
		//			"long row content",
		//			[]string{"Name", "Description", "Type", "Required", "Default"},
		//			[][]string{{"command", "A command to be executed inside the container to assess its health. Each space delimited token of the command is a separate array element. Commands exiting 0 are considered to be successful probes, whilst all other exit codes are considered failures.", "yes", "hello", "yep"}},
		//		},
		//		{
		//			"missing row content",
		//			[]string{"Name", "Description", "Type", "Required", "Default"},
		//			[][]string{{"command", "A command to be executed inside the container to assess its health. Each space delimited token of the command is a separate array element. Commands exiting 0 are considered to be successful probes, whilst all other exit codes are considered failures.", "yes", "", ""}},
		//		},
		{
			"long header content, long and short rows",
			[]string{"Destination", "Why are you going on this trip? Is it a hot or cold climate?", "Affordability"},
			[][]string{
				{"Mexico", "I want to go somewhere hot, dry, and affordable. Mexico has really good food, just don't drink tap water!", "$"},
				{"New York", "I'm thinking about going during the Christmas season to check out Rockefeller center. Might be cold though...", "$$$"},
				{"California", "", "$$$"},
			},
		},
		{
			"Long text, different languages",
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
		table := New().
			Headers(tc.headers...).
			Rows(tc.data...)
			//		StyleFunc(func(_, col int) lipgloss.Style {
			//			if len(data[col]) > 25 {
			//				return lipgloss.NewStyle().Width(30)
			//			}
			//			return lipgloss.NewStyle()
			//		})
		//		table.Row("command", lipgloss.NewStyle().Width(60).Render("A command to be executed inside the container to assess its health. Each space delimited token of the command is a separate array element. Commands exiting 0 are considered to be successful probes, whilst all other exit codes are considered failures."), "yes")
		table.Width(80)

		t.Log(lipgloss.Width(table.String()))
		t.Log("\n" + table.String() + "\n")
	}
}

func TestContentWrapping_WithPadding(t *testing.T) {
	tests := []struct {
		name    string
		headers []string
		data    [][]string
	}{
		{
			"long row content",
			[]string{"Name", "Description", "Type", "Required", "Default"},
			[][]string{{"command", "A command to be executed inside the container to assess its health. Each space delimited token of the command is a separate array element. Commands exiting 0 are considered to be successful probes, whilst all other exit codes are considered failures.", "yes", "hello", "yep"}},
		},
		{
			"missing row content",
			[]string{"Name", "Description", "Type", "Required", "Default"},
			[][]string{{"command", "A command to be executed inside the container to assess its health. Each space delimited token of the command is a separate array element. Commands exiting 0 are considered to be successful probes, whilst all other exit codes are considered failures.", "yes", "", ""}},
		},
		{
			"long header content, long and short rows",
			[]string{"Destination", "Why are you going on this trip? Is it a hot or cold climate?", "Affordability"},
			[][]string{
				{"Mexico", "I want to go somewhere hot, dry, and affordable. Mexico has really good food, just don't drink tap water!", "$"},
				{"New York", "I'm thinking about going during the Christmas season to check out Rockefeller center. Might be cold though...", "$$$"},
				{"California", "", "$$$"},
			},
		},
		{
			"Long text, different languages",
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

	for _, tc := range tests {
		table := New().
			Headers(tc.headers...).
			Rows(tc.data...).
			StyleFunc(func(_, col int) lipgloss.Style {
				return lipgloss.NewStyle().Padding(0, 1)
			})
		table.Width(80)

		t.Log(lipgloss.Width(table.String()))
		t.Log("\n" + table.String() + "\n")
	}
}

func TestContentWrapping_WithMargins(t *testing.T) {
	tests := []struct {
		name    string
		headers []string
		data    [][]string
	}{
		{
			"long row content",
			[]string{"Name", "Description", "Type", "Required", "Default"},
			[][]string{{"command", "A command to be executed inside the container to assess its health. Each space delimited token of the command is a separate array element. Commands exiting 0 are considered to be successful probes, whilst all other exit codes are considered failures.", "yes", "hello", "yep"}},
		},
		{
			"missing row content",
			[]string{"Name", "Description", "Type", "Required", "Default"},
			[][]string{{"command", "A command to be executed inside the container to assess its health. Each space delimited token of the command is a separate array element. Commands exiting 0 are considered to be successful probes, whilst all other exit codes are considered failures.", "yes", "", ""}},
		},
		{
			"long header content, long and short rows",
			[]string{"Destination", "Why are you going on this trip? Is it a hot or cold climate?", "Affordability"},
			[][]string{
				{"Mexico", "I want to go somewhere hot, dry, and affordable. Mexico has really good food, just don't drink tap water!", "$"},
				{"New York", "I'm thinking about going during the Christmas season to check out Rockefeller center. Might be cold though...", "$$$"},
				{"California", "", "$$$"},
			},
		},
		{
			"Long text, different languages",
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
		table := New().
			Headers(tc.headers...).
			Rows(tc.data...).
			StyleFunc(func(row, col int) lipgloss.Style {
				return lipgloss.NewStyle().Margin(0, 4)
			})
		table.Width(80)

		t.Log(lipgloss.Width(table.String()))
		t.Log("\n" + table.String() + "\n")
	}
}

func TestContentWrapping_ColumnWidth(t *testing.T) {
	tests := []struct {
		name    string
		headers []string
		data    [][]string
	}{
		{
			"long row content",
			[]string{"Name", "Description", "Type", "Required", "Default"},
			[][]string{{"command", "A command to be executed inside the container to assess its health. Each space delimited token of the command is a separate array element. Commands exiting 0 are considered to be successful probes, whilst all other exit codes are considered failures.", "yes", "hello", "yep"}},
		},
		{
			"missing row content",
			[]string{"Name", "Description", "Type", "Required", "Default"},
			[][]string{{"command", "A command to be executed inside the container to assess its health. Each space delimited token of the command is a separate array element. Commands exiting 0 are considered to be successful probes, whilst all other exit codes are considered failures.", "yes", "", ""}},
		},
		{
			"long header content, long and short rows",
			[]string{"Destination", "Why are you going on this trip? Is it a hot or cold climate?", "Affordability"},
			[][]string{
				{"Mexico", "I want to go somewhere hot, dry, and affordable. Mexico has really good food, just don't drink tap water!", "$"},
				{"New York", "I'm thinking about going during the Christmas season to check out Rockefeller center. Might be cold though...", "$$$"},
				{"California", "", "$$$"},
			},
		},
		{
			"Long text, different languages",
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
		table := New().
			Headers(tc.headers...).
			Rows(tc.data...).
			StyleFunc(func(row, col int) lipgloss.Style {
				if row == 0 && col == 1 {
					return lipgloss.NewStyle().Width(15)
				}
				return lipgloss.NewStyle()
			})
		table.Width(80)
		t.Log(lipgloss.Width(table.String()))
		t.Log("\n" + table.String() + "\n")
	}
}

func TestCarriageReturn(t *testing.T) {
	data := [][]string{
		{"a0", "b0", "c0", "d0"},
		{"a1", "b1.0\r\nb1.1\r\nb1.2\r\nb1.3\r\nb1.4\r\nb1.5\r\nb1.6", "c1", "d1"},
		{"a2", "b2", "c2", "d2"},
		{"a3", "b3", "c3", "d3"},
	}
	table := New().Rows(data...).Border(lipgloss.NormalBorder())
	got := table.String()
	want := `┌──┬────┬──┬──┐
│a0│b0  │c0│d0│
│a1│b1.0│c1│d1│
│  │b1.1│  │  │
│  │b1.2│  │  │
│  │b1.3│  │  │
│  │b1.4│  │  │
│  │b1.5│  │  │
│  │b1.6│  │  │
│a2│b2  │c2│d2│
│a3│b3  │c3│d3│
└──┴────┴──┴──┘`

	if got != want {
		t.Logf("detailed view...\ngot:\n%q\nwant:\n%q", got, want)
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
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

func debug(s string) string {
	return strings.ReplaceAll(s, " ", ".")
}

func stripString(str string) string {
	s := ansi.Strip(str)
	ss := strings.Split(s, "\n")

	var lines []string
	for _, l := range ss {
		trim := strings.TrimRightFunc(l, unicode.IsSpace)
		lines = append(lines, trim)
	}

	return strings.Join(lines, "\n")
}

// Examples

func ExampleWrap() {
	tests := []struct {
		name    string
		headers []string
		data    [][]string
	}{
		{
			"long row content",
			[]string{"Name", "Description", "Type", "Required", "Default"},
			[][]string{{"command", "A command to be executed inside the container to assess its health. Each space delimited token of the command is a separate array element. Commands exiting 0 are considered to be successful probes, whilst all other exit codes are considered failures.", "yes", "hello", "yep"}},
		},
		{
			"missing row content",
			[]string{"Name", "Description", "Type", "Required", "Default"},
			[][]string{{"command", "A command to be executed inside the container to assess its health. Each space delimited token of the command is a separate array element. Commands exiting 0 are considered to be successful probes, whilst all other exit codes are considered failures.", "yes", "", ""}},
		},
		{
			"long header content, long and short rows",
			[]string{"Destination", "Why are you going on this trip? Is it a hot or cold climate?", "Affordability"},
			[][]string{
				{"Mexico", "I want to go somewhere hot, dry, and affordable. Mexico has really good food, just don't drink tap water!", "$"},
				{"New York", "I'm thinking about going during the Christmas season to check out Rockefeller center. Might be cold though...", "$$$"},
				{"California", "", "$$$"},
			},
		},
		{
			"Long text, different languages",
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
		table := New().
			Headers(tc.headers...).
			Rows(tc.data...).
			StyleFunc(func(row, col int) lipgloss.Style {
				//				if col == 1 {
				//					return lipgloss.NewStyle().MaxWidth(10)
				//				}
				return DefaultStyles(row, col)
			}).Wrap(false)
		table.Width(80)

		fmt.Println(table.String())

		// Output:
		//
	}
}

// . TODO test custom column widths defined by StyleFunc.
func ExampleStyleFuncCustomColumns() {}
