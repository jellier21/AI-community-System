package AgentService

import (
	"Service/Agent/AgentService/SQLservice"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
)

func AgentSQLService(r *gin.Engine, db *sql.DB) {
	agent := r.Group("/agent")
	medical := agent.Group("/medical")
	medical.POST("/add", func(context *gin.Context) {
		var req struct {
			Account   string `json:"account"`
			Symptom   string `json:"symptom"`
			Diagnosis string `json:"diagnosis"`
		}
		if err := context.ShouldBindJSON(&req); err != nil {
			context.JSON(400, gin.H{"error": err.Error()})
			return
		}
		_account_ := req.Account
		_symptom_ := req.Symptom
		_diagnosis_ := req.Diagnosis
		log.Printf("Account: %s, Symptom: %s, Diagnosis: %s", req.Account, _symptom_, _diagnosis_)

		err := SQLservice.InsertMedicalLog(db, _account_, _symptom_, _diagnosis_)
		if err != nil {
			context.JSON(500, gin.H{"error": err.Error()})
		} else {
			context.JSON(200, gin.H{"message": "Medical log added successfully"})
		}
	})
	medical.GET("/logs", func(context *gin.Context) {
		_account_ := context.Query("account")
		medicalLogs, err := SQLservice.GetMedicalLog(db, _account_)
		if err != nil {
			context.JSON(500, gin.H{"error": err.Error()})
		} else {
			context.JSON(200, gin.H{"medicalLogs": medicalLogs})
		}
	})
}
