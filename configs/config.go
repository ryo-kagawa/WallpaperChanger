package configs

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Input struct {
		ImagePath string `yaml:"imagePath"`
		ImageList []struct {
			X uint64 `yaml:"x"`
			Y uint64 `yaml:"y"`
			W uint64 `yaml:"w"`
			H uint64 `yaml:"h"`
		} `yaml:"imageList"`
	} `yaml:"input"`
	Output struct {
		OutputFilePath string `yaml:"outputFilePath"`
	} `yaml:"output"`
}

func LoadConfig(filePath string) (Config, error) {
	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		err = errors.New("設定ファイル" + filePath + "が見つかりません")
		return Config{}, err
	}

	config := Config{}
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		err = errors.New("設定ファイル" + filePath + "が内容が不正です")
		return Config{}, err
	}

	return config, nil
}

// TODO: バリデート作成
func (c Config) Validate() error {
	return nil
}
