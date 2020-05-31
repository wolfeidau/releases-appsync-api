# releases-appsync-api

This project provides a simple working AppSync API with Go code generated from a schema using [github.com/99designs/gqlgen](https://github.com/99designs/gqlgen) storing data in DynamoDB.

Authentication for this project is provided by [Okta](https://www.okta.com) which is integrated with AppSync using OpenID.

# Development

Create a `.envrc` file and use [direnv](https://direnv.net/) to switch environments:

```
export AWS_PROFILE=whatever
export AWS_REGION=us-east-1
export PACKAGE_BUCKET=deploy-lambda-us-east-1
# Used in development only
export RAW_EVENT_LOGGING=true
export OPENID_CONNECT_ISSUER=https://dev-xxxxxx.oktapreview.com
export OPENID_CONNECT_CLIENTID=1234567890abcdef1234
```

# License

This application is released under Apache 2.0 license by [Mark Wolfe](https://www.wolfe.id.au).
