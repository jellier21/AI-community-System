package SQLservice

import (
	"database/sql"
	"time"
)

type Medical struct {
	Symptom   string `json:"symptom"`
	Diagnosis string `json:"diagnosis"`
	CreateAt  string `json:"create_at"`
}

func InsertMedicalLog(db *sql.DB, account, symptom string, diagnosis string) error {
	_, err := db.Exec("INSERT INTO medical (account, symptom, diagnosis, create_at) VALUES (?, ?, ?, ?)", account, symptom, diagnosis, time.Now().Format("2006-01-02 15:04:05"))
	return err
}
func GetMedicalLog(db *sql.DB, account string) ([]*Medical, error) {
	rows, err := db.Query("SELECT symptom, diagnosis, create_at FROM medical WHERE account = ?", account)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	medicalLogs := []*Medical{}
	for rows.Next() {
		medicalLog := &Medical{}
		err := rows.Scan(&medicalLog.Symptom, &medicalLog.Diagnosis, &medicalLog.CreateAt)
		if err != nil {
			return nil, err
		}
		medicalLogs = append(medicalLogs, medicalLog)
	}
	return medicalLogs, nil
}
