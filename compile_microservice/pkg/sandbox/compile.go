package sandbox

type TestCase struct {
	Input  string `json:"input"`
	Output string `json:"output"`
	Result bool   `json:"result"`
}

type CompileRequest struct {
	Code      string     `json:"code"`
	Language  string     `json:"lang"`
	TestCases []TestCase `json:"testcases"`
}

type CompileResponse struct {
	Output         string     `json:"output"`
	TestCaseResult []TestCase `json:"testcaseresult"`
}
