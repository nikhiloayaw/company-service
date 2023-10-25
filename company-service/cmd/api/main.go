package main

import (
	"company-service/pkg/config"
	"company-service/pkg/di"
	"log"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	srv, err := di.InitializeAPI(cfg)
	if err != nil {
		log.Fatalf("failed to initialize api: %v", err)
	}

	if err = srv.Start(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

// var (
// 	TotalEmployees = 500
// )

// func main() {

// 	randomGenerator, err := random.NewRandomGenerator()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	usecase := usecase.NewCompanyUseCase(randomGenerator)

// 	company := usecase.Create()

// 	// fmt.Printf("company: %+v\n", company)

// 	data, err := json.Marshal(company)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	file, err := os.Create("./company.json")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	_, err = file.Write(data)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
