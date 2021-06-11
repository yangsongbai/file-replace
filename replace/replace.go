package replace

import (
	"flag"
	"github.com/yangsongbai/file-replace/config"
	"github.com/yangsongbai/file-replace/tools"
	"strings"
)

type App struct {
	uppercaseName  string
	lowercaseName  string
}

func NewApp(name string) *App {
	return &App{uppercaseName: strings.ToUpper(tools.TrimSpaces(name)),
		lowercaseName: strings.ToLower(tools.TrimSpaces(name))}
}

func (app App) GetAppLowercaseName() string {
	return app.lowercaseName
}

func (app App) GetUppercaseName() string {
	return app.uppercaseName
}

const (
	CONTAIN = ""
	START = "\""
	END = "\""
	REPLACE=""
	FILE_NAME_START=""
	FILE_NAME_END = "conf"
	DEFAULT_DIR = "/etc/logstash"
	//DEFAULT_DIR = "/Users/yangsongbai/workspace/github/file-replace/config"
)

func (app *App) parseFlag()*config.Config	{
	config := config.NewConfig()
	//本地文件配置
	flag.StringVar(&config.ReplaceInfo.Contain, "contain", CONTAIN, "要替换的内容包含的子字符串")
	flag.StringVar(&config.ReplaceInfo.Start, "start", START, "以什么开头")
	flag.StringVar(&config.ReplaceInfo.End, "end", END, "以什么结尾")
	flag.StringVar(&config.ReplaceInfo.Replace, "replace", REPLACE, "匹配的文件名字")

	flag.StringVar(&config.FileInfo.Contain, "file_name_contain", "", "文件名字包含的子字符串")
	flag.StringVar(&config.FileInfo.Start, "file_name_start", "", "文件名以什么开头")
	flag.StringVar(&config.FileInfo.End, "file_name_end", FILE_NAME_END, "文件名以什么结尾")
	flag.StringVar(&config.FileInfo.Dir, "dir", DEFAULT_DIR, "扫描的目录")
	flag.BoolVar(&config.FileInfo.Recursion, "recursion", true, "是否递归扫描")

	flag.Parse()
	return config
}



func (app *App) Init(setUp func()(), check func()) {
	concurrent:= app.parseFlag()
	config.SetConfig(concurrent)

}

func (app App) Start() {
	tools.ReplaceFile()

}
