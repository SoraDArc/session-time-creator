package models

import (
	"coursach/db"
	"database/sql"
	"errors"
	"reflect"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type MyObject interface{}

func HandleError(err error) {
	if err != nil {
		panic(errors.New("ошибка в models/orm.go").Error() + err.Error())
	}
}

// Returns []string - fields and values of struct.
func GetFieldsAndValues(obj MyObject) ([]string, []string) {
	var resFields []string
	var resValues []string
	t := reflect.TypeOf(obj)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		resFields = append(resFields, strings.ToLower(field.Name))
		fieldValue := reflect.ValueOf(obj).FieldByName(field.Name)
		var temp string
		if field.Type.Kind() == reflect.Int {
			temp = strconv.FormatInt(fieldValue.Int(), 10)
		} else if field.Type.Kind() == reflect.Struct {
			help := fieldValue.FieldByName("Id").Int()
			temp = strconv.FormatInt(help, 10)
		} else if field.Type.Kind() == reflect.String {
			temp = fieldValue.String()
		} else {
			help := fieldValue.Interface().(int64)
			temp = strconv.FormatInt(help, 10)
		}
		resValues = append(resValues, temp)
	}
	return resFields, resValues
}

// Returns query for Insert to DB.
func BuildInsert(tableName string, obj MyObject) string {
	resValues := ""
	resFields := ""
	Fields, Values := GetFieldsAndValues(obj)
	for i := 1; i < len(Fields); i++ {
		if i < len(Fields)-1 {
			resFields += Fields[i] + ","
			resValues += "'" + Values[i] + "',"
		} else {
			resFields += Fields[i]
			resValues += "'" + Values[i] + "'"
		}
	}
	resValues = "(" + resValues + ")"
	resFields = "(" + resFields + ")"
	result := "INSERT INTO " + tableName + " " + resFields + " VALUES " + resValues + " returning id"
	return result
}

// Returns query for Update to DB.
func BuildUpdate(tableName string, obj MyObject) string {
	Fields, Values := GetFieldsAndValues(obj)
	temp := ""
	for i := 1; i < len(Fields); i++ {
		if i < len(Fields)-1 {
			temp += Fields[i] + "=" + "'" + Values[i] + "',"
		} else {
			temp += Fields[i] + "=" + "'" + Values[i] + "'"
		}
	}
	res := "UPDATE " + tableName + " SET " + temp + " WHERE " + Fields[0] + "=" + Values[0]
	return res
}

// Save if id is zero else update object in DB.
func Save(tableName string, obj MyObject) (int, error) {
	id := reflect.ValueOf(obj).FieldByName("Id").Int()
	var err error
	if id == 0 {
		err = db.DbConnect.QueryRow(BuildInsert(tableName, obj)).Scan(&id)
		if err != nil {
			panic(errors.New("models/orm.go:Save: Ошибка в выполнение insert"))
		}
		return int(id), nil
	} else {
		res, err := db.DbConnect.Exec(BuildUpdate(tableName, obj))
		if err != nil {
			panic(errors.New("models/orm.go:Save: Ошибка в выполнение update"))
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			return 0, errors.New("такого объекта нет")
		} else {
			return int(id), nil
		}
	}
}

// Search object in table by field and value.
func Find(tableName, Field, value string) (MyObject, error) {
	//Фабрика объектов
	obj, err := FactoryObjByName(tableName)
	if err != nil {
		return nil, err
	}

	//Создаём select запрос по условию
	resColumns := ""
	columns, _ := GetFieldsAndValues(obj)
	for i := 0; i < len(columns); i++ {
		if i < len(columns)-1 {
			resColumns += columns[i] + ","
		} else {
			resColumns += columns[i]
		}
	}
	query := "SELECT " + resColumns + " FROM " + tableName + " WHERE " + Field + "='" + value + "'"
	// Массив названий полей
	var Fields []string
	t := reflect.TypeOf(obj)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		Fields = append(Fields, field.Name)
	}
	//Массив с ссылками для считывания результата после запроса
	values := make([]interface{}, len(Fields))
	for i := range Fields {
		var value interface{}
		values[i] = &value
	}
	err = db.DbConnect.QueryRow(query).Scan(values...)
	if err != nil {
		return nil, errors.New("данного объекта не существует")
	}

	//Создаём и заполняем новый объект - newObj
	newObj := reflect.New(reflect.TypeOf(obj)).Elem()
	for i, field := range Fields {
		fieldValue := newObj.FieldByName(field)
		if fieldValue.IsValid() && fieldValue.CanSet() {
			if fieldValue.Kind() == reflect.Struct {
				tn, err := FactoryTableNameByField(field)
				if err != nil {
					panic(errors.New("models/orm.go:Find: не существует таблицы для этого поля ").Error() + field)
				}
				idInter := (*(values[i].(*interface{}))).(int64)
				id := strconv.Itoa(int(idInter))
				tempObj, _ := Find(tn, "id", id)
				fieldValue.Set(reflect.ValueOf(tempObj))

			} else {
				valuePtr := values[i].(*interface{})
				fieldValue.Set(reflect.ValueOf(*valuePtr))
			}
		}
	}
	return newObj.Interface(), nil
}

// Delete object from DB. Return error if object not exists.
func Delete(tableName string, obj MyObject) error {
	var err error
	id := reflect.ValueOf(obj).FieldByName("Id").Int()
	if id != 0 {
		var result sql.Result
		query := "DELETE FROM " + tableName + " WHERE id=$1"
		result, err = db.DbConnect.Exec(query, id)
		if err == nil {
			affRows, _ := result.RowsAffected()
			if affRows == 0 {
				err = errors.New("такого объекта не существует, поэтому не удалено")
			}
		}
	}
	return err
}

//Return list of objects from DB.

func AllInstitutes() []Institutes {
	var result []Institutes
	rows, err := db.DbConnect.Query("SELECT id,title,short_title FROM institutes")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		temp := Institutes{}
		err := rows.Scan(&temp.Id, &temp.Title, &temp.Short_Title)
		if err != nil {
			panic(err)
		}
		result = append(result, temp)
	}
	return result
}

func AllLecturers() []Lecturers {
	var result []Lecturers
	rows, err := db.DbConnect.Query("SELECT id,short_name,surname,name,patronymic FROM lecturers")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		temp := Lecturers{}
		err := rows.Scan(&temp.Id, &temp.Short_Name, &temp.Surname, &temp.Name, &temp.Patronymic)
		if err != nil {
			panic(err)
		}
		result = append(result, temp)
	}
	return result
}

func AllCabinets() []Cabinets {
	var result []Cabinets
	rows, err := db.DbConnect.Query("SELECT id,title,cabinets_type FROM cabinets")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		temp := Cabinets{}
		err := rows.Scan(&temp.Id, &temp.Title, &temp.Cabinets_Type)
		if err != nil {
			panic(err)
		}
		result = append(result, temp)
	}
	return result
}

func AllTrainingDirections() []TrainingDirections {
	var result []TrainingDirections
	rows, err := db.DbConnect.Query("SELECT id,num,code,title,level_of_education,institutes_id FROM training_directions")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		temp := TrainingDirections{}
		var i_id int
		err := rows.Scan(&temp.Id, &temp.Num, &temp.Code, &temp.Title, &temp.Level_Of_Education, &i_id)
		if err != nil {
			panic(err)
		}
		help, _ := Find("institutes", "id", strconv.Itoa(i_id))
		temp.Institutes_Id = help.(Institutes)
		result = append(result, temp)
	}
	return result
}

func AllAcademicSubjects() []AcademicSubjects {
	var result []AcademicSubjects
	rows, err := db.DbConnect.Query("SELECT id,title,training_directions_id FROM academic_subjects")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		temp := AcademicSubjects{}
		var as_id int
		err := rows.Scan(&temp.Id, &temp.Title, &as_id)
		if err != nil {
			panic(err)
		}
		help, _ := Find("training_directions", "id", strconv.Itoa(as_id))
		temp.Training_Directions_Id = help.(TrainingDirections)
		result = append(result, temp)
	}
	return result
}

func AllGroupOfStudents() []GroupOfStudents {
	var result []GroupOfStudents
	rows, err := db.DbConnect.Query("SELECT id,title,students_count,training_directions_id FROM group_of_students")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		temp := GroupOfStudents{}
		var td_id int
		err := rows.Scan(&temp.Id, &temp.Title, &temp.Students_Count, &td_id)
		if err != nil {
			panic(err)
		}
		help, _ := Find("training_directions", "id", strconv.Itoa(td_id))
		temp.Training_Directions_Id = help.(TrainingDirections)
		result = append(result, temp)
	}
	return result
}

func AllLeadSubjects() []LeadSubjects {
	var result []LeadSubjects
	rows, err := db.DbConnect.Query("SELECT id,academic_subjects_id,course,group_of_students_id,lecturers_id FROM lead_subjects")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		temp := LeadSubjects{}
		var as_id, gos_id, lect_id int
		err := rows.Scan(&temp.Id, &as_id, &temp.Course, &gos_id, &lect_id)
		if err != nil {
			panic(err)
		}
		help, _ := Find("academic_subjects", "id", strconv.Itoa(as_id))
		temp.Academic_Subjects_Id = help.(AcademicSubjects)
		help, _ = Find("group_of_students", "id", strconv.Itoa(gos_id))
		temp.Group_Of_Students_Id = help.(GroupOfStudents)
		help, _ = Find("lecturers", "id", strconv.Itoa(lect_id))
		temp.Lecturers_Id = help.(Lecturers)
		result = append(result, temp)
	}
	return result
}

func AllResultDaysExams() []ResultDaysExams {
	var result []ResultDaysExams
	rows, err := db.DbConnect.Query("SELECT id,lead_subjects_id,consultation,exam,has_date FROM result_days_exams")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		temp := ResultDaysExams{}
		var ls_id int
		err := rows.Scan(&temp.Id, &ls_id, &temp.Consultation, &temp.Exam, &temp.Has_Date)
		if err != nil {
			panic(err)
		}
		help, _ := Find("lead_subjects", "id", strconv.Itoa(ls_id))
		temp.Lead_Subjects_Id = help.(LeadSubjects)
		result = append(result, temp)
	}
	return result
}
