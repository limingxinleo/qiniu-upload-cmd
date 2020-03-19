package cmd

import (
	"fmt"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"io/ioutil"
	"os"
	"strings"
)

var accessKey = ""
var secretKey = ""
var bucket = ""
var PthSep = string(os.PathSeparator)

var rootCmd = &cobra.Command{
	Use:   "push [dirName] [targetName]",
	Short: "推送 文件 到 七牛云存储",
	Long:  `推送 文件 到 七牛云存储`,
	Args:  cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		dirName := args[1]
		targetName := args[2]
		files := getFiles(dirName)
		var target string
		for _, file := range files {
			target = strings.Replace(file, dirName, targetName, 1)
			upload(file, target)
			fmt.Println(fmt.Sprintf("文件 %s 上传到 %s 对应文件名 %s", file, bucket, target))
		}
	},
}

func getFiles(dirName string) []string {
	var files []string
	dir, err := ioutil.ReadDir(dirName)
	if err != nil {
		fmt.Println(err)
		return files
	}

	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			files = append(files, getFiles(dirName+PthSep+fi.Name())...)
		} else {
			files = append(files, dirName+PthSep+fi.Name())
		}
	}

	return files
}

func upload(localFile string, target string) {
	putPolicy := storage.PutPolicy{
		Scope:      fmt.Sprintf("%s:%s", bucket, target),
		InsertOnly: 0,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuadong
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	err := formUploader.PutFile(context.Background(), &ret, upToken, target, localFile, nil)
	fmt.Println(ret)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
