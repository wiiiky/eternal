package file

import (
	"eternal/filemanager"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"net/http"
)

/* 获取当前支持的国家 */
func UploadFile(ctx echo.Context) error {
	fm := filemanager.Get()

	file, err := ctx.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	key, err := fm.Save(src)
	if err != nil {
		log.Error("File save failed:", err)
		return filemanager.ErrFileSaveFailure
	}

	fileInfo := &FileInfo{
		ID: key,
	}
	return ctx.JSON(http.StatusOK, fileInfo)
}
