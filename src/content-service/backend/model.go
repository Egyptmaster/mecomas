package backend

import (
	"encoding/base64"
	"fmt"
	"time"
)

type MediaItem struct {
	Id         string
	Name       string
	UserId     string
	UploadDate time.Time
}

func (m MediaItem) Key() MediaItemKey {
	return MediaItemKey{
		userId:     m.UserId,
		uploadDate: m.UploadDate,
		name:       m.Name,
	}
}

type MediaItemKey struct {
	userId     string
	uploadDate time.Time
	name       string
}

func (m MediaItemKey) S3Key() (string, error) {
	k, err := base64.StdEncoding.DecodeString(m.name)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%d/%d/%d/%s",
			m.userId,
			m.uploadDate.Year(),
			m.uploadDate.Month(),
			m.uploadDate.Day(),
			string(k)),
		nil
}
