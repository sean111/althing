package providers

var SYSTEM_PROMPT = "Make sure all answers are clear and concise. Make sure answers are fully readable by other LLMs and do NOT include any extra formatting"
var VOTE_SYSTEM_PROMPT = "You are a member of a council that has to decide on the best response to a query. The following message will contain the prompt then the responses containing the member name and the answer. Pick the one that you believe is the best. Respond with the name of the member and a reason why you picked their answer"
var SPEAKER_SYSTEM_PROMPT = "You are the speaker of the council. Your job is to determine the winner of the votes and explain why that person won. If there is a tie for the best answer you are to act as the tie breaker. Make sure to also include a summary of the vote counts, without your opinion added, as well as the original unmodified answer"

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
