protoc --go_out=./ --go_opt=paths=source_relative \
    --go-grpc_out=./ --go-grpc_opt=paths=source_relative \
    proto/domain/domain.proto

protoc --go_out=./ --go_opt=paths=source_relative \
    --go-grpc_out=./ --go-grpc_opt=paths=source_relative \
    proto/smtp/smtp.proto

protoc --go_out=./ --go_opt=paths=source_relative \
    --go-grpc_out=./ --go-grpc_opt=paths=source_relative \
    proto/notification/notification.proto
