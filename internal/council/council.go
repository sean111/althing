package council

import (
	"context"
	"fmt"
	"sean111/althing/internal/formatting"
	"sean111/althing/internal/providers"
	"sync"

	"github.com/spf13/viper"
)

var conncil *Council

type Council struct {
	Members []Member
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

	tally := make(map[string]int)

	for vote := range votes {
		if vote.Error != nil {
			fmt.Printf("[%s]: Error ... %v\n", formatting.MemberNameStyle.Render(vote.Name), formatting.ErrorStyle.Render(vote.Error.Error()))
			continue
		}
		fmt.Printf("[%s] %s\n", formatting.MemberNameStyle.Render(vote.Name), formatting.ResponseStyle.Render(vote.Response))

		tally[vote.Response]++
	}

	winner := 0
	var winnerName string
	for member, count := range tally {
		fmt.Printf("%s: %d\n", member, count)
		if count > winner {
			winner = count
			winnerName = member
		}
		fmt.Printf("Member: %s Count: %d\n\n", member, count)
	}

	fmt.Printf("\nWinner: %s\n", winnerName)
}
