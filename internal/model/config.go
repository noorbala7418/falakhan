package model

type Config struct {
	Listen        string         `yaml:"listen_port"`
	LogLevel      string         `yaml:"log_level"`
	S3Credentials []S3Credential `yaml:"s3_credentials"`
	Routes        []struct {
		Name     string `yaml:"name"`
		S3Config string `yaml:"s3_config"`
	} `yaml:"routes"`
}
