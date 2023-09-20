package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func main() {
	re := lipgloss.NewRenderer(os.Stdout)
	baseStyle := re.NewStyle().Padding(0, 1)
	headerStyle := baseStyle.Copy().Foreground(lipgloss.Color("#F8D030"))
	selectedStyle := baseStyle.Copy().Foreground(lipgloss.Color("#F8D030")).Background(lipgloss.Color("#282828"))
	typeColors := map[string]lipgloss.AdaptiveColor{
		"Grass":    {Light: "#78C850", Dark: "#78C850"},
		"Bug":      {Light: "#A8B820", Dark: "#A8B820"},
		"Fire":     {Light: "#F08030", Dark: "#F08030"},
		"Water":    {Light: "#6890F0", Dark: "#6890F0"},
		"Normal":   {Light: "#A8A878", Dark: "#A8A878"},
		"Poison":   {Light: "#A040A0", Dark: "#A040A0"},
		"Flying":   {Light: "#A890F0", Dark: "#A890F0"},
		"Electric": {Light: "#F8D030", Dark: "#F8D030"},
		"Ground":   {Light: "#E0C068", Dark: "#E0C068"},
	}

	headers := []any{"#", "Name", "Type 1", "Type 2", "Japanese", "Official Rom."}
	data := [][]any{
		{1, "Bulbasaur", "Grass", "Poison", "フシギダネ", "Bulbasaur"},
		{2, "Ivysaur", "Grass", "Poison", "フシギソウ", "Ivysaur"},
		{3, "Venusaur", "Grass", "Poison", "フシギバナ", "Venusaur"},
		{4, "Charmander", "Fire", "", "ヒトカゲ", "Hitokage"},
		{5, "Charmeleon", "Fire", "", "リザード", "Lizardo"},
		{6, "Charizard", "Fire", "Flying", "リザードン", "Lizardon"},
		{7, "Squirtle", "Water", "", "ゼニガメ", "Zenigame"},
		{8, "Wartortle", "Water", "", "カメール", "Kameil"},
		{9, "Blastoise", "Water", "", "カメックス", "Kamex"},
		{10, "Caterpie", "Bug", "", "キャタピー", "Caterpie"},
		{11, "Metapod", "Bug", "", "トランセル", "Trancell"},
		{12, "Butterfree", "Bug", "Flying", "バタフリー", "Butterfree"},
		{13, "Weedle", "Bug", "Poison", "ビードル", "Beedle"},
		{14, "Kakuna", "Bug", "Poison", "コクーン", "Cocoon"},
		{15, "Beedrill", "Bug", "Poison", "スピアー", "Spear"},
		{16, "Pidgey", "Normal", "Flying", "ポッポ", "Poppo"},
		{17, "Pidgeotto", "Normal", "Flying", "ピジョン", "Pigeon"},
		{18, "Pidgeot", "Normal", "Flying", "ピジョット", "Pigeot"},
		{19, "Rattata", "Normal", "", "コラッタ", "Koratta"},
		{20, "Raticate", "Normal", "", "ラッタ", "Ratta"},
		{21, "Spearow", "Normal", "Flying", "オニスズメ", "Onisuzume"},
		{22, "Fearow", "Normal", "Flying", "オニドリル", "Onidrill"},
		{23, "Ekans", "Poison", "", "アーボ", "Arbo"},
		{24, "Arbok", "Poison", "", "アーボック", "Arbok"},
		{25, "Pikachu", "Electric", "", "ピカチュウ", "Pikachu"},
		{26, "Raichu", "Electric", "", "ライチュウ", "Raichu"},
		{27, "Sandshrew", "Ground", "", "サンド", "Sand"},
		{28, "Sandslash", "Ground", "", "サンドパン", "Sandpan"},
	}

	CapitalizeHeaders := func(data []any) []any {
		for i := range data {
			data[i] = strings.ToUpper(data[i].(string))
		}
		return data
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(re.NewStyle().Foreground(lipgloss.Color("238"))).
		Headers(CapitalizeHeaders(headers)...).
		Rows(data...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return headerStyle
			}

			if data[row-1][1] == "Pikachu" {
				return selectedStyle
			}

			switch col {
			case 2, 3: // Type 1 + 2
				color := typeColors[fmt.Sprint(data[row-1][col])]
				return baseStyle.Copy().Foreground(color)
			}

			switch {
			case row%2 == 0:
				return baseStyle.Copy().Foreground(lipgloss.Color("245"))
			case row%2 == 1:
				return baseStyle.Copy().Foreground(lipgloss.Color("247"))
			}

			return baseStyle
		})
	fmt.Println(t)
}
