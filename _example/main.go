package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/xbmlz/kitx/config"
	"github.com/xbmlz/kitx/db"
	"github.com/xbmlz/kitx/ginx"
	"github.com/xbmlz/kitx/log"
	"github.com/xbmlz/kitx/server"
)

type TODO struct {
	db.BaseModel
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func main() {
	var cfg struct {
		Server struct {
			Port int `mapstructure:"port"`
		} `mapstructure:"server"`
		DBMap map[string]db.Config `mapstructure:"db"`
	}

	err := config.MustLoad("./config.example.yaml").Parse(&cfg)
	if err != nil {
		log.Error("Failed to parse config: %v", err)
		return
	}

	db.Register(cfg.DBMap)

	db.Get("default").AutoMigrate(&TODO{})

	r := ginx.New()

	// curl -X POST -H "Content-Type: application/json" -d '{"title": "Buy groceries", "description": "Buy groceries for the week", "completed": false}' http://localhost:8080/todos
	r.POST("/todos", func(c *gin.Context) {
		var todo TODO
		if !ginx.BindJSON(c, &todo) {
			return
		}
		db.Get("default").Create(&todo)
		ginx.ResponseOk(c, todo)
	})

	// curl -X GET http://localhost:8080/todos
	r.GET("/todos", func(c *gin.Context) {
		var todos []TODO
		db.Get("default").Find(&todos)
		ginx.ResponseOk(c, todos)
	})

	log.Info("server started on port %d", cfg.Server.Port)
	server.Run(fmt.Sprintf(":%d", cfg.Server.Port), r)
}
