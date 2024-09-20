package main

import (
	"encoding/xml"
	"fmt"

	"os"

	"github.com/GuohaoMa/tucowDemo/database"
	"github.com/GuohaoMa/tucowDemo/handlers"
	"github.com/GuohaoMa/tucowDemo/model"
	"github.com/GuohaoMa/tucowDemo/validation"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	err := validation.Validate("data/exampleTest.xml")
	if err != nil {
		fmt.Println("Error in validating XML file:", err)
		return
	}

	xmlFile, err := os.Open("data/exampleTest.xml")
	if err != nil {
		fmt.Println("Error opening XML file:", err)
		return
	}
	defer xmlFile.Close()

	graph := model.Graph{Db: database.Db}
	xmlParserDecoder := xml.NewDecoder(xmlFile)
	xmlParserDecoder.Decode(&graph)

	err = graph.Create()
	if err != nil {
		fmt.Println("Error in saving graph to database:", err)
	}

	// register gin server and run
	var r = gin.New()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))
	graph2 := model.Graph{Db: database.Db, Id: graph.Id}
	r.POST("/graphs/paths", handlers.FindPathHandler(&graph2))
	r.Run(":" + "8080")
}
