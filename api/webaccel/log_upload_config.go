package webaccel

// LogUploadConfig ログアップロード設定
type LogUploadConfig struct {
	Bucket          string `json:","`
	Endpoint        string `json:"," validate:",equal=https://s3.isk01.sakurastorage.jp"`
	Region          string `json:"," validate:",equal=jp-north-1"`
	AccessKeyID     string `json:","`
	SecretAccessKey string `json:","`
	Status          string `json:"," validate:"oneof=enabled disabled"`
}
