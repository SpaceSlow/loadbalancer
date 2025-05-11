package config

type DBConfig struct {
	DBMS string `yaml:"dbms"`
	DSN  string `yaml:"dsn"`
}
