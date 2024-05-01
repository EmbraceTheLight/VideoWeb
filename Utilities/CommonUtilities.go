package Utilities

import (
	"VideoWeb/define"
	"errors"
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
func CheckPicExt(filename string) error {
	if _, ok := define.PicExtCheck[path.Ext(filename)]; !ok {
		return errors.New("图片格式错误或不支持此图片格式")
	}
	return nil
}

// CheckVideoExt 检查视频扩展名
func CheckVideoExt(filename string) error {
	if _, ok := define.VideoExtCheck[path.Ext(filename)]; !ok {
		return errors.New("视频格式错误或不支持此图片格式")
	}
	return nil
}

func ReadFileContent(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	ret, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
