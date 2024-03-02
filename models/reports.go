package models

type Report struct {
	Row                   int      `json:"row"`
	MainUploadedVariation string   `json:"main_uploaded_variation"`
	MainExistingVariation string   `json:"main_existing_variation"`
	MainSymbol            string   `json:"main_symbol"`
	MainAfVcf             float64  `json:"main_af_vcf"`
	MainDp                float64  `json:"main_dp"`
	Details2Provean       string   `json:"details2_provean"`
	Details2DannScore     *float64 `json:"details2_dann_score"`
	LinksMondo            string   `json:"links_mondo"`
	LinksPhenoPubmed      string   `json:"links_pheno_pubmed"`
}

type FilterRequest struct {
	Filters  map[string]interface{} `json:"filters"`
	Ordering []map[string]string    `json:"ordering"`
}
type Response struct {
	Page     int      `json:"page"`
	PageSize int      `json:"page_size"`
	Count    int      `json:"count"`
	Results  []Report `json:"results"`
}

// Security
var AllowedColumns = map[string]bool{
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

var AllowedDirections = map[string]bool{
	"ASC":  true,
	"DESC": true,
}
