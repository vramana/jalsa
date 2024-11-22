package lsp

type InitializeRequest struct {
	Request
	Params *InitializeParams `json:"params"`
}

type InitializeParams struct {
	ProcessID  *int  `json:"processId,omitempty"`
	ClientInfo *Info `json:"clientInfo,omitempty"`
}

type Info struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

type InitializeResponse struct {
	Response
	Result *InitializeResult `json:"result,omitempty"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   *Info              `json:"serverInfo,omitempty"`
}

type ServerCapabilities struct {
	TextDocumentSync    TextDocumentSyncOptions `json:"textDocumentSync"`
	DiagnosticsProvider DiagnosticsOptions      `json:"diagnosticsProvider"`
}

type DiagnosticsOptions struct {
	Identifier            string `json:"identifier"`
	InterFileDependencies bool   `json:"interFileDependencies"`
	WorkspaceDiagnostics  bool   `json:"workspaceDiagnostics"`
}

type TextDocumentSyncOptions struct {
	OpenClose bool        `json:"openClose"`
	Change    int         `json:"change"`
	Save      SaveOptions `json:"save"`
}

// TODO: Check why doesn't this work?
type SaveOptions struct {
	IncludeText bool `json:"incudeText"`
}

func NewInitializeResponse(ID *int) *InitializeResponse {
	response := &InitializeResponse{
		Response: Response{RPC: "2.0", ID: ID},
		Result: &InitializeResult{
			ServerInfo: &Info{Name: "jalsa", Version: "0.0.1"},
			Capabilities: ServerCapabilities{
				TextDocumentSync: TextDocumentSyncOptions{
					OpenClose: true,
					Change:    1,
					Save:      SaveOptions{IncludeText: true},
				},
				DiagnosticsProvider: DiagnosticsOptions{
					Identifier:            "jalsa",
					InterFileDependencies: false,
					WorkspaceDiagnostics:  false,
				},
			},
		},
	}

	return response
}
