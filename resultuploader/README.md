# result uploader

Using provided credentials

```go
import (
    	"github.com/aws/aws-sdk-go-v2/credentials"
        "os"
        "github.com/garnet-org/pkg/resultuploader"
)

func main() {
    accessKey:= os.Getenv("AWS_ACCESS_KEY_ID")
    secretKey: = os.Getenv("AWS_SECRET_ACCESS_KEY")
    region := os.Getenv("AWS_REGION")
    s3Creds := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")
    cfg := aws.Config{
        Region:      region,
        Credentials: credentials,
    }

    s3Uploader := resultuploader.NewS3Uploader(cfg, "my-garnet-bucket-name")
    s3Uploader.WithCustomEndpoint("http://127.0.0.1:9000")
}
```

Using container credentials (e.g: on ECS)

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
