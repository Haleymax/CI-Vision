package autoload

type SystemConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Language string `mapstructure:"language"`
}
