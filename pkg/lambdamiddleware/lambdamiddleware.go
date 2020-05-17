package lambdamiddleware

import (
	"github.com/aws/aws-lambda-go/lambdacontext"
	"os"
)

// update the common stuff from environment variables using lambda context
func AddFieldsFromLambdaEnv(fields map[string]interface{}) {
	// https://docs.aws.amazon.com/lambda/latest/dg/configuration-envvars.html#configuration-envvars-runtime
	//
	// These don't change so may as well do it once here, also clobber any other values passed in just in case
	fields["aws_region"] = os.Getenv("AWS_REGION")
	fields["function_name"] = lambdacontext.FunctionName
	fields["function_version"] = lambdacontext.FunctionVersion
}
