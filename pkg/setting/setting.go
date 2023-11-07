package setting

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
	"webhook/pkg/util"
)

type AppSetting struct {
	Port         string
	Mode         string
	MaxHeader    int
	ReadTimeOut  time.Duration
	WriteTimeOut time.Duration
}

type PageSetting struct {
	PageSize int
	MaxPage  int
}
type RedisSetting struct {
	Address  string
	UserName string
	Pass     string
	DB       int
	Topic    string
}

type MySQLSetting struct {
	Address     string
	User        string
	Pass        string
	Db          string
	MaxIdleConn int
	MaxOpenConn int
}

type KafkaSetting struct {
	Address       string
	WebHookTopic  string
	DurationTopic string
	GeneralTopic  string
	WebHookGroup  string
	DurationGroup string
	GeneralGroup  string
}

type ServerConfig struct {
	App   AppSetting
	Page  PageSetting
	Redis RedisSetting
	Mysql MySQLSetting
	Kafka KafkaSetting
}

var (
	Svc = &ServerConfig{}
)

func LoadEnv() {
	util.SonarServer = os.Getenv("SONAR_SERVER")
	if util.SonarServer == "" {
		panic("env[SONAR_SERVER] absent")
	}
	util.SonarUsername = os.Getenv("SONAR_USERNAME")
	if util.SonarUsername == "" {
		panic("env[SONAR_USERNAME] absent")
	}
	util.SonarPassword = os.Getenv("SONAR_PASSWORD")
	if util.SonarPassword == "" {
		panic("env[SONAR_PASSWORD] absent")
	}

}

func LoadServerConfig() {
	viper.SetConfigName("server")
	viper.AddConfigPath("conf")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			log.Println("LoadServerConfig.config file not found...")
		} else {
			log.Printf("LoadServerConfig.config parse fail:%s", err)
		}
		log.Fatal(err)
	}
	if err := viper.Unmarshal(Svc); err != nil {
		log.Printf("LoadServerConfig.config file mapping ServerConfig fail:%s", err)
		log.Fatal(err)
	}

	var prettyJSON bytes.Buffer
	marshal, _ := json.Marshal(Svc)
	if err := json.Indent(&prettyJSON, marshal, "", " "); err != nil {
		log.Fatalf("json.Indent.err:%s", err.Error())
	}
	log.Printf("****LoadConfig:\n%s", prettyJSON.String())
}
