package main

import (
	"bufio"
	"encoding/json"
	"io"
	"os"

	"jalsa/lsp"
	"jalsa/rpc"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	writer := os.Stdout

	server := lsp.NewServer()

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, msg, err := rpc.DecodeMessage(msg)
		if err != nil {
			server.Logger.Printf("Got an error %s", err)
			continue
		}
		handleMessage(msg, writer, method, server)
	}
}

func handleMessage(msg []byte, writer io.Writer, method string, server *lsp.Server) {
	server.Logger.Printf("Method: %s", method)
	switch method {
	case "initialize":
		request := new(lsp.InitializeRequest)
		if err := json.Unmarshal(msg, request); err != nil {
			server.Logger.Printf("We could not parse initialize %s", err)
			return
		}
		server.Logger.Printf("Connected to client %s %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)

		response := lsp.NewInitializeResponse(request.ID)
		writeMessage(writer, response)
	case "textDocument/didOpen":
		notification := new(lsp.DidOpenTextDocumentNotification)
		if err := json.Unmarshal(msg, notification); err != nil {
			server.Logger.Printf("We could not parse %s %s", method, err)
			return
		}

		server.Files[notification.Params.TextDocument.URI] = notification.Params.TextDocument.Text

		diagnosticsNotification := server.Analyze(notification.Params.TextDocument.URI)
		writeMessage(writer, diagnosticsNotification)

	case "textDocument/didChange":
		notification := new(lsp.DidChangeTextDocumentNotification)
		if err := json.Unmarshal(msg, notification); err != nil {
			server.Logger.Printf("We could not parse %s %s", method, err)
			return
		}

		server.Files[notification.Params.TextDocument.URI] = notification.Params.ContentChanges[0].Text
	case "textDocument/didSave":
		notification := new(lsp.DidSaveTextDocumentNotification)
		if err := json.Unmarshal(msg, notification); err != nil {
			server.Logger.Printf("We could not parse %s %s", method, err)
			return
		}

		writeMessage(writer, server.CachedDiagnostics(notification.Params.TextDocument.URI))
		diagnosticsNotification := server.Analyze(notification.Params.TextDocument.URI)
		writeMessage(writer, diagnosticsNotification)

	default:
		server.Logger.Println(string(msg))
	}

}

func writeMessage(writer io.Writer, message any) {
	data := rpc.EncodeMessage(message)
	writer.Write([]byte(data))
}
