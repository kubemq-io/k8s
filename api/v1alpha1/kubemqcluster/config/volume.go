package config

type VolumeConfig struct {
	// +optional
	Size string `json:"size,omitempty"`

	// +optional
	StorageClass string `json:"storageClass,omitempty"`
}
