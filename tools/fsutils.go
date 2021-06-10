//File fsutils.go
//@Author Songbai Yang
//@Date 2019/7/31
package tools

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)


type FileInfo struct {
	BaseDir      string       `json:"base_dir,omitempty"`
	SubDir      string       `json:"sub_dir,omitempty"`
	FileName string       `json:"file_name,omitempty"`
	FileSize int64        `json:"file_size,omitempty"`
	Err      error        `json:"error,omitempty"`
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
func GetAllFile(baseDir ,subPath,pattern,fileName,end string,recursion bool) []*FileInfo {
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
			ffinfos := GetAllFile(baseDir, next,pattern,fileName, runtime.GOOS,recursion)
			if  ffinfos != nil && len(ffinfos) > 0 {
				fileInfos = append(fileInfos, ffinfos...)
			}
			continue
		}
		if pattern != "" || fileName != "" {
			if (pattern != "" && pattern != "\"\""&& strings.Contains(fi.Name(), pattern)) || (fileName != "" && fileName != "\"\""&& fileName == fi.Name()) {
					fmt.Println(fi.Name())
					fileInfos = append(fileInfos, NewFileInfo(baseDir, subPath, fi.Name(), fi.Size(), nil))
				}
		} else {
			fileInfos = append(fileInfos, NewFileInfo(baseDir, subPath, fi.Name(), fi.Size(), nil))
		}
	}
	return fileInfos
}


func ReplaceFile(baseDir ,subPath,osSystem,contain,replace string,recursion bool) {
	files:=GetAllFile(baseDir,subPath,contain,replace,osSystem,recursion)
	for _,file:=range files {
		pathFile:=file.BaseDir;
		if file.SubDir!="" {
			pathFile+="/"+file.SubDir
		}
		pathFile+=file.FileName
		ReplaceFileContent(pathFile,contain,replace)
	}
}

func ReplaceFileContent(fileName,contain,replace string) []string {
	in, err := os.Open(fileName)
	if err != nil {
		fmt.Println("open file fail:", err)
		os.Exit(-1)
	}
	defer in.Close()
	out, err := os.OpenFile(fileName+".mdf", os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		fmt.Println("Open write file fail:", err)
		os.Exit(-1)
	}
	defer out.Close()

	br := bufio.NewReader(in)
	index := 1
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("read err:", err)
			os.Exit(-1)
		}
		source:=string(line)
		if strings.Contains(source,contain) {
			newLine := strings.Replace(source, source, replace, -1)
			_, err = out.WriteString(newLine + "\n")
			if err != nil {
				fmt.Println("write to file fail:", err)
				os.Exit(-1)
			}
		}

		fmt.Println("done ", index)
		index++
	}

	return nil
}