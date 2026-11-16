package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"tas-app/internal/application/usecase"
)

var (
	errExit   = errors.New("exit")
	errCancel = errors.New("cancel")
)

const sectionWidth = 52

type StudentAdder interface {
	Execute(usecase.AddStudentCommand) error
}

type CriteriaAdder interface {
	Execute(usecase.AddCriteriaCommand) error
}

type StudentEvaluator interface {
	Execute(usecase.StudentEvaluationCommand) error
}

type StudentsExporter interface {
	ExportStudents(path string) error
}

type EvaluationRemover interface {
	Execute(usecase.RemoveEvaluationCommand) error
}

type Dependencies struct {
	AddCriteria      CriteriaAdder
	AddStudent       StudentAdder
	Evaluate         StudentEvaluator
	RemoveEvaluation EvaluationRemover
	GetReport        usecase.ReportReader
	ExportStudents   StudentsExporter
}

type CLI struct {
	deps Dependencies

	reader    *bufio.Scanner
	writer    io.Writer
	errWriter io.Writer
}

func (c *CLI) Run() {
	if err := c.mainMenu(); err != nil && !errors.Is(err, errExit) {
		c.printError(err)
	}
}

func (c *CLI) mainMenu() error {
	for {
		c.clearScreen()
		c.printMainMenu()

		choice, err := c.askInt("Выберите пункт: ")
		if err != nil {
			return err
		}

		switch choice {
		case 1:
			if err := c.studentsMenu(); err != nil {
				return err
			}
		case 2:
			if err := c.criteriaMenu(); err != nil {
				return err
			}
		case 3:
			if err := c.runAction(c.evaluateScreen); err != nil {
				return err
			}
		case 4:
			if err := c.runAction(c.removeEvaluationScreen); err != nil {
				return err
			}
		case 5:
			if err := c.runAction(c.scoresTableScreen); err != nil {
				return err
			}
		case 6:
			if err := c.runAction(c.totalsScreen); err != nil {
				return err
			}
		case 7:
			if err := c.exportMenu(); err != nil {
				return err
			}
		case 0:
			fmt.Fprintln(c.writer, "до встречи!")
			return nil
		default:
			c.invalidChoice(choice)
		}
	}
}

func (c *CLI) studentsMenu() error {
	for {
		c.clearScreen()
		c.printSection("Студенты")
		fmt.Fprintln(c.writer, "  1. Добавить студента")
		fmt.Fprintln(c.writer, "  2. Показать всех")
		fmt.Fprintln(c.writer, "  0. Назад")

		choice, err := c.askInt("Выберите пункт: ")
		if err != nil {
			return err
		}

		switch choice {
		case 1:
			if err := c.runAction(c.addStudent); err != nil {
				return err
			}
		case 2:
			if err := c.runAction(c.listStudents); err != nil {
				return err
			}
		case 0:
			return nil
		default:
			c.invalidChoice(choice)
		}
	}
}

func (c *CLI) criteriaMenu() error {
	for {
		c.clearScreen()
		c.printSection("Критерии")
		fmt.Fprintln(c.writer, "  1. Добавить критерий")
		fmt.Fprintln(c.writer, "  2. Показать все")
		fmt.Fprintln(c.writer, "  0. Назад")

		choice, err := c.askInt("Выберите пункт: ")
		if err != nil {
			return err
		}

		switch choice {
		case 1:
			if err := c.runAction(c.addCriteria); err != nil {
				return err
			}
		case 2:
			if err := c.runAction(c.listCriteria); err != nil {
				return err
			}
		case 0:
			return nil
		default:
			c.invalidChoice(choice)
		}
	}
}

func (c *CLI) exportMenu() error {
	for {
		c.clearScreen()
		c.printSection("Экспорт")

		fmt.Fprintln(c.writer, "  1. Экспорт студентов (.csv или .json)")
		fmt.Fprintln(c.writer, "  0. Назад")

		choice, err := c.askInt("Выберите пункт: ")
		if err != nil {
			return err
		}

		switch choice {
		case 1:
			if err := c.runAction(c.exportStudents); err != nil {
				return err
			}
		case 0:
			return nil
		default:
			c.invalidChoice(choice)
		}
	}
}

func (c *CLI) addStudent() error {
	c.clearScreen()
	c.printSection("Новый студент")

	id, err := c.askInt("ID студента: ")
	if err != nil {
		return err
	}

	name, err := c.askNonEmpty("Имя и фамилия: ")
	if err != nil {
		return err
	}

	if err := c.deps.AddStudent.Execute(usecase.AddStudentCommand{ID: id, Name: name}); err != nil {
		return err
	}

	fmt.Fprintf(c.writer, "\n[OK] Студент %q (id=%d) добавлен.\n", name, id)
	return nil
}

func (c *CLI) addCriteria() error {
	c.clearScreen()
	c.printSection("Новый критерий")

	name, err := c.askNonEmpty("Название: ")
	if err != nil {
		return err
	}

	maxScore, err := c.askInt("Максимальный балл: ")
	if err != nil {
		return err
	}
	if maxScore <= 0 {
		return errors.New("максимальный балл должен быть больше нуля")
	}

	weight, err := c.askFloat("Вес (0..1): ")
	if err != nil {
		return err
	}
	if weight < 0 || weight > 1 {
		return fmt.Errorf("вес должен быть от 0 до 1, получено: %.2f", weight)
	}

	err = c.deps.AddCriteria.Execute(usecase.AddCriteriaCommand{
		Name:        name,
		MaxScoreInt: maxScore,
		WeightFloat: float64(weight),
	})
	if err != nil {
		return err
	}

	fmt.Fprintf(c.writer, "\n[OK] Критерий %q добавлен (макс=%d, вес=%.2f).\n", name, maxScore, weight)
	return nil
}

func (c *CLI) listStudents() error {
	rows, err := c.reportStudents()
	if err != nil {
		return err
	}

	c.clearScreen()
	c.printSection("Список студентов")
	if len(rows) == 0 {
		fmt.Fprintln(c.writer, "Список пуст.")
		return nil
	}

	headers := []string{"ID", "Имя", "Итог"}
	data := make([][]string, 0, len(rows))
	for _, row := range rows {
		data = append(data, []string{strconv.Itoa(row.ID), row.Name, formatTotal(row)})
	}

	c.printTable(headers, data)
	fmt.Fprintf(c.writer, "\nВсего: %d\n", len(rows))
	return nil
}

func (c *CLI) listCriteria() error {
	criteria, err := c.reportCriteria()
	if err != nil {
		return err
	}

	c.clearScreen()
	c.printSection("Список критериев")
	if len(criteria) == 0 {
		fmt.Fprintln(c.writer, "Список пуст.")
		return nil
	}

	headers := []string{"Название", "Макс. балл", "Вес"}
	data := make([][]string, 0, len(criteria))
	for _, row := range criteria {
		data = append(data, []string{row.Name, strconv.Itoa(row.MaxScore), fmt.Sprintf("%.2f", row.Weight)})
	}

	c.printTable(headers, data)
	fmt.Fprintf(c.writer, "\nВсего: %d\n", len(criteria))
	return nil
}

func (c *CLI) evaluateScreen() error {
	students, err := c.reportStudents()
	if err != nil {
		return err
	}
	if len(students) == 0 {
		return errors.New("сначала добавьте хотя бы одного студента")
	}

	criteria, err := c.reportCriteria()
	if err != nil {
		return err
	}
	if len(criteria) == 0 {
		return errors.New("сначала добавьте хотя бы один критерий")
	}

	c.clearScreen()
	c.printSection("Выставить оценку")
	fmt.Fprintln(c.writer, "Выберите студента:")
	student, err := c.selectStudent(students)
	if err != nil {
		return err
	}

	c.clearScreen()
	c.printSection("Выставить оценку")
	fmt.Fprintf(c.writer, "Студент: %s (id=%d)\n\n", student.Name, student.ID)
	fmt.Fprintln(c.writer, "Выберите критерий:")
	criterion, err := c.selectCriteria(criteria)
	if err != nil {
		return err
	}

	score, err := c.askScore(criterion.MaxScore)
	if err != nil {
		return err
	}

	err = c.deps.Evaluate.Execute(usecase.StudentEvaluationCommand{
		StudentID:    student.ID,
		CriteriaName: criterion.Name,
		Score:        score,
	})
	if err != nil {
		return err
	}

	fmt.Fprintf(c.writer, "\n[OK] %s: %s = %d.\n", student.Name, criterion.Name, score)
	return nil
}

func (c *CLI) removeEvaluationScreen() error {
	students, err := c.reportStudents()
	if err != nil {
		return err
	}
	if len(students) == 0 {
		return errors.New("сначала добавьте хотя бы одного студента")
	}

	criteria, err := c.reportCriteria()
	if err != nil {
		return err
	}
	if len(criteria) == 0 {
		return errors.New("сначала добавьте хотя бы один критерий")
	}

	c.clearScreen()
	c.printSection("Удалить оценку")
	fmt.Fprintln(c.writer, "Выберите студента:")
	student, err := c.selectStudent(students)
	if err != nil {
		return err
	}

	c.clearScreen()
	c.printStudentDetails(student, criteria)

	var scored []usecase.CriteriaDTO
	for _, criterion := range criteria {
		if _, ok := student.Scores[criterion.Name]; ok {
			scored = append(scored, criterion)
		}
	}
	if len(scored) == 0 {
		return fmt.Errorf("у студента %s нет оценок", student.Name)
	}

	fmt.Fprintln(c.writer)
	fmt.Fprintln(c.writer, "Выберите критерий для удаления:")
	criterion, err := c.selectCriteria(scored)
	if err != nil {
		return err
	}

	confirmed, err := c.askConfirm(fmt.Sprintf("Удалить оценку %q у %s?", criterion.Name, student.Name))
	if err != nil {
		return err
	}
	if !confirmed {
		return nil
	}

	err = c.deps.RemoveEvaluation.Execute(usecase.RemoveEvaluationCommand{
		StudentID:    student.ID,
		CriteriaName: criterion.Name,
	})
	if err != nil {
		return err
	}

	fmt.Fprintf(c.writer, "\n[OK] Оценка по критерию %q удалена.\n", criterion.Name)
	return nil
}

func (c *CLI) scoresTableScreen() error {
	students, err := c.reportStudents()
	if err != nil {
		return err
	}

	criteria, err := c.reportCriteria()
	if err != nil {
		return err
	}

	c.clearScreen()
	c.printSection("Таблица оценок")
	if len(students) == 0 {
		fmt.Fprintln(c.writer, "Нет студентов.")
		return nil
	}
	if len(criteria) == 0 {
		fmt.Fprintln(c.writer, "Нет критериев.")
		return nil
	}

	headers := []string{"ID", "Имя"}
	for _, criterion := range criteria {
		headers = append(headers, criterion.Name)
	}
	headers = append(headers, "Итог")

	data := make([][]string, 0, len(students))
	for _, student := range students {
		row := []string{strconv.Itoa(student.ID), student.Name}
		for _, criterion := range criteria {
			row = append(row, scoreOf(student, criterion.Name))
		}
		row = append(row, formatTotal(student))
		data = append(data, row)
	}

	c.printTable(headers, data)
	return nil
}

func (c *CLI) totalsScreen() error {
	students, err := c.reportStudents()
	if err != nil {
		return err
	}

	c.clearScreen()
	c.printSection("Итоги по студентам")
	if len(students) == 0 {
		fmt.Fprintln(c.writer, "Нет студентов.")
		return nil
	}

	sort.SliceStable(students, func(i, j int) bool {
		if students[i].Total == students[j].Total {
			return students[i].ID < students[j].ID
		}
		return students[i].Total > students[j].Total
	})

	headers := []string{"Место", "ID", "Имя", "Итог"}
	data := make([][]string, 0, len(students))
	for i, student := range students {
		data = append(data, []string{strconv.Itoa(i + 1), strconv.Itoa(student.ID), student.Name, formatTotal(student)})
	}

	c.printTable(headers, data)
	return nil
}

func (c *CLI) exportStudents() error {
	path, err := c.askPath("Путь к файлу (укажите расширение .csv или .json)", "students.csv")
	if err != nil {
		return err
	}

	if err := c.deps.ExportStudents.ExportStudents(path); err != nil {
		return err
	}

	fmt.Fprintf(c.writer, "\n[OK] Данные успешно экспортированы в %s.\n", path)
	return nil
}

func (c *CLI) reportStudents() ([]usecase.StudentDTO, error) {
	return c.deps.GetReport.GetStudents()
}

func (c *CLI) reportCriteria() ([]usecase.CriteriaDTO, error) {
	return c.deps.GetReport.GetCriteria()
}

func (c *CLI) selectStudent(students []usecase.StudentDTO) (usecase.StudentDTO, error) {
	for i, student := range students {
		fmt.Fprintf(c.writer, "  %d. %s (id=%d)\n", i+1, student.Name, student.ID)
	}
	fmt.Fprintln(c.writer, "  0. Назад")

	for {
		choice, err := c.askInt("Студент: ")
		if err != nil {
			return usecase.StudentDTO{}, err
		}
		if choice == 0 {
			return usecase.StudentDTO{}, errCancel
		}
		if choice >= 1 && choice <= len(students) {
			return students[choice-1], nil
		}
		c.printError(fmt.Errorf("выберите число от 0 до %d", len(students)))
	}
}

func (c *CLI) selectCriteria(criteria []usecase.CriteriaDTO) (usecase.CriteriaDTO, error) {
	for i, criterion := range criteria {
		fmt.Fprintf(c.writer, "  %d. %s (макс. %d)\n", i+1, criterion.Name, criterion.MaxScore)
	}
	fmt.Fprintln(c.writer, "  0. Назад")

	for {
		choice, err := c.askInt("Критерий: ")
		if err != nil {
			return usecase.CriteriaDTO{}, err
		}
		if choice == 0 {
			return usecase.CriteriaDTO{}, errCancel
		}
		if choice >= 1 && choice <= len(criteria) {
			return criteria[choice-1], nil
		}
		c.printError(fmt.Errorf("выберите число от 0 до %d", len(criteria)))
	}
}

func (c *CLI) printStudentDetails(student usecase.StudentDTO, criteria []usecase.CriteriaDTO) {
	c.printSection("Студент")
	fmt.Fprintf(c.writer, "ID: %d\n", student.ID)
	fmt.Fprintf(c.writer, "Имя: %s\n", student.Name)
	fmt.Fprintln(c.writer, "Баллы:")

	if len(criteria) == 0 {
		names := make([]string, 0, len(student.Scores))
		for name := range student.Scores {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			fmt.Fprintf(c.writer, "  %s: %d\n", name, student.Scores[name])
		}
	} else {
		for _, criterion := range criteria {
			fmt.Fprintf(c.writer, "  %s: %s / %d\n", criterion.Name, scoreOf(student, criterion.Name), criterion.MaxScore)
		}
	}

	fmt.Fprintf(c.writer, "Итог: %s\n", formatTotal(student))
}

func (c *CLI) runAction(action func() error) error {
	if err := action(); err != nil {
		if errors.Is(err, errExit) {
			return err
		}
		if errors.Is(err, errCancel) {
			return nil
		}
		c.printError(err)
	}
	c.pause()
	return nil
}

func (c *CLI) askInt(label string) (int, error) {
	for {
		line, err := c.readLine(label)
		if err != nil {
			return 0, err
		}
		if line == "" {
			c.printError(errors.New("введите целое число"))
			continue
		}
		value, err := strconv.Atoi(line)
		if err != nil {
			c.printError(fmt.Errorf("%q — не целое число", line))
			continue
		}
		return value, nil
	}
}

func (c *CLI) askFloat(label string) (float64, error) {
	for {
		line, err := c.readLine(label)
		if err != nil {
			return 0, err
		}
		if line == "" {
			c.printError(errors.New("введите число"))
			continue
		}
		value, err := strconv.ParseFloat(strings.ReplaceAll(line, ",", "."), 64)
		if err != nil {
			c.printError(fmt.Errorf("%q — не число", line))
			continue
		}
		return value, nil
	}
}

func (c *CLI) askScore(maxScore int) (int, error) {
	for {
		score, err := c.askInt(fmt.Sprintf("Балл (0..%d): ", maxScore))
		if err != nil {
			return 0, err
		}
		// повторная проверка для удобства чтобы не возвращать ошибку через Use Case
		if score < 0 || score > maxScore {
			c.printError(fmt.Errorf("балл должен быть от 0 до %d", maxScore))
			continue
		}
		return score, nil
	}
}

func (c *CLI) askNonEmpty(label string) (string, error) {
	for {
		line, err := c.readLine(label)
		if err != nil {
			return "", err
		}
		if line != "" {
			return line, nil
		}
		c.printError(errors.New("значение не может быть пустым"))
	}
}

func (c *CLI) askConfirm(label string) (bool, error) {
	for {
		line, err := c.readLine(label + " (д/н): ")
		if err != nil {
			return false, err
		}
		switch strings.ToLower(line) {
		case "д", "да", "y", "yes":
			return true, nil
		case "н", "нет", "n", "no":
			return false, nil
		default:
			c.printError(errors.New("ответьте «д» или «н»"))
		}
	}
}

func (c *CLI) askPath(label string, defaultPath string) (string, error) {
	line, err := c.readLine(fmt.Sprintf("%s [%s]: ", label, defaultPath))
	if err != nil {
		return "", err
	}
	if line == "" {
		return defaultPath, nil
	}
	return line, nil
}

func (c *CLI) readLine(label string) (string, error) {
	fmt.Fprint(c.writer, label)

	if !c.reader.Scan() {
		fmt.Fprintln(c.writer)
		if err := c.reader.Err(); err != nil {
			return "", err
		}
		return "", errExit
	}

	return strings.TrimSpace(c.reader.Text()), nil
}

func (c *CLI) pause() {
	fmt.Fprint(c.writer, "\nНажмите Enter, чтобы продолжить...")
	c.reader.Scan()
}

func (c *CLI) clearScreen() {
	fmt.Fprint(c.writer, "\033[H\033[2J")
}

func (c *CLI) invalidChoice(choice int) {
	c.printError(fmt.Errorf("нет такого пункта: %d", choice))
	c.pause()
}

func (c *CLI) printError(err error) {
	fmt.Fprintf(c.errWriter, "[!] %v\n", err)
}

func (c *CLI) printSection(title string) {
	fmt.Fprintf(c.writer, "── %s ", title)
	used := utf8.RuneCountInString(title) + 4
	if used < sectionWidth {
		fmt.Fprint(c.writer, strings.Repeat("─", sectionWidth-used))
	}
	fmt.Fprintln(c.writer)
}

func (c *CLI) printMainMenu() {
	fmt.Fprintln(c.writer, "╔════════════════════════════════════════════╗")
	fmt.Fprintln(c.writer, "║           TAS — Оценка студентов           ║")
	fmt.Fprintln(c.writer, "╠════════════════════════════════════════════╣")
	fmt.Fprintln(c.writer, "║  1. Студенты                               ║")
	fmt.Fprintln(c.writer, "║  2. Критерии                               ║")
	fmt.Fprintln(c.writer, "║  3. Выставить оценку                       ║")
	fmt.Fprintln(c.writer, "║  4. Удалить оценку                         ║")
	fmt.Fprintln(c.writer, "║  5. Таблица оценок                         ║")
	fmt.Fprintln(c.writer, "║  6. Итоги по студентам                     ║")
	fmt.Fprintln(c.writer, "║  7. Экспорт                                ║")
	fmt.Fprintln(c.writer, "║  0. Выход                                  ║")
	fmt.Fprintln(c.writer, "╚════════════════════════════════════════════╝")
}

func (c *CLI) printTable(headers []string, rows [][]string) {
	widths := make([]int, len(headers))
	for i, header := range headers {
		widths[i] = utf8.RuneCountInString(header)
	}
	for _, row := range rows {
		for i := range headers {
			if i >= len(row) {
				continue
			}
			if width := utf8.RuneCountInString(row[i]); width > widths[i] {
				widths[i] = width
			}
		}
	}

	printRow := func(cells []string) {
		for i := range widths {
			cell := ""
			if i < len(cells) {
				cell = cells[i]
			}
			fmt.Fprint(c.writer, cell)
			padding := widths[i] - utf8.RuneCountInString(cell)
			if i < len(widths)-1 {
				padding += 2
			}
			if padding > 0 {
				fmt.Fprint(c.writer, strings.Repeat(" ", padding))
			}
		}
		fmt.Fprintln(c.writer)
	}

	printRow(headers)

	separator := make([]string, len(widths))
	for i, width := range widths {
		separator[i] = strings.Repeat("-", width)
	}
	printRow(separator)

	for _, row := range rows {
		printRow(row)
	}
}

func scoreOf(student usecase.StudentDTO, criteriaName string) string {
	score, ok := student.Scores[criteriaName]
	if !ok {
		return "-"
	}
	return strconv.Itoa(score)
}

func formatTotal(student usecase.StudentDTO) string {
	if student.MaxTotal > 0 {
		return fmt.Sprintf("%.2f / %.2f", student.Total, student.MaxTotal)
	}
	return fmt.Sprintf("%.2f", student.Total)
}
