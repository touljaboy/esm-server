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
// TODO this code also begins to get pretty repetitive, maybe there is a way to generalize the functions?
// TODO - implement a nice project structure
// TODO need way better error messages to get sent, because this fucking sucks dude, no logs, no anything to debug
// TODO adding, updating, deleting a Client
// TODO adding, updating, deleting a Project
// TODO adding, updating, deleting a Skill to an Employee
// TODO adding, updating, deleting an Employee to a Project

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

func addSkill(context *gin.Context) {
	var skill instances.Skill
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

func updateSkill(context *gin.Context) {
	strId := context.Params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	currSkill, err := sqlGetSkill(id)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	if err := context.BindJSON(&currSkill); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	result, err := sqlUpdateSkill(currSkill)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"err": err})
		return
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

func sqlUpdateSkill(skill instances.Skill) (int64, error) {
	result, err := db.Exec(
		"UPDATE Skills SET skill_id=?, skill_class=?, skill=? WHERE skill_id = ?",
		skill.SkillId, skill.SkillClass, skill.Skill, skill.SkillId)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

// We use Skill struct which also contains skill level, as it is usually associated with an Employee.
// In this case however, we only want to see what Skills are available in database, thus skill level is nil
func sqlGetSkills() ([]instances.Skill, error) {
	var skills []instances.Skill

	rows, err := db.Query("SELECT skill_id, skill_class, skill FROM Skills")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var skill instances.Skill

		if err := rows.Scan(&skill.SkillId, &skill.SkillClass, &skill.Skill); err != nil {
			return nil, err
		}

		skills = append(skills, skill)
	}
	return skills, nil
}

func sqlAddSkill(skill instances.Skill) (int, error) {
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

func sqlGetSkill(id int64) (instances.Skill, error) {
	var skill instances.Skill
	row := db.QueryRow("SELECT * FROM Skills WHERE skill_id=?", id)
	if err := row.Scan(&skill.SkillId, &skill.SkillClass, &skill.Skill); err != nil {
		return instances.Skill{}, err
	}
	return skill, nil
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
func sqlGetAllProjects() ([]instances.Project, error) {
	var projects []instances.Project

	rows, err := db.Query("SELECT * FROM projects")
	if err != nil {
		return nil, fmt.Errorf("sqlGetAllProjects: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var project instances.Project
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

func sqlGetProject(id int64) (instances.Project, error) {
	var proj instances.Project

	row := db.QueryRow("SELECT * FROM Projects WHERE project_id = ?", id)
	if err := row.Scan(&proj.ProjectId, &proj.ClientId, &proj.FocusArea, &proj.Description, &proj.IsSecret); err != nil {
		return instances.Project{}, err
	}
	return proj, nil
}

func sqlGetAllClients() ([]instances.Client, error) {
	var clients []instances.Client

	rows, err := db.Query("SELECT * FROM Clients")
	if err != nil {
		return nil, fmt.Errorf("sqlGetAllClients: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var client instances.Client
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

func sqlGetClient(id int64) (instances.Client, error) {
	var client instances.Client
	row := db.QueryRow("SELECT * FROM Clients WHERE id = ?", id)
	if err := row.Scan(&client.ID, &client.Name, &client.Description); err != nil {
		return instances.Client{}, err
	}
	return client, nil
}

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
	empStore, err := NewMySQLEmployeeStore(cfg)
	if err != nil {
		log.Fatal(err)
	}

	empHandler := NewEmployeeHandler(empStore)
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

	router.GET("/v1/projects", getProjects)
	router.GET("/v1/projects/:id", getProject)

	router.GET("/v1/clients", getClients)
	router.GET("/v1/clients/:id", getClient)

	router.GET("/v1/skills", getSkills)
	router.GET("/v1/skills/:id", getSkill)
	router.POST("/v1/skills", addSkill)
	router.PUT("/v1/skills/:id", updateSkill)
	router.DELETE("/v1/skills/:id", deleteSkill)

	router.Run("localhost:9090")

}
