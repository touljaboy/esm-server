package instances

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

type EmployeeSkill struct {
	SkillId    int64 `json:"skill_id"`
	SkillLevel int64 `json:"skill_level"`
}

type EmployeeProject struct {
	ProjectId   int64  `json:"project_id"`
	ProjectRole string `json:"project_role"`
}
