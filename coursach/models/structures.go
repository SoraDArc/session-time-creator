package models

//Имена таблиц -> всё с прописными  и если строк из нескольких слов, то разделяется "_";
type Institutes struct {
	Id          int64
	Title       string
	Short_Title string
}

type Lecturers struct {
	Id         int64
	Short_Name string
	Surname    string
	Name       string
	Patronymic string
}

type Cabinets struct {
	Id            int64
	Title         string
	Cabinets_Type string
}

type TrainingDirections struct {
	Id                 int64
	Num                int64
	Code               string
	Title              string
	Level_Of_Education string
	Institutes_Id      Institutes
}

type AcademicSubjects struct {
	Id                     int64
	Title                  string
	Training_Directions_Id TrainingDirections
}

type GroupOfStudents struct {
	Id                     int64
	Title                  string
	Students_Count         int64
	Training_Directions_Id TrainingDirections
}

type LeadSubjects struct {
	Id                   int64
	Academic_Subjects_Id AcademicSubjects
	Course               int64
	Group_Of_Students_Id GroupOfStudents
	Lecturers_Id         Lecturers
	Comment              string
}

//добавить строку требования requirement string и cabinets_id

type ResultDaysExams struct {
	Id               int64
	Lead_Subjects_Id LeadSubjects
	Consultation     string
	Exam             string
	Has_Date         string
}

// Result DAys Exams - List_of_exams
//2 массива из расписания в бд записать как Наименование string  ПоставленныеДни string(JSON) НеПоставленныеДни string(JSON)

//Пользователь
