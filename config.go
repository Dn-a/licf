package main

// Initialize Configuration
type appApp struct {
	Name string `yaml:"name"`
}

type appLogo struct {
	Show bool `yaml:"show"`
}

type appFile struct {
	Recursion       bool     `yaml:"recursion"`
	Print           bool     `yaml:"print"`
	SupportedFormat []string `yaml:"supported-format"`
	Ignore          []string `yaml:"ignore"`
	Only            []string `yaml:"only"`
}

type appConfig struct {
	App  appApp  `yaml:"app"`
	Logo appLogo `yaml:"logo"`
	File appFile `yaml:"file"`
}

func (f *appFile) CheckIgnore(str string) bool {
	for _, v := range f.Ignore {
		if v == str {
			return true
		}
	}
	return false
}

func (f *appFile) CheckFormat(str string) bool {
	for _, v := range f.SupportedFormat {
		if v == str {
			return true
		}
	}
	return false
}

func (f *appFile) CheckOnly(str string) bool {
	if len(f.Only) == 0 {
		return true
	}
	for _, v := range f.Only {
		if v == str {
			return true
		}
	}
	return false
}

var config *appConfig

// Initialize with singleton pattern
func initConfig() {
	if config == nil {
		config = new(appConfig)
		config.File.SupportedFormat = []string{"json"}
		unmarshalFile("config.yaml", &config)
	}
}

func GetConfig() *appConfig {
	initConfig()
	return config
}
