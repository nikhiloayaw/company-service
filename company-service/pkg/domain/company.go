package domain

import "time"

/*
** Important Note
** 500 - 100 employees
** Department like HR should contains less members
 */

type Company struct {
	Name     string
	Industry string
	CEO      string
	// EstablishedYear time.Time
	Departments    []Department
	TotalEmployees int
}

type Department struct {
	Name           string
	Leader         Employee
	Teams          []Team
	TotalEmployees int
}

type Team struct {
	Name           string
	Manager        Employee
	Members        []Employee
	TotalEmployees int
}

type Employee struct {
	ID       string
	Name     string
	Age      int
	Email    string
	Salary   float64
	HireDate time.Time
	Role     string
}
