//Created by Goland
//@User: lenora
//@Date: 2021/3/8
//@Time: 10:56 上午
package conf

type ServerConfig struct {
	Debug          bool               `yaml:"debug"`
	JWTSecret      string             `required:"true" yaml:"jwt_secret"`
	BaseUrl        string             `yaml:"base_url"`
	DbConfig       *DBConfig          `required:"true" yaml:"mysql"`
}

type DBConfig struct {
	User     string `default:"root" yaml:"user"`
	Password string `default:"" yaml:"password"`
	Name     string `yaml:"ip"`
	Port     uint   `default:"3306" yaml:"port"`
	DbName   string `required:"true" yaml:"db_name"`
	Charset  string `default:"utf8" yaml:"charset"`
	MaxIdle  int    `default:"10" yaml:"max_idle"`
	MaxOpen  int    `default:"50" yaml:"max_open"`
	LogMode  bool   `yaml:"log_mode"`
	Loc      string `required:"true" yaml:"loc"`
}
