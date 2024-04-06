package api

import (
	"coursach/models"
	"encoding/json"
	"errors"
	"strconv"
)

// Save objects to DB and return saved object.
func saveObjToDB(data *json.Decoder, objectType string) (interface{}, error) {
	var err error
	var id int
	switch objectType {
	case "institutes":
		obj := models.Institutes{}
		err = data.Decode(&obj)
		if err == nil {
			id, err = models.Save(objectType, obj)
			obj.Id = int64(id)
		}
		return obj, err
	case "lecturers":
		obj := models.Lecturers{}
		err = data.Decode(&obj)
		if err == nil {
			id, err = models.Save(objectType, obj)
			obj.Id = int64(id)
		}
		return obj, err
	case "cabinets":
		obj := models.Cabinets{}
		err = data.Decode(&obj)
		if err == nil {
			id, err = models.Save(objectType, obj)
			obj.Id = int64(id)
		}
		return obj, err
	case "training_directions":
		obj := models.TrainingDirections{}
		err = data.Decode(&obj)
		if err == nil {
			id, err = models.Save(objectType, obj)
			obj.Id = int64(id)
		}
		return obj, err
	case "academic_subjects":
		obj := models.AcademicSubjects{}
		err = data.Decode(&obj)
		if err == nil {
			id, err = models.Save(objectType, obj)
			obj.Id = int64(id)
		}
		return obj, err
	case "group_of_students":
		obj := models.GroupOfStudents{}
		err = data.Decode(&obj)
		if err == nil {
			id, err = models.Save(objectType, obj)
			obj.Id = int64(id)
		}
		return obj, err
	case "lead_subjects":
		obj := models.LeadSubjects{}
		err = data.Decode(&obj)
		if err == nil {
			id, err = models.Save(objectType, obj)
			obj.Id = int64(id)
		}
		return obj, err
	case "result_days_exams":
		obj := models.ResultDaysExams{}
		err = data.Decode(&obj)
		if err == nil {
			id, err = models.Save(objectType, obj)
			obj.Id = int64(id)
		}
		return obj, err
	default:
		return nil, errors.New("api/dbInterface.go: saveObjToDB: Неверный тип данных")
	}
}

// Delete object from DB and return deleted object. Return error if object does not exists.
func deleteObjFromDB(objectType string, id int) (interface{}, error) {
	obj, err := models.Find(objectType, "id", strconv.Itoa(id))
	if err != nil {
		return nil, err
	}
	err = models.Delete(objectType, obj)
	return obj, err
}
