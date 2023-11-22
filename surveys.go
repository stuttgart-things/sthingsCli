/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/allegro/bigcache/v3"
)

type RenderSurvey struct {
	SingleInputSurvey  func(string) string
	SingleSelectSurvey func(string) string
}

func AskSingleSelectQuestion(questionText string, options []string) string {

	qs := []*survey.Question{
		{
			Name: "question",
			Prompt: &survey.Select{
				Message: questionText,
				Options: options,
			},
		},
	}

	answers := struct {
		Question string `survey:"question"`
	}{}

	err := survey.Ask(qs, &answers)

	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	return answers.Question
}

func AskSingleInputQuestion(questionText string, defaultAnswer string) string {

	qs := []*survey.Question{
		{
			Name: "question",
			Prompt: &survey.Input{
				Message: questionText,
				Default: defaultAnswer,
			},
		},
	}

	answers := struct {
		Question string `survey:"question"`
	}{}

	err := survey.Ask(qs, &answers)

	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	return answers.Question
}

func AskMultiSelectQuestion(questionText string, options []string) []string {

	qs := []*survey.Question{
		{
			Name: "questionMultiSelect",
			Prompt: &survey.MultiSelect{
				Message: questionText,
				Options: options,
			},
		},
	}

	selectedAnswers := []string{}
	err := survey.Ask(qs, &selectedAnswers)

	if err != nil {
		fmt.Println(err.Error())
	}

	return selectedAnswers
}

func RenderTemplateSurvey(templateContent string, globalValues map[string]interface{}) (string, map[string]interface{}) {
	var buf bytes.Buffer

	// INITIALIZE CACHE
	cache, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))

	survey := RenderSurvey{
		SingleInputSurvey: func(defaultKey string) (answer string) {

			// SET DEFAULT VALUES
			values := []string{defaultKey, "-"}
			cacheKey := defaultKey

			fmt.Println("DEFAULT KEY:", defaultKey)

			if globalValues[defaultKey] != nil && strings.Contains(globalValues[defaultKey].(string), "|") {
				values = strings.Split(globalValues[defaultKey].(string), "|")

				// GET REUSABLE VARIABLES AND SET CACHE KEY
				if strings.Contains(values[1], "+") {
					valuesVar := strings.Split(globalValues[defaultKey].(string), "+")
					cacheKey = valuesVar[1]
				}

				// GET KEY FROM CACHE
				cachedEntry, _ := cache.Get(cacheKey)
				if len(cachedEntry) != 0 {
					values[1] = string(cachedEntry)
					fmt.Println("CACHED ENTRY:", string(cachedEntry))
				} else {
					values[1] = "PLEASE OVERWRITE"
				}

			}

			// ASK QUESTION
			answer = AskSingleInputQuestion("Enter "+values[0]+":", values[1])

			// SET VALUE TO CACHE
			cache.Set(values[0], []byte(answer))
			cache.Set(cacheKey, []byte(answer))

			return
		},
		SingleSelectSurvey: func(defaultValues string) (answer string) {
			values := []string{"value"}

			if strings.Contains(globalValues[defaultValues].(string), "|") {
				values = strings.Split(globalValues[defaultValues].(string), "|")
			}

			answer = AskSingleSelectQuestion("Select "+values[0]+":", strings.Split(values[1], ";"))

			globalValues[globalValues[defaultValues].(string)] = answer

			return
		},
	}

	tmpl, err := template.New("rendering").Parse(templateContent)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(&buf, survey)

	if err != nil {
		log.Fatalf("execution: %s", err)
	}

	return buf.String(), globalValues

}
