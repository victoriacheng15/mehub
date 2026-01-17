package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type NavItem struct {
	Href string `yaml:"href"`
	Text string `yaml:"text"`
}

type NavigationConfig struct {
	Header []NavItem `yaml:"header"`
	Footer []NavItem `yaml:"footer"`
}

type Social struct {
	Name string `yaml:"name"`
	Href string `yaml:"href"`
	Icon string `yaml:"icon"`
}

type Project struct {
	Title            string `yaml:"title"`
	ShortDescription string `yaml:"shortDescription"`
	Link             string `yaml:"link"`
	Techs            string `yaml:"techs"`
}

type Skill struct {
	Name string `yaml:"name"`
	Icon string `yaml:"icon"`
}

type NowCategory struct {
	Title string `yaml:"title"`
	Items string `yaml:"items"`
}

type AboutConfig struct {
	Image      string `yaml:"image"`
	Paragraphs string `yaml:"paragraphs"`
}

type NowConfig struct {
	LastUpdated string        `yaml:"lastUpdated"`
	Categories  []NowCategory `yaml:"categories"`
}

type SiteMetadata struct {
	URL         string      `yaml:"url"`
	Title       string      `yaml:"title"`
	Name        string      `yaml:"name"`
	Slogan      string      `yaml:"slogan"`
	Description string      `yaml:"description"`
	About       AboutConfig `yaml:"about"`
	Now         NowConfig   `yaml:"now"`
}

type SiteConfig struct {
	Site       SiteMetadata     `yaml:"site"`
	Navigation NavigationConfig `yaml:"navigation"`
	Socials    []Social         `yaml:"socials"`
	Projects   []Project        `yaml:"projects"`
	Skills     []Skill          `yaml:"skills"`
}

func LoadConfig(configDir string) (*SiteConfig, error) {
	var config SiteConfig

	files := []string{
		"site.yaml",
		"navigation.yaml",
		"socials.yaml",
		"projects.yaml",
		"skills.yaml",
	}

	for _, file := range files {
		path := filepath.Join(configDir, file)
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, err
		}
	}

	return &config, nil
}
