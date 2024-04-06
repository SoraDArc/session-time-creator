package seed

import (
	"coursach/models"
	"errors"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/xuri/excelize/v2"
)

type MyObjects interface{}

// Read xslx file.
func GetDataFromXSLX(fileName string, sheet string) [][]string {
	//Open File
	pwd, _ := os.Getwd()
	filePath := pwd + "\\data\\" + fileName
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// Get all the rows in the sheet.
	rows, err := file.GetRows(sheet)
	if err != nil {
		log.Fatal(err)
	}

	return rows
}

// Load objects from xlsx file and create nessesary structs from "structures.go". Returns []Objects, TableName, error.
func LoadObjects(name string) ([]MyObjects, string, error) {
	switch name {
	case "Institutes":
		rows := GetDataFromXSLX("SessionTime2024.xlsx", "Факультеты")
		//Create []obj
		var Institute []MyObjects
		for _, row := range rows {
			tempInst := models.Institutes{
				Title:       row[0],
				Short_Title: row[1],
			}
			Institute = append(Institute, tempInst)
		}

		return Institute, "institutes", nil
	case "Lecturers":
		rows := GetDataFromXSLX("SessionTime2024.xlsx", "Преподаватели")
		//Create obj
		var Lecturer []MyObjects
		for _, row := range rows {
			tempLect := models.Lecturers{
				Short_Name: row[0],
				Surname:    row[1],
				Name:       row[2],
				Patronymic: row[3],
			}
			Lecturer = append(Lecturer, tempLect)
		}

		return Lecturer, "lecturers", nil
	case "Cabinets":
		rows := GetDataFromXSLX("SessionTime2024.xlsx", "Аудитории")
		//Create obj
		var Cabinet []MyObjects
		for _, row := range rows {
			tempCab := models.Cabinets{
				Title:         row[0],
				Cabinets_Type: row[1],
			}
			Cabinet = append(Cabinet, tempCab)
		}

		return Cabinet, "cabinets", nil
	case "TrainingDirections":
		rows := GetDataFromXSLX("SessionTime2024.xlsx", "Направления подготовки")
		//Create obj
		var TrainingDirection []MyObjects
		for _, row := range rows {
			num, _ := strconv.Atoi(row[0])
			tempTD := models.TrainingDirections{
				Num:                int64(num),
				Code:               row[1],
				Title:              row[2],
				Level_Of_Education: row[3],
			}
			help, _ := models.Find("institutes", "short_title", row[4])
			tempTD.Institutes_Id = help.(models.Institutes)
			TrainingDirection = append(TrainingDirection, tempTD)
		}

		return TrainingDirection, "training_directions", nil
	case "AcademicSubjects":
		rows := GetDataFromXSLX("SessionTime2024.xlsx", "Предметы")
		//Create obj
		var AcademicSubject []MyObjects
		for _, row := range rows {
			tempAS := models.AcademicSubjects{
				Title: row[0],
			}
			help, _ := models.Find("training_directions", "id", row[2])
			tempAS.Training_Directions_Id = help.(models.TrainingDirections)
			AcademicSubject = append(AcademicSubject, tempAS)
		}

		return AcademicSubject, "academic_subjects", nil
	case "GroupOfStudents":
		rows := GetDataFromXSLX("SessionTime2024.xlsx", "Группы")
		//Create obj
		var GroupOfStudent []MyObjects
		for _, row := range rows {
			StudCnt, _ := strconv.Atoi(row[1])
			tempGOS := models.GroupOfStudents{
				Title:          row[0],
				Students_Count: int64(StudCnt),
			}
			help, _ := models.Find("training_directions", "id", row[2])
			tempGOS.Training_Directions_Id = help.(models.TrainingDirections)
			GroupOfStudent = append(GroupOfStudent, tempGOS)
		}

		return GroupOfStudent, "group_of_students", nil
	case "LeadSubjects":
		rows := GetDataFromXSLX("SessionTime2024.xlsx", "Ведут предмет")
		//Create obj
		var LeadSubject []MyObjects
		for _, row := range rows {
			tempLS := models.LeadSubjects{}
			help, _ := models.Find("academic_subjects", "title", row[0])
			tempLS.Academic_Subjects_Id = help.(models.AcademicSubjects)
			Crs, _ := strconv.Atoi(row[1])
			tempLS.Course = int64(Crs)
			help, _ = models.Find("group_of_students", "title", row[2])
			tempLS.Group_Of_Students_Id = help.(models.GroupOfStudents)
			help, _ = models.Find("lecturers", "short_name", row[3])
			tempLS.Lecturers_Id = help.(models.Lecturers)
			if len(row) > 4 {
				tempLS.Comment = row[4]
			}
			LeadSubject = append(LeadSubject, tempLS)
		}

		return LeadSubject, "lead_subjects", nil
	default:
		return nil, "Error", errors.New("seed:fromxslstodb:LoadObjects: Неверный тип ")
	}
}
