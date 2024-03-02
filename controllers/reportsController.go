package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"massivebio/database"
	"massivebio/models"
	"massivebio/utils"
	"net/http"
	"strings"
)

func MassiveFilter(w http.ResponseWriter, r *http.Request) {
	// Just for Pagination without filter (only GET)------------

	if r.Method == http.MethodGet {
		//GET METHOD---------------------------------------------------------
		// 'page' ve 'page_size' control (not allowed string and negative numbers)
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
		//POST METHOD---------------------------------------------------------
		// Pagination + Filter (post + get)----------

		// FilterRequest, represents incoming data
		type FilterRequest struct {
			Filters  map[string]interface{} `json:"filters"`
			Ordering []map[string]string    `json:"ordering"`
		}

		// Security
		var allowedColumns = map[string]bool{
			"main_uploaded_variation": true,
			"main_existing_variation": true,
			"main_symbol":             true,
			"main_af_vcf":             true,
			"main_dp":                 true,
			"details2_provean":        true,
			"details2_dann_score":     true,
			"links_mondo":             true,
			"links_pheno_pubmed":      true,
			"row":                     true,
		}

		var allowedDirections = map[string]bool{
			"ASC":  true,
			"DESC": true,
		}

		var req FilterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.SendJSONError(w, "400 - Bad Request", http.StatusBadRequest)
			return
		}
		type Response struct {
			Page     int             `json:"page"`
			PageSize int             `json:"page_size"`
			Count    int             `json:"count"`
			Results  []models.Report `json:"results"`
		}

		// filters:{}
		var args []interface{}
		var whereClauses []string
		for col, value := range req.Filters {
			if !allowedColumns[col] {
				utils.SendJSONError(w, "Invalid column in filters", http.StatusBadRequest)
				return
			}
			whereClause := fmt.Sprintf("%s = $%d", col, len(args)+1)
			whereClauses = append(whereClauses, whereClause)
			args = append(args, value)
		}

		// ordering:[]
		var orderClauses []string
		for _, orderMap := range req.Ordering {
			for col, dir := range orderMap {
				if !allowedColumns[col] || !allowedDirections[dir] {
					utils.SendJSONError(w, "Invalid column or direction in ordering", http.StatusBadRequest)
					return
				}
				orderClause := fmt.Sprintf("%s %s", col, dir)
				orderClauses = append(orderClauses, orderClause)
			}
		}

		// Create Where and Order By Sentences
		where := ""
		if len(whereClauses) > 0 {
			where = "WHERE " + strings.Join(whereClauses, " AND ")
		}

		order := ""
		if len(orderClauses) > 0 {
			order = "ORDER BY " + strings.Join(orderClauses, ", ")
		}

		// 'page' ve 'page_size' control (not allowed string and negative numbers)
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

		//Total Records
		countQuery := fmt.Sprintf("SELECT COUNT(*) FROM report_output %s", where)
		var totalRecords int
		queryTotalNumber := database.DB.QueryRow(countQuery, args...)
		queryTotalNumber.Scan(&totalRecords)
		if queryTotalNumber.Err() != nil {
			utils.SendJSONError(w, "Database query error: Cannot get total record count", http.StatusInternalServerError)
			return
		}
		// Create Query
		//query := fmt.Sprintf("SELECT * FROM report_output %s %s", where, order)
		offset := (page - 1) * pageSize
		query := fmt.Sprintf("SELECT * FROM report_output %s %s LIMIT $%d OFFSET $%d", where, order, len(args)+1, len(args)+2)
		args = append(args, pageSize, offset)

		// Response
		rows, err := database.DB.Query(query, args...)
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

		response := Response{
			Page:     page,
			PageSize: pageSize,
			Count:    totalRecords,
			Results:  reports,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}

}
