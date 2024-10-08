package main

import (
	"esmAPI/pkg/instances"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type EmployeeHandler struct {
	store employeeStore
}

type SkillHandler struct {
	store skillStore
}

type ProjectHandler struct {
	store projectStore
}

type ClientHandler struct {
	store clientStore
}

// NewEmployeeHandler - constructor
func NewEmployeeHandler(store employeeStore) *EmployeeHandler {
	return &EmployeeHandler{
		store: store,
	}
}

func (h EmployeeHandler) addEmployee(context *gin.Context) {
	var emp instances.Employee
	if err := context.BindJSON(&emp); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	result, err := h.store.Add(emp)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	context.IndentedJSON(http.StatusCreated, gin.H{"rows_affected": result})
}

// full entry update done by id
func (h EmployeeHandler) updateEmployee(context *gin.Context) {
	strId := context.Params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	currEmployee, err := h.store.Get(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	if err := context.BindJSON(&currEmployee); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	result, err := h.store.Update(id, currEmployee)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{"rows_affected": result})
}

func (h EmployeeHandler) deleteEmployee(context *gin.Context) {
	strId := context.Params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.store.Delete(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{"rows_affected": result})
}

func (h EmployeeHandler) getEmployee(context *gin.Context) {
	strId := context.Params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	employee, err := h.store.Get(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	context.IndentedJSON(http.StatusOK, employee)
}

func (h EmployeeHandler) getEmployees(context *gin.Context) {
	employees, err := h.store.List()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.IndentedJSON(http.StatusOK, employees)
}

func (h EmployeeHandler) getFullEmployees(context *gin.Context) {
	fullEmployees, err := h.store.ListFull()
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	context.IndentedJSON(http.StatusOK, fullEmployees)
}

func (h EmployeeHandler) getFullEmployee(context *gin.Context) {
	strId := context.Params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fullEmployee, err := h.store.GetFull(id)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	context.IndentedJSON(http.StatusOK, fullEmployee)
}

func (h EmployeeHandler) addSkill(context *gin.Context) {
	strId := context.Params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var empSkill instances.EmployeeSkill
	if err := context.BindJSON(&empSkill); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	fmt.Printf("empdId:%d, skillId: %d, skillLevel: %d", id, empSkill.SkillId, empSkill.SkillLevel)
	result, err := h.store.AddSkill(id, empSkill.SkillId, empSkill.SkillLevel)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	context.IndentedJSON(http.StatusCreated, gin.H{"rows_affected": result})
}

// NewSkillHandler - constructor
func NewSkillHandler(store skillStore) *SkillHandler {
	return &SkillHandler{
		store: store,
	}
}

func (h EmployeeHandler) deleteSkill(context *gin.Context) {
	strId := context.Params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var empSkill instances.EmployeeSkill
	if err := context.BindJSON(&empSkill); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	result, err := h.store.DeleteSkill(id, empSkill.SkillId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{"rows_affected": result})
}

func (h EmployeeHandler) updateSkill(context *gin.Context) {
	strId := context.Params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var empSkill instances.EmployeeSkill
	if err := context.BindJSON(&empSkill); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	result, err := h.store.UpdateSkill(id, empSkill.SkillId, empSkill.SkillLevel)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{"rows_affected": result})
}

func (h EmployeeHandler) addProject(context *gin.Context) {
	strId := context.Params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var empProject instances.EmployeeProject
	if err := context.BindJSON(&empProject); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	result, err := h.store.AddProject(empProject.ProjectId, id, empProject.ProjectRole)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	context.IndentedJSON(http.StatusCreated, gin.H{"rows_affected": result})
}

func (h EmployeeHandler) updateProject(context *gin.Context) {
	strId := context.Params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var empProject instances.EmployeeProject
	if err := context.BindJSON(&empProject); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	result, err := h.store.UpdateProject(empProject.ProjectId, id, empProject.ProjectRole)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{"rows_affected": result})
}

func (h EmployeeHandler) deleteProject(context *gin.Context) {
	strId := context.Params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var empProject instances.EmployeeProject
	if err := context.BindJSON(&empProject); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	result, err := h.store.DeleteProject(empProject.ProjectId, id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{"rows_affected": result})
}

func (h SkillHandler) getSkills(context *gin.Context) {
	skills, err := h.store.List()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	context.IndentedJSON(http.StatusOK, skills)
}

func (h SkillHandler) getSkill(context *gin.Context) {
	strId := context.Params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	skill, err := h.store.Get(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	context.IndentedJSON(http.StatusOK, skill)
}

func (h SkillHandler) addSkill(context *gin.Context) {
	var skill instances.Skill
	if err := context.BindJSON(&skill); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	result, err := h.store.Add(skill)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	context.IndentedJSON(http.StatusCreated, gin.H{"rows_affected": result})
}

func (h SkillHandler) updateSkill(context *gin.Context) {
	strId := context.Params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	currSkill, err := h.store.Get(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	if err := context.BindJSON(&currSkill); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	result, err := h.store.Update(id, currSkill)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{"rows_affected": result})
}

func (h SkillHandler) deleteSkill(context *gin.Context) {
	strId := context.Params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.store.Delete(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{"rows_affected": result})
}

// NewProjectHandler - constructor
func NewProjectHandler(store projectStore) *ProjectHandler {
	return &ProjectHandler{
		store: store,
	}
}

func (h ProjectHandler) getProjects(context *gin.Context) {
	projects, err := h.store.List()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	context.IndentedJSON(http.StatusOK, projects)
}

func (h ProjectHandler) getProject(context *gin.Context) {
	strId := context.Params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	project, err := h.store.Get(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	context.IndentedJSON(http.StatusOK, project)
}

func (h ProjectHandler) addProject(context *gin.Context) {
	var project instances.Project
	if err := context.BindJSON(&project); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	result, err := h.store.Add(project)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	context.IndentedJSON(http.StatusCreated, gin.H{"rows_affected": result})
}

func (h ProjectHandler) updateProject(context *gin.Context) {
	strId := context.Params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	proj, err := h.store.Get(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	if err := context.BindJSON(&proj); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	result, err := h.store.Update(id, proj)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{"rows_affected": result})
}

func (h ProjectHandler) deleteProject(context *gin.Context) {
	strId := context.Params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.store.Delete(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{"rows_affected": result})
}

// NewClientHandler - constructor
func NewClientHandler(store clientStore) *ClientHandler {
	return &ClientHandler{
		store: store,
	}
}

func (h ClientHandler) getClients(context *gin.Context) {
	clients, err := h.store.List()
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	context.IndentedJSON(http.StatusOK, clients)
}

func (h ClientHandler) getClient(context *gin.Context) {
	strId := context.Params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	client, err := h.store.Get(id)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	context.IndentedJSON(http.StatusOK, client)
}

func (h ClientHandler) addClient(context *gin.Context) {
	var client instances.Client
	if err := context.BindJSON(&client); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	result, err := h.store.Add(client)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	context.IndentedJSON(http.StatusCreated, gin.H{"rows_affected": result})
}

func (h ClientHandler) updateClient(context *gin.Context) {
	strId := context.Params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	client, err := h.store.Get(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	if err := context.BindJSON(&client); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	result, err := h.store.Update(id, client)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{"rows_affected": result})
}

func (h ClientHandler) deleteClient(context *gin.Context) {
	strId := context.Params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.store.Delete(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{"rows_affected": result})
}
