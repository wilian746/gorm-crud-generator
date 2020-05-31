package main

import (
	"fmt"
	"github.com/wilian746/go-generator/pkg/standart-gorm/configs"
	_ "github.com/wilian746/go-generator/pkg/standart-gorm/docs" // docs is generatedon base of repository
	"github.com/wilian746/go-generator/pkg/standart-gorm/internal/entities/product"
	"github.com/wilian746/go-generator/pkg/standart-gorm/internal/routes"
	"github.com/wilian746/go-generator/pkg/standart-gorm/pkg/repository/adapter"
	"github.com/wilian746/go-generator/pkg/standart-gorm/pkg/repository/database"
	"log"
	"net/http"
)

// @title Swagger Example API EEEEEEYYYY
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	configs := config.GetConfig()
	entity := &product.Product{}

	connection := database.GetConnection(configs.Dialect, configs.DatabaseURI)
	connection.Table(entity.TableName()).AutoMigrate(entity)
	repository := adapter.NewAdapter(connection)

	port := fmt.Sprintf(":%v", configs.Port)
	router := routes.NewRouter().SetRouters(repository)
	log.Println("service running on port ", port)
	log.Println("swagger running on url: ", fmt.Sprintf("http://localhost:%v/swagger/index.html", configs.Port))

	server := http.ListenAndServe(port, router)
	log.Fatal(server)
}
