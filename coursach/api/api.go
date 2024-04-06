package api

import (
	"coursach/db"
	"coursach/models"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// Return one object from table by id. Example url: http://localhost:8080/getObject/?id=N&objectType=Name
func getObject(w http.ResponseWriter, r *http.Request) {
	id_param := r.URL.Query().Get("id")
	id, err := strconv.Atoi(id_param)
	if err != nil {
		jsonData, _ := json.Marshal(errors.New("параметр id должен быть не отрицательным и целым").Error())
		w.Write(jsonData)
		return
	}
	objectType := r.URL.Query().Get("objectType")
	object, err := models.FactoryObjByNameAndId(objectType, id)
	if err != nil {
		jsonData, _ := json.Marshal(err.Error())
		w.Write(jsonData)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(object)
	if err != nil {
		panic(err)
	}
}

// Return all objects from table. Example url: http://localhost:8080/getObjectsList/?objectType=Name
func getObjectsList(w http.ResponseWriter, r *http.Request) {
	objectType := r.URL.Query().Get("objectType")
	objects, err := models.FactoryObjList(objectType)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(objects)
	if err != nil {
		panic(err)
	}
}

// Save if without id else update object to DB. Example url: http://localhost:8080/saveObject/?objectType=Name
func saveObject(w http.ResponseWriter, r *http.Request) {
	objectType := r.URL.Query().Get("objectType")
	object, err := saveObjToDB(json.NewDecoder(r.Body), objectType)
	if err != nil {
		jsonData, _ := json.Marshal(err.Error())
		w.Write(jsonData)
		return
	}
	err = json.NewEncoder(w).Encode(object)
	if err != nil {
		panic(errors.New("api/api.go: saveObject: ошибка конвертации в json файл"))
	}
}

// Delete object from DB by id. Example url: http://localhost:8080/deleteObject/?id=N&objectType=Name
func deleteObject(w http.ResponseWriter, r *http.Request) {
	id_param := r.URL.Query().Get("id")
	id, err := strconv.Atoi(id_param)
	if err != nil {
		jsonData, _ := json.Marshal(errors.New("параметр id должен быть не отрицательным числом"))
		w.Write(jsonData)
		return
	}
	objectType := r.URL.Query().Get("objectType")
	object, err := deleteObjFromDB(objectType, id)
	if err != nil {
		jsonData, _ := json.Marshal(err)
		w.Write(jsonData)
		return
	}
	err = json.NewEncoder(w).Encode(object)
	if err != nil {
		panic(errors.New("api/api.go: deleteObject: ошибка конвертации в json файл"))
	}
}

// Start api.
func Start() {
	db.Connect("postgres", "12345678", "sessiontime")
	defer db.Close()

	http.HandleFunc("/getObject/", getObject)
	http.HandleFunc("/getObjectsList/", getObjectsList)
	http.HandleFunc("/saveObject/", saveObject)
	http.HandleFunc("/deleteObject/", deleteObject)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		panic(errors.New("api/api.go: start:ошибка запуска http"))
	}
}
