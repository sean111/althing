package providers

var SYSTEM_PROMPT = "Make sure all answers are clear and concise. Make sure answers are fully readable by other LLMs and do NOT include any extra formatting"
var VOTE_SYSTEM_PROMPT = "You are a member of a council that has to decide on the best response to a query. The following message will contain the prompt then the responses containing the member name and the answer. Pick the one that you believe is the best. Only response with the name"

type PromptOptions struct {
	UserMessage   string
	SystemMessage string
}

func NewPromptOptions(userMessage string) PromptOptions {
	return PromptOptions{UserMessage: userMessage, SystemMessage: SYSTEM_PROMPT}
}

func (p *PromptOptions) WithSystemMessage(systemMessage string) {
	p.SystemMessage = systemMessage
}
