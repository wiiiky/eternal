package filemanager

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
)

/* 文件信息 */
type FileStat struct {
	ID   string `json:"id"`
	Size int64  `json:"size"`
}

var (
	ErrFileNotFound    = errors.New("File Not Found")
	ErrFileSaveFailure = errors.New("File Save failure")
)

const (
	FM_TYPE_LOCAL = "local"
)

type FileManager interface {
	/* 保存文件，返回文件ID和错误信息 */
	Save(io.Reader) (string, error)

	/* 读取文件 */
	Read(string, io.Writer) error

	/* 获取我文件信息 */
	Stat(string) (*FileStat, error)
}

var fileManager FileManager = nil

func Init() {
	ftype := viper.GetString("filemanager.type")
	if ftype == FM_TYPE_LOCAL {
		initLocalFileManager()
	} else {
		log.Fatal("Unknown file manager type:", ftype)
	}
}

func initLocalFileManager() {
	savePath := viper.GetString("filemanager.local.path")
	fm, err := NewLocalFileManager(savePath)
	if err != nil {
		log.Fatal("Initialize local file manager failed:", err)
	}
	fileManager = fm
}

func Get() FileManager {
	return fileManager
}
