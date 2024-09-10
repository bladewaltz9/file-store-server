package mq

type FileTransferMessage struct {
	FileID    int
	LocalFile string
	ObjectKey string
}
