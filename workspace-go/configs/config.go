package configs

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ImagePath     string      `yaml:"imagePath"`
	RectangleList []Rectangle `yaml:"rectangleList"`
}

func LoadConfig(filePath string) (Config, error) {
	buf, err := os.ReadFile(filePath)
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
