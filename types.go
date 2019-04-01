package main

// FactSheet represents a Fact Sheet in the form it is stored in the json file
type FactSheet map[string]interface{}

// KubernetesNodeInfo holds meta information about a kubernetes cluster
type KubernetesNodeInfo struct {
	DataCenter       string
	AvailabilityZone string
	NumberNodes      string
	TypeNodes        []string
}

// Content holds paths for the different json files
type Content struct {
	Manifest   string `json:"manifest,omitempty"`
	FactSheets string `json:"factSheets,omitempty"`
	Relations  string `json:"relations,omitempty"`
}

// Manifest describes metadata of the archives content
type Manifest struct {
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	ImportSettings ImportSettings         `json:"importSettings,omitempty"`
}

// ImportSettings defines a set of import settings
type ImportSettings struct {
	CreateUnknownEntities           bool `json:"createUnknownEntities"`
	UpdateKnownEntities             bool `json:"updateKnownEntities"`
	ReplaceMultiValueFields         bool `json:"replaceMultiValueFields"`
	ResetUnreferencedFields         bool `json:"resetUnreferencedFields"`
	DeleteUnreferencedEntities      bool `json:"deleteUnreferencedEntities"`
	DeleteUnreferencedRelationships bool `json:"deleteUnreferencedRelationships"`
}
