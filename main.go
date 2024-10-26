package main

import (
	"bufio"
	"encoding/json"
	"io"
	"os"

	"jalsa/lsp"
	"jalsa/rpc"
)

// func test() {
// 	// Iterate over os.Args slice and print each argument
//
// 	db, err := sql.Open("sqlite3", "./test.db")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()
//
// 	arg := os.Args[1]
// 	b, err := os.ReadFile(arg)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
//
// 	fmt.Printf("Reading file: %s\n", arg)
//
// 	contents := string(b)
// 	sentences := strings.Split(contents, "\n")
//
// 	frontMatter := false
//
// 	for _, sentence := range sentences {
// 		if sentence == "+++" {
// 			frontMatter = !frontMatter
// 			continue
// 		}
// 		if sentence == "" {
// 			continue
// 		}
// 		fmt.Printf("Checking sentence: %s\n", sentence)
// 		h := hash(sentence)
// 		var result string
// 		err = db.QueryRow("SELECT correction FROM sentences WHERE sentence_hash = ?", h).Scan(&result)
//
// 		if err != nil && err != sql.ErrNoRows {
// 			fmt.Printf("Error: %v\n", err)
// 			continue
// 		}
//
// 		if result != "" {
// 			fmt.Printf("Result: %s\n", result)
// 			continue
// 		}
//
// 		fmt.Printf("Not found in cache\n")
//
// 		result, err = checkSentence(sentence)
// 		if err != nil {
// 			fmt.Printf("Error: %v\n", err)
// 			continue
// 		}
// 		fmt.Printf("Result: %s\n", result)
//
// 		time.Sleep(100 * time.Millisecond)
//
// 		_, err = db.Exec("INSERT INTO sentences (sentence_hash, sentence, correction) VALUES (?, ?, ?)", hash(sentence), sentence, result)
// 	}
// }

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
		server.Analyze(notification.Params.TextDocument.URI)

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
		server.Analyze(notification.Params.TextDocument.URI)
	default:
		server.Logger.Println(string(msg))
	}

}

func writeMessage(writer io.Writer, message any) {
	data := rpc.EncodeMessage(message)
	writer.Write([]byte(data))
}
