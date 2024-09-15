package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
)

//Define structs to be used for representing the db data
//For now, I assume that struct EmployeeFull will be the "highest in hierarchy"

type Skill struct {
	SkillClass string `json:"skill_class"`
	Skill      string `json:"skill"`
	SkillLevel int    `json:"skill_level"`
}
type Client struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Project struct {
	ProjectId   int64  `json:"project_id"`
	ClientId    int    `json:"client_id"`
	FocusArea   string `json:"focus_area"`
	Description string `json:"description"`
	IsSecret    bool   `json:"isSecret"`
}

// ProjectFull is used to combine Project with Employee to add an EmployeeRole. This way, the EmployeeFull can have
// all the information combined about an Employee
type ProjectFull struct {
	Project      Project `json:"project"`
	EmployeeRole string  `json:"employee_role"`
}

type Employee struct {
	EmployeeId int64  `json:"employee_id"`
	Name       string `json:"name"`
	Lastname   string `json:"lastname"`
	FocusArea  string `json:"focus_area"`
}

type EmployeeFull struct {
	Employee Employee      `json:"employee"`
	Skills   []Skill       `json:"skills"`
	Projects []ProjectFull `json:"projects"`
}

//TODO better table design idea - create a new table EmployeeSkills in sql. When querying from EmployeeSkills joined with Employee, fill in the entire, singular Employee struct

// TODO this code also begins to get pretty repetitive, maybe there is a way to generalize the functions?
func getEmployees(context *gin.Context) {
	employees, err := sqlGetAllEmployees()
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.IndentedJSON(http.StatusOK, employees)
}

func getProjects(context *gin.Context) {
	projects, err := sqlGetAllProjects()
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	context.IndentedJSON(http.StatusOK, projects)
}

func getClients(context *gin.Context) {
	clients, err := sqlGetAllClients()
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	context.IndentedJSON(http.StatusOK, clients)
}

// TODO the following code seems very similar. Try and find a way to generalize such code
func sqlGetAllProjects() ([]Project, error) {
	var projects []Project

	rows, err := db.Query("SELECT * FROM projects")
	if err != nil {
		return nil, fmt.Errorf("sqlGetAllProjects: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var project Project
		if err := rows.Scan(&project.ProjectId, &project.ClientId, &project.FocusArea, &project.Description, &project.IsSecret); err != nil {
			return nil, fmt.Errorf("sqlGetAllProjects: %v", err)
		}
		projects = append(projects, project)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("sqlGetAllProjects: %v", err)
	}
	return projects, nil
}

func sqlGetAllClients() ([]Client, error) {
	var clients []Client

	rows, err := db.Query("SELECT * FROM Clients")
	if err != nil {
		return nil, fmt.Errorf("sqlGetAllClients: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var client Client
		if err := rows.Scan(&client.ID, &client.Name, &client.Description); err != nil {
			return nil, fmt.Errorf("sqlGetAllClients: %v", err)
		}
		clients = append(clients, client)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("sqlGetAllClients: %v", err)
	}
	return clients, nil
}

func sqlGetAllEmployees() ([]Employee, error) {
	var employees []Employee

	rows, err := db.Query("SELECT * FROM Employees")
	if err != nil {
		return nil, fmt.Errorf("sqlGetAllEmployees %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var emp Employee
		if err := rows.Scan(&emp.EmployeeId, &emp.Name, &emp.Lastname, &emp.FocusArea); err != nil {
			return nil, fmt.Errorf("sqlGetAllEmployees %v", err)
		}
		employees = append(employees, emp)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("sqlGetAllEmployees %v", err)
	}
	return employees, nil
}

var db *sql.DB

func main() {
	// Capture connection properties.
	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "esmdb",
		AllowNativePasswords: true,
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	//Configure endpoints
	router := gin.Default()
	router.GET("/employees", getEmployees)
	router.GET("/projects", getProjects)
	router.GET("/clients", getClients)
	router.Run("localhost:9090")

}
