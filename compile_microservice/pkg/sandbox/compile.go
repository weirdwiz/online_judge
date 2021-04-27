package sandbox

type TestCase struct {
	Input  string `json:"input"`
	Output string `json:"output"`
	Result bool   `json:"result"`
}

type CompileRequest struct {
	Code     string   `json:"code"`
	Language string   `json:"lang"`
	TestCase TestCase `json:"testcase"`
}

type CompileResponse struct {
	Output string `json:"output"`
}
