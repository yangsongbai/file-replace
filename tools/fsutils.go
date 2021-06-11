//File fsutils.go
//@Author Songbai Yang
//@Date 2019/7/31
package tools

import (
	"bufio"
	"fmt"
	"github.com/yangsongbai/file-replace/config"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)


type FileInfo struct {
	BaseDir      string       `json:"base_dir,omitempty"`
	SubDir       string       `json:"sub_dir,omitempty"`
	FileName     string       `json:"file_name,omitempty"`
	FileSize     int64        `json:"file_size,omitempty"`
	Err          error        `json:"error,omitempty"`
}

func NewFileInfo(BaseDir,SubDir,FileName string,FileSize int64,Err error) *FileInfo {
	return &FileInfo{BaseDir:BaseDir,SubDir:SubDir,FileName:FileName,FileSize:FileSize,Err:Err}
}


// FileExists check if the path are exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// FilePutContent put string to file
func FilePutContent(file string, content string) (int, error) {
	fs, e := os.Create(file)
	if e != nil {
		return 0, e
	}
	defer fs.Close()
	return fs.WriteString(content)
}

// JoinPath return joined file path
func JoinPath(filenames ...string) string {

	hasSlash := false
	result := ""
	for _, str := range filenames {
		currentHasSlash := false
		if len(result) > 0 {
			currentHasSlash = strings.HasPrefix(str, "/")
			if hasSlash && currentHasSlash {
				str = strings.TrimLeft(str, "/")
			}
			if !(hasSlash || currentHasSlash) {
				str = "/" + str
			}
		}
		hasSlash = strings.HasSuffix(str, "/")
		result += str
	}
	return result
}

/**
 递归获取该目录下的所有文件及其文件大小
 */
func GetAllFile(baseDir ,subPath,end string,recursion bool) []*FileInfo {
	path:=baseDir
	if subPath !="" {
		path += GetSystemDirSplit(runtime.GOOS)+subPath
	}
	rd, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Errorf("%s",err)
		return nil
	}
	var (
		fileInfos []*FileInfo
	)
	if len(rd)>0 {
		fileInfos = make([]*FileInfo, 0)
	}
	for _, fi := range rd {
		if "." == fi.Name() || ".." == fi.Name(){continue}
		if recursion && fi.IsDir() {
			next := fi.Name()
			if subPath !="" {
				next = subPath+GetSystemDirSplit(runtime.GOOS)+next
			}
			ffinfos := GetAllFile(baseDir, next, end,recursion)
			if  ffinfos != nil && len(ffinfos) > 0 {
				fileInfos = append(fileInfos, ffinfos...)
			}
			continue
		}
		if end!="" && strings.HasSuffix(fi.Name(),end){
			fileInfos = append(fileInfos, NewFileInfo(baseDir, subPath, fi.Name(), fi.Size(), nil))
			continue
		}
	}
	return fileInfos
}


func ReplaceFile() {
	config:=config.GetConfig()
	replace:=config.ReplaceInfo.Replace
	files:=GetAllFile(config.FileInfo.Dir,config.FileInfo.SubDir,config.FileInfo.End,config.FileInfo.Recursion)
	for _,file:=range files {
		pathFile:=file.BaseDir;
		if file.SubDir!="" {
			pathFile+="/"+file.SubDir
		}
		pathFile+="/"+file.FileName
		ReplaceFileContent(pathFile,replace,config.ReplaceInfo)
	}
}

func ReplaceFileContent(pathFile ,replace string,replaceInfo *config.ReplaceInfo  ) []string {
	fmt.Println(replaceInfo.Contain)
	fmt.Println(pathFile)
	fmt.Println(replace)
	lineCount := getFileline(pathFile)
	 b := getFileContent(pathFile, replace, replaceInfo,lineCount)
	WriteContentToFile(pathFile, b)
	return nil
}

func WriteContentToFile(pathFile string, b strings.Builder) {
	in, err := os.Open(pathFile)
	if err != nil {
		fmt.Println("open file fail:", err)
		os.Exit(-1)
	}
	defer in.Close()
	out, err := os.OpenFile(pathFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
	if err != nil {
		fmt.Println("Open write file fail:", err)
		os.Exit(-1)
	}
	defer out.Close()
	_, err = out.WriteString(b.String())
	if err != nil {
		fmt.Println("write to file fail:", err)
		os.Exit(-1)
	}
}

func getFileContent(pathFile string, replace string, replaceInfo *config.ReplaceInfo, lineCount int) (strings.Builder) {
	in, err := os.Open(pathFile)
	if err != nil {
		fmt.Println("open file fail:", err)
		os.Exit(-1)
	}
	defer in.Close()
	out, err := os.OpenFile(pathFile, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		fmt.Println("Open write file fail:", err)
		os.Exit(-1)
	}
	defer out.Close()
	br := bufio.NewReader(in)
	var b strings.Builder
	index := 1
	for index <= lineCount {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("read err:", err)
			os.Exit(-1)
		}
		source := string(line)
		newLine := source
		if strings.Contains(source, replaceInfo.Contain) && strings.Contains(source, "bootstrap_servers") {
			start := strings.Index(newLine, replaceInfo.Start) + 1
			end := strings.LastIndex(newLine, replaceInfo.End)
			old := string([]byte(newLine)[start:end])
			newLine = strings.Replace(source, old, replace, -1)
		}
		b.Write([]byte(newLine + "\n"))
		fmt.Println("done ", index)
		index++
	}
	fmt.Println(b.String())
	return  b
}

func getFileline(pathFile string) int {
	in, err := os.Open(pathFile)
	if err != nil {
		fmt.Println("open file fail:", err)
		os.Exit(-1)
	}
	fileScanner := bufio.NewScanner(in)
	lineCount := 0
	for fileScanner.Scan() {
		lineCount++
	}
	fmt.Println("number of lines:", lineCount)
	defer in.Close()
	return lineCount
}