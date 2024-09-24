package main

import (
	"esmAPI/pkg/instances"
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
	result, err := h.store.Update(currEmployee)
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
	}
	result, err := h.store.Delete(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"err": err})
	}
	context.IndentedJSON(http.StatusOK, gin.H{"rows_affected": result})
}

func (h EmployeeHandler) getEmployee(context *gin.Context) {
	strId := context.Params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

// NewSkillHandler - constructor
func NewSkillHandler(store skillStore) *SkillHandler {
	return &SkillHandler{
		store: store,
	}
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
	result, err := h.store.Update(currSkill)
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
	}
	result, err := h.store.Delete(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"err": err})
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
	}
	project, err := h.store.Get(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	context.IndentedJSON(http.StatusOK, project)
}
