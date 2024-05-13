package backend

import (
	"cid.com/content-service/common/retry"
	"cid.com/content-service/common/secrets"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log/slog"
	"time"
)

type FileStorageConfiguration struct {
	retry.Settings `yaml:",inline"`
	Url            string
	Bucket         string
	Expire         time.Duration
}

type FileStorage interface {
	GeneratePreSignedUrl(key MediaItemKey) (string, error)
}

type fileStorage struct {
	cfg    FileStorageConfiguration
	bucket bucket
}

func NewFileStorage(cfg FileStorageConfiguration, vault secrets.Secrets) (FileStorage, error) {
	slog.Info("establishing connection to file storage")

	auth, err := vault.Get("s3")
	if err != nil {
		return nil, err
	}

	retryer := client.DefaultRetryer{
		NumMaxRetries: int(cfg.Retry),
		MinRetryDelay: cfg.Delay,
	}
	cred := credentials.NewStaticCredentials(auth.User, auth.Password, "")
	sess := session.Must(session.NewSession())
	awsCfg := &aws.Config{
		Credentials:      cred,
		Region:           aws.String(endpoints.UsWest1RegionID),
		S3ForcePathStyle: aws.Bool(true), // required to work with minIO
		Endpoint:         aws.String(cfg.Url),
		Retryer:          retryer,
	}
	conn := s3.New(sess, awsCfg)
	b := bucket{
		conn: conn,
		name: cfg.Bucket,
	}

	// verify if bucket exists or can be created
	if exists, err := b.exists(); err != nil {
		slog.Error("the request for existing buckets failed", slog.String("err", err.Error()))
		return nil, err
	} else if !exists {
		slog.Info(fmt.Sprintf("the bucket %s doesn't exists, so we trying to create it", b.name))
		if err = b.create(); err != nil {
			slog.Error(fmt.Sprintf("the creation of the bucket %s failed", b.name), slog.String("err", err.Error()))
			return nil, err
		}
	} else {
		slog.Info(fmt.Sprintf("the bucket %s already exists", b.name))
	}

	fs := &fileStorage{
		bucket: b,
		cfg:    cfg,
	}

	return fs, nil
}

func (fs fileStorage) GeneratePreSignedUrl(key MediaItemKey) (string, error) {
	return fs.bucket.generatePreSignedUrl(key, fs.cfg.Expire)
}

type bucket struct {
	name string
	conn *s3.S3
}

func (b bucket) exists() (bool, error) {
	resp, err := b.conn.ListBuckets(nil)
	if err != nil {
		return false, err
	}

	for _, item := range resp.Buckets {
		if aws.StringValue(item.Name) == b.name {
			return true, nil
		}
	}
	return false, nil
}

func (b bucket) create() error {
	req := &s3.CreateBucketInput{
		Bucket: aws.String(b.name),
	}
	_, err := b.conn.CreateBucket(req)
	return err
}

func (b bucket) generatePreSignedUrl(key MediaItemKey, expire time.Duration) (string, error) {
	s3Key, err := key.S3Key()
	if err != nil {
		return "", err
	}
	req, _ := b.conn.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(b.name),
		Key:    aws.String(s3Key),
	})
	url, err := req.Presign(expire)
	if err != nil {
		slog.Error(fmt.Sprintf("generation of pre-signed url for %s failed", s3Key), slog.String("err", err.Error()))
	}
	return url, err
}
