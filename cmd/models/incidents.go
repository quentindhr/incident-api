package models

import (
	"database/sql"
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
	rows, err := database.DB.Query("SELECT id, title, description, status, severity, created_at, resolved_at FROM incidents")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var incidents []Incident
	for rows.Next() {
		var incident Incident
		err := rows.Scan(&incident.ID, &incident.Title, &incident.Description, &incident.Severity, &incident.Status, &incident.CreatedAt, &incident.ResolvedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		incidents = append(incidents, incident)
	}

	c.JSON(http.StatusOK, gin.H{"incidents": incidents})
}

func GetIncidentByID(c *gin.Context) {
	id := c.Param("id")
	query := `SELECT id, title, description, severity, status, created_at, resolved_at FROM incidents WHERE id = ?`
	row := database.DB.QueryRow(query, id)

	var incident Incident
	err := row.Scan(&incident.ID, &incident.Title, &incident.Description, &incident.Severity, &incident.Status, &incident.CreatedAt, &incident.ResolvedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Incident not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"incident": incident})
}

func CreateIncident(c *gin.Context) {
	var incident Incident
	if err := c.ShouldBindJSON(&incident); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO incidents (title, description, severity, status, created_at, resolved_at) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := database.DB.Exec(query, incident.Title, incident.Description, incident.Severity, incident.Status, incident.CreatedAt, incident.ResolvedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	incident.ID = int(id)
	incident.CreatedAt = time.Now()

	c.JSON(http.StatusCreated, gin.H{"incident": incident})
}

func DeleteIncident(c *gin.Context) {
	id := c.Param("id")
	query := `DELETE FROM incidents WHERE id = ?`
	result, err := database.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Incident not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Incident deleted successfully"})
}

func UpdateIncident(c *gin.Context) {
	id := c.Param("id")
	var incident Incident
	if err := c.ShouldBindJSON(&incident); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `UPDATE incidents SET title = ?, description = ?, severity = ?, status = ? WHERE id = ?`
	result, err := database.DB.Exec(query, incident.Title, incident.Description, incident.Severity, incident.Status, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Incident not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Incident updated successfully"})
}

func UpdateIncidentStatus(c *gin.Context) {
	id := c.Param("id")
	var incident Incident
	if err := c.ShouldBindJSON(&incident); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `UPDATE incidents SET status = ?, resolved_at = ? WHERE id = ?`
	resolvedAt := time.Now()
	result, err := database.DB.Exec(query, incident.Status, resolvedAt, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Incident not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Incident status updated successfully"})
}
