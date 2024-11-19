package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss/v2"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path"
)

const failMessage = "usage: sixels <jpg or png file>"

func main() {
	// Pull a jpg or png from OS args
	if len(os.Args) != 2 {
		log.Fatalln(failMessage)
		return
	}

	fileInfo, err := os.Stat(os.Args[1])
	if err != nil {
		log.Fatalln(failMessage)
		return
	}

	ext := path.Ext(fileInfo.Name())
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		log.Fatalln(failMessage)
		return
	}

	// Even on latest Windows Terminal preview, MINGW requires this
	lipgloss.EnableLegacyWindowsANSI(os.Stdout)
	lipgloss.EnableLegacyWindowsANSI(os.Stdin)

	bgColor, _ := lipgloss.BackgroundColor(os.Stdin, os.Stdout)

	imgFile, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_ = imgFile.Close()
	}()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Fatalln(err)
	}

	sixelImage := lipgloss.Sixel(img)
	fmt.Println(lipgloss.NewStyle().Background(bgColor).RenderSixelImage(sixelImage))
}
