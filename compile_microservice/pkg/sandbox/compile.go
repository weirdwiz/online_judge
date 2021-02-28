package sandbox

type CompileRequest struct {
	Code     string `json:"code"`
	Language string `json:"lang"`
}

type CompileResponse struct {
	Output string `json:"output"`
}
