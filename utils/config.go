package utils

import (
	"encoding/json"
	"io"
	"os"
)

var Config struct {
	Database struct {
		RediDB struct {
			Host     string
			Port     int
			Login    string
			Password string
			Database string
		}
		MongoDB struct {
			Endpoint   string
			AuthSource string
			Login      string
			Password   string
			Database   string
		}
		Postgres struct {
			ConnectionString string
		}
		YDB struct {
			ConnectionString string
			Token            string
			UseMetadataAuth  bool
		}
	}
}

func init() {
	f := Must(os.Open("db_config.json"))
	Must0(json.Unmarshal(Must(io.ReadAll(f)), &Config))

}
