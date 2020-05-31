package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog"
	"github.com/wolfeidau/dynastore"
	lmw "github.com/wolfeidau/lambda-go-extras/middleware"
	"github.com/wolfeidau/realworld-appsync-ddb/internal/flags"
	"github.com/wolfeidau/realworld-appsync-ddb/internal/graph"
	"github.com/wolfeidau/realworld-appsync-ddb/internal/handlers"
	"github.com/wolfeidau/realworld-appsync-ddb/pkg/lambdamiddleware/dump"
	"github.com/wolfeidau/realworld-appsync-ddb/pkg/lambdamiddleware/zlog"
)

var (
	buildDate string
	commit    string

	apiFlags flags.API
)

func main() {
	kong.Parse(&apiFlags,
		kong.Vars{"version": fmt.Sprintf("%s_%s", commit, buildDate)}, // bind a var for version
	)

	session := dynastore.New()

	resolvers := graph.NewResolvers(session.Table(apiFlags.ReleaseTable).Partition("releases"))

	rp := handlers.NewRequestProcessor(resolvers)

	ch := lmw.New(zlog.WithConfig(zlog.Config{
		Level:  zerolog.DebugLevel,
		Output: os.Stderr,
	}), dump.Middleware()).ThenFunc(rp.Handler)

	lambda.StartHandler(ch)
}
