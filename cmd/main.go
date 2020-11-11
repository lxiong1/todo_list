package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/rs/cors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

var database, _ = gorm.Open("mysql", "root:root@/todolist?charset=utf8&parseTime=True&loc=Local")

type Todo struct {
	Id          int `gorm:"primary_key"`
	Description string
	Completed   bool
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

func main() {
	defer database.Close()

	setupDatabase()
	setupServer()
}

func Health(writer http.ResponseWriter, request *http.Request) {
	log.Info("Server health is UP")

	writer.Header().Set("Content-Type", "application/json")
	io.WriteString(writer, `{"health": UP}`)
}

func CreateTodo(writer http.ResponseWriter, request *http.Request) {
	description := request.FormValue("description")
	todo := &Todo{Description: description, Completed: false}
	database.Create(&todo)

	log.WithFields(log.Fields{"description": description}).Info("Created todo in database")

	result := database.Last(todo)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(result.Value)
}

func UpdateTodoCompletion(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, _ := strconv.Atoi(vars["id"])
	todoExists := canFindTodoById(id)

	writer.Header().Set("Content-Type", "application/json")

	if todoExists == false {
		message := "Could not update todo because id given does not exist in database"
		log.WithFields(log.Fields{"Id": id}).Info(message)
		io.WriteString(writer, fmt.Sprintf(`{"updated": false, "message": "%v"}`, message))
	} else {
		completed, _ := strconv.ParseBool(request.FormValue("completed"))
		todo := &Todo{}
		database.First(&todo, id)
		todo.Completed = completed
		database.Save(&todo)

		log.WithFields(log.Fields{"id": todo.Id, "description": todo.Description, "completed": todo.Completed}).Info("Updated Todo")
		io.WriteString(writer, `{"updated": true}`)
	}
}

func DeleteTodo(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, _ := strconv.Atoi(vars["id"])
	todoExists := canFindTodoById(id)

	writer.Header().Set("Content-Type", "application/json")

	if todoExists == false {
		message := "Could not delete todo because id given does not exist in database"
		log.WithFields(log.Fields{"Id": id}).Info(message)
		io.WriteString(writer, fmt.Sprintf(`{"deleted": false, "message": "%v"}`, message))
	} else {
		todo := &Todo{}
		database.First(&todo, id)
		database.Delete(&todo)

		log.WithFields(log.Fields{"id": todo.Id, "description": todo.Description, "completed": todo.Completed}).Info("Deleted todo in database")
		io.WriteString(writer, `{"deleted": true}`)
	}
}

func GetTodos(writer http.ResponseWriter, request *http.Request) {
	var todoList []Todo
	todos := database.Find(&todoList).Value
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(todos)

	log.Info("Retrieved todos")
}

func GetCompletedTodos(writer http.ResponseWriter, request *http.Request) {
	completedTodos := findTodosByCompletion(true)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(completedTodos)

	log.Info("Retrieved completed todos")
}

func GetIncompletedTodos(writer http.ResponseWriter, request *http.Request) {
	incompleteTodos := findTodosByCompletion(false)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(incompleteTodos)

	log.Info("Retrieved incompleted todos")
}

func findTodosByCompletion(completed bool) interface{} {
	var todoList []Todo
	todos := database.Where("completed = ?", completed).Find(&todoList).Value

	return todos
}

func canFindTodoById(id int) bool {
	todo := &Todo{}
	result := database.First(&todo, id)

	if result.Error != nil {
		log.WithFields(log.Fields{"Id": id}).Info("Todo not found in database")

		return false
	}

	log.WithFields(log.Fields{"Id": id}).Info("Todo found in database")

	return true
}

func setupDatabase() {
	database.Debug().DropTableIfExists(&Todo{})
	database.Debug().AutoMigrate(&Todo{})
}

func setupServer() {
	log.Info("Starting up server")

	router := mux.NewRouter()
	router.HandleFunc("/health", Health).Methods("GET")
	router.HandleFunc("/todos", GetTodos).Methods("GET")
	router.HandleFunc("/todos/complete", GetCompletedTodos).Methods("GET")
	router.HandleFunc("/todos/incomplete", GetIncompletedTodos).Methods("GET")
	router.HandleFunc("/todo", CreateTodo).Methods("POST")
	router.HandleFunc("/todo/{id}", UpdateTodoCompletion).Methods("POST")
	router.HandleFunc("/todo/{id}", DeleteTodo).Methods("DELETE")

	handler := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "DELETE"},
	}).Handler(router)

	http.ListenAndServe(":8000", handler)
}
