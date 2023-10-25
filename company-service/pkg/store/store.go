package store

import (
	"company-service/pkg/utils"
	"fmt"
	"sync"
)

var (
	// DepartmentData = []Department{
	// 	{
	// 		Name:  "Development",
	// 		Roles: []string{"Software Developer", "QA Engineer", "DevOps Specialist"},
	// 		Teams: []string{"Frontend Team", "Backend Team", "QA Team"},
	// 	},
	// 	{
	// 		Name:  "HR",
	// 		Roles: []string{"HR Manager", "Recruiter", "Employee Relations Specialist"},
	// 		Teams: []string{"Recruitment Team", "Employee Relations Team"},
	// 	},
	// 	{
	// 		Name:  "Sales",
	// 		Roles: []string{"Sales Manager", "Account Executive"},
	// 		Teams: []string{},
	// 	},
	// 	{
	// 		Name:  "Marketing",
	// 		Roles: []string{"Marketing Manager", "Content Specialist"},
	// 		Teams: []string{},
	// 	},
	// }

	// DepartmentRoles = map[string][]string{
	// 	"Development": {"Software Developer", "QA Engineer", "DevOps Specialist"},
	// 	"HR":          {"HR Manager", "Recruiter", "Employee Relations Specialist"},
	// 	// "Sales":       {"Sales Manager", "Account Executive"},
	// 	// "Marketing":   {"Marketing Manager", "Content Specialist"},
	// }
	//
	mu sync.RWMutex

	// department and percentage should be calculated logically by ourself
	departments = []Department{ // distributing 100 percentage to everyone
		{
			Name:                  "Development",
			PercentageOfEmployees: 75,
		},
		{
			Name:                  "HR",
			PercentageOfEmployees: 5,
		},
		{
			Name:                  "Sales",
			PercentageOfEmployees: 10,
		},
		{
			Name:                  "Marketing",
			PercentageOfEmployees: 10,
		},
	}

	departmentTeamNames = map[string][]string{
		"Development": {"Frontend Team", "Backend Team", "QA Team"},
		"HR":          {"Recruitment Team", "Employee Relations Team"},
		"Sales":       {"Sales Team"},
		"Marketing":   {"Marketing Team"},
	}

	teamRoles = map[string][]string{
		// Development
		"Frontend Team": {"Frontend Developer", "UI/UX Designer", "Quality Assurance Engineer"},
		"Backend Team":  {"Backend Developer", "Database Administrator", "DevOps Engineer"},
		"QA Team":       {"Quality Assurance Engineer", "Test Automation Engineer"},
		// HR
		"Recruitment Team":        {"HR Manager", "Recruiter", "HR Coordinator"},
		"Employee Relations Team": {"HR Manager", "Employee Relations Specialist"},
		// Sales
		"Sales Team": {"Sales Manager", "Account Executive"},
		// Marketing
		"Marketing Team": {"Marketing Manager", "Content Specialist"},
	}

	salaries = []float64{20000, 30000, 40000, 50000, 60000, 70000, 80000, 90000, 10000}
)

type Department struct {
	Name                  string
	PercentageOfEmployees int
}

// to get a copy of departments(to avoid passing the actual department to avoid mutation)
func GetAllDepartments() []Department {

	d := make([]Department, len(departments))
	// copy the departments
	copy(d, departments)

	return d
}

func GetRandomRoleByTeamName(team string) string {

	mu.RLock()
	defer mu.RUnlock()
	// get the role of the team
	roles, ok := teamRoles[team]

	if !ok {
		fmt.Println("invalid department got")
		return ""
	}
	// return a random role from the roles
	return roles[utils.GetRandomIndex(len(roles))]
}

func GetRandomTeamByDepartmentName(department string) string {

	mu.RLock()
	defer mu.RUnlock()
	// get the teams of the department
	teams, ok := departmentTeamNames[department]

	if !ok {
		fmt.Println("invalid department got")
		return ""
	}
	// return a random teams from the teams
	return teams[utils.GetRandomIndex(len(teams))]
}

func GetRandomSalary() float64 {
	mu.RLock()
	defer mu.RUnlock()

	return salaries[utils.GetRandomIndex(len(salaries))]
}
