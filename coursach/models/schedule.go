package models

import (
	"errors"
	"fmt"
	"time"
)

// Contains LeadSubj for giving date.
type accordence struct {
	Id               int
	Lead_Subjects_Id LeadSubjects
	Consultation     time.Time
	Exam             time.Time
	HasDate          bool
}

// Makes a schedule for the session time. Returns two lists - LeadSubjects with succesful date in session time and LeadSubjects without dates.
func MakeSchedule(dateStart, dateEnd string) ([]ResultDaysExams, []ResultDaysExams) {
	//Проверка на правильность даты
	layout := "02.01.2006" // - Формат
	time_start, err := time.Parse(layout, dateStart)
	if err != nil {
		panic(errors.New("schedule.go: makeschedule: time_start - invalid format").Error() + err.Error())
	}
	time_end, err := time.Parse(layout, dateEnd)
	if err != nil {
		panic(errors.New("schedule.go: makeschedule: time_end - invalid format").Error() + err.Error())
	}

	//Для проверок: в день не больше 14 экзаменов. Почему 14? Потому что кабинетов нетб больше.
	FullDay := make(map[time.Time]int) // - map[date]count
	for i := time_start; !i.After(time_end); i = i.AddDate(0, 0, 1) {
		FullDay[i] = 0
	}

	//Все сортированные данные из таблицы lead_subjects
	CurrentData := (&LeadSubjects{}).GetCurrentData()

	//Составление расписания
	Busy := make(map[string]bool) // - map[Name+Date]bool
	TempResult := []accordence{}
	CurGroup := GroupOfStudents{Id: 0, Title: "00000", Students_Count: 0}
	for index := 0; index < len(CurrentData); index++ {
		if CurGroup == CurrentData[index].Group_Of_Students_Id {
			TempResult = append(TempResult, accordence{Lead_Subjects_Id: CurrentData[index]})
			continue
		}
		CurGroup = CurrentData[index].Group_Of_Students_Id
		for i := time_start; !i.After(time_end); {
			if FullDay[i] > 14 || i.Format(layout) == "12.06.2024" {
				i = i.AddDate(0, 0, 1)
				continue
			}
			helpAccord := accordence{}
			var weekday string
			helpAccord.Lead_Subjects_Id = CurrentData[index]

			//день консультации
			weekday = fmt.Sprintf("%v", i.Weekday())
			if weekday == "Sunday" || i.Format(layout) == "12.06.2024" {
				i = i.AddDate(0, 0, 1)
			}
			helpAccord.Consultation = i
			FullDay[i]++
			i = i.AddDate(0, 0, 1)
			if _, inMap := Busy[helpAccord.Lead_Subjects_Id.Lecturers_Id.Short_Name+" "+i.Format(layout)]; inMap {
				helpAccord.Consultation = helpAccord.Consultation.AddDate(0, 0, 1)
				i = i.AddDate(0, 0, 1)
			}

			//день экзамена
			weekday = fmt.Sprintf("%v", i.Weekday())
			if weekday == "Sunday" || i.Format(layout) == "12.06.2024" {
				i = i.AddDate(0, 0, 1)
			}
			helpAccord.Exam = i
			i = i.AddDate(0, 0, 2)

			if !helpAccord.Exam.After(time_end) {
				helpAccord.HasDate = true
			}

			if _, inMap := Busy[helpAccord.Lead_Subjects_Id.Lecturers_Id.Short_Name+" "+helpAccord.Exam.Format(layout)]; !inMap && helpAccord.HasDate {
				Busy[helpAccord.Lead_Subjects_Id.Lecturers_Id.Short_Name+" "+helpAccord.Exam.Format(layout)] = true
			}

			TempResult = append(TempResult, helpAccord)

			if index+1 < len(CurrentData) {
				if CurGroup != CurrentData[index+1].Group_Of_Students_Id || helpAccord.Exam.After(time_end) {
					break
				}
				if i.After(time_end) {
					break
				}
				index++
			} else {
				break
			}
		}
	}
	//Перебор на выставленные дни и нет
	GoodResult := []ResultDaysExams{}
	BadResult := []ResultDaysExams{}
	for _, tempRes := range TempResult {
		if tempRes.HasDate {
			temp := ResultDaysExams{
				Lead_Subjects_Id: tempRes.Lead_Subjects_Id,
				Consultation:     tempRes.Consultation.Format(layout),
				Exam:             tempRes.Exam.Format(layout),
				Has_Date:         "Да",
			}
			GoodResult = append(GoodResult, temp)
		} else {
			tempRes.Consultation = time.Time{}
			tempRes.Exam = time.Time{}
			temp := ResultDaysExams{
				Lead_Subjects_Id: tempRes.Lead_Subjects_Id,
				Consultation:     tempRes.Consultation.Format(layout),
				Exam:             tempRes.Exam.Format(layout),
				Has_Date:         "Нет",
			}
			BadResult = append(BadResult, temp)
		}
	}
	return GoodResult, BadResult
}

/*
для будущей доработки выбора свободного препода
func main() {
	var sl []int = []int{1, 2, 3, 4, 5, 6}
	f1 := &sl[0]
	f2 := &sl[5]
	*f1, *f2 = *f2, *f1
	fmt.Println(sl)
}
func main() {
	var sl []int = []int{1, 2, 3, 4, 5, 6}
	mp := make(map[string](*int))
	mp["Первый"] = &sl[1]
	//	mp["Второй"] = &sl[5]
	fmt.Println(sl)
	for key, value := range mp {
		fmt.Println(key, *value)
	}
	f := &sl[5]
	if k, inMap := mp["Первый"]; inMap {
		*k, *f = *f, *k
	}
	fmt.Println(sl)
	for key, value := range mp {
		fmt.Println(key, *value)
	}
}

*/
