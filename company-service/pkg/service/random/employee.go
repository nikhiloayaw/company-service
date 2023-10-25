package random

import (
	"company-service/pkg/domain"
	"company-service/pkg/store"
	"company-service/pkg/utils"
	"math/rand"
	"strings"
	"time"
)

const (
	employeeMinAge = 18
	employeeMaxAge = 50
	minSalary      = 150000
	maxSalary      = 500000
)

var (
	startTime = time.Date(2014, 11, 12, 0, 0, 0, 0, time.Local)
)

func (r *randomGenerator) CreateEmployee(team string) domain.Employee {

	// select a random name from names
	name := r.getRandomName()

	return domain.Employee{
		ID:    utils.GenerateUUID(),
		Name:  name,
		Email: createEmail(name),
		Age:   utils.GetIntBetween(employeeMinAge, employeeMaxAge),
		Role:  store.GetRandomRoleByTeamName(team), // get a random role according to the team
		// needs to be change
		Salary:   store.GetRandomSalary(),
		HireDate: utils.GetTimeBetween(startTime, time.Now()),
	}

}

func createEmail(name string) string {

	// change the name to lowercase
	name = strings.ToLower(name)

	numChar := []byte("123456789")

	extraLetters := 4

	// cover the name to byte slice as email
	email := []byte(name)

	for i := 1; i <= extraLetters; i++ {
		// take a random number character from numChan and append to email
		email = append(email, numChar[rand.Intn(len(numChar))])
	}

	mailEnd := []byte("@email.com")
	// add the email ending to the email
	email = append(email, mailEnd...)

	// convert the mail slice to string and return
	return string(email)
}
