package storage

import (
	"bytes"
	"fmt"
	"time"

	storage_go "github.com/supabase-community/storage-go"
)

type StorageService struct {
	client       *storage_go.Client
	BucketCV     string
	BucketAvatar string
}

func NewStorageService(supabaseURL, serviceRoleKey, bucketCV, bucketAvatar string) *StorageService {
	client := storage_go.NewClient(supabaseURL+"/storage/v1", serviceRoleKey, nil)
	return &StorageService{
		client:       client,
		BucketCV:     bucketCV,
		BucketAvatar: bucketAvatar,
	}
}

func (s *StorageService) UploadCV(userID string, file []byte, contentType string) (string, error) {
	fileName := fmt.Sprintf("%s/cv_%d.pdf", userID, time.Now().Unix())
	upsert := true
	_, err := s.client.UploadFile(s.BucketCV, fileName, bytes.NewReader(file), storage_go.FileOptions{
		ContentType: &contentType,
		Upsert:      &upsert,
	})
	if err != nil {
		return "", err
	}
	res := s.client.GetPublicUrl(s.BucketCV, fileName)
	return res.SignedURL, nil
}

func (s *StorageService) UploadAvatar(userID string, file []byte, contentType string) (string, error) {
	fileName := fmt.Sprintf("%s/avatar_%d", userID, time.Now().Unix())
	upsert := true
	_, err := s.client.UploadFile(s.BucketAvatar, fileName, bytes.NewReader(file), storage_go.FileOptions{
		ContentType: &contentType,
		Upsert:      &upsert,
	})
	if err != nil {
		return "", err
	}
	res := s.client.GetPublicUrl(s.BucketAvatar, fileName)
	return res.SignedURL, nil
}

func (s *StorageService) DeleteFile(bucket, fileName string) error {
	_, err := s.client.RemoveFile(bucket, []string{fileName})
	return err
}
