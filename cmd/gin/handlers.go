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

// NewEmployeeHandler - constructor
func NewEmployeeHandler(store employeeStore) *EmployeeHandler {
	return &EmployeeHandler{
		store: store,
	}
}

func (h EmployeeHandler) addEmployee(context *gin.Context) {
	var emp instances.Employee
	if err := context.BindJSON(&emp); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	result, err := h.store.Add(emp)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	context.IndentedJSON(http.StatusCreated, gin.H{"rows_affected": result})
}

// full entry update done by id
func (h EmployeeHandler) updateEmployee(context *gin.Context) {
	strId := context.Params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	currEmployee, err := h.store.Get(id)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	if err := context.BindJSON(&currEmployee); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	result, err := h.store.Update(currEmployee)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{"rows_affected": result})
}

func (h EmployeeHandler) deleteEmployee(context *gin.Context) {
	id := context.Params.ByName("id")
	result, err := h.store.Delete(id)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"err": err})
	}
	context.IndentedJSON(http.StatusOK, gin.H{"rows_affected": result})
}

func (h EmployeeHandler) getEmployee(context *gin.Context) {
	strId := context.Params.ByName("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	employee, err := h.store.Get(id)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	context.IndentedJSON(http.StatusOK, employee)
}

func (h EmployeeHandler) getEmployees(context *gin.Context) {
	employees, err := h.store.List()
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.IndentedJSON(http.StatusOK, employees)
}
