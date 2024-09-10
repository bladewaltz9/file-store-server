package config

const (
	// Directory path
	FileStoreDir = "/home/bladewaltz/data/files/"
	FileChunkDir = "/home/bladewaltz/data/chunks/"

	MaxUploadSize = 32 << 20 // 32MB

	// SSL cert and key
	CertFile = "/etc/apache2/ssl/bladewaltz.cn.crt"
	KeyFile  = "/etc/apache2/ssl/bladewaltz.cn.key"
)
