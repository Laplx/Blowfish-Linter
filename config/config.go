package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	TitleMode   string
	DateMode    string
	TagsMode    []string
	Authors     []string
	Category    string
	DraftMode   string
	InlineKatex string
	MultiKatex  string
}

var Cfg Config

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$USERPROFILE/.config/bflint")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic("配置文件读取失败: " + err.Error())
	}
	Cfg.TitleMode = viper.GetString("Front Matter.title")
	Cfg.DateMode = viper.GetString("Front Matter.date")
	Cfg.TagsMode = viper.GetStringSlice("Front Matter.tags")
	Cfg.Authors = viper.GetStringSlice("Front Matter.authors")
	Cfg.Category = viper.GetString("Front Matter.category")
	Cfg.DraftMode = viper.GetString("Front Matter.draft")
	Cfg.InlineKatex = viper.GetString("Katex.inline")
	Cfg.MultiKatex = viper.GetString("Katex.multiline")
}
