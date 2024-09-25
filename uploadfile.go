package uploadfilesforwin

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/delong666888/uploadFilesGo/config"
	"github.com/fsnotify/fsnotify"
)

func uploadFile(filename string, apiUrl string, payload *config.Payload) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	/*// 创建 JSON 部分
	jsonData := config.Payload{
		Token:   config.TOKEN,
		FileTxt: "123123123",
	}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %v", err)
	}

	// 创建 JSON 部分
	jsonPart, err := writer.CreateFormField("data")
	if err != nil {
		return fmt.Errorf("error creating form field for JSON: %v", err)
	}

	_, err = jsonPart.Write(jsonBytes)
	if err != nil {
		return fmt.Errorf("error writing JSON data: %v", err)
	}*/
	err = writer.WriteField("token", payload.Token)
	if err != nil {
		return fmt.Errorf("error writing token field: %v", err)
	}

	err = writer.WriteField("file_txt", payload.FileTxt)
	if err != nil {
		return fmt.Errorf("error writing file_txt field: %v", err)
	}
	// 创建文件部分
	filePart, err := writer.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		return fmt.Errorf("error creating form file: %v", err)
	}

	_, err = io.Copy(filePart, file)
	if err != nil {
		return fmt.Errorf("error copying file: %v", err)
	}

	// 关闭 multipart writer
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("error closing multipart writer: %v", err)
	}

	// 发送 API 请求
	req, err := http.NewRequest("POST", apiUrl, body)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		resBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %v", err)
		}
		return fmt.Errorf("response Status: %s, Response Status Code: %d, Response Message: %s", resp.Status, resp.StatusCode, string(resBody))
	}

	return nil
}

// 移动文件到其他目录
func moveFile(source, targetDir string) error {
	filename := filepath.Base(source)

	targetFile := filepath.Join(targetDir, filename)
	return os.Rename(source, targetFile)
}

// 检查如果已经有存在的文件，则上传并移动文件到目标目录
func CheckExistingFiles(uploadConfig *config.UploadConfig) {
	err := filepath.Walk(uploadConfig.SourceDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		//处理文件，忽略目录
		if !info.IsDir() {
			// upload file
			fmt.Println("CheckExistingFiles - Uploading:", path)
			err := uploadFile(path, uploadConfig.ApiUrl, uploadConfig.Payload)
			if err != nil {
				return fmt.Errorf("CheckExistingFiles - Upload faild: %v", err)
			}
			fmt.Println("Upload Sucessful for", path)
			err = moveFile(path, uploadConfig.TargetDir)
			if err != nil {
				return fmt.Errorf("CheckExistingFiles - error moving file: %v", err)
			}
			fmt.Println("moved Sucessful for", path)
		}

		return nil
	})
	if err != nil {
		log.Println("Error scanning existing files:", err)
	}
}

// 通过插件fsnotify监听来源目录，如果有创建的文件，则上传文件和移动文件到新目录
func WatchDirectory(uploadConfig *config.UploadConfig) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// 处理创建文件事件
				if event.Op&fsnotify.Create == fsnotify.Create {
					fmt.Println("New File detected:", event.Name)
					if strings.HasSuffix(event.Name, "jpg") || strings.HasSuffix(event.Name, ".png") {
						fmt.Println("Uploading:", event.Name)
						err := uploadFile(event.Name, uploadConfig.ApiUrl, uploadConfig.Payload)
						if err != nil {
							log.Println("Upload faild:", err)
							return
						}

						fmt.Println("Upload Sucessful for", event.Name)
						err = moveFile(event.Name, uploadConfig.TargetDir)
						if err != nil {
							log.Println("error moving file:", err)
						}
						fmt.Println("moved Sucessful for", event.Name)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Error:", err)
			}
		}
	}()

	err = watcher.Add(uploadConfig.SourceDir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
