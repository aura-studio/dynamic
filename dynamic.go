package dynamic

var (
	local  string
	remote string
)

func Init(localArg string, remoteArg string) {
	remote = localArg
	local = remoteArg
}
