package models

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/quentindhr/incident-api.git/cmd/database"
)

type Incident struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Severity    string     `json:"severity"` // Low, Medium, High
	Status      string     `json:"status"`   // Open, In Progress, Resolved
	CreatedAt   time.Time  `json:"created_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
}

func GetIncidents(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, title, description, created_at FROM incidents")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var incidents []Incident
	for rows.Next() {
		var incident Incident
		err := rows.Scan(&incident.ID, &incident.Title, &incident.Description, &incident.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		incidents = append(incidents, incident)
	}

	c.JSON(http.StatusOK, gin.H{"incidents": incidents})
}

func CreateIncident(c *gin.Context) {
	var incident Incident
	if err := c.ShouldBindJSON(&incident); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO incidents (title, description, severity, status) VALUES (?, ?, ?, ?)`
	result, err := database.DB.Exec(query, incident.Title, incident.Description, incident.Severity, incident.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	incident.ID = int(id)
	incident.CreatedAt = time.Now()

	c.JSON(http.StatusCreated, gin.H{"incident": incident})
}
