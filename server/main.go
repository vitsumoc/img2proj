package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/api/imgs", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": getImages(),
		})
	})
	r.Run("0.0.0.0:9090") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

type projImg struct {
	Name     string    `json:"name"`     // 图片名称
	DirName  string    `json:"dirName"`  // 图片所在文件夹名称
	DirPath  string    `json:"dirPath"`  // 图片所在文件夹完整路径
	ImgPath  string    `json:"imgPath"`  // 图片相对文件夹路径
	UpdateAt time.Time `json:"updateAt"` // 图片更新时间
}

func getImages() []projImg {
	// 获取当前目录
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("获取当前目录失败:", err)
		return nil
	}
	fmt.Println("当前目录:", dir)

	// 记录自己所在的文件夹名称，进入上一级目录，上级目录中所有的文件夹形成列表（不包括自己）
	// 获取当前目录的文件夹名称
	currentDir := filepath.Base(dir)

	// 获取上级目录路径
	parentDir := filepath.Dir(dir)

	// 读取上级目录中的所有文件和文件夹
	entries, err := os.ReadDir(parentDir)
	if err != nil {
		fmt.Println("读取上级目录失败:", err)
		return nil
	}

	// 遍历并记录文件夹(排除当前文件夹)
	var folders []string
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != currentDir {
			folders = append(folders, entry.Name())
		}
	}

	// 文件夹名称与图片列表的映射
	folder2Imgs := make(map[string][]projImg)

	// 遍历每个文件夹
	for _, folder := range folders {
		// 构建文件夹完整路径
		folderPath := filepath.Join(parentDir, folder)

		// 递归遍历文件夹获取图片
		images := getImagesFromDir(folderPath, folder, folderPath)
		folder2Imgs[folder] = images
	}

	// 选择 folder2Imgs 中每个文件夹中 updateAt 最晚的图片，返回图片列表
	var latestImgs []projImg
	for _, imgs := range folder2Imgs {
		if len(imgs) == 0 {
			continue
		}
		latestImg := imgs[0]
		for _, img := range imgs {
			if img.UpdateAt.After(latestImg.UpdateAt) {
				latestImg = img
			}
		}
		latestImgs = append(latestImgs, latestImg)
	}

	return latestImgs
}

// 递归遍历文件夹获取图片
func getImagesFromDir(folderRaw string, baseFolder string, dirPath string) []projImg {
	var images []projImg

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("读取文件夹 %s 失败: %v\n", dirPath, err)
		return images
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// 递归遍历子文件夹
			subPath := filepath.Join(dirPath, entry.Name())
			subImages := getImagesFromDir(folderRaw, baseFolder, subPath)
			images = append(images, subImages...)
		} else {
			// 检查文件扩展名是否为图片
			ext := filepath.Ext(entry.Name())
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
				// 获取基础文件夹的完整路径
				baseFolderPath := filepath.Join(filepath.Dir(dirPath), baseFolder)
				// 获取图片相对于folderRaw的路径
				relPath, _ := filepath.Rel(folderRaw, filepath.Join(dirPath, entry.Name()))

				info, _ := entry.Info()
				img := projImg{
					Name:     entry.Name(),
					DirName:  baseFolder,
					DirPath:  baseFolderPath,
					ImgPath:  relPath,
					UpdateAt: info.ModTime(),
				}
				images = append(images, img)
			}
		}
	}

	return images
}
