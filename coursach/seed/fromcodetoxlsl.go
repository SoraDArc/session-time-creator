package seed

import (
	"coursach/db"
	"coursach/models"
	"errors"
	"fmt"

	"github.com/xuri/excelize/v2"
)

func createAndFillInExcel(data []models.ResultDaysExams, name string) {
	file := excelize.NewFile()
	file.SetSheetName("Sheet1", "Расписание сессии")

	headers := []string{"Предмет", "Курс", "Группа", "Преподаватель", "Консультация", "Экзамен", "Комментарий"}
	//SetCellValue(1,2,3) - 1 аргумент это название листа
	// 2 аргумент это сссылка на клетку (Типо А1)
	// 3 аргумент это значение
	for i, header := range headers {
		file.SetCellValue("Расписание сессии", fmt.Sprintf("%s%d", string(rune(65+i)), 1), header)
	}

	for i, obj := range data {
		dataRow := i + 2
		academic_subj := obj.Lead_Subjects_Id.Academic_Subjects_Id.Title
		course := obj.Lead_Subjects_Id.Course
		group := obj.Lead_Subjects_Id.Group_Of_Students_Id.Title
		lect := obj.Lead_Subjects_Id.Lecturers_Id.Short_Name
		consultation := obj.Consultation
		exam := obj.Exam
		comment := obj.Lead_Subjects_Id.Comment
		file.SetCellValue("Расписание сессии", fmt.Sprintf("%s%d", string(rune(65+0)), dataRow), academic_subj)
		file.SetCellValue("Расписание сессии", fmt.Sprintf("%s%d", string(rune(65+1)), dataRow), course)
		file.SetCellValue("Расписание сессии", fmt.Sprintf("%s%d", string(rune(65+2)), dataRow), group)
		file.SetCellValue("Расписание сессии", fmt.Sprintf("%s%d", string(rune(65+3)), dataRow), lect)
		file.SetCellValue("Расписание сессии", fmt.Sprintf("%s%d", string(rune(65+4)), dataRow), consultation)
		file.SetCellValue("Расписание сессии", fmt.Sprintf("%s%d", string(rune(65+5)), dataRow), exam)
		file.SetCellValue("Расписание сессии", fmt.Sprintf("%s%d", string(rune(65+6)), dataRow), comment)
	}

	if err := file.SaveAs("./readymadeschedule/raspisanie-2024" + name + ".xlsx"); err != nil {
		panic(errors.New("seed: fromcodetoxlsl.go: createAndFillInExcel: Ошибка сохранения эксель файла"))
	}
}

func ToExcelFile() {
	db.Connect("postgres", "12345678", "sessiontime")
	defer db.Close()

	GoodResult, BadResult := models.MakeSchedule("10.06.2024", "06.07.2024")
	createAndFillInExcel(GoodResult, "good")
	createAndFillInExcel(BadResult, "bad")
}
