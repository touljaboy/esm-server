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

type Employee struct {
	EntryId    int64  `json:"entry_id"`
	EmployeeId int64  `json:"employee_id"`
	Name       string `json:"name"`
	Lastname   string `json:"lastname"`
	FocusArea  string `json:"focus_area"`
	SkillClass string `json:"skill_class"`
	Skill      string `json:"skill"`
	SkillLevel int    `json:"skill_level"`
}

type Client struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
	IsSecret    bool   `json:"isSecret"`
}

type Projects struct {
	ID          int64  `json:"id"`
	ClientID    int    `json:"client_id"`
	FocusArea   string `json:"focus_area"`
	Description string `json:"description"`
	IsSecret    bool   `json:"isSecret"`
}

func getEmployees(context *gin.Context) {
	employees, err := sqlGetAllEmployees()
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	context.IndentedJSON(http.StatusOK, employees)
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
		if err := rows.Scan(&emp.EntryId, &emp.EmployeeId, &emp.Name, &emp.Lastname, &emp.FocusArea, &emp.SkillClass, &emp.Skill, &emp.SkillLevel); err != nil {
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
	router.Run("localhost:9090")

}
