package main

import (
	"github.com/gin-gonic/gin"
	"github.com/quentindhr/incident-api.git/cmd/database"
	"github.com/quentindhr/incident-api.git/cmd/models"
)

func main() {
	database.InitDB("incidents.db")
	r := gin.Default()

	r.GET("/incidents", models.GetIncidents)
	r.GET("/incidents/:id", models.GetIncidentByID)
	r.POST("/incidents", models.CreateIncident)
	r.DELETE("/incidents/:id", models.DeleteIncident)
	r.PUT("/incidents/:id", models.UpdateIncident)
	r.PATCH("/incidents/:id/resolve", models.UpdateIncidentStatus)
	{
		r.Run()

	}
}
