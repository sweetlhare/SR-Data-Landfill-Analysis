package cdnclient

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"svalka-service/internal/config"
	logicentities "svalka-service/internal/logic/entities"
	logicinterfaces "svalka-service/internal/logic/interfaces"
	"svalka-service/pkg/closer"
	"time"

	"github.com/jlaffaye/ftp"
)

type cdnClient struct {
	client *ftp.ServerConn
	config config.FtpClientConfig
}

// NewCdnClient ...
func NewCdnClient(_ context.Context) (logicinterfaces.CdnClient, error) {
	config := getFtpClientConfig()
	client, err := ftp.Connect(config.FtpAddr())
	if err != nil {
		return nil, err
	}

	err = client.Login(config.FtpUSER(), config.FtpPassword())
	if err != nil {
		return nil, err
	}
	cdnClient := cdnClient{
		client: client,
		config: config,
	}

	closer.Add(cdnClient.Close)

	return cdnClient, nil
}

// Close ...
func (c cdnClient) Close() error {
	return c.client.Quit()
}

// SaveImage ...
func (c cdnClient) SaveImages(images ...logicentities.File) (result *logicentities.CdnResult, err error) {
	result = &logicentities.CdnResult{}
	for _, image := range images {
		filePath, err := c.UploadFileToFTP(image)
		if err != nil {
			return nil, err
		}
		result.RawImagesPaths = append(result.RawImagesPaths, filePath)
	}
	return result, nil
}

// UploadFileToFTP ...
func (c cdnClient) UploadFileToFTP(image logicentities.File) (filePath string, err error) {

	file, err := image.Open()
	if err != nil {
		return "", err
	}

	defer file.Close()

	imageType := strings.TrimPrefix(image.Header.Get("Content-Type"), "image/")
	filePath = fmt.Sprintf("%s/%s.%s", c.config.FtpPATH(), GenerateUniqueFileName(), imageType)
	err = c.client.Stor(filePath, file)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

// GenerateUniqueFileName ...
func GenerateUniqueFileName() string {
	rand.Seed(time.Now().UnixNano())

	randomNumber := rand.Intn(900000) + 100000

	timestamp := time.Now().UnixNano()

	uniqueFileName := fmt.Sprintf("%d_%d", timestamp, randomNumber)

	return uniqueFileName
}
