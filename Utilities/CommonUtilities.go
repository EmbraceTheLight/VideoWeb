package Utilities

import (
	"VideoWeb/define"
	"io"
	"mime/multipart"
	"os"
	"path"
)

// WriteToNewFile 写入新文件,src为源文件,dst为目标文件路径
func WriteToNewFile(src *multipart.FileHeader, dst string) error {
	srcFile, err := src.Open()
	defer srcFile.Close()
	if err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	defer dstFile.Close()
	if err != nil {
		return err
	}

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	return nil
}

// CheckPicExt 检查图片扩展名
func CheckPicExt(filename string) bool {
	_, ok := define.PicExtCheck[path.Ext(filename)]
	return ok
}

// CheckVideoExt 检查视频扩展名
func CheckVideoExt(filename string) bool {
	_, ok := define.VideoExtCheck[path.Ext(filename)]
	return ok
}
