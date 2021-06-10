package main

import "github.com/yangsongbai/file-replace/replace"

func main() {
	app := replace.NewApp("smart_upload")
	app.Init(nil, nil)
	app.Start()
}
