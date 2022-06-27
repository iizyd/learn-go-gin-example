package upload

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"github.com/iamzhiyudong/go-gin-example/pkg/file"
	"github.com/iamzhiyudong/go-gin-example/pkg/logging"
	"github.com/iamzhiyudong/go-gin-example/pkg/setting"
	"github.com/iamzhiyudong/go-gin-example/pkg/util"
)

// 获取图片完整 url
func GetImageFullUrl(name string) string {
	return setting.AppSetting.ImagePrefixUrl + "/" + GetImagePath() + name
}

// 获取 md5 处理后的文件名
func GetImageName(name string) string {
	ext := path.Ext(name)
	fileName := strings.TrimSuffix(name, ext) // 去除文件名后缀
	fileName = util.EncodeMD5(fileName)       // 名字 md5 加密

	return fileName + ext
}

// 获取图片存储路径
func GetImagePath() string {
	return setting.AppSetting.ImageSavePath
}

// 获取图片完整的存储路径
func GetImageFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetImagePath()
}

// 检查后缀是否符合规范
func CheckImageExt(fileName string) bool {
	ext := file.GetExt(fileName)
	for _, allowExt := range setting.AppSetting.ImageAllowExts {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}

	return false
}

// 检查文件大小是否符合
func CheckImageSize(f multipart.File) bool {
	size, err := file.GetSize(f)
	if err != nil {
		log.Println(err)
		logging.Warn(err)
		return false
	}

	return size <= setting.AppSetting.ImageMaxSize
}

// 检查是否存在、权限
func CheckImage(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}

	err = file.IsNotExistMkDir(dir + "/" + src)
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err: %v", err)
	}

	perm := file.CheckPermission(src)
	if perm == true {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	return nil
}
