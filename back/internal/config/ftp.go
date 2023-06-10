package config

import (
	"os"

	"github.com/pkg/errors"
)

var _ FtpClientConfig = (*ftpClientConfig)(nil)

const (
	ftpAddrEnv     = "FTP_ADDR"
	ftpPasswordEnv = "FTP_PASSWORD"
	ftpUSEREnv     = "FTP_USER"
	ftpPATHEnv     = "FTP_PATH"
)

type FtpClientConfig interface {
	FtpAddr() string
	FtpPassword() string
	FtpUSER() string
	FtpPATH() string
}

type ftpClientConfig struct {
	ftpAddr     string
	ftpPassword string
	ftpUSER     string
	ftpPATH     string
}

func NewFtpClientConfig() (*ftpClientConfig, error) {
	ftpAddr := os.Getenv(ftpAddrEnv)
	if ftpAddr == "" {
		return nil, errors.New("ftp client: addr not found")
	}

	ftpPassword := os.Getenv(ftpPasswordEnv)
	if ftpPassword == "" {
		return nil, errors.New("ftp client: password not found")
	}

	ftpUSER := os.Getenv(ftpUSEREnv)
	if ftpUSER == "" {
		return nil, errors.New("ftp client: user not found")
	}

	ftpPATH := os.Getenv(ftpPATHEnv)
	if ftpUSER == "" {
		return nil, errors.New("ftp client: path not found")
	}

	return &ftpClientConfig{
		ftpAddr:     ftpAddr,
		ftpPassword: ftpPassword,
		ftpUSER:     ftpUSER,
		ftpPATH:     ftpPATH,
	}, nil
}

// ftpAddr ...
func (cfg *ftpClientConfig) FtpAddr() string {
	return cfg.ftpAddr
}

// ftpPassword ...
func (cfg *ftpClientConfig) FtpPassword() string {
	return cfg.ftpPassword
}

// ftpUSER ...
func (cfg *ftpClientConfig) FtpUSER() string {
	return cfg.ftpUSER
}

// ftpPATH ...
func (cfg *ftpClientConfig) FtpPATH() string {
	return cfg.ftpPATH
}
