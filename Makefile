.PHONY: err
err:
	protoc --proto_path=./pkg/errx \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:./pkg/errx \
		   --go-errors_out=paths=source_relative:./pkg/errx \
 	       --go-http_out=paths=source_relative:./pkg/errx \
 	       --go-grpc_out=paths=source_relative:./pkg/errx \
	       --openapi_out=fq_schema_naming=true,default_response=false:. \
	       error_reason.proto