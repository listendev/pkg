# result uploader


```go
import (
    	"github.com/aws/aws-sdk-go-v2/credentials"
)

func example() {
    accessKey:= os.Getenv("AWS_ACCESS_KEY_ID")
    secretKey: = os.Getenv("AWS_SECRET_ACCESS_KEY")
    region := os.Getenv("AWS_REGION")
    s3Endpoint := "" // empty for default (aws upstream, you can put your minio URL here)
    s3Creds := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")
    s3Uploader := resultuploader.NewS3Uploader(region, "my-garnet-bucket-name", s3Creds, s3Endpoint)
}
```
