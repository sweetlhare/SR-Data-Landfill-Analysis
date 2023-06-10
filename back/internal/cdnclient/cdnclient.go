package cdnclient

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"svalka-service/internal/config"
	logicentities "svalka-service/internal/logic/entities"
	logicinterfaces "svalka-service/internal/logic/interfaces"
	"time"

	"github.com/jlaffaye/ftp"

	"github.com/nfnt/resize"
)

type cdnClient struct {
	config config.FtpClientConfig
}

// NewCdnClient ...
func NewCdnClient(_ context.Context) (logicinterfaces.CdnClient, error) {
	config := getFtpClientConfig()

	cdnClient := cdnClient{
		config: config,
	}

	return cdnClient, nil
}

// SaveImage ...
func (c cdnClient) SaveImages(maxSize int, images ...logicentities.File) (result *logicentities.CdnResult, err error) {
	result = &logicentities.CdnResult{}
	for _, image := range images {
		filePath, err := c.UploadFileToFTP(image, maxSize)
		if err != nil {
			return nil, err
		}
		result.RawImagesPaths = append(result.RawImagesPaths, filePath)
	}
	return result, nil
}

// UploadFileToFTP ...
func (c cdnClient) UploadFileToFTP(image logicentities.File, maxSize int) (filePath string, err error) {

	imageType := strings.TrimPrefix(image.Header.Get("Content-Type"), "image/")
	filePath = fmt.Sprintf("%s/%s.%s", c.config.FtpPATH(), GenerateUniqueFileName(), imageType)

	// err = SaveToLocalStorage(filePath, image)
	// if err != nil {
	// 	return "", err
	// }

	err = c.SaveToCDN(filePath, image, maxSize)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
func (c cdnClient) SaveToCDN(filePath string, image logicentities.File, maxSize int) error {
	client, err := ftp.Connect(c.config.FtpAddr())
	if err != nil {
		return err
	}
	err = client.Login(c.config.FtpUSER(), c.config.FtpPassword())
	if err != nil {
		return err
	}

	file, err := image.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	// newFile, err := Compress(file, maxSize)
	// if err != nil {
	// 	return err
	// }
	err = client.Stor(filePath, file)
	if err != nil {
		return err
	}
	client.Quit()
	return nil
}

func Compress(file io.Reader, maxSize int) (res *bytes.Buffer, err error) {
	img, imgType, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	if width > height {
		width = maxSize
		height = 0
	} else {
		width = 0
	}

	resizedImg := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

	extension := strings.TrimPrefix(imgType, "image/")

	buffer := &bytes.Buffer{}
	_, err = io.Copy(buffer, file)
	if err != nil {
		return nil, err
	}

	switch extension {
	case "jpeg":
		jpeg.Encode(buffer, resizedImg, nil)
	case "png":
		png.Encode(buffer, resizedImg)
	}
	return buffer, nil
}

func SaveToLocalStorage(filePath string, image logicentities.File) error {
	file, err := image.Open()
	if err != nil {
		return err
	}

	defer file.Close()

	fileNew, err := os.Create(filePath)
	if err != nil {
		// Проверяем, что ошибка не связана с отсутствием директории
		if os.IsNotExist(err) {
			dirPath := filepath.Dir(filePath)
			err := os.MkdirAll(dirPath, fs.ModeExclusive) // Создаем директорию рекурсивно
			if err != nil {
				return fmt.Errorf("Ошибка при создании директории:%s", err.Error())
			}

			// Пытаемся создать файл снова
			file, err = os.Create(filePath)
			if err != nil {
				return fmt.Errorf("Ошибка при создании файла:%s", err.Error())
			}
		} else {
			return fmt.Errorf("Ошибка при создании файла:%s", err.Error())
		}
	}
	defer file.Close()

	fmt.Println("Файл успешно создан:", filePath)

	// Копирование содержимого ответа в файл
	_, err = io.Copy(fileNew, file)
	if err != nil {
		return fmt.Errorf("Ошибка при сохранении файла:%s", err.Error())
	}

	return nil
}

// GenerateUniqueFileName ...
func GenerateUniqueFileName() string {
	rand.Seed(time.Now().UnixNano())

	randomNumber := rand.Intn(900000) + 100000

	timestamp := time.Now().UnixNano()

	uniqueFileName := fmt.Sprintf("%d_%d", timestamp, randomNumber)

	return uniqueFileName
}
