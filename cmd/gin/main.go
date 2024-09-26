package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
)

// TODO introduce a logger
// TODO tests can be written in .http format
// TODO also, consider using a router gorilla/mux
// TODO ids should probably be a uint
// TODO need way better error messages to get sent, because this fucking sucks dude, no logs, no anything to debug
// TODO adding, updating, deleting a Client
// TODO adding, updating, deleting a Skill to an Employee
// TODO adding, updating, deleting an Employee to a Project

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
	empStore, err := NewEmployeeStore(cfg)
	if err != nil {
		log.Fatal(err)
	}
	skillStore, err := NewSkillStore(cfg)
	if err != nil {
		log.Fatal(err)
	}
	projectStore, err := NewProjectStore(cfg)
	if err != nil {
		log.Fatal(err)
	}
	clientStore, err := NewClientStore(cfg)
	if err != nil {
		log.Fatal(err)
	}
	// create handlers
	empHandler := NewEmployeeHandler(empStore)
	skillHandler := NewSkillHandler(skillStore)
	projectHandler := NewProjectHandler(projectStore)
	clientHandler := NewClientHandler(clientStore)
	//Configure endpoints
	router := gin.Default()
	router.Routes()
	router.GET("/v1/employees", empHandler.getEmployees)
	router.GET("/v1/employees/:id", empHandler.getEmployee)
	router.POST("/v1/employees", empHandler.addEmployee)
	router.PUT("/v1/employees/:id", empHandler.updateEmployee)
	router.DELETE("/v1/employees/:id", empHandler.deleteEmployee)

	router.GET("/v1/fullEmployees", empHandler.getFullEmployees)
	router.GET("/v1/fullEmployees/:id", empHandler.getFullEmployee)
	//special endpoints
	router.POST("/v1/skills/employees/:id", empHandler.addSkill)
	router.DELETE("/v1/skills/employees/:id", empHandler.deleteSkill)
	router.PUT("/v1/skills/employees/:id", empHandler.updateSkill)
	router.POST("/v1/projects/employees/:id", empHandler.addProject)
	router.DELETE("/v1/projects/employees/:id", empHandler.deleteProject)
	router.PUT("/v1/projects/employees/:id", empHandler.updateProject)

	router.GET("/v1/projects", projectHandler.getProjects)
	router.GET("/v1/projects/:id", projectHandler.getProject)
	router.POST("/v1/projects", projectHandler.addProject)
	router.PUT("v1/projects/:id", projectHandler.updateProject)
	router.DELETE("v1/projects/:id", projectHandler.deleteProject)

	router.GET("/v1/clients", clientHandler.getClients)
	router.GET("/v1/clients/:id", clientHandler.getClient)
	router.POST("/v1/clients", clientHandler.addClient)
	router.PUT("v1/clients/:id", clientHandler.updateClient)
	router.DELETE("v1/clients/:id", clientHandler.deleteClient)

	router.GET("/v1/skills", skillHandler.getSkills)
	router.GET("/v1/skills/:id", skillHandler.getSkill)
	router.POST("/v1/skills", skillHandler.addSkill)
	router.PUT("/v1/skills/:id", skillHandler.updateSkill)
	router.DELETE("/v1/skills/:id", skillHandler.deleteSkill)

	router.Run("localhost:9090")

}
