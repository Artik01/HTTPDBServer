package main

import (
	"net/http"
	"io"
	"fmt"
	"encoding/json"
	"encoding/xml"
	"strings"
)

var db DB

var GlobalMutex chan int = make(chan int, 1)

type (
	Person struct {
		Name         string `json:"name" xml:"name"`
		Surname      string `json:"surname" xml:"surname"`
		PersonalCode string `json:"personalCode" xml:"personalCode"`
	}

	Teacher struct {
		mutex	  chan int
		ID        string   `json:"id" xml:"id"`
		Subject   string   `json:"subject" xml:"subject"`
		Salary    float64  `json:"salary" xml:"salary"`
		Classroom []string `json:"classroom" xml:"classroom>value"`
		Person    `json:"person"`
	}

	Student struct {
		mutex  chan int
		ID     string `json:"id" xml:"id"`
		Class  string `json:"class" xml:"class"`
		Person `json:"person"`
	}

	Staff struct {
		mutex	  chan int
		ID        string  `json:"id" xml:"id"`
		Salary    float64 `json:"salary" xml:"salary"`
		Classroom string  `json:"classroom" xml:"classroom"`
		Phone     string  `json:"phone" xml:"phone"`
		Person    `json:"person"`
	}
	DB []GeneralObject
)

var FirstFreeId int = 1

type Action struct {
	Action  string `json:"action" xml:"action"`
	ObjName string `json:"object" xml:"object"`
}
type DefinedAction interface {
	GetFromJSON([]byte)
	GetFromXML([]byte)
	Process(db *DB) string
}
type GeneralObject interface {
	GetCreateAction() DefinedAction
	GetUpdateAction() DefinedAction
	GetReadAction() DefinedAction
	GetDeleteAction() DefinedAction
	Print()
	GetId() string
}

func (db DB) GetIndex(id string) int {
	for i, p := range db {
		if p.GetId() == id {
			return i
		}
	}
	return -1
}

//Teacher:
func (t Teacher) GetCreateAction() DefinedAction {
	return &CreateTeacher{}
}

type CreateTeacher struct {
	T Teacher `json:"data" xml:"data"`
}

func (action *CreateTeacher) GetFromJSON(rawData []byte) {
	err := json.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action *CreateTeacher) GetFromXML(rawData []byte) {
	err := xml.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action CreateTeacher) Process(db *DB) string {
	<-GlobalMutex
	action.T.ID=fmt.Sprint(FirstFreeId)
	action.T.mutex = make(chan int, 1)
	action.T.mutex <- 1
	FirstFreeId++
	*db = append(*db, action.T)
	GlobalMutex <- 1
	return "Teacher created successfully\n"
}

func (t Teacher) GetUpdateAction() DefinedAction {
	return &UpdateTeacher{}
}

type UpdateTeacher struct {
	T Teacher `json:"data" xml:"data"`
}

func (action *UpdateTeacher) GetFromJSON(rawData []byte) {
	err := json.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action *UpdateTeacher) GetFromXML(rawData []byte) {
	err := xml.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action UpdateTeacher) Process(db *DB) string {
	<-GlobalMutex
	GlobalMutex <- 1
	id := action.T.GetId()
	<-((*db)[db.GetIndex(id)]).(Teacher).mutex
	action.T.mutex = make(chan int, 1)
	(*db)[db.GetIndex(id)] = action.T
	((*db)[db.GetIndex(id)]).(Teacher).mutex <- 1
	return "Teacher updated successfully\n"
}

func (t Teacher) GetReadAction() DefinedAction {
	return &ReadTeacher{}
}

type ReadTeacher struct {
	Data struct {
		ID string `json:"id" xml:"id"`
	} `json:"data" xml:"data"`
}

func (action *ReadTeacher) GetFromJSON(rawData []byte) {
	err := json.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action *ReadTeacher) GetFromXML(rawData []byte) {
	err := xml.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action ReadTeacher) Process(db *DB) string {
	<-GlobalMutex
	GlobalMutex <- 1
	<-((*db)[db.GetIndex(action.Data.ID)]).(Teacher).mutex
	t := ((*db)[db.GetIndex(action.Data.ID)]).(Teacher)
	((*db)[db.GetIndex(action.Data.ID)]).(Teacher).mutex <- 1
	return fmt.Sprintf("ID:%s\tName:%s\tSurname:%s\tSalary:%.2f\tSubject:%s\tClassroom:%v\n", t.ID, t.Name, t.Surname, t.Salary, t.Subject, t.Classroom)
}

func (t Teacher) GetDeleteAction() DefinedAction {
	return &DeleteTeacher{}
}

type DeleteTeacher struct {
	Data struct {
		ID string `json:"id" xml:"id"`
	} `json:"data" xml:"data"`
}

func (action *DeleteTeacher) GetFromJSON(rawData []byte) {
	err := json.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action *DeleteTeacher) GetFromXML(rawData []byte) {
	err := xml.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action DeleteTeacher) Process(db *DB) string {
	<-GlobalMutex
	for i, p := range *db {
		if p.GetId() == action.Data.ID {
			*db = append((*db)[:i], (*db)[i+1:]...)
		}
	}
	GlobalMutex <- 1
	return "Teacher deleted successfully\n"
}
func (t Teacher) Print() {
	fmt.Printf("ID:%s\tName:%s\tSurname:%s\tSalary:%.2f\tSubject:%s\tClassroom:%v\n", t.ID, t.Name, t.Surname, t.Salary, t.Subject, t.Classroom)
}

func (t Teacher) GetId() string {
	return t.ID
}

//Student:
func (s Student) GetCreateAction() DefinedAction {
	return &CreateStudent{}
}

type CreateStudent struct {
	S Student `json:"data" xml:"data"`
}

func (action *CreateStudent) GetFromJSON(rawData []byte) {
	err := json.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action *CreateStudent) GetFromXML(rawData []byte) {
	err := xml.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action CreateStudent) Process(db *DB) string {
	<-GlobalMutex
	action.S.ID=fmt.Sprint(FirstFreeId)
	action.S.mutex = make(chan int, 1)
	action.S.mutex <- 1
	FirstFreeId++
	*db = append(*db, action.S)
	GlobalMutex <- 1
	return "Student created successfully\n"
}

func (s Student) GetUpdateAction() DefinedAction {
	return &UpdateStudent{}
}

type UpdateStudent struct {
	S Student `json:"data" xml:"data"`
}

func (action *UpdateStudent) GetFromJSON(rawData []byte) {
	err := json.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action *UpdateStudent) GetFromXML(rawData []byte) {
	err := xml.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action UpdateStudent) Process(db *DB) string {
	<-GlobalMutex
	GlobalMutex <- 1
	id := action.S.GetId()
	<-((*db)[db.GetIndex(id)]).(Student).mutex
	action.S.mutex = make(chan int, 1)
	(*db)[db.GetIndex(id)] = action.S
	((*db)[db.GetIndex(id)]).(Student).mutex <- 1
	return "Student updated successfully\n"
}

func (s Student) GetReadAction() DefinedAction {
	return &ReadStudent{}
}

type ReadStudent struct {
	Data struct {
		ID string `json:"id" xml:"id"`
	} `json:"data" xml:"data"`
}

func (action *ReadStudent) GetFromJSON(rawData []byte) {
	err := json.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action *ReadStudent) GetFromXML(rawData []byte) {
	err := xml.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action ReadStudent) Process(db *DB) string {
	<-GlobalMutex
	GlobalMutex <- 1
	<-((*db)[db.GetIndex(action.Data.ID)]).(Student).mutex
	s := ((*db)[db.GetIndex(action.Data.ID)]).(Student)
	((*db)[db.GetIndex(action.Data.ID)]).(Student).mutex <- 1
	return fmt.Sprintf("ID:%s\tName:%s\tSurname:%s\tClass:%s\n", s.ID, s.Name, s.Surname, s.Class)
}

func (s Student) GetDeleteAction() DefinedAction {
	return &DeleteStudent{}
}

type DeleteStudent struct {
	Data struct {
		ID string `json:"id" xml:"id"`
	} `json:"data" xml:"data"`
}

func (action *DeleteStudent) GetFromJSON(rawData []byte) {
	err := json.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action *DeleteStudent) GetFromXML(rawData []byte) {
	err := xml.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action DeleteStudent) Process(db *DB) string {
	<-GlobalMutex
	for i, p := range *db {
		if p.GetId() == action.Data.ID {
			*db = append((*db)[:i], (*db)[i+1:]...)
		}
	}
	GlobalMutex <- 1
	return "Student deleted successfully\n"
}
func (s Student) Print() {
	fmt.Printf("ID:%s\tName:%s\tSurname:%s\tClass:%s\n", s.ID, s.Name, s.Surname, s.Class)
}

func (s Student) GetId() string {
	return s.ID
}

//Staff:
func (s Staff) GetCreateAction() DefinedAction {
	return &CreateStaff{}
}

type CreateStaff struct {
	S Staff `json:"data" xml:"data"`
}

func (action *CreateStaff) GetFromJSON(rawData []byte) {
	err := json.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action *CreateStaff) GetFromXML(rawData []byte) {
	err := xml.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action CreateStaff) Process(db *DB) string {
	<-GlobalMutex
	action.S.ID=fmt.Sprint(FirstFreeId)
	action.S.mutex = make(chan int, 1)
	action.S.mutex <- 1
	FirstFreeId++
	*db = append(*db, action.S)
	GlobalMutex <- 1
	return "Staff created successfully\n"
}

func (s Staff) GetUpdateAction() DefinedAction {
	return &UpdateStaff{}
}

type UpdateStaff struct {
	S Staff `json:"data" xml:"data"`
}

func (action *UpdateStaff) GetFromJSON(rawData []byte) {
	err := json.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action *UpdateStaff) GetFromXML(rawData []byte) {
	err := xml.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action UpdateStaff) Process(db *DB) string {
	<-GlobalMutex
	GlobalMutex <- 1
	id := action.S.GetId()
	<-((*db)[db.GetIndex(id)]).(Staff).mutex
	action.S.mutex = make(chan int, 1)
	(*db)[db.GetIndex(id)] = action.S
	((*db)[db.GetIndex(id)]).(Staff).mutex <- 1
	return "Staff updated successfully\n"
}

func (s Staff) GetReadAction() DefinedAction {
	return &ReadStaff{}
}

type ReadStaff struct {
	Data struct {
		ID string `json:"id" xml:"id"`
	} `json:"data" xml:"data"`
}

func (action *ReadStaff) GetFromJSON(rawData []byte) {
	err := json.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action *ReadStaff) GetFromXML(rawData []byte) {
	err := xml.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action ReadStaff) Process(db *DB) string {
	<-GlobalMutex
	GlobalMutex <- 1
	<-((*db)[db.GetIndex(action.Data.ID)]).(Staff).mutex
	s := ((*db)[db.GetIndex(action.Data.ID)]).(Staff)
	((*db)[db.GetIndex(action.Data.ID)]).(Staff).mutex <- 1
	return fmt.Sprintf("ID:%s\tName:%s\tSurname:%s\tSalary:%.2f\tClassroom:%s\tPhone:%s\n", s.ID, s.Name, s.Surname, s.Salary, s.Classroom, s.Phone)
}

func (s Staff) GetDeleteAction() DefinedAction {
	return &DeleteStaff{}
}

type DeleteStaff struct {
	Data struct {
		ID string `json:"id" xml:"id"`
	} `json:"data" xml:"data"`
}

func (action *DeleteStaff) GetFromJSON(rawData []byte) {
	err := json.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action *DeleteStaff) GetFromXML(rawData []byte) {
	err := xml.Unmarshal(rawData, action)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func (action DeleteStaff) Process(db *DB) string {
	<-GlobalMutex
	for i, p := range *db {
		if p.GetId() == action.Data.ID {
			*db = append((*db)[:i], (*db)[i+1:]...)
		}
	}
	GlobalMutex <- 1
	return "Staff deleted successfully\n"
}
func (s Staff) Print() {
	fmt.Printf("ID:%s\tName:%s\tSurname:%s\tSalary:%.2f\tClassroom:%s\tPhone:%s\n", s.ID, s.Name, s.Surname, s.Salary, s.Classroom, s.Phone)
}

func (s Staff) GetId() string {
	return s.ID
}

func Handler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		var ids string
		for _, obj := range db {
			ids += obj.GetId() + " "
		}
		io.WriteString(w, ids)
	} else if req.Method == "POST" {
		data, err := io.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {return }
		res := db.UseAction(data)
		io.WriteString(w, res)
	} else if req.Method == "DELETE" {
		id, err := io.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {return }
		<-GlobalMutex
		for i, p := range db {
			if p.GetId() == string(id) {
				db = append(db[:i], db[i+1:]...)
			}
		}
		GlobalMutex <- 1
		io.WriteString(w, "Object deleted successfully\n")
	} else {
		//w.WriteHeader(405)
		io.WriteString(w, "Unknown cmd\n")
	}
	
}

func (db *DB) UseAction(data []byte) string {
	var FType string
	if strings.HasPrefix(string(data), "{") {
		FType="json"
	} else if strings.HasPrefix(string(data), "<") {
		FType="xml"
	} else {
		return "Unsuported file type\n"
	}
	
	var act Action
	var err error
	if FType == "json" {
		err = json.Unmarshal(data, &act)
	} else if FType == "xml" {
		err = xml.Unmarshal(data, &act)
	}
	if err != nil {
		return fmt.Sprintln(err)
	}
	
	var obj GeneralObject
	switch act.ObjName {
		case "Teacher":
			obj = &Teacher{}
		case "Student":
			obj = &Student{}
		case "Staff":
			obj = &Staff{}
		default:
			return fmt.Sprintln("unknown object",act.ObjName)
	}
	var toDo DefinedAction
	
	switch act.Action {
		case "create":
			toDo = obj.GetCreateAction()
		case "update":
			toDo = obj.GetUpdateAction()
		case "read":
			toDo = obj.GetReadAction()
		case "delete":
			toDo = obj.GetDeleteAction()
		default:
			return fmt.Sprintln("unknown action",act.Action)
	}
	
	if FType == "json" {
		toDo.GetFromJSON(data)
	} else if FType == "xml" {
		toDo.GetFromXML(data)
	}
	return toDo.Process(db)
}

func min(a, b int) int {
	if a < b {return a}
	return b
}

func main() {
	GlobalMutex <- 1
	http.HandleFunc("/", Handler)
	
	err := http.ListenAndServe(":8080", nil)
	panic(err)
}

