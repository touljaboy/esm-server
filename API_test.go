package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

//Testing each endpoint (testing handlers seems to be more logical and convenient way to go)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

// Test adding and deleting an employee in one go, its just logical, since I gotta delete him from the db anyway
func TestCRUDEmployee(t *testing.T) {
	//initialize database connection
	initDB()
	mockResponse := `{
    "rows_affected": 1
}`

	eng := SetUpRouter()

	// POST /employees TEST
	//sample employee to add
	emp := Employee{
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
	eng.POST("/employees", addEmployee)
	req, _ := http.NewRequest("POST", "/employees", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()
	eng.ServeHTTP(w, req)
	//test that status is OK
	assert.Equal(t, http.StatusCreated, w.Code)
	//test that 1 row was affected
	assert.Equal(t, mockResponse, w.Body.String())

	// PUT /employees TEST
	empUpt := Employee{
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
	eng.PUT("/employees/:id", updateEmployee)
	req, _ = http.NewRequest("PUT", "/employees/0", bytes.NewBuffer(jsonData))
	w = httptest.NewRecorder()
	eng.ServeHTTP(w, req)
	//test that status is OK
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, mockResponse, w.Body.String())

	// GET /employees TEST
	eng.GET("/employees", getEmployees)
	req, _ = http.NewRequest("GET", "/employees", nil)
	w = httptest.NewRecorder()
	eng.ServeHTTP(w, req)
	//test that status is OK
	assert.Equal(t, http.StatusOK, w.Code)

	// GET /employees:id TEST.
	eng.GET("/employees/:id", getEmployee)
	req, _ = http.NewRequest("GET", "/employees/0", nil)
	w = httptest.NewRecorder()
	eng.ServeHTTP(w, req)

	//test that status is OK
	assert.Equal(t, http.StatusOK, w.Code)

	// DELETE /employees/:id TEST
	eng.DELETE("/employees/:id", deleteEmployee)
	req, _ = http.NewRequest("DELETE", "/employees/0", nil)
	w = httptest.NewRecorder()
	eng.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, mockResponse, w.Body.String())

}

func TestCRUDFullEmployee(t *testing.T) {
	//initialize the db
	initDB()

	eng := SetUpRouter()

	// GET /fullEmployees TEST
	eng.GET("/fullEmployees", getFullEmployees)
	req, _ := http.NewRequest("GET", "/fullEmployees", nil)
	w := httptest.NewRecorder()
	eng.ServeHTTP(w, req)
	//test that status is OK
	assert.Equal(t, http.StatusOK, w.Code)

	// GET /fullEmployees:id TEST. If db has no values, this will fail
	eng.GET("/fullEmployees/:id", getFullEmployee)
	req, _ = http.NewRequest("GET", "/fullEmployees/1", nil)
	w = httptest.NewRecorder()
	eng.ServeHTTP(w, req)

	//test that status is OK
	assert.Equal(t, http.StatusOK, w.Code)

	//TODO add create, update and delete tests
}

func TestCRUDProject(t *testing.T) {
	//initialize the db
	initDB()

	eng := SetUpRouter()

	// GET /projects TEST
	eng.GET("/projects", getProjects)
	req, _ := http.NewRequest("GET", "/projects", nil)
	w := httptest.NewRecorder()
	eng.ServeHTTP(w, req)
	//test that status is OK
	assert.Equal(t, http.StatusOK, w.Code)

	// GET /projects/:id TEST. If db has no values, this will fail
	eng.GET("/projects/:id", getProject)
	req, _ = http.NewRequest("GET", "/projects/1", nil)
	w = httptest.NewRecorder()
	eng.ServeHTTP(w, req)

	//test that status is OK
	assert.Equal(t, http.StatusOK, w.Code)

	//TODO add create, update and delete tests
}

func TestCRUDClient(t *testing.T) {
	//initialize the db
	initDB()

	eng := SetUpRouter()

	// GET /clients TEST
	eng.GET("/clients", getClients)
	req, _ := http.NewRequest("GET", "/clients", nil)
	w := httptest.NewRecorder()
	eng.ServeHTTP(w, req)
	//test that status is OK
	assert.Equal(t, http.StatusOK, w.Code)

	// GET /clients/:id TEST. If db has no values, this will fail
	eng.GET("/clients/:id", getClient)
	req, _ = http.NewRequest("GET", "/clients/1", nil)
	w = httptest.NewRecorder()
	eng.ServeHTTP(w, req)

	//test that status is OK
	assert.Equal(t, http.StatusOK, w.Code)

	//TODO add create, update and delete tests
}

func TestCRUDSkill(t *testing.T) {
	//initialize the db
	initDB()
	mockResponse := `{
    "rows_affected": 1
}`
	eng := SetUpRouter()

	// ADD /skills TEST
	skill := Skill{
		SkillId:    0,
		SkillClass: "Language",
		Skill:      "Polish",
	}
	jsonData, err := json.Marshal(skill)
	if err != nil {
		t.Error(err)
	}
	eng.POST("/employees", addSkill)
	req, _ := http.NewRequest("POST", "/employees", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()
	eng.ServeHTTP(w, req)
	//test that status is OK
	assert.Equal(t, http.StatusCreated, w.Code)
	//test that 1 row was affected
	assert.Equal(t, mockResponse, w.Body.String())

	// GET /skills TEST
	eng.GET("/skills", getSkills)
	req, _ = http.NewRequest("GET", "/skills", nil)
	w = httptest.NewRecorder()
	eng.ServeHTTP(w, req)
	//test that status is OK
	assert.Equal(t, http.StatusOK, w.Code)

	// GET /skills/:id TEST. If db has no values, this will fail
	eng.GET("/skills/:id", getSkill)
	req, _ = http.NewRequest("GET", "/skills/0", nil)
	w = httptest.NewRecorder()
	eng.ServeHTTP(w, req)

	//test that status is OK
	assert.Equal(t, http.StatusOK, w.Code)

	// DELETE /skills/:id TEST
	eng.DELETE("/skills/:id", deleteSkill)
	req, _ = http.NewRequest("DELETE", "/skills/0", nil)
	w = httptest.NewRecorder()
	eng.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, mockResponse, w.Body.String())

	//TODO add update test
}
