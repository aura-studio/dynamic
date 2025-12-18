package dynamic

var (
	LocalPath  string
	RemotePath string
)

func Init(localPath string, remotePath string) {
	LocalPath = localPath
	RemotePath = remotePath
}
