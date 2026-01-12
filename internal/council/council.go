package council

import (
	"context"
	"fmt"
	"sean111/althing/internal/formatting"
	"sean111/althing/internal/providers"
	"sean111/althing/internal/tools"
	"sync"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/viper"
)

var conncil *Council

type Council struct {
	Members []Member
	Tools   []tools.Tool
}

type MemberResponse struct {
	Name     string
	Response string
	Error    error
}

type Member struct {
	Name       string         `mapstructure:"name"`
	ProviderId string         `mapstructure:"provider"`
	Model      string         `mapstructure:"model"`
	Provider   MemberProvider `mapstructure:"-"`
}

type MemberProvider interface {
	Prompt(ctx context.Context, options providers.PromptOptions) (string, error)
}

func Init() {
	conncil = &Council{}

	// Setup tools
	//tools.ToolList = map[string]tools.Tool{
	//	"web_search": tools.NewSearch(),
	//}

	// Set up members
	if err := viper.UnmarshalKey("members", &conncil.Members); err != nil {
		fmt.Printf("Error loading members: %v\n", err)
	}

	for i := range conncil.Members {
		member := &conncil.Members[i]
		//fmt.Printf("%v\n", member)
		providerInfo := viper.GetStringMap(fmt.Sprintf("providers.%s", member.ProviderId))
		//fmt.Printf("%v\n", providerInfo)
		switch providerInfo["type"].(string) {
		case "openai":
			providerOptions := providers.OpenAIProviderOptions{
				APIKey: providerInfo["api_key"].(string),
				URL:    providerInfo["url"].(string),
				Model:  member.Model,
			}
			provider, err := providers.CreateOpenAIProvider(providerOptions)

			if err != nil {
				panic(err)
			}
			member.Provider = provider
		case "gemini":
			providerOptions := providers.GeminiProviderOptions{
				APIKey: providerInfo["api_key"].(string),
				Model:  member.Model,
			}
			provider, err := providers.CreateGeminiProvider(providerOptions)

			if err != nil {
				panic(err)
			}
			member.Provider = provider
		case "openrouter":
			providerOptions := providers.OpenRouterProviderOptions{
				Token: providerInfo["token"].(string),
				Title: providerInfo["title"].(string),
				URL:   providerInfo["url"].(string),
				Model: member.Model,
			}

			provider, err := providers.CreateOpenRouterProvider(providerOptions)

			if err != nil {
				panic(err)
			}
			member.Provider = provider
		}
	}
}

func getCouncil() *Council {
	return conncil
}

func Run(prompt string) {
	council := getCouncil()
	var wg sync.WaitGroup

	fmt.Println(formatting.HeaderStyle.Render("Running Council..."))
	results := make(chan MemberResponse, len(council.Members))
	promptOptions := providers.NewPromptOptions(prompt)
	for _, member := range council.Members {
		wg.Add(1)
		go memberWorker(member, promptOptions, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	votePrompt := fmt.Sprintf("Prompt: %s\n", prompt)

	for res := range results {
		if res.Error != nil {
			fmt.Printf("[%s]: %s\n", formatting.MemberNameStyle.Render(res.Name), formatting.ErrorStyle.Render(res.Error.Error()))
		} else {
			fmt.Printf("[%s]: %s\n", formatting.MemberNameStyle.Render(res.Name), formatting.ResponseStyle.Render(res.Response))
			votePrompt += fmt.Sprintf("Member: %s\n Response: %s\n\n", res.Name, res.Response)
		}
	}
	formatting.Line()
	vote(votePrompt)

}

func memberWorker(member Member, promptOptions providers.PromptOptions, results chan<- MemberResponse, wg *sync.WaitGroup) {
	defer wg.Done()
	var result MemberResponse
	result.Name = member.Name
	response, err := member.Provider.Prompt(context.Background(), promptOptions)
	if err != nil {
		result.Error = err
	} else {
		result.Response = response
	}
	results <- result
}

func vote(prompt string) {
	council := getCouncil()
	var wg sync.WaitGroup
	votes := make(chan MemberResponse, len(council.Members))
	promptOptions := providers.NewPromptOptions(prompt)
	promptOptions.WithSystemMessage(providers.VOTE_SYSTEM_PROMPT)

	for _, member := range council.Members {
		wg.Add(1)
		go memberWorker(member, promptOptions, votes, &wg)
	}

	go func() {
		wg.Wait()
		close(votes)
	}()

	fmt.Println(formatting.HeaderStyle.Render("Voting..."))

	memberVotes := ""

	for vote := range votes {
		if vote.Error != nil {
			fmt.Printf("[%s]: Error ... %v\n", formatting.MemberNameStyle.Render(vote.Name), formatting.ErrorStyle.Render(vote.Error.Error()))
			continue
		}
		fmt.Printf("[%s] %s\n", formatting.MemberNameStyle.Render(vote.Name), formatting.ResponseStyle.Render(vote.Response))
		memberVotes += fmt.Sprintf("%s: %s\n", vote.Name, vote.Response)
	}

	votePrompt := fmt.Sprintf("\nMember answers:\n%s\nMember votes:%s", prompt, memberVotes)

	speakerResponse, err := speaker(votePrompt)

	if err != nil {
		fmt.Printf("Error getting speaker: %v\n", err)
	} else {
		var output string
		formatting.Line()
		renderedResponse, err := glamour.Render(speakerResponse, "dark")
		if err != nil {
			output = speakerResponse
		} else {
			output = renderedResponse
		}
		fmt.Printf("%s\n\n%s\n", formatting.MemberNameStyle.Render("Speaker:"), output)
	}
}

func getMember(name string) (Member, error) {
	council := getCouncil()
	for _, member := range council.Members {
		if member.Name == name {
			return member, nil
		}
	}
	return Member{}, fmt.Errorf("member not found")
}

func speaker(votes string) (string, error) {
	speaker, err := getMember(viper.GetString("speaker"))
	if err != nil {
		return "", err
	}
	promptOptions := providers.NewPromptOptions(votes)
	promptOptions.WithSystemMessage(providers.SPEAKER_SYSTEM_PROMPT)
	response, err := speaker.Provider.Prompt(context.Background(), promptOptions)

	if err != nil {
		return "", err
	}

	return response, nil
}
