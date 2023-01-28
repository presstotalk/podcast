package storage

// TODO: onedrive apis are fucking annoying

type Storage interface {
	UploadFromURL(destPath string, url string) error
}
