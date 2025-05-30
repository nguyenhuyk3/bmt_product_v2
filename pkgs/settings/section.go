package settings

type Config struct {
	Server         serverSetting
	ServiceSetting serviceSetting
}

type serviceSetting struct {
	PostgreSql   postgreSetting `mapstructure:"database"`
	KafkaSetting kafkaSetting   `mapstructure:"kafka"`
	RedisSetting redisSetting   `mapstructure:"redis"`
	S3Setting    s3Setting      `mapstructure:"s3"`
}

type serverSetting struct {
	ServerPort    string `mapstructure:"SERVER_PORT"`
	RPCServerPort string `mapstructure:"RPC_SERVER_PORT"`
	XInternalCall string `mapstructure:"X_INTERNAL_CALL"`
}

type postgreSetting struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username,omitempty"`
	Password        string `mapstructure:"password,omitempty"`
	DbName          string `mapstructure:"db_name"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

type kafkaSetting struct {
	KafkaBroker_1 string `mapstructure:"kafka_broker_1"`
	KafkaBroker_2 string `mapstructure:"kafka_broker_2"`
	KafkaBroker_3 string `mapstructure:"kafka_broker_3"`
}

type redisSetting struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username,omitempty"`
	Password string `mapstructure:"password,omitempty"`
	Database int    `mapstructure:"database,omitempty"`
}

type s3Setting struct {
	AwsAccessKeyId       string `mapstructure:"aws_access_key_id"`
	AwsSercetAccessKeyId string `mapstructure:"aws_sercet_access_key_id"`
	AwsRegion            string `mapstructure:"aws_region"`
	FilmBucketName       string `mapstructure:"film_bucket_name"`
}
