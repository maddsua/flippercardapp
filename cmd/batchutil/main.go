package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	db_model "github.com/maddsua/flippercardapp/db/model"
	"github.com/maddsua/flippercardapp/utils"
)

func main() {

	srcPath := flag.String("src", "data/table.csv", "sets source data file")
	deckName := flag.String("name", "", "set deck name")

	flag.Parse()

	srcFile, err := os.Open(*srcPath)
	if err != nil {
		fmt.Println("Can't open source data file:", err)
		os.Exit(1)
	}

	defer srcFile.Close()

	reader := csv.NewReader(srcFile)

	var mapper *utils.RecordMapper

	var cards []db_model.CardNodeContent

	for {

		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("csv.Reader.Read:", err)
			os.Exit(1)
		}

		if mapper == nil {
			mapper = utils.NewRecordMapper(row)
			continue
		}

		title, _ := mapper.Get(row, "title")
		if title == "" {
			fmt.Println("WARN: Empty row title; Skipped")
			continue
		}

		question, _ := mapper.Get(row, "textarea")
		if question == "" {
			fmt.Println("WARN: Empty question; Skipped")
			continue
		}

		answer, _ := mapper.Get(row, "answer_text")
		if answer == "" {
			fmt.Println("WARN: Empty answer; Skipped")
			continue
		}

		frontContent := []db_model.CardContentNode{
			db_model.NewCardTitleNode(title),
			db_model.NewCardTextNode(question),
		}

		if opts, _ := mapper.Get(row, "poll_options"); opts != "" {

			var poll db_model.CardPollElement

			for idx, option := range strings.Split(opts, ",") {

				if option = strings.TrimSpace(option); option == "" {
					fmt.Printf("WARN: Empty poll option")
					continue
				}

				poll.Content = append(poll.Content, db_model.CardPollElementOptionNode{
					Value:    option,
					IsAnswer: idx == 0,
				})
			}

			if len(poll.Content) >= 2 {
				frontContent = append(frontContent, db_model.CardContentNode{
					Element: poll,
				})
			}
		}

		cards = append(cards, db_model.CardNodeContent{

			Front: db_model.CardContentFace{
				Content: frontContent,
			},

			Back: db_model.CardContentFace{
				Content: []db_model.CardContentNode{
					db_model.NewCardTextNode(answer),
				},
			},
		})
	}

	dstPath := strings.TrimSuffix(*srcPath, path.Ext(*srcPath)) + "-output.json"

	output, err := os.Create(dstPath)
	if err != nil {
		fmt.Println("os.Create:", err)
		os.Exit(1)
	}

	defer output.Close()

	name := strings.TrimSpace(*deckName)
	if name == "" {
		name = strings.TrimSuffix(path.Base(*srcPath), path.Ext(*srcPath))
	}

	json.NewEncoder(output).Encode(DeckTemplateBundle{
		Content: DeckTemplate{
			Name:  name,
			Cards: cards,
		},
	})
}

type DeckTemplateBundle struct {
	Content DeckTemplate `json:"content"`
}

type DeckTemplate struct {
	Name        string                     `json:"name"`
	Description string                     `json:"description,omitempty"`
	Cards       []db_model.CardNodeContent `json:"cards"`
}

func (node DeckTemplateBundle) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":    "flippercardtemplate",
		"content": node.Content,
	})
}
