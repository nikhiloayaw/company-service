package utils

type EmployeeCount struct {
	MinEmployeesInATeam int
	MaxEmployeesInATeam int
}

func GetEmployeeCountInATeamBasedOnDepartment(departmentEmployeePercentage int) EmployeeCount {

	var empCount EmployeeCount

	switch {
	case departmentEmployeePercentage <= 10: // department like HR
		empCount.MinEmployeesInATeam = 4
		empCount.MaxEmployeesInATeam = 8
	case departmentEmployeePercentage <= 40:
		empCount.MinEmployeesInATeam = 9
		empCount.MaxEmployeesInATeam = 15
	default: // department like Development
		empCount.MinEmployeesInATeam = 16
		empCount.MaxEmployeesInATeam = 25
	}

	return empCount
}
