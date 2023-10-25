package usecase

import (
	"company-service/pkg/domain"
	"company-service/pkg/models"
	"company-service/pkg/service/random"
	"company-service/pkg/store"
	"company-service/pkg/usecase/interfaces"
	"company-service/pkg/utils"
	"fmt"
)

type companyUseCase struct {
	randomGen random.RandomGenerator
}

func NewCompanyUseCase(randomGen random.RandomGenerator) interfaces.CompanyUseCase {

	return &companyUseCase{
		randomGen: randomGen,
	}
}

type departmentDetails struct {
	name                string
	requiredEmployees   int
	minEmployeesInATeam int
	maxEmployeesInATeam int
}

func (c *companyUseCase) Create(companyReq models.CompanyRequest) domain.Company {

	// create a copy of departments
	departments := store.GetAllDepartments()

	company := domain.Company{
		Name:           companyReq.Name,
		CEO:            companyReq.CEO,
		Industry:       "IT Industry",
		TotalEmployees: companyReq.TotalEmployees,
		Departments:    make([]domain.Department, len(departments)), // create a department slice according to the length of department
	}

	// create each department for company
	for i, department := range departments {

		// calculate the minimum and maximum employees in a team by the percentage of employees in the department
		empCount := utils.GetEmployeeCountInATeamBasedOnDepartment(department.PercentageOfEmployees)
		// calculate the required employees for the department based on the total employees needed and employee percentage for the department.
		requiredEmployeesForDepartment := (companyReq.TotalEmployees / 100) * department.PercentageOfEmployees

		// create a department details
		depDetails := departmentDetails{
			name:                department.Name,
			requiredEmployees:   requiredEmployeesForDepartment,
			minEmployeesInATeam: empCount.MinEmployeesInATeam,
			maxEmployeesInATeam: empCount.MaxEmployeesInATeam,
		}
		// create a new department
		company.Departments[i] = c.getDepartment(depDetails)
	}

	return company
}

func (c *companyUseCase) getDepartment(depDetails departmentDetails) domain.Department {

	department := domain.Department{
		Name:  depDetails.name,
		Teams: []domain.Team{},
	}

	var (
		teamName          string
		teamEmployeeCount int
	)

	// keep create team until the needed employees count is comes to zero
	for depDetails.requiredEmployees > 0 {
		// select a random team name for this department, from available team for this department
		teamName = store.GetRandomTeamByDepartmentName(depDetails.name)
		// if the total employees need to create is less than max employees in a team, then set the teamEmployeeCount as the available employees count
		if depDetails.requiredEmployees <= depDetails.maxEmployeesInATeam {
			teamEmployeeCount = depDetails.requiredEmployees

		} else { // else select a random employees count for the team based on the department min and max employees in a team
			teamEmployeeCount = utils.GetIntBetween(depDetails.minEmployeesInATeam, depDetails.maxEmployeesInATeam)
		}
		// update the department employees count
		department.TotalEmployees += teamEmployeeCount
		// update the required employees count
		depDetails.requiredEmployees -= teamEmployeeCount

		// create a new team and append it to department team
		team := c.createTeam(teamName, teamEmployeeCount)
		department.Teams = append(department.Teams, team)
	}

	return department
}

func (c *companyUseCase) createTeam(teamName string, totalEmployees int) domain.Team {

	done := make(chan struct{})

	empChan := Generator(done, teamName, c.randomGen.CreateEmployee)

	// get an employee as manager
	manager := <-empChan
	// update manger role to team manager
	manager.Role = fmt.Sprintf("%s Manager", teamName)
	team := domain.Team{
		Name:           teamName,
		Manager:        manager,                                   // take the first employee created as manager
		Members:        make([]domain.Employee, totalEmployees-1), // create a slice of employees,
		TotalEmployees: totalEmployees,
	}

	// fill employees for team
	for i := range team.Members {
		team.Members[i] = <-empChan
	}

	// close the done channel to let the generator to close as well.
	close(done)

	return team
}

// This function is generic function, which a take a done channel,callback function and return a channel.
// The the callback function return values can be receive through return channel until the done channel close
func Generator[T any](done <-chan struct{}, team string, fn func(team string) T) <-chan T {

	stream := make(chan T)

	go func() {
		defer close(stream)
		for {
			select {
			case <-done:
				return
			case stream <- fn(team):
			}
		}
	}()

	return stream
}
