package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog"
	lmw "github.com/wolfeidau/lambda-go-extras/middleware"
	"github.com/wolfeidau/realworld-appsync-ddb/internal/handlers"
	"github.com/wolfeidau/realworld-appsync-ddb/pkg/lambdamiddleware/dump"
	"github.com/wolfeidau/realworld-appsync-ddb/pkg/lambdamiddleware/zlog"
)

func main() {

	rp := &handlers.RequestProcessor{}

	ch := lmw.New(zlog.WithConfig(zlog.Config{
		Level:  zerolog.DebugLevel,
		Output: os.Stderr,
	}), dump.Middleware()).ThenFunc(rp.Handler)

	lambda.StartHandler(ch)
}
