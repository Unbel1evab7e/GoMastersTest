package main

import (
	_ "GoMastersTest/docs"
	"GoMastersTest/middleware"
	"GoMastersTest/user/delivery/http"
	"GoMastersTest/user/repository"
	"GoMastersTest/user/usecase"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"
	"time"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		fmt.Println("Service RUN on DEBUG mode")
	}
}

// @title           Swagger Example API
// @version         1.0
// @description     This is a MasterGo Test.

// @host      localhost:8080
// @BasePath  /api/v1

func main() {
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	connection := fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s sslmode=disable", dbName, dbUser, dbPass, dbHost, dbPort)
	dsn := fmt.Sprintf("%s", connection)
	dbConn, err := sql.Open(`postgres`, dsn)
	if err != nil && viper.GetBool("debug") {
		fmt.Println(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	r := gin.Default()
	middL := middleware.InitMiddleware()
	r.Use(middL.Logger())
	userRepo := repository.NewPostgreUserRepository(dbConn)

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	uc := usecase.NewUserUseCase(userRepo, timeoutContext)
	http.NewUserHandler(r, uc)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(viper.GetString("server.address"))
}
