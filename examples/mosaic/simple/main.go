package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/mosaic"
)

func main() {
	dogImg, err := loadImage("./pekinas.jpg")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	m := mosaic.New().Scale(mosaic.Center)

	fmt.Println(lipgloss.JoinVertical(lipgloss.Right, lipgloss.JoinHorizontal(lipgloss.Center, m.Render(dogImg))))
}

func loadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return jpeg.Decode(f)
}
