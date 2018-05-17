package filemanager

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

type LocalFileManager struct {
	Path string
}

func NewLocalFileManager(savepath string) (*LocalFileManager, error) {
	if err := os.MkdirAll(savepath, 0755); err != nil {
		return nil, err
	}
	return &LocalFileManager{
		Path: savepath,
	}, nil
}

/* 保存文件 */
func (fm *LocalFileManager) Save(r io.Reader) (string, error) {
	tmpfile, err := ioutil.TempFile(fm.Path, "eternal-local")
	if err != nil {
		return "", err
	}
	defer tmpfile.Close()
	defer os.Remove(tmpfile.Name())
	buf := make([]byte, 1024)
	h := md5.New()
	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}
		if _, err := tmpfile.Write(buf[:n]); err != nil {
			return "", err
		}
		if _, err := h.Write(buf[:n]); err != nil {
			return "", err
		}
	}
	key := fmt.Sprintf("%x", h.Sum(nil))
	filePath := path.Join(fm.Path, key)
	if err := os.Rename(tmpfile.Name(), filePath); err != nil {
		return "", err
	}
	return key, nil
}

/* 读取文件 */
func (fm *LocalFileManager) Read(key string) (io.Reader, error) {
	savepath := path.Join(fm.Path, key)
	file, err := os.Open(savepath)
	if err != nil {
		return nil, ErrFileNotFound
	}
	return file, nil
}

/* 获取文件信息 */
func (fm *LocalFileManager) Stat(key string) (*FileStat, error) {
	savepath := path.Join(fm.Path, key)
	stat, err := os.Stat(savepath)
	if err != nil {
		return nil, err
	}
	return &FileStat{
		ID:   key,
		Size: stat.Size(),
	}, nil
}
