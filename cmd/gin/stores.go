package main

//TODO maybe there is a way to at least make the constructor generic?
import (
	"database/sql"
	"esmAPI/pkg/instances"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
)

// data store interface for employee
type employeeStore interface {
	Add(emp instances.Employee) (int, error)
	Get(employeeId int64) (emp instances.Employee, err error)
	List() ([]instances.Employee, error)
	Update(currId int64, emp instances.Employee) (int64, error)
	Delete(employeeId int64) (int64, error)
	GetFull(employeeId int64) (emp instances.EmployeeFull, err error)
	ListFull() ([]instances.EmployeeFull, error)
	AddSkill(employeeId int64, skillId int64, skillLevel int64) (int64, error)
	DeleteSkill(employeeId int64, skillId int64) (int64, error)
	UpdateSkill(employeeId int64, skillId int64, skillLevel int64) (int64, error)
	AddProject(projectId int64, employeeId int64, employeeRole string) (int64, error)
	DeleteProject(projectId int64, employeeId int64) (int64, error)
	UpdateProject(projectId int64, employeeId int64, employeeRole string) (int64, error)
	//TODO associate a project with an employee
}

type skillStore interface {
	Add(skill instances.Skill) (int, error)
	Get(skillId int64) (emp instances.Skill, err error)
	List() ([]instances.Skill, error)
	Update(currId int64, skill instances.Skill) (int64, error)
	Delete(skillId int64) (int64, error)
}

type projectStore interface {
	Add(proj instances.Project) (int, error)
	Get(projId int64) (proj instances.Project, err error)
	List() ([]instances.Project, error)
	Update(currId int64, proj instances.Project) (int64, error)
	Delete(projId int64) (int64, error)
}

type clientStore interface {
	Add(client instances.Client) (int, error)
	Get(clientId int64) (client instances.Client, err error)
	List() ([]instances.Client, error)
	Update(currId int64, client instances.Client) (int64, error)
	Delete(clientId int64) (int64, error)
}

type MySQLEmployeeStore struct {
	db *sql.DB
}

func NewEmployeeStore(cfg mysql.Config) (*MySQLEmployeeStore, error) {

	// Get a database handle.
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
	return &MySQLEmployeeStore{db: db}, nil
}

func (s *MySQLEmployeeStore) Add(emp instances.Employee) (int, error) {
	result, err := s.db.Exec(
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

func (s *MySQLEmployeeStore) Delete(employeeId int64) (int64, error) {
	result, err := s.db.Exec("DELETE FROM Employees WHERE employee_id=?", employeeId)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (s *MySQLEmployeeStore) Update(currId int64, emp instances.Employee) (int64, error) {
	result, err := s.db.Exec(
		"UPDATE Employees SET name=?, lastname=?, focus_area=?, email=? WHERE employee_id = ?",
		emp.Name, emp.Lastname, emp.FocusArea, emp.Email, currId)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (s *MySQLEmployeeStore) Get(employeeId int64) (instances.Employee, error) {
	var emp instances.Employee

	row := s.db.QueryRow("SELECT * FROM Employees WHERE employee_id = ?", employeeId)
	if err := row.Scan(&emp.EmployeeId, &emp.Name, &emp.Lastname, &emp.FocusArea, &emp.Email); err != nil {
		return instances.Employee{}, err
	}
	return emp, nil
}

func (s *MySQLEmployeeStore) List() ([]instances.Employee, error) {
	var employees []instances.Employee

	rows, err := s.db.Query("SELECT * FROM Employees")
	if err != nil {
		return nil, fmt.Errorf("sqlGetAllEmployees %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var emp instances.Employee
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

type MySQLSkillStore struct {
	db *sql.DB
}

func NewSkillStore(cfg mysql.Config) (*MySQLSkillStore, error) {

	// Get a database handle.
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
	return &MySQLSkillStore{db: db}, nil
}

func (s *MySQLSkillStore) Delete(id int64) (int64, error) {
	result, err := s.db.Exec("DELETE FROM Skills WHERE skill_id=?", id)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (s *MySQLSkillStore) Update(currId int64, skill instances.Skill) (int64, error) {
	result, err := s.db.Exec(
		"UPDATE Skills SET skill_id=?, skill_class=?, skill=? WHERE skill_id = ?",
		skill.SkillId, skill.SkillClass, skill.Skill, currId)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

// We use Skill struct which also contains skill level, as it is usually associated with an Employee.
// In this case however, we only want to see what Skills are available in database, thus skill level is nil
func (s *MySQLSkillStore) List() ([]instances.Skill, error) {
	var skills []instances.Skill

	rows, err := s.db.Query("SELECT skill_id, skill_class, skill FROM Skills")
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

func (s *MySQLSkillStore) Add(skill instances.Skill) (int, error) {
	result, err := s.db.Exec(
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

func (s *MySQLSkillStore) Get(id int64) (instances.Skill, error) {
	var skill instances.Skill
	row := s.db.QueryRow("SELECT * FROM Skills WHERE skill_id=?", id)
	if err := row.Scan(&skill.SkillId, &skill.SkillClass, &skill.Skill); err != nil {
		return instances.Skill{}, err
	}
	return skill, nil
}

type MySQLProjectStore struct {
	db *sql.DB
}

func NewProjectStore(cfg mysql.Config) (*MySQLProjectStore, error) {
	// Get a database handle.
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
	return &MySQLProjectStore{db: db}, nil
}

func (s *MySQLProjectStore) List() ([]instances.Project, error) {
	var projects []instances.Project

	rows, err := s.db.Query("SELECT * FROM projects")
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

func (s *MySQLProjectStore) Get(id int64) (instances.Project, error) {
	var proj instances.Project

	row := s.db.QueryRow("SELECT * FROM Projects WHERE project_id = ?", id)
	if err := row.Scan(&proj.ProjectId, &proj.ClientId, &proj.FocusArea, &proj.Description, &proj.IsSecret); err != nil {
		return instances.Project{}, err
	}
	return proj, nil
}

func (s *MySQLProjectStore) Add(proj instances.Project) (int, error) {
	result, err := s.db.Exec("INSERT INTO Projects (project_id, client_id, focus_area, description, isSecret)"+
		" VALUES(?, ?, ?, ?, ?)", proj.ProjectId, proj.ClientId, proj.FocusArea, proj.Description, proj.IsSecret)
	if err != nil {
		return -1, err
	}
	id, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}
	return int(id), nil
}

func (s *MySQLProjectStore) Update(currId int64, proj instances.Project) (int64, error) {
	result, err := s.db.Exec(
		"UPDATE Projects SET project_id=?, client_id=?, focus_area=?, description=?, isSecret=? WHERE project_id = ?",
		proj.ProjectId, proj.ClientId, proj.FocusArea, proj.Description, proj.IsSecret, currId)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}
func (s *MySQLProjectStore) Delete(projId int64) (int64, error) {
	result, err := s.db.Exec("DELETE FROM Projects WHERE project_id = ?", projId)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

type MySQLClientStore struct {
	db *sql.DB
}

func NewClientStore(cfg mysql.Config) (*MySQLClientStore, error) {
	// Get a database handle.
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
	return &MySQLClientStore{db: db}, nil
}

func (s *MySQLClientStore) List() ([]instances.Client, error) {
	var clients []instances.Client

	rows, err := s.db.Query("SELECT * FROM Clients")
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

func (s *MySQLClientStore) Get(id int64) (instances.Client, error) {
	var client instances.Client
	row := s.db.QueryRow("SELECT * FROM Clients WHERE id = ?", id)
	if err := row.Scan(&client.ID, &client.Name, &client.Description); err != nil {
		return instances.Client{}, err
	}
	return client, nil
}

func (s *MySQLClientStore) Add(client instances.Client) (int, error) {
	result, err := s.db.Exec("INSERT INTO Clients (id, name, description)"+
		" VALUES(?, ?, ?)", client.ID, client.Name, client.Description)
	if err != nil {
		return -1, err
	}
	id, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}
	return int(id), nil
}

func (s *MySQLClientStore) Update(currId int64, client instances.Client) (int64, error) {
	result, err := s.db.Exec(
		"UPDATE Clients SET id=?, name=?, description=? WHERE id = ?",
		client.ID, client.Name, client.Description, currId)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}
func (s *MySQLClientStore) Delete(clientId int64) (int64, error) {
	result, err := s.db.Exec("DELETE FROM Clients WHERE id = ?", clientId)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (s *MySQLEmployeeStore) ListFull() ([]instances.EmployeeFull, error) {
	var employeesFull []instances.EmployeeFull

	//first, get all the employees
	employees, err := s.List()
	if err != nil {
		return nil, fmt.Errorf("sqlGetAllProjects: %v", err)
	}

	//iterate through each employee and find associated projects and skills. Then append employeesFull
	for _, employee := range employees {
		employeeFull, err := s.GetFull(employee.EmployeeId)
		if err != nil {
			return nil, fmt.Errorf("sqlGetFullEmployeeById: %v", err)
		}
		employeesFull = append(employeesFull, employeeFull)
	}

	return employeesFull, nil
}

func (s *MySQLEmployeeStore) GetFull(id int64) (instances.EmployeeFull, error) {
	employee, err := s.Get(id)
	if err != nil {
		return instances.EmployeeFull{}, err
	}

	var employeeFull instances.EmployeeFull
	var skills []instances.Skill
	var projects []instances.ProjectFull

	//find associated skills
	rows, err := s.db.Query("SELECT s.skill_id,s.skill_class, s.skill, e.skill_level FROM EmployeeSkills AS e "+
		"INNER JOIN Skills AS s ON e.skill_id = s.skill_id WHERE employee_id = ?", employee.EmployeeId)
	if err != nil {
		return instances.EmployeeFull{}, fmt.Errorf("sqlGetFullEmployees: %v", err)
	}
	for rows.Next() {
		var skill instances.Skill
		if err := rows.Scan(&skill.SkillId, &skill.SkillClass, &skill.Skill, &skill.SkillLevel); err != nil {
			return instances.EmployeeFull{}, fmt.Errorf("sqlGetFullEmployees: %v", err)
		}
		skills = append(skills, skill)
	}

	//find associate projects
	rows, err = s.db.Query("SELECT a.*, b.employee_role FROM Projects AS a "+
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

func (s *MySQLEmployeeStore) AddSkill(employeeId int64, skillId int64, skillLevel int64) (int64, error) {
	result, err := s.db.Exec("INSERT INTO EmployeeSkills (employee_id, skill_id, skill_level) VALUES(?,?,?)",
		employeeId, skillId, skillLevel)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (s *MySQLEmployeeStore) DeleteSkill(employeeId int64, skillId int64) (int64, error) {
	result, err := s.db.Exec("DELETE FROM EmployeeSkills WHERE employee_id=? AND skill_id = ?",
		employeeId, skillId)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (s *MySQLEmployeeStore) UpdateSkill(employeeId int64, skillId int64, skillLevel int64) (int64, error) {
	results, err := s.db.Exec("UPDATE EmployeeSkills SET skill_level=? WHERE employee_id=? AND skill_id=?",
		skillLevel, employeeId, skillId)
	if err != nil {
		return -1, err
	}
	return results.RowsAffected()
}

func (s *MySQLEmployeeStore) AddProject(projectId int64, employeeId int64, projectRole string) (int64, error) {
	result, err := s.db.Exec("INSERT INTO ProjectDetails (project_id, employee_id, employee_role) VALUES (?,?,?)",
		projectId, employeeId, projectRole)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}
func (s *MySQLEmployeeStore) UpdateProject(projectId int64, employeeId int64, projectRole string) (int64, error) {
	result, err := s.db.Exec("UPDATE ProjectDetails SET employee_role=? WHERE project_id=? AND employee_id=?",
		projectRole, projectId, employeeId)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}
func (s *MySQLEmployeeStore) DeleteProject(projectId int64, employeeId int64) (int64, error) {
	result, err := s.db.Exec("DELETE FROM ProjectDetails WHERE project_id=? AND employee_id=?",
		projectId, employeeId)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}
