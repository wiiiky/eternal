package file

import (
	"eternal/errors"
	"eternal/filemanager"
	fileModel "eternal/model/file"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"net/http"
)

/* 获取当前支持的国家 */
func UploadFile(ctx echo.Context) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	fm := filemanager.Get()
	key, err := fm.Save(src)
	if err != nil {
		log.Error("File save failed:", err)
		return filemanager.ErrFileSaveFailure
	}
	log.Infof("%v", file.Header)

	contentType := file.Header.Get("Content-type")
	fileInfo, err := fileModel.SaveFile(key, contentType)
	if err != nil {
		log.Error("SaveFile failed:", err)
		return err
	}

	return ctx.JSON(http.StatusOK, fileInfo)
}

func DownloadFile(ctx echo.Context) error {
	fileID := ctx.Param("id")

	fileInfo, err := fileModel.GetFile(fileID)
	if err != nil {
		return err
	} else if fileInfo == nil {
		return errors.ErrFileNotFound
	}

	fm := filemanager.Get()
	r, err := fm.Read(fileID)
	if err != nil {
		log.Warn("Read failed:", err)
		return errors.ErrFileNotFound
	}

	return ctx.Stream(http.StatusOK, fileInfo.ContentType, r)
}
