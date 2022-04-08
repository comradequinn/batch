build :
	@ go build 

test :
	@ go test ./...

exec : 
	@ cd ./cmd/example && go run example.go