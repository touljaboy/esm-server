package main

import (
	"database/sql"
	"esmAPI/pkg/instances"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"strconv"
)

// data store interface for employee
type employeeStore interface {
	Add(emp instances.Employee) (int, error)
	Get(employeeId int64) (emp instances.Employee, err error)
	List() ([]instances.Employee, error)
	Update(emp instances.Employee) (int64, error)
	Delete(employeeId string) (int64, error)
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

func (s *MySQLEmployeeStore) Delete(employeeId string) (int64, error) {
	id, err := strconv.ParseInt(employeeId, 10, 64)
	if err != nil {
		return -1, err
	}
	result, err := s.db.Exec("DELETE FROM Employees WHERE employee_id=?", id)
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
