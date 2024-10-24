package lsp

type InitializeRequest struct {
	Request
	Params *InitializeParams `json:"params"`
}

type InitializeParams struct {
	ProcessID  *int   `json:"processId",omitempty`
	ClientInfo *Info  `json:"clientInfo",omitempty`
	RootUri    string `json:"rootUri",omitempty`
}

type Info struct {
	Name    string `json:"name"`
	Version string `json:"version",omitempty`
}

type InitializeResponse struct {
	Response
	Result *InitializeResult `json:"result",omitempty`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   *Info              `json:"serverInfo",omitempty`
}

type ServerCapabilities struct {
}

func NewInitializeResponse(ID *int) *InitializeResponse {
	return nil
}
