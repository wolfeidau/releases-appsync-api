package dump

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/wolfeidau/lambda-go-extras/lambdaextras"
)

func Middleware() func(next lambda.Handler) lambda.Handler {

	return func(next lambda.Handler) lambda.Handler {
		return lambdaextras.HandlerFunc(func(ctx context.Context, payload []byte) ([]byte, error) {

			fmt.Println(string(payload))

			result, err := next.Invoke(ctx, payload)

			fmt.Println(string(result))

			return result, err
		})

	}
}
