// config/opentsdb.go
package config

// OpentsdbConfig holds configuration for OpenTSDB
type OpentsdbConfig struct {
	Host string
	Port int
}

// NewOpentsdbConfig creates a new instance of OpentsdbConfig with default values
func NewOpentsdbConfig() *OpentsdbConfig {
	return &OpentsdbConfig{
		Host: "localhost",
		Port: 4242,
	}
}
