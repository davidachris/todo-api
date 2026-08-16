# todo-api

- Run `make run` to start the API
- Run `make build` to compile the Lambda bootstrap zip
- Run `make clean` to clean up all folder and files created by make commands

S3 is used only when `BUCKET_NAME` is set and `AWS_DISABLED` is not `true`/`1`/`yes`. Otherwise the API uses SQLite only (default `/tmp/todos.db`, override with `SQLITE_PATH`). `go run .` works with no AWS env vars.

```
Add AWS Layer: arn:aws:lambda:us-east-1:753240598075:layer:LambdaAdapterLayerX86:22
Add AWS Layer: arn:aws:lambda:us-east-1:753240598075:layer:LambdaAdapterLayerArm64
Git page awslabs: https://github.com/awslabs/aws-lambda-web-adapter
```
