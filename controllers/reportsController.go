package controllers

import (
	"database/sql"
	"encoding/json"
	"massivebio/database"
	"massivebio/models"
	"massivebio/utils"
	"net/http"
)

func MassiveFilter(w http.ResponseWriter, r *http.Request) {
	// Just for Pagination without filter (only GET)

	if r.Method == http.MethodGet {

		// 'page' ve 'page_size' control
		// Show page (def:1)
		page, err := utils.ValidateQueryParam(r, "page", 1)
		if err != nil || page < 1 {
			utils.SendJSONError(w, "Invalid 'page' query parameter", http.StatusBadRequest)
			return
		}
		// Show data per page (def:10)
		pageSize, err := utils.ValidateQueryParam(r, "page_size", 10)
		if err != nil || pageSize <= 0 {
			utils.SendJSONError(w, "Invalid 'page_size' query parameter", http.StatusBadRequest)
			return
		}

		// Show DATA
		offset := (page - 1) * pageSize
		query := "SELECT * FROM report_output ORDER BY row ASC LIMIT $1 OFFSET $2"

		rows, err := database.DB.Query(query, pageSize, offset)
		if err != nil {
			utils.SendJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var reports []models.Report
		for rows.Next() {
			var report models.Report
			var details2DannScore sql.NullFloat64
			err := rows.Scan(&report.Row, &report.MainUploadedVariation, &report.MainExistingVariation, &report.MainSymbol, &report.MainAfVcf, &report.MainDp, &report.Details2Provean, &details2DannScore, &report.LinksMondo, &report.LinksPhenoPubmed)
			if err != nil {
				utils.SendJSONError(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if details2DannScore.Valid {
				report.Details2DannScore = &details2DannScore.Float64
			}
			reports = append(reports, report)
		}

		if len(reports) == 0 {
			utils.SendJSONError(w, "No records found for the requested page.", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reports)

	} else if r.Method == http.MethodPost {
		// Pagination + Filter (post + get)
		utils.SendJSONError(w, "Invalid 'page_size' query parameter", http.StatusBadRequest)
	}

}
