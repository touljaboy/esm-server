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
	EmployeeRole string  `json:"employee_role"`
	Project      Project `json:"project"`
}

type Employee struct {
	EmployeeId int64  `json:"employee_id"`
	Name       string `json:"name"`
	Lastname   string `json:"lastname"`
	FocusArea  string `json:"focus_area"`
	Email      string `json:"email"`
}

type EmployeeFull struct {
	Employee Employee      `json:"employee"`
	Skills   []Skill       `json:"skills"`
	Projects []ProjectFull `json:"projects"`
}

// TODO this code also begins to get pretty repetitive, maybe there is a way to generalize the functions?
// TODO - implement a nice project structure
// TODO need way better error messages to get sent, because this fucking sucks dude, no logs, no anything to debug
// TODO: updating, deleting an Employee
// TODO adding, updating, deleting a Skill
// TODO adding, updating, deleting a Skill to an Employee
// TODO adding, updating, deleting a Client
// TODO adding, updating, deleting a Project
// TODO adding, updating, deleting an Employee to a Project
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

func getFullEmployees(context *gin.Context) {
	fullEmployees, err := sqlGetFullEmployees()
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	context.IndentedJSON(http.StatusOK, fullEmployees)
}

func getSkills(context *gin.Context) {
	skills, err := sqlGetSkills()
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	context.IndentedJSON(http.StatusOK, skills)
}

func addEmployee(context *gin.Context) {
	var emp Employee
	if err := context.BindJSON(&emp); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	id, err := sqlAddEmp(emp)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"isAdded": id})
		return
	}
	context.IndentedJSON(http.StatusCreated, gin.H{"isAdded": id})
}

func sqlAddEmp(emp Employee) (int, error) {
	result, err := db.Exec(
		"INSERT INTO Employees (employee_id, name, lastname, focus_area, email) VALUES (?,?,?,?,?)",
		emp.EmployeeId, emp.Name, emp.Lastname, emp.FocusArea, emp.Email)
	if err != nil {
		return -1, err
	}
	id, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}
	return int(id), nil
}

// We use Skill struct which also contains skill level, as it is usually associated with an Employee.
// In this case however, we only want to see what Skills are available in database, thus skill level is nil
func sqlGetSkills() ([]Skill, error) {
	var skills []Skill

	rows, err := db.Query("SELECT skill_class, skill FROM Skills")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var skill Skill

		if err := rows.Scan(&skill.SkillClass, &skill.Skill); err != nil {
			return nil, err
		}

		skills = append(skills, skill)
	}
	return skills, nil
}

func sqlGetFullEmployees() ([]EmployeeFull, error) {
	var employeesFull []EmployeeFull

	//first, get all the employees
	employees, err := sqlGetAllEmployees()
	if err != nil {
		return nil, fmt.Errorf("sqlGetAllProjects: %v", err)
	}

	//iterate through each employee and find associated projects and skills. Then append employeesFull
	for _, employee := range employees {
		var employeeFull EmployeeFull
		var skills []Skill
		var projects []ProjectFull

		//find associated skills
		rows, err := db.Query("SELECT s.skill_class, s.skill, e.skill_level FROM EmployeeSkills AS e "+
			"INNER JOIN Skills AS s ON e.skill_id = s.skill_id WHERE employee_id = ?", employee.EmployeeId)
		if err != nil {
			return nil, fmt.Errorf("sqlGetFullEmployees: %v", err)
		}
		for rows.Next() {
			var skill Skill
			if err := rows.Scan(&skill.SkillClass, &skill.Skill, &skill.SkillLevel); err != nil {
				return nil, fmt.Errorf("sqlGetFullEmployees: %v", err)
			}
			skills = append(skills, skill)
		}

		//find associate projects
		rows, err = db.Query("SELECT a.*, b.employee_role FROM Projects AS a "+
			"INNER JOIN ProjectDetails as b  ON a.project_id = b.project_id WHERE employee_id = ?", employee.EmployeeId)
		if err != nil {
			return nil, fmt.Errorf("sqlGetFullEmployees: %v", err)
		}
		for rows.Next() {
			var projectFull ProjectFull

			if err := rows.Scan(&projectFull.Project.ProjectId,
				&projectFull.Project.ClientId, &projectFull.Project.FocusArea,
				&projectFull.Project.Description, &projectFull.Project.IsSecret, &projectFull.EmployeeRole); err != nil {
				return nil, fmt.Errorf("sqlGetFullEmployees: %v", err)
			}
			projects = append(projects, projectFull)
		}
		employeeFull.Employee = employee
		employeeFull.Skills = skills
		employeeFull.Projects = projects

		employeesFull = append(employeesFull, employeeFull)
	}

	return employeesFull, nil
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
		if err := rows.Scan(&emp.EmployeeId, &emp.Name, &emp.Lastname, &emp.FocusArea, &emp.Email); err != nil {
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
	router.GET("/fullEmployees", getFullEmployees)
	router.GET("/projects", getProjects)
	router.GET("/clients", getClients)
	router.GET("/skills", getSkills)
	router.POST("/employees", addEmployee)
	router.Run("localhost:9090")

}
