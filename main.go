package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// 1. Модель (состояние нашего приложения)
type model struct {
	files        []string // Список путей к картинкам
	cursor       int      // На какой картинке сейчас курсор
	preview      string   // ASCII-арт текущей картинки
	loading      bool     // Идет ли сейчас рендеринг картинки
	width, height int     // Размеры терминала
}

// Сообщение, которое приходит, когда chafa закончил работу
type previewMsg string

// 2. Инициализация (что делать при старте)
func initialModel(wallpaperDir string) model {
	// Читаем папку с обоями
	var files []string
	entries, err := os.ReadDir(wallpaperDir)
	if err != nil {
		fmt.Printf("Ошибка чтения папки: %v\n", err)
		os.Exit(1)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			// Берем только картинки
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp" {
				files = append(files, filepath.Join(wallpaperDir, entry.Name()))
			}
		}
	}

	return model{
		files:   files,
		cursor:  0,
		loading: true,
	}
}

// Функция, которая асинхронно вызывает chafa
func renderPreview(path string) tea.Cmd {
	return func() tea.Msg {
		// Вызываем chafa. --size 40x20 ограничивает размер картинки
		cmd := exec.Command("chafa", "--format", "symbols", "--size", "40x20", path)
		out, err := cmd.Output()
		if err != nil {
			return previewMsg("Ошибка рендеринга: " + err.Error())
		}
		return previewMsg(string(out))
	}
}

// 3. Функция Update (реакция на нажатия и сообщения)
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// При старте сразу рендерим первую картинку
		if len(m.files) > 0 {
			return m, renderPreview(m.files[m.cursor])
		}

	case previewMsg:
		// Пришел результат от chafa
		m.preview = string(msg)
		m.loading = false
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				m.loading = true
				return m, renderPreview(m.files[m.cursor])
			}
		case "down", "j":
			if m.cursor < len(m.files)-1 {
				m.cursor++
				m.loading = true
				return m, renderPreview(m.files[m.cursor])
			}
		case "enter":
			// Тут будет логика смены обоев
		}
	}
	return m, nil
}

// 4. Функция View (что и как рисовать на экране)
func (m model) View() string {
	var b strings.Builder

	// Заголовок
	b.WriteString("=== wltui ===\n\n")

	// Левая колонка (список)
	var list strings.Builder
	for i, file := range m.files {
		cursor := "  "
		if m.cursor == i {
			cursor = "> "
		}
		fileName := filepath.Base(file)
		if m.cursor == i {
			list.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(cursor + fileName + "\n"))
		} else {
			list.WriteString(cursor + fileName + "\n")
		}
	}

	// Правая колонка (превью)
	var previewStr string
	if m.loading {
		previewStr = "Загрузка предпросмотра..."
	} else {
		previewStr = m.preview
	}

	// Склеиваем левую и правую часть
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, list.String(), "   |   ", previewStr))
	b.WriteString("\n\nНажмите q для выхода, стрелки для навигации")

	return b.String()
}

// Точка входа
func main() {
	// Укажи свою папку с обоями
	wallpaperDir := "/home/umka/wallpapers" 
	
	p := tea.NewProgram(initialModel(wallpaperDir), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Ошибка запуска: %v\n", err)
		os.Exit(1)
	}
}