package storage

import (
	"fmt"
	"github.com/donknap/dpanel/common/function"
	"github.com/donknap/dpanel/common/service/acme"
	"github.com/we7coreteam/w7-rangine-go/v2/pkg/support/facade"
	"log/slog"
	"os"
	"path/filepath"
)

type Local struct {
}

func (self Local) Delete(name string) error {
	err := os.Remove(self.GetRealPath(name))
	return err
}

func (self Local) GetSaveRootPath() string {
	return filepath.Join(self.GetStorageLocalPath(), "storage")
}

func (self Local) GetRealPath(name string) string {
	return filepath.Join(self.GetStorageLocalPath(), "storage", name)
}

func (self Local) GetStorageCertPath() string {
	return filepath.Join(self.GetStorageLocalPath(), "cert")
}

func (self Local) GetComposePath() string {
	return filepath.Join(self.GetStorageLocalPath(), "compose")
}

func (self Local) GetStorePath() string {
	return filepath.Join(self.GetStorageLocalPath(), "store")
}

func (self Local) GetLicenseFilePath() string {
	return filepath.Join(self.GetStorageLocalPath(), "dpanel.lic")
}

func (self Local) GetSshKnownHostsPath() string {
	return filepath.Join(self.GetStorageLocalPath(), "known_hosts")
}

func (self Local) GetScriptTemplatePath() string {
	return filepath.Join(self.GetStorageLocalPath(), "script")
}

func (self Local) GetBackupPath() string {
	return filepath.Join(self.GetStorageLocalPath(), "backup")
}

func (self Local) GetStorageLocalPath() string {
	if facade.GetConfig() == nil {
		slog.Debug("storage local path empty")
		return ""
	}
	path := facade.GetConfig().GetString("storage.local.path")
	if path == "" {
		panic("storage.local.path empty")
	}
	return facade.GetConfig().GetString("storage.local.path")
}

func (self Local) GetNginxSettingPath() string {
	return fmt.Sprintf("%s/nginx/proxy_host/", self.GetStorageLocalPath())
}

func (self Local) GetNginxCertPath() string {
	if override := os.Getenv(acme.EnvOverrideConfigHome); override != "" {
		return override
	}
	return fmt.Sprintf("%s/cert/", self.GetStorageLocalPath())
}

func (self Local) CreateTempFile(name string) (*os.File, error) {
	if name == "" {
		return os.CreateTemp(self.GetSaveRootPath(), "dpanel-temp-")
	}
	_ = os.MkdirAll(filepath.Dir(filepath.Join(self.GetSaveRootPath(), name)), os.ModePerm)
	return os.Create(filepath.Join(self.GetSaveRootPath(), name))
}

func (self Local) CreateTempDir(name string) (string, error) {
	if name == "" {
		return os.MkdirTemp(self.GetSaveRootPath(), "dpanel-temp-")
	}
	path := filepath.Dir(filepath.Join(self.GetSaveRootPath(), name))
	err := os.MkdirAll(path, os.ModePerm)
	return path, err
}

func (self Local) SaveUploadImage(uploadFileName, newFileNamePrefix string, appendRandomString bool) string {
	// 删除旧的前缀文件
	rootPath := filepath.Join(self.GetSaveRootPath(), "image")
	if matches, err := filepath.Glob(filepath.Join(rootPath, newFileNamePrefix+"*")); err == nil {
		for _, match := range matches {
			_ = os.Remove(match)
		}
	}
	var newFileName string
	if appendRandomString {
		newFileName = fmt.Sprintf("%s-%s.png", newFileNamePrefix, function.GetRandomString(5))
	} else {
		newFileName = fmt.Sprintf("%s.png", newFileNamePrefix)
	}

	newBgFile := filepath.Join(rootPath, newFileName)
	_ = os.MkdirAll(filepath.Dir(newBgFile), 0777)
	_ = os.Rename(
		filepath.Join(self.GetSaveRootPath(), uploadFileName),
		newBgFile,
	)
	return "/dpanel/static/image/" + newFileName
}
