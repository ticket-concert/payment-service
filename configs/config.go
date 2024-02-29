package configs

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var Cfg Config

type Config struct {
	ServiceName       string           `envconfig:"service_name"`
	ServiceVersion    string           `envconfig:"service_version"`
	ServicePort       string           `envconfig:"service_port"`
	ServiceEnv        string           `envconfig:"service_env"`
	HttpServer        HttpServerConfig `envconfig:"http_server"`
	Logger            LoggerConfig     `envconfig:"logger"`
	Database          DatabaseConfig   `envconfig:"database"`
	Redis             RedisConfig      `envconfig:"redis"`
	MongoDB           MongoDBConfig    `envconfig:"mongo"`
	APMElastic        APMElasticConfig `envconfig:"apm"`
	Datadog           DatadogConfig    `envconfig:"datadog"`
	Kafka             KafkaConfig      `envconfig:"kafka"`
	Jwt               JwtConfig        `envconfig:"jwt"`
	UsernameBasicAuth string           `envconfig:"username_basic_auth"`
	PasswordBasicAuth string           `envconfig:"password_basic_auth"`
	ShutDownDelay     string           `envconfig:"shutdown_delay"`
	SecretHashPass    string           `envconfig:"secret_hash_pass"`
	IdHash            string           `envconfig:"id_hash"`
	Midtrans          MidtransConfig   `envconfig:"midtrans"`
}

type HttpServerConfig struct {
	Host string `envconfig:"http_server_host"`
	Port string `envconfig:"http_server_port"`
}

type LoggerConfig struct {
	IsVerbose       string `envconfig:"logger_is_verbose"`
	LoggerCollector string `envconfig:"logger_logger_collector"`
}

type DatabaseConfig struct {
	Host         string `envconfig:"database_host"`
	Port         string `envconfig:"database_port"`
	Username     string `envconfig:"database_username"`
	Password     string `envconfig:"database_password"`
	DBName       string `envconfig:"database_db_name"`
	SSL          string `envconfig:"database_ssl"`
	SchemaName   string `envconfig:"database_schema_name"`
	MaxIdleConns string `envconfig:"database_max_idle_conns"`
	MaxOpenConns string `envconfig:"database_max_open_conns"`
	Timeout      string `envconfig:"database_timeout"`
}

type MongoDBConfig struct {
	MongoMasterDBUrl string `envconfig:"mongo_master_database_url"`
	MongoSlaveDBUrl  string `envconfig:"mongo_slave_database_url"`
	MongoPoolSize    string `envconfig:"mongo_pool_size"`
}

type RedisConfig struct {
	RedisDB        string `envconfig:"redis_db"`
	RedisHost      string `envconfig:"redis_host"`
	RedisPort      string `envconfig:"redis_port"`
	RedisPassword  string `envconfig:"redis_password"`
	RedisAppConfig string `envconfig:"redis_app_config"`
}

type APMElasticConfig struct {
	APMUrl         string `envconfig:"apm_url"`
	APMSecretToken string `envconfig:"apm_secret_token"`
}

type DatadogConfig struct {
	DatadogHost    string `envconfig:"datadog_host"`
	DatadogPort    string `envconfig:"datadog_port"`
	DatadogEnabled string `envconfig:"datadog_enabled"`
	DatadogEnv     string `envconfig:"datadog_env"`
	DatadogService string `envconfig:"datadog_service"`
}

type KafkaConfig struct {
	KafkaUrl      string `envconfig:"kafka_url"`
	KafkaUsername string `envconfig:"kafka_username"`
	KafkaPassword string `envconfig:"kafka_password"`
}

type JwtConfig struct {
	JwtPrivateKey        string `envconfig:"private_key"`
	JwtPublicKey         string `envconfig:"public_key"`
	JwtRefreshPrivateKey string `envconfig:"private_key_refresh"`
	JwtRefreshPublicKey  string `envconfig:"public_key_refresh"`
}

type MidtransConfig struct {
	BaseUrl   string `envconfig:"midtrans_base_url"`
	BasicAuth string `envconfig:"midtrans_basic_auth"`
}

func InitConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
	err = envconfig.Process("fww_core", &Cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &Cfg
}

func GetConfig() *Config {
	return &Cfg
}
