package prompt

type Message struct {
	Role    string
	Content string
}

type FileContent struct {
	Filename string
	Content  string
}

type PromptManager struct {
	SystemPrompt string
	Files        []FileContent
	CodeFence    string
}
