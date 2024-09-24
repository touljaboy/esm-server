package main

import (
	"esmAPI/pkg/instances"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"strconv"
)

// TODO introduce handlers for each struct type of the REST API
// TODO introduce a logger
// TODO tests can be written in .http format
// TODO also, consider using a router gorilla/mux
// TODO ids should probably be a uint
// TODO need way better error messages to get sent, because this fucking sucks dude, no logs, no anything to debug
// TODO adding, updating, deleting a Client
// TODO adding, updating, deleting a Skill to an Employee
// TODO adding, updating, deleting an Employee to a Project

func getFullEmployees(context *gin.Context) {
	fullEmployees, err := sqlGetFullEmployees()
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	context.IndentedJSON(http.StatusOK, fullEmployees)
}

func getFullEmployee(context *gin.Context) {
	strId := context.Params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	fullEmployee, err := sqlGetFullEmployeeById(id)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	context.IndentedJSON(http.StatusOK, fullEmployee)
}
func sqlGetFullEmployees() ([]instances.EmployeeFull, error) {
	var employeesFull []instances.EmployeeFull

	//first, get all the employees
	employees, err := sqlGetAllEmployees()
	if err != nil {
		return nil, fmt.Errorf("sqlGetAllProjects: %v", err)
	}

	//iterate through each employee and find associated projects and skills. Then append employeesFull
	for _, employee := range employees {
		employeeFull, err := sqlGetFullEmployeeById(employee.EmployeeId)
		if err != nil {
			return nil, fmt.Errorf("sqlGetFullEmployeeById: %v", err)
		}
		employeesFull = append(employeesFull, employeeFull)
	}

	return employeesFull, nil
}

func sqlGetFullEmployeeById(id int64) (instances.EmployeeFull, error) {
	employee, err := sqlGetEmployeeById(id)
	if err != nil {
		return instances.EmployeeFull{}, err
	}

	var employeeFull instances.EmployeeFull
	var skills []instances.Skill
	var projects []instances.ProjectFull

	//find associated skills
	rows, err := db.Query("SELECT s.skill_class, s.skill, e.skill_level FROM EmployeeSkills AS e "+
		"INNER JOIN Skills AS s ON e.skill_id = s.skill_id WHERE employee_id = ?", employee.EmployeeId)
	if err != nil {
		return instances.EmployeeFull{}, fmt.Errorf("sqlGetFullEmployees: %v", err)
	}
	for rows.Next() {
		var skill instances.Skill
		if err := rows.Scan(&skill.SkillClass, &skill.Skill, &skill.SkillLevel); err != nil {
			return instances.EmployeeFull{}, fmt.Errorf("sqlGetFullEmployees: %v", err)
		}
		skills = append(skills, skill)
	}

	//find associate projects
	rows, err = db.Query("SELECT a.*, b.employee_role FROM Projects AS a "+
		"INNER JOIN ProjectDetails as b  ON a.project_id = b.project_id WHERE employee_id = ?", employee.EmployeeId)
	if err != nil {
		return instances.EmployeeFull{}, fmt.Errorf("sqlGetFullEmployees: %v", err)
	}
	for rows.Next() {
		var projectFull instances.ProjectFull

		if err := rows.Scan(&projectFull.Project.ProjectId,
			&projectFull.Project.ClientId, &projectFull.Project.FocusArea,
			&projectFull.Project.Description, &projectFull.Project.IsSecret, &projectFull.EmployeeRole); err != nil {
			return instances.EmployeeFull{}, fmt.Errorf("sqlGetFullEmployees: %v", err)
		}
		projects = append(projects, projectFull)
	}
	employeeFull.Employee = employee
	employeeFull.Skills = skills
	employeeFull.Projects = projects

	return employeeFull, nil
}

// TODO the following code seems very similar. Try and find a way to generalize such code

func main() {
	// Capture connection properties.
	// TODO read cfg from a separate file in gitignore
	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "esmdb",
		AllowNativePasswords: true,
	}

	// create stores
	empStore, err := NewMySQLEmployeeStore(cfg)
	if err != nil {
		log.Fatal(err)
	}
	skillStore, err := NewMySQLSkillStore(cfg)
	if err != nil {
		log.Fatal(err)
	}
	projectStore, err := NewProjectStore(cfg)
	// create handlers
	empHandler := NewEmployeeHandler(empStore)
	skillHandler := NewSkillHandler(skillStore)
	projectHandler := NewProjectHandler(projectStore)
	//Configure endpoints
	router := gin.Default()
	router.Routes()
	router.GET("/v1/employees", empHandler.getEmployees)
	router.GET("/v1/employees/:id", empHandler.getEmployee)
	router.POST("/v1/employees", empHandler.addEmployee)
	router.PUT("/v1/employees/:id", empHandler.updateEmployee)
	router.DELETE("/v1/employees/:id", empHandler.deleteEmployee)

	router.GET("/v1/fullEmployees", getFullEmployees)
	router.GET("/v1/fullEmployees/:id", getFullEmployee)

	router.GET("/v1/projects", projectHandler.getProjects)
	router.GET("/v1/projects/:id", projectHandler.getProject)

	router.GET("/v1/clients", getClients)
	router.GET("/v1/clients/:id", getClient)

	router.GET("/v1/skills", skillHandler.getSkills)
	router.GET("/v1/skills/:id", skillHandler.getSkill)
	router.POST("/v1/skills", skillHandler.addSkill)
	router.PUT("/v1/skills/:id", skillHandler.updateSkill)
	router.DELETE("/v1/skills/:id", skillHandler.deleteSkill)

	router.Run("localhost:9090")

}
