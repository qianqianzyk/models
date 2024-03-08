package userController

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"image/jpeg"
	"io"
	"live-chat/app/apiException"
	"live-chat/app/configSecret"
	"live-chat/app/utils"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func UploadImage(c *gin.Context) {
	// 保存图片文件
	file, err := c.FormFile("img")
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	// 创建临时目录
	tmp, err := os.MkdirTemp("", "tmp")
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	defer func() {
		if err := os.RemoveAll(tmp); err != nil {
			fmt.Println("Error remove destination tmp:", err)
		}
	}()
	// 在临时目录中创建临时文件
	tmpFile := filepath.Join(tmp, file.Filename)
	f, err := os.Create(tmpFile)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("Error closing destination f:", err)
		}
	}()
	// 将上传的文件保存到临时文件中
	src, err := file.Open()
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	defer func() {
		if err := src.Close(); err != nil {
			fmt.Println("Error closing destination src:", err)
		}
	}()
	_, err = io.Copy(f, src)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	// 判断文件的MIME类型是否为图片
	mime, err := mimetype.DetectFile(tmpFile)
	if err != nil || !strings.HasPrefix(mime.String(), "image/") {
		utils.JsonErrorResponse(c, apiException.PictureError)
		return
	}
	// 保存原始图片
	filename := uuid.New().String() + ".jpg"
	dst := "./static/" + filename
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	// 转换图像为JPG格式并压缩
	jpgFile := filepath.Join(tmp, "compressed.jpg")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := convertAndCompressImage(dst, jpgFile); err != nil {
			fmt.Println("Error converting and compressing image:", err)
		}
	}()
	wg.Wait()
	err = os.Rename(jpgFile, dst)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	key, _ := configSecret.GetUrlKey()
	url := key + filename
	utils.JsonSuccessResponse(c, gin.H{
		"img": url,
	})
}

// 用于转换和压缩图像的函数
func convertAndCompressImage(srcPath, dstPath string) error {
	srcImg, err := imaging.Open(srcPath)
	if err != nil {
		return err
	}
	// 调整图像大小（根据需要进行调整）
	resizedImg := resize.Resize(300, 0, srcImg, resize.Lanczos3)
	// 创建新的JPG文件
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := dstFile.Close(); err != nil {
			fmt.Println("Error closing destination dstFile:", err)
		}
	}()
	// 以JPG格式保存调整大小的图像，并设置压缩质量为90
	err = jpeg.Encode(dstFile, resizedImg, &jpeg.Options{Quality: 90})
	if err != nil {
		return err
	}
	return nil
}
