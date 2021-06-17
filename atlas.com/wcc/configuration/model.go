package configuration

type Model struct {
	TimeoutTaskInterval int64 `yaml:"timeoutTaskInterval"`
	TimeoutDuration     int64 `yaml:"timeoutDuration"`
}
