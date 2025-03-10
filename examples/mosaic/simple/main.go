package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/v2/mosaic"
)

func main() {
	dogImg, err := loadImage("./pekinas.jpg")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	m := mosaic.New().Width(80).Scale(1)

	fmt.Println(lipgloss.JoinVertical(lipgloss.Right, lipgloss.JoinHorizontal(lipgloss.Center, m.Render(dogImg))))
}

func loadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()
	if err != nil {
		return nil, err
	}
	return jpeg.Decode(f)
}
