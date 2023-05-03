# result uploader

```go
import (
    	"github.com/aws/aws-sdk-go-v2/credentials"
        "github.com/aws/aws-sdk-go-v2/config"
        "os"
        "github.com/garnet-org/pkg/resultuploader"
)


func main() {
    cfg, _ := config.LoadDefaultConfig(context.TODO())
    s3Uploader := resultuploader.NewS3Uploader(cfg, "my-garnet-bucket-name")
}

```


During local development (e.g with minio), you can export the following env variables

```
export AWS_ACCESS_KEY_ID=user
export AWS_SECRET_ACCESS_KEY=password
```

and instruct the s3Uploader to use the minio server

```go
s3Uploader.WithCustomEndpoint("http://127.0.0.1:9000")
```
