// responses - пакет, содержащий структуры ответов на Api запросы
package responses

// AssetLinksResponse - структура для ответа на запрос
// для создании AssetLinks для связи с мобильным приложением
type AssetLinksResponse []AssetLink

type AssetLink struct {
	Relation []string    `json:"relation"`
	Target   AssetTarget `json:"target"`
}

type AssetTarget struct {
	Namespace              string   `json:"namespace"`
	PackageName            string   `json:"package_name"`
	Sha256CertFingerprints []string `json:"sha256_cert_fingerprints"`
}
