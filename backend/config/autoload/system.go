package autoload

type SystemConfig struct {
	Port     int    `mapstructure:"port"`
	Language string `mapstructure:"language"`
}
