package config

type VolumeConfig struct {
	// +optional
	Size string `json:"size,omitempty" yaml:"size,omitempty"`

	// +optional
	StorageClass string `json:"storageClass,omitempty" yaml:"storageClass,omitempty"`
}
