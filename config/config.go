package config

type ReplaceInfo struct {
	Contain       string      `json:"contain,omitempty"`
	Start         string      `json:"start,omitempty"`
	End           string      `json:"end,omitempty"`
	Replace       string     `json:"replace,omitempty"`
}

type FileMatchInfo struct {
	Contain       string      `json:"contain,omitempty"`
	Start         string      `json:"start,omitempty"`
	End           string      `json:"end,omitempty"`
	Dir           string     `json:"dir,omitempty"`
	SubDir        string     `json:"sub_dir,omitempty"`
	FileSize      int64      `json:"file_size,omitempty"`
	Error         error      `json:"error,omitempty"`
	Recursion      bool      `json:"recursion,omitempty"`
}

var config *Config

func SetConfig(conf *Config)  {
	config = conf
}

func GetConfig() *Config  {
	  return  config
}


type Config struct {
	FileInfo       *FileMatchInfo      `json:"file_info,omitempty"`
	ReplaceInfo    *ReplaceInfo    `json:"replace_info,omitempty"`
}

func NewConfig() *Config {
	return &Config{FileInfo:NewFileMatchInfo(),ReplaceInfo: NewReplaceInfo()}
}

func NewReplaceInfo() *ReplaceInfo {
	return &ReplaceInfo{}
}

func NewFileMatchInfo()*FileMatchInfo {
	return &FileMatchInfo{}
}