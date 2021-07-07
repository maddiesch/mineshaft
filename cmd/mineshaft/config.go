package main

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config provides the configuration for the mineshaft application
type Config struct {
	Machine struct {
		JavaCommand   string `toml:"java_command" validate:"required,file"`
		WorkingDir    string `toml:"working_directory" validate:"required,startswith=/"`
		PreAllocation string `toml:"pre_allocation" validate:"required"`
		MaxAllocation string `toml:"max_allocation" validate:"required"`
	} `toml:"machine" validate:"required"`
	Server struct {
		FileName  string `toml:"jar_file_name" validate:"required,endswith=.jar"`
		SourceURL string `toml:"jar_file_url" validate:"required,url"`
		Port      string `toml:"port" validate:"required,numeric"`
	} `toml:"server" validate:"required"`
}

func (c Config) JarPath() string {
	return filepath.Join(c.Machine.WorkingDir, c.Server.FileName)
}

type cfgContext struct {
	Pwd string
}

func loadConfig(filePath, workingDir string) (Config, error) {
	var config Config

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return config, err
	}

	ctx := cfgContext{
		Pwd: workingDir,
	}
	tmpl, err := template.New("config").Parse(string(content))
	if err != nil {
		return config, err
	}

	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, ctx); err != nil {
		return config, err
	}

	if err := toml.Unmarshal(buffer.Bytes(), &config); err != nil {
		return config, err
	}

	if config.Server.Port == "" {
		config.Server.Port = "25565"
	}

	if err := Validate.Struct(config); err != nil {
		return config, err
	}

	return config, nil
}

func fileExists(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		panic("Schrodinger's File!")
	}
}
