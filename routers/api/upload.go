package api

import (
	"JRYY/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func UploadFile(g *gin.Context) {
	file, err := g.FormFile("file")
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintln(err),
		})
	}
	src, err := file.Open()
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintln(err),
		})
		return
	}
	defer src.Close()
	t := time.Now()
	nameList := strings.Split(file.Filename, ".")
	if len(nameList) != 2 {
		g.String(http.StatusBadRequest, "请检查文件名是否有非法字符，Error Code：201")
		return
	}
	fileName := string(nameList[0] + "-" + t.Format("20060102150405") + "." + nameList[1])
	//fmt.Println(fileName)
	dst, err := os.Create(utils.BASE_PATH + "file/" + fileName)
	//dst, err := os.Create("/Users/engine/Dropbox/Codes/go/src/EEB/upload/" + fileName)
	//fmt.Println(dst)
	if err != nil {
		g.String(http.StatusBadRequest, "上传文件失败，Error Code：202")
		return
	}
	defer dst.Close()
	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		g.String(http.StatusBadRequest, "上传文件失败，Error Code：203")
		return
	}
	g.String(http.StatusOK, fileName)
	return
}

func UploadImage(g *gin.Context) {
	file, err := g.FormFile("file")
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintln(err),
		})
	}
	src, err := file.Open()
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintln(err),
		})
		return
	}
	defer src.Close()
	t := time.Now()
	nameList := strings.Split(file.Filename, ".")
	if len(nameList) != 2 {
		g.String(http.StatusBadRequest, "请检查文件名是否有非法字符，Error Code：201")
		return
	}
	fileName := string(nameList[0] + "-" + t.Format("20060102150405") + "." + nameList[1])
	//fmt.Println(fileName)
	dst, err := os.Create(utils.BASE_PATH + "image/" + fileName)
	//dst, err := os.Create("/Users/engine/Dropbox/Codes/go/src/EEB/upload/" + fileName)
	//fmt.Println(dst)
	if err != nil {
		g.String(http.StatusBadRequest, "上传图片失败，Error Code：202")
		return
	}
	defer dst.Close()
	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		g.String(http.StatusBadRequest, "上传图片失败，Error Code：203")
		return
	}
	g.String(http.StatusOK, fileName)
	return
}
