package models

import (
	"coursach/db"
	"errors"
	"strconv"
)

// Return sorted list by count and frequently encountered  groups.
func (ls *LeadSubjects) GetCurrentData() []LeadSubjects {
	query := "SELECT ls.id, ls.academic_subjects_id, ls.course, ls.group_of_students_id, ls.lecturers_id, ls.comment FROM lead_subjects ls INNER JOIN (SELECT COUNT(*) AS the_count, group_of_students_id FROM lead_subjects GROUP BY group_of_students_id) t ON ls.group_of_students_id = t.group_of_students_id ORDER BY t.the_count DESC, ls.group_of_students_id DESC"
	rows, err := db.DbConnect.Query(query)
	if err != nil {
		panic(errors.New("models/methods.go: GetCurrentData: не удалось выполнить запрос"))
	}
	var LS []LeadSubjects
	for rows.Next() {
		temp := LeadSubjects{}
		var as_id, gos_id, lect_id int
		var comm string
		err := rows.Scan(&temp.Id, &as_id, &temp.Course, &gos_id, &lect_id, &comm)
		if err != nil {
			panic(errors.New("models/methods.go: GetCurrentData: Ошибка чтения данных из rows"))
		}
		help, _ := Find("academic_subjects", "id", strconv.Itoa(as_id))
		temp.Academic_Subjects_Id = help.(AcademicSubjects)
		help, _ = Find("group_of_students", "id", strconv.Itoa(gos_id))
		temp.Group_Of_Students_Id = help.(GroupOfStudents)
		help, _ = Find("lecturers", "id", strconv.Itoa(lect_id))
		temp.Lecturers_Id = help.(Lecturers)
		temp.Comment = comm
		LS = append(LS, temp)
	}
	return LS
}
