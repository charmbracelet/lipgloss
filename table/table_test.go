package table

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
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

func TestTableBorder(t *testing.T) {
	table := New().
		Border(lipgloss.DoubleBorder()).
		StyleFunc(TableStyle).
		Headers("LANGUAGE", "FORMAL", "INFORMAL").
		Row("Chinese", "Nǐn hǎo", "Nǐ hǎo").
		Row("French", "Bonjour", "Salut").
		Row("Japanese", "こんにちは", "やあ").
		Row("Russian", "Zdravstvuyte", "Privet").
		Row("Spanish", "Hola", "¿Qué tal?")

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
	rows := [][]any{
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
	rows := [][]any{
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
	rows := [][]any{
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
	rows := [][]any{
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
	table := New().
		Border(lipgloss.NormalBorder()).
		StyleFunc(TableStyle).
		Headers("LANGUAGE", "FORMAL", "INFORMAL").
		Row("Chinese", "Nǐn hǎo", "Nǐ hǎo").
		Row("French", "Bonjour", "Salut").
		Row("Japanese", "こんにちは", "やあ").
		Row("Russian", "Zdravstvuyte", "Privet").
		Row("Spanish", "Hola", "¿Qué tal?").
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
	table := New().
		Border(lipgloss.NormalBorder()).
		StyleFunc(TableStyle).
		Headers("LANGUAGE", "FORMAL", "INFORMAL").
		Row("Chinese", "Nǐn hǎo", "Nǐ hǎo").
		Row("French", "Bonjour", "Salut").
		Row("Japanese", "こんにちは", "やあ").
		Row("Russian", "Zdravstvuyte", "Privet").
		Row("Spanish", "Hola", "¿Qué tal?").
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
	table := New().
		Border(lipgloss.NormalBorder()).
		StyleFunc(TableStyle).
		Headers("LANGUAGE", "FORMAL", "INFORMAL").
		Row("Chinese", "Nǐn hǎo", "Nǐ hǎo").
		Row("French", "Bonjour", "Salut").
		Row("Japanese", "こんにちは", "やあ").
		Row("Russian", "Zdravstvuyte", "Privet").
		Row("Spanish", "Hola", "¿Qué tal?").
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
	table := New().
		Border(lipgloss.NormalBorder()).
		StyleFunc(TableStyle).
		BorderRow(true).
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
			return lipgloss.NewStyle().Bold(true).Padding(0, 1)
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
			return lipgloss.NewStyle().Bold(true).Padding(0, 1)
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
	rows := [][]any{
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
	rows := [][]any{
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
	rows := [][]any{
		{"Kini", 40, "New York"},
		{"Eli", 30, "London"},
		{"Iris", 20, "Paris"},
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
	rows := [][]any{
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
	rows := [][]any{
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
	rows := [][]any{
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
	rows := [][]any{
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

func debug(s string) string {
	return strings.ReplaceAll(s, " ", ".")
}
