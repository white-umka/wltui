package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

// создаем структура model
type model struct {
	choices  []string // список элементов для выбора
	cursor   int      // индекс нынешнего элемента
	selected string   // то что было при нажатии enter
}

// инициализируем model
func initialModel(wallpaperNames []string) model {
	return model{
		choices: wallpaperNames,
		cursor:  0, // стартуем с 1 элемента
	}
}

// Обязательный технический метод для Bubble Tea (в нашем случае пустой)
func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			m.selected = m.choices[m.cursor]
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	// Собираем интерфейс в одну большую строку
	s := "Выбери картинку стрелочками (Enter для подтверждения):\n\n"

	for i, choice := range m.choices {
		// Если индекс элемента совпадает с курсором — рисуем указатель
		cursorStr := " "
		if m.cursor == i {
			cursorStr = "->" // Наш курсор
		}

		// Форматируем строку: "[>] название_файла"
		s += fmt.Sprintf("%s %s\n", cursorStr, choice)
	}

	s += "\nНажми Q или Ctrl+C для выхода.\n"
	return s
}

func timg(wallDirPath string, pictureName string) {
	picturePath := filepath.Join(wallDirPath, pictureName) // склеивание ~/wallpapers + wallpaper.png

	cmd := exec.Command("timg", "-g50x50", picturePath)

	// связываем стандартные потоки дочернего процесса с нашим терминалом (ВАЖНО)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	wallDirPath := "/home/umka/wallpapers"

	// заходим в дерикторию где хранятся обои
	err := os.Chdir(wallDirPath)
	if err != nil {
		panic(err)
	}

	files, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	wallpaperNames := []string{}

	// добавляем имя каждого файла в слайс
	for _, file := range files {
		wallpaperNames = append(wallpaperNames, file.Name())
	}

	p := tea.NewProgram(initialModel(wallpaperNames))

	finalModel, err := p.Run()
	if err != nil {
		log.Fatalf("Ошибка запуска TUI: %v", err)
	}

	// Приводим интерфейс к нашей структуре model, чтобы достать выбранное значение
	if m, ok := finalModel.(model); ok && m.selected != "" {
		fmt.Printf("\nТы выбрал файл: %s\n", m.selected)
		timg(wallDirPath, m.selected)
	}
}
