package file

/* 文件信息 */
type FileStat struct {
	ID string `json:"id"`
	Size uint64 `json:"size"`
}

type FileManager interface{
	/* 保存文件，返回文件ID和错误信息 */
	Save(io.Reader) (string, error)

	/* 读取文件 */
	Read(string, io.Writer) error

	/* 获取我文件信息 */
	Stat(string) (*FileStat, error)
}