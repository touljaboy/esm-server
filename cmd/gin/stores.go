package main

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
	Update(emp instances.Employee) (int64, error)
	Delete(employeeId int64) (int64, error)
}

type skillStore interface {
	Add(skill instances.Skill) (int, error)
	Get(skillId int64) (emp instances.Skill, err error)
	List() ([]instances.Skill, error)
	Update(skill instances.Skill) (int64, error)
	Delete(skillId int64) (int64, error)
}

type projectStore interface {
	Add(proj instances.Project) (int, error)
	Get(projId int64) (emp instances.Project, err error)
	List() ([]instances.Project, error)
	Update(proj instances.Project) (int64, error)
	Delete(projId int64) (int64, error)
}

type MySQLEmployeeStore struct {
	db *sql.DB
}

func NewMySQLEmployeeStore(cfg mysql.Config) (*MySQLEmployeeStore, error) {

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

func (s *MySQLEmployeeStore) Update(emp instances.Employee) (int64, error) {
	result, err := s.db.Exec(
		"UPDATE Employees SET name=?, lastname=?, focus_area=?, email=? WHERE employee_id = ?",
		emp.Name, emp.Lastname, emp.FocusArea, emp.Email, emp.EmployeeId)
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

func NewMySQLSkillStore(cfg mysql.Config) (*MySQLSkillStore, error) {

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

func (s *MySQLSkillStore) Update(skill instances.Skill) (int64, error) {
	result, err := s.db.Exec(
		"UPDATE Skills SET skill_id=?, skill_class=?, skill=? WHERE skill_id = ?",
		skill.SkillId, skill.SkillClass, skill.Skill, skill.SkillId)
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
	//TODO
	return 0, nil
}

func (s *MySQLProjectStore) Update(proj instances.Project) (int64, error) {
	//TODO
	return 0, nil
}
func (s *MySQLProjectStore) Delete(projId int64) (int64, error) {
	//TODO
	return 0, nil
}
