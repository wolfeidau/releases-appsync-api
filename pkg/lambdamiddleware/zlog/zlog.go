package zlog

import (
	"context"
	"io"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/lambda-go-extras/lambdaextras"
	"github.com/wolfeidau/realworld-appsync-ddb/pkg/lambdamiddleware"
)

type Config struct {
	Level  zerolog.Level
	Output io.Writer
	Fields map[string]interface{}
}

func WithConfig(cfg Config) func(next lambda.Handler) lambda.Handler {
	if cfg.Fields == nil {
		cfg.Fields = map[string]interface{}{}
	}

	lambdamiddleware.AddFieldsFromLambdaEnv(cfg.Fields)

	return func(next lambda.Handler) lambda.Handler {
		return lambdaextras.HandlerFunc(func(ctx context.Context, payload []byte) ([]byte, error) {
			lc, _ := lambdacontext.FromContext(ctx)
			zl := zerolog.New(cfg.Output).Level(cfg.Level).With().Caller().
				Fields(cfg.Fields).
				Str("aws_request_id", lc.AwsRequestID).
				Str("amzn_trace_id", os.Getenv("_X_AMZN_TRACE_ID")).
				Logger()

			// inject the logger into the context and pass it down the chain
			return next.Invoke(zl.WithContext(ctx), payload)
		})
	}

}

// WithContext this will add fields to the logger stored in the context enabling you to
// include more context as it flows down to the next layer of your service.
func WithContext(ctx context.Context, fields map[string]interface{}) context.Context {
	l := log.Ctx(ctx).With().Fields(fields).Logger()
	return l.WithContext(ctx)
}
