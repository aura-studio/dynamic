package dynamic

var (
	localPath  string
	remotePath string
)

func Init(argLocalPath string, argRemotePath string) {
	localPath = argLocalPath
	remotePath = argRemotePath
}
