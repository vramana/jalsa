set dotenv-load

host FILE :
  go run main.go {{FILE}}

watch :
  go build main.go && go test ./rpc
