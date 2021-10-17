package model

type Config struct {
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Domain    string `yaml:"domain"`
	Record    string `yaml:"record"`
	RegionID  string `yaml:"region_id"`
}
