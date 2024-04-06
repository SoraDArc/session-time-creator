package models

import (
	"errors"
	"strconv"
)

type DBObject interface{}

// Return models by tableName and error if table no exist.
func FactoryObjByName(tableName string) (interface{}, error) {
	switch tableName {
	case "academic_subjects":
		return AcademicSubjects{}, nil
	case "cabinets":
		return Cabinets{}, nil
	case "group_of_students":
		return GroupOfStudents{}, nil
	case "institutes":
		return Institutes{}, nil
	case "lead_subjects":
		return LeadSubjects{}, nil
	case "lecturers":
		return Lecturers{}, nil
	case "training_directions":
		return TrainingDirections{}, nil
	default:
		return nil, errors.New("такой таблицы не существует")
	}
}

// Return tableName by field name and error if table by field name not exist.
func FactoryTableNameByField(field string) (string, error) {
	switch field {
	case "Institutes_Id":
		return "institutes", nil
	case "Training_Directions_Id":
		return "training_directions", nil
	case "Academic_Subjects_Id":
		return "academic_subjects", nil
	case "Lecturers_Id":
		return "lecturers", nil
	case "Group_Of_Students_Id":
		return "group_of_students", nil
	case "Cabinets_Id":
		return "cabinets", nil
	default:
		return "", errors.New("таких данных не существует")
	}
}

// Return object by tableName and id and error if table not exist.
func FactoryObjByNameAndId(objectType string, id int) (DBObject, error) {
	obj, err := Find(objectType, "id", strconv.Itoa(id))
	return obj, err
}

// Return all objects from DB table and error if table not exist.
func FactoryObjList(objectType string) ([]DBObject, error) {
	switch objectType {
	case "institutes":
		list := AllInstitutes()
		result := make([]DBObject, len(list))
		for i, v := range list {
			result[i] = v
		}
		return result, nil
	case "lecturers":
		list := AllLecturers()
		result := make([]DBObject, len(list))
		for i, v := range list {
			result[i] = v
		}
		return result, nil
	case "cabinets":
		list := AllCabinets()
		result := make([]DBObject, len(list))
		for i, v := range list {
			result[i] = v
		}
		return result, nil
	case "training_directions":
		list := AllTrainingDirections()
		result := make([]DBObject, len(list))
		for i, v := range list {
			result[i] = v
		}
		return result, nil
	case "academic_subjects":
		list := AllAcademicSubjects()
		result := make([]DBObject, len(list))
		for i, v := range list {
			result[i] = v
		}
		return result, nil
	case "group_of_students":
		list := AllGroupOfStudents()
		result := make([]DBObject, len(list))
		for i, v := range list {
			result[i] = v
		}
		return result, nil
	case "lead_subjects":
		list := AllLeadSubjects()
		result := make([]DBObject, len(list))
		for i, v := range list {
			result[i] = v
		}
		return result, nil
	case "result_days_exams":
		list := AllResultDaysExams()
		result := make([]DBObject, len(list))
		for i, v := range list {
			result[i] = v
		}
		return result, nil
	default:
		return nil, errors.New("models/factory.go:FactoryObjList: Неверный тип данных")
	}
}
