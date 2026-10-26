package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	WallDirPath := "/home/umka/wallpapers"
	
	// заходим в дерикторию где хранятся обои
	err := os.Chdir(WallDirPath)
	if err != nil {
		panic(err)
	}

	// список файлов в дериктории
	files, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	// выводим поштучно
	for _, file := range files {
		fmt.Println(file.Name())
	}

	// ввод юзером названия конкретного изображения 
	var PictureName string
	fmt.Println("enter name of picture")
	fmt.Scan(&PictureName)

	picturePath := filepath.Join(WallDirPath, PictureName) // склеивание ~/wallpapers + wallpaper.png
	cmd := exec.Command("timg", "-g50x50", picturePath)

	// КРИТИЧЕСКИ ВАЖНО
	// связываем потоки дочернего процесса с терминалом 
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run();
	if err != nil {
		log.Fatal(err)
	}
}