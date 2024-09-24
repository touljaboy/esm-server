package main

import (
	"bytes"
	"encoding/json"
	"esmAPI/pkg/instances"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

//Testing each endpoint (testing handlers seems to be more logical and convenient way to go)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

// Test adding and deleting an employee in one go, its just logical, since I gotta delete him from the db anyway
func TestCRUDEmployee(t *testing.T) {
	//initialize the empHandler
	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "esmdb",
		AllowNativePasswords: true,
	}
	empStore, err := NewEmployeeStore(cfg)
	if err != nil {
		log.Fatal(err)
	}
	empHandler := NewEmployeeHandler(empStore)

	mockResponse := `{
    "rows_affected": 1
}`

	eng := SetUpRouter()

	// POST /employees TEST
	//sample employee to add
	emp := instances.Employee{
		EmployeeId: 0,
		Name:       "Marco",
		Lastname:   "Plathweith",
		FocusArea:  "Singer",
		Email:      "marco.plathweith@company.co",
	}
	jsonData, err := json.Marshal(emp)
	if err != nil {
		t.Error(err)
	}
	eng.POST("/employees", empHandler.addEmployee)
	req, _ := http.NewRequest("POST", "/employees", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()
	eng.ServeHTTP(w, req)
	//test that status is OK
	assert.Equal(t, http.StatusCreated, w.Code)
	//test that 1 row was affected
	assert.Equal(t, mockResponse, w.Body.String())

	// PUT /employees TEST
	empUpt := instances.Employee{
		EmployeeId: 0,
		Name:       "Pavel",
		Lastname:   "Markovitz",
		FocusArea:  "Dancer",
		Email:      "Pavel.Markovitz@company.co",
	}
	jsonData, err = json.Marshal(empUpt)
	if err != nil {
		t.Error(err)
	}
	eng.PUT("/employees/:id", empHandler.updateEmployee)
	req, _ = http.NewRequest("PUT", "/employees/0", bytes.NewBuffer(jsonData))
	w = httptest.NewRecorder()
	eng.ServeHTTP(w, req)
	//test that status is OK
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, mockResponse, w.Body.String())

	// GET /employees TEST
	eng.GET("/employees", empHandler.getEmployees)
	req, _ = http.NewRequest("GET", "/employees", nil)
	w = httptest.NewRecorder()
	eng.ServeHTTP(w, req)
	//test that status is OK
	assert.Equal(t, http.StatusOK, w.Code)

	// GET /employees:id TEST.
	eng.GET("/employees/:id", empHandler.getEmployee)
	req, _ = http.NewRequest("GET", "/employees/0", nil)
	w = httptest.NewRecorder()
	eng.ServeHTTP(w, req)

	//test that status is OK
	assert.Equal(t, http.StatusOK, w.Code)

	// DELETE /employees/:id TEST
	eng.DELETE("/employees/:id", empHandler.deleteEmployee)
	req, _ = http.NewRequest("DELETE", "/employees/0", nil)
	w = httptest.NewRecorder()
	eng.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, mockResponse, w.Body.String())

}

func TestCRUDFullEmployee(t *testing.T) {

	//initialize the empFullHandler
	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "esmdb",
		AllowNativePasswords: true,
	}
	empStore, err := NewEmployeeStore(cfg)
	if err != nil {
		log.Fatal(err)
	}
	empFullStore, err := NewEmployeeFullStore(cfg)
	if err != nil {
		log.Fatal(err)
	}
	empFullHandler := NewEmployeeFullHandler(empFullStore, empStore)

	eng := SetUpRouter()

	// GET /fullEmployees TEST
	eng.GET("/fullEmployees", empFullHandler.getFullEmployees)
	req, _ := http.NewRequest("GET", "/fullEmployees", nil)
	w := httptest.NewRecorder()
	eng.ServeHTTP(w, req)
	//test that status is OK
	assert.Equal(t, http.StatusOK, w.Code)

	// GET /fullEmployees:id TEST. If db has no values, this will fail
	eng.GET("/fullEmployees/:id", empFullHandler.getFullEmployee)
	req, _ = http.NewRequest("GET", "/fullEmployees/1", nil)
	w = httptest.NewRecorder()
	eng.ServeHTTP(w, req)

	//test that status is OK
	assert.Equal(t, http.StatusOK, w.Code)

	//TODO add create, update and delete tests
}

func TestCRUDProject(t *testing.T) {
	//initialize the projHandler
	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "esmdb",
		AllowNativePasswords: true,
	}
	projStore, err := NewProjectStore(cfg)
	if err != nil {
		log.Fatal(err)
	}
	projHandler := NewProjectHandler(projStore)

	eng := SetUpRouter()

	// GET /projects TEST
	eng.GET("/projects", projHandler.getProjects)
	req, _ := http.NewRequest("GET", "/projects", nil)
	w := httptest.NewRecorder()
	eng.ServeHTTP(w, req)
	//test that status is OK
	assert.Equal(t, http.StatusOK, w.Code)

	// GET /projects/:id TEST. If db has no values, this will fail
	eng.GET("/projects/:id", projHandler.getProject)
	req, _ = http.NewRequest("GET", "/projects/1", nil)
	w = httptest.NewRecorder()
	eng.ServeHTTP(w, req)

	//test that status is OK
	assert.Equal(t, http.StatusOK, w.Code)

	//TODO add create, update and delete tests
}

func TestCRUDClient(t *testing.T) {
	//initialize the clientHandler
	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "esmdb",
		AllowNativePasswords: true,
	}
	clientStore, err := NewClientStore(cfg)
	if err != nil {
		log.Fatal(err)
	}
	clientHandler := NewClientHandler(clientStore)
	eng := SetUpRouter()

	// GET /clients TEST
	eng.GET("/clients", clientHandler.getClients)
	req, _ := http.NewRequest("GET", "/clients", nil)
	w := httptest.NewRecorder()
	eng.ServeHTTP(w, req)
	//test that status is OK
	assert.Equal(t, http.StatusOK, w.Code)

	// GET /clients/:id TEST. If db has no values, this will fail
	eng.GET("/clients/:id", clientHandler.getClient)
	req, _ = http.NewRequest("GET", "/clients/1", nil)
	w = httptest.NewRecorder()
	eng.ServeHTTP(w, req)

	//test that status is OK
	assert.Equal(t, http.StatusOK, w.Code)

	//TODO add create, update and delete tests
}

func TestCRUDSkill(t *testing.T) {
	//initialize the skillHandler
	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "esmdb",
		AllowNativePasswords: true,
	}
	skillStore, err := NewSkillStore(cfg)
	if err != nil {
		log.Fatal(err)
	}
	skillHandler := NewSkillHandler(skillStore)
	mockResponse := `{
    "rows_affected": 1
}`
	eng := SetUpRouter()

	// ADD /skills TEST
	skill := instances.Skill{
		SkillId:    0,
		SkillClass: "Language",
		Skill:      "Polish",
	}
	jsonData, err := json.Marshal(skill)
	if err != nil {
		t.Error(err)
	}
	eng.POST("/employees", skillHandler.addSkill)
	req, _ := http.NewRequest("POST", "/employees", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()
	eng.ServeHTTP(w, req)
	//test that status is OK
	assert.Equal(t, http.StatusCreated, w.Code)
	//test that 1 row was affected
	assert.Equal(t, mockResponse, w.Body.String())

	// GET /skills TEST
	eng.GET("/skills", skillHandler.getSkills)
	req, _ = http.NewRequest("GET", "/skills", nil)
	w = httptest.NewRecorder()
	eng.ServeHTTP(w, req)
	//test that status is OK
	assert.Equal(t, http.StatusOK, w.Code)

	// GET /skills/:id TEST. If db has no values, this will fail
	eng.GET("/skills/:id", skillHandler.getSkill)
	req, _ = http.NewRequest("GET", "/skills/0", nil)
	w = httptest.NewRecorder()
	eng.ServeHTTP(w, req)

	//test that status is OK
	assert.Equal(t, http.StatusOK, w.Code)

	// PUT /skills/:id TEST
	eng.PUT("/skills/:id", skillHandler.updateSkill)
	skill.Skill = "Mandarin"
	jsonData, err = json.Marshal(skill)
	if err != nil {
		t.Error(err)
	}
	req, err = http.NewRequest("PUT", "/skills/0", bytes.NewBuffer(jsonData))
	w = httptest.NewRecorder()
	eng.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, mockResponse, w.Body.String())

	// DELETE /skills/:id TEST
	eng.DELETE("/skills/:id", skillHandler.deleteSkill)
	req, _ = http.NewRequest("DELETE", "/skills/0", nil)
	w = httptest.NewRecorder()
	eng.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, mockResponse, w.Body.String())

}
