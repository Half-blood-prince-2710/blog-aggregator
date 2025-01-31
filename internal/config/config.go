package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DbUrl string `json:"db_url"`
	Username string `json:"current_user_name"`
}
const (
	configFileName = ".gatorconfig.json"
)
func getConfigFilePath() (string , error) {
		val ,err:=os.UserHomeDir()
	if err!=nil {
		return "" , err
	}
	
	filePath := filepath.Join(val,configFileName)
	return filePath,nil
}

func Read() (Config,error){
	var cfg Config

	filePath, err := getConfigFilePath()
	// fmt.Print("file: ",filePath)
	data,err:=os.ReadFile(filePath)
	if err !=nil {
		return cfg,err
	}
	
	// fmt.Print("val",val)
	
	err =json.Unmarshal(data,&cfg)
	if err!=nil {
		return cfg ,err
	}


	return cfg,nil
}

func (cfg *Config) SetUser(name string) (error) {
	cfg.Username = name

	data , err := json.Marshal(cfg)
	if err!=nil {
		return err
	}
	str, err:=getConfigFilePath()
	err = os.WriteFile(str,data,0664)
	if err!=nil {
		return err
	}
	return nil
}