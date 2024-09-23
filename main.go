package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"strconv"
)

//Define structs to be used for representing the db data
//For now, I assume that struct EmployeeFull will be the "highest in hierarchy", combining all data

type Skill struct {
	SkillId    int    `json:"skill_id"`
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

// TODO ids should probably be a uint
// TODO version your API, for now it's v1
// TODO write some tests
// TODO this code also begins to get pretty repetitive, maybe there is a way to generalize the functions?
// TODO - implement a nice project structure
// TODO need way better error messages to get sent, because this fucking sucks dude, no logs, no anything to debug
// TODO adding, updating, deleting a Skill
// TODO adding, updating, deleting a Client
// TODO adding, updating, deleting a Project
// TODO adding, updating, deleting a Skill to an Employee
// TODO adding, updating, deleting an Employee to a Project

func getEmployees(context *gin.Context) {
	employees, err := sqlGetAllEmployees()
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.IndentedJSON(http.StatusOK, employees)
}

func getEmployee(context *gin.Context) {
	strId := context.Params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	employee, err := sqlGetEmployeeById(id)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	context.IndentedJSON(http.StatusOK, employee)
}

func getProjects(context *gin.Context) {
	projects, err := sqlGetAllProjects()
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	context.IndentedJSON(http.StatusOK, projects)
}

func getProject(context *gin.Context) {
	strId := context.Params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	project, err := sqlGetProject(id)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	context.IndentedJSON(http.StatusOK, project)
}

func getClients(context *gin.Context) {
	clients, err := sqlGetAllClients()
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	context.IndentedJSON(http.StatusOK, clients)
}

func getClient(context *gin.Context) {
	strId := context.Params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	client, err := sqlGetClient(id)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	context.IndentedJSON(http.StatusOK, client)
}

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

func getSkills(context *gin.Context) {
	skills, err := sqlGetSkills()
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	context.IndentedJSON(http.StatusOK, skills)
}

func getSkill(context *gin.Context) {
	strId := context.Params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	skill, err := sqlGetSkill(id)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	context.IndentedJSON(http.StatusOK, skill)
}

func addEmployee(context *gin.Context) {
	var emp Employee
	if err := context.BindJSON(&emp); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	result, err := sqlAddEmp(emp)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	context.IndentedJSON(http.StatusCreated, gin.H{"rows_affected": result})
}

func addSkill(context *gin.Context) {
	var skill Skill
	if err := context.BindJSON(&skill); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	result, err := sqlAddSkill(skill)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	context.IndentedJSON(http.StatusCreated, gin.H{"rows_affected": result})
}

// full entry update done by id
func updateEmployee(context *gin.Context) {
	strId := context.Params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	currEmployee, err := sqlGetEmployeeById(id)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	if err := context.BindJSON(&currEmployee); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	result, err := sqlUpdateEmployee(currEmployee)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{"rows_affected": result})
}

func updateSkill(context *gin.Context) {
	strId := context.Params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	currEmployee, err := sqlGetEmployeeById(id)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	if err := context.BindJSON(&currEmployee); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	result, err := sqlUpdateEmployee(currEmployee)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{"rows_affected": result})
}

func deleteEmployee(context *gin.Context) {
	id := context.Params.ByName("id")
	result, err := sqlDeleteEmployee(id)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"err": err})
	}
	context.IndentedJSON(http.StatusOK, gin.H{"rows_affected": result})
}

func deleteSkill(context *gin.Context) {
	id := context.Params.ByName("id")
	result, err := sqlDeleteSkill(id)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"err": err})
	}
	context.IndentedJSON(http.StatusOK, gin.H{"rows_affected": result})
}

func sqlDeleteSkill(strId string) (int64, error) {
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		return -1, err
	}
	result, err := db.Exec("DELETE FROM Skills WHERE skill_id=?", id)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func sqlDeleteEmployee(strId string) (int64, error) {
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		return -1, err
	}
	result, err := db.Exec("DELETE FROM Employees WHERE employee_id=?", id)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func sqlUpdateEmployee(emp Employee) (int64, error) {
	result, err := db.Exec(
		"UPDATE Employees SET name=?, lastname=?, focus_area=?, email=? WHERE employee_id = ?",
		emp.Name, emp.Lastname, emp.FocusArea, emp.Email, emp.EmployeeId)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
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

func sqlAddSkill(skill Skill) (int, error) {
	result, err := db.Exec(
		"INSERT INTO Skills (skill_id, skill_class, skill) VALUES (?,?,?)",
		skill.SkillId, skill.SkillClass, skill.Skill)
	if err != nil {
		return -1, err
	}
	id, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}
	return int(id), nil
}

func sqlGetSkill(id int64) (Skill, error) {
	var skill Skill
	row := db.QueryRow("SELECT * FROM Skills WHERE skill_id=?", id)
	if err := row.Scan(&skill.SkillId, &skill.SkillClass, &skill.Skill); err != nil {
		return Skill{}, err
	}
	return skill, nil
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
		employeeFull, err := sqlGetFullEmployeeById(employee.EmployeeId)
		if err != nil {
			return nil, fmt.Errorf("sqlGetFullEmployeeById: %v", err)
		}
		employeesFull = append(employeesFull, employeeFull)
	}

	return employeesFull, nil
}

func sqlGetFullEmployeeById(id int64) (EmployeeFull, error) {
	employee, err := sqlGetEmployeeById(id)
	if err != nil {
		return EmployeeFull{}, err
	}

	var employeeFull EmployeeFull
	var skills []Skill
	var projects []ProjectFull

	//find associated skills
	rows, err := db.Query("SELECT s.skill_class, s.skill, e.skill_level FROM EmployeeSkills AS e "+
		"INNER JOIN Skills AS s ON e.skill_id = s.skill_id WHERE employee_id = ?", employee.EmployeeId)
	if err != nil {
		return EmployeeFull{}, fmt.Errorf("sqlGetFullEmployees: %v", err)
	}
	for rows.Next() {
		var skill Skill
		if err := rows.Scan(&skill.SkillClass, &skill.Skill, &skill.SkillLevel); err != nil {
			return EmployeeFull{}, fmt.Errorf("sqlGetFullEmployees: %v", err)
		}
		skills = append(skills, skill)
	}

	//find associate projects
	rows, err = db.Query("SELECT a.*, b.employee_role FROM Projects AS a "+
		"INNER JOIN ProjectDetails as b  ON a.project_id = b.project_id WHERE employee_id = ?", employee.EmployeeId)
	if err != nil {
		return EmployeeFull{}, fmt.Errorf("sqlGetFullEmployees: %v", err)
	}
	for rows.Next() {
		var projectFull ProjectFull

		if err := rows.Scan(&projectFull.Project.ProjectId,
			&projectFull.Project.ClientId, &projectFull.Project.FocusArea,
			&projectFull.Project.Description, &projectFull.Project.IsSecret, &projectFull.EmployeeRole); err != nil {
			return EmployeeFull{}, fmt.Errorf("sqlGetFullEmployees: %v", err)
		}
		projects = append(projects, projectFull)
	}
	employeeFull.Employee = employee
	employeeFull.Skills = skills
	employeeFull.Projects = projects

	return employeeFull, nil
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

func sqlGetProject(id int64) (Project, error) {
	var proj Project

	row := db.QueryRow("SELECT * FROM Projects WHERE project_id = ?", id)
	if err := row.Scan(&proj.ProjectId, &proj.ClientId, &proj.FocusArea, &proj.Description, &proj.IsSecret); err != nil {
		return Project{}, err
	}
	return proj, nil
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

func sqlGetClient(id int64) (Client, error) {
	var client Client
	row := db.QueryRow("SELECT * FROM Clients WHERE id = ?", id)
	if err := row.Scan(&client.ID, &client.Name, &client.Description); err != nil {
		return Client{}, err
	}
	return client, nil
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

func sqlGetEmployeeById(id int64) (Employee, error) {
	var emp Employee

	row := db.QueryRow("SELECT * FROM Employees WHERE employee_id = ?", id)
	if err := row.Scan(&emp.EmployeeId, &emp.Name, &emp.Lastname, &emp.FocusArea, &emp.Email); err != nil {
		return Employee{}, err
	}
	return emp, nil
}

func initDB() {
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
}

var db *sql.DB

func main() {
	initDB()

	//Configure endpoints
	router := gin.Default()
	router.GET("/v1/employees", getEmployees)
	router.GET("/v1/employees/:id", getEmployee)
	router.POST("/v1/employees", addEmployee)
	router.PUT("/v1/employees/:id", updateEmployee)
	router.DELETE("/v1/employees/:id", deleteEmployee)

	router.GET("/v1/fullEmployees", getFullEmployees)
	router.GET("/v1/fullEmployees/:id", getFullEmployee)

	router.GET("/v1/projects", getProjects)
	router.GET("/v1/projects/:id", getProject)

	router.GET("/v1/clients", getClients)
	router.GET("/v1/clients/:id", getClient)

	router.GET("/v1/skills", getSkills)
	router.GET("/v1/skills/:id", getSkill)
	router.POST("/v1/skills", addSkill)
	router.DELETE("/v1/skills/:id", deleteSkill)

	router.Run("localhost:9090")

}
