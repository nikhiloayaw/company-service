package random

import (
	"company-service/pkg/domain"
	"company-service/pkg/file"
	"company-service/pkg/utils"
	"sync"
)

const (
	namesSheetFileName = "./inputs/names.xls"
	namesSheetName     = "names"
)

type RandomGenerator interface {
	CreateEmployee(team string) domain.Employee
}

type randomGenerator struct {
	mu    sync.RWMutex
	names []string
}

func NewRandomGenerator() (RandomGenerator, error) {

	// get the names for random generation
	names, err := file.GetAllNamesFromSheet(namesSheetFileName, namesSheetName)
	if err != nil {
		return nil, err
	}

	return &randomGenerator{
		names: names,
		mu:    sync.RWMutex{},
	}, nil
}

// To get a random name from random generator with thread safe
func (r *randomGenerator) getRandomName() string {

	// lock the mutex for read.
	r.mu.RLock()
	defer r.mu.RUnlock()

	// return a random name from the name
	return r.names[utils.GetRandomIndex(len(r.names))]
}
