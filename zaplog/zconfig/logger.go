package zconfig

type LoggerConfig struct {
	Name  string   `json:"name" yaml:"name"`
	Cores []string `json:"devices" yaml:"devices"`
}
