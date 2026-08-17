package exporter

import (
	"fmt"
	"time"

	"github.com/PashakArt/go_lift_backend/services/workout-service/internal/domain"
	"github.com/xuri/excelize/v2"
)

type WorkoutReportSet struct {
	SetNumber    int
	ExerciseName string
	Weight       *float64
	Reps         *int32
	DurationSec  *int32
	DistanceM    *int32
}

type WorkoutReportSession struct {
	StartedAt   time.Time
	EndedAt     *time.Time
	SessionType domain.SessionType
	Exercises   []WorkoutReportSet
}

type ExcelExporter struct{}

func NewExcelExporter() *ExcelExporter {
	return &ExcelExporter{}
}

func (e *ExcelExporter) GenerateWorkoutReport(sessions []WorkoutReportSession) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Тренировки"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)
	_ = f.DeleteSheet("Sheet1")

	sessionHeaderStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "#1F4E78", Size: 11},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#D9E1F2"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
	})

	tableHeaderStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "#FFFFFF", Size: 10},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#1F4E78"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	centerStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	currentRow := 1

	for _, session := range sessions {
		durationStr := "в процессе"
		if session.EndedAt != nil {
			dur := session.EndedAt.Sub(session.StartedAt)
			hours := int(dur.Hours())
			mins := int(dur.Minutes()) % 60
			if hours > 0 {
				durationStr = fmt.Sprintf("%dч %dмин", hours, mins)
			} else {
				durationStr = fmt.Sprintf("%dмин", mins)
			}
		}

		sessionTitle := fmt.Sprintf("🟢 Тренировка: %s · %s (Длительность: %s)",
			humanizeSessionType(session.SessionType),
			session.StartedAt.Format("02.01.2006 в 15:04"),
			durationStr,
		)

		// Объединяем ячейки от A до F (всего 6 колонок)
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", currentRow), sessionTitle)
		f.MergeCell(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("F%d", currentRow))
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("F%d", currentRow), sessionHeaderStyle)
		f.SetRowHeight(sheetName, currentRow, 26)
		currentRow++

		headers := []string{"№ Подхода", "Название упражнения", "Кол-во (повторы)", "Вес (кг)", "Время (сек)", "Дистанция (м)"}
		for colIdx, text := range headers {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, currentRow)
			f.SetCellValue(sheetName, cell, text)
		}
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("F%d", currentRow), tableHeaderStyle)
		f.SetRowHeight(sheetName, currentRow, 20)
		currentRow++

		for _, set := range session.Exercises {
			// A: Номер подхода
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", currentRow), set.SetNumber)
			f.SetCellStyle(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("A%d", currentRow), centerStyle)

			// B: Название упражнения
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", currentRow), set.ExerciseName)

			// C: Повторы
			if set.Reps != nil {
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", currentRow), *set.Reps)
			} else {
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", currentRow), "—")
			}
			f.SetCellStyle(sheetName, fmt.Sprintf("C%d", currentRow), fmt.Sprintf("C%d", currentRow), centerStyle)

			// D: Вес
			if set.Weight != nil {
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", currentRow), *set.Weight)
			} else {
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", currentRow), "—")
			}
			f.SetCellStyle(sheetName, fmt.Sprintf("D%d", currentRow), fmt.Sprintf("D%d", currentRow), centerStyle)

			// E: Время
			if set.DurationSec != nil {
				f.SetCellValue(sheetName, fmt.Sprintf("E%d", currentRow), *set.DurationSec)
			} else {
				f.SetCellValue(sheetName, fmt.Sprintf("E%d", currentRow), "—")
			}
			f.SetCellStyle(sheetName, fmt.Sprintf("E%d", currentRow), fmt.Sprintf("E%d", currentRow), centerStyle)

			// F: Дистанция
			if set.DistanceM != nil {
				f.SetCellValue(sheetName, fmt.Sprintf("F%d", currentRow), *set.DistanceM)
			} else {
				f.SetCellValue(sheetName, fmt.Sprintf("F%d", currentRow), "—")
			}
			f.SetCellStyle(sheetName, fmt.Sprintf("F%d", currentRow), fmt.Sprintf("F%d", currentRow), centerStyle)

			currentRow++
		}

		currentRow++
	}

	// Настройка ширины 6 колонок (A-F)
	cols := []string{"A", "B", "C", "D", "E", "F"}
	defaultWidths := []float64{12, 32, 16, 12, 14, 15}
	for i, col := range cols {
		_ = f.SetColWidth(sheetName, col, col, defaultWidths[i])
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func humanizeSessionType(st domain.SessionType) string {
	switch string(st) {
	case "classic":
		return "Классическая"
	case "calisthenics":
		return "Калистеника"
	case "cardio":
		return "Кардио"
	default:
		if string(st) != "" {
			return string(st)
		}
		return "Силовая"
	}
}
