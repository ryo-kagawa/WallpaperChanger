package configs

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ImagePath     string `yaml:"imagePath"`
	RectangleList []struct {
		X      uint64 `yaml:"x"`
		Y      uint64 `yaml:"y"`
		Width  uint64 `yaml:"width"`
		Height uint64 `yaml:"height"`
	} `yaml:"rectangleList"`
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
