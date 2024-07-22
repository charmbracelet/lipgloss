package table

import (
	"strings"
	"testing"
	"unicode"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

var TableStyle = func(row, col int) lipgloss.Style {
	switch {
	case row == 0:
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
			case row == 0:
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

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
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
		if row == 0 {
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
		if row == 0 {
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

func TestTableGlow(t *testing.T) {
	table := New().
		Border(lipgloss.NormalBorder()).
		BorderTop(true).
		BorderRight(true).
		BorderBottom(true).
		BorderLeft(true).
		Width(80).
		StyleFunc(func(row, col int) lipgloss.Style {
			st := lipgloss.NewStyle()
			st = st.Margin(0, 1)
			if row == 0 {
				st = st.Bold(true)
			}

			switch col {
			case 0, 1, 2:
				st = st.Align(lipgloss.Right).MarginLeft(0)
			case 3:
				st = st.Align(lipgloss.Left).MarginRight(0)
			default:
				st = st.Align(lipgloss.Center)
			}

			return st
		}).
		Headers("Dec", "Hex", "Oct", "Char", "Description").
		Row("001", "0x00", "0o000", "NUL", "null").
		Row("002", "0x01", "0o001", "SOH", "startofheading").
		Row("003", "0x02", "0o002", "STX", "startoftext").
		Row("004", "0x03", "0o003", "ETX", "endoftext").
		Row("005", "0x04", "0o004", "EOT", "endoftransmission").
		Row("006", "0x05", "0o005", "ENQ", "enquiry")

	expected := strings.TrimSpace(`
┌───────────┬────────────┬─────────────┬────────────┬──────────────────────────┐
│       Dec │        Hex │         Oct │ Char       │        Description       │
├───────────┼────────────┼─────────────┼────────────┼──────────────────────────┤
│       001 │       0x00 │       0o000 │ NUL        │            null          │
│       002 │       0x01 │       0o001 │ SOH        │       startofheading     │
│       003 │       0x02 │       0o002 │ STX        │        startoftext       │
│       004 │       0x03 │       0o003 │ ETX        │         endoftext        │
│       005 │       0x04 │       0o004 │ EOT        │     endoftransmission    │
│       006 │       0x05 │       0o005 │ ENQ        │          enquiry         │
└───────────┴────────────┴─────────────┴────────────┴──────────────────────────┘
`)

	if table.String() != expected {
		t.Fatalf("expected:\n\n%s\n\ngot:\n\n%s", expected, table.String())
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
