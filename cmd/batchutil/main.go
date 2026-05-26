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
		} else if len(row) == 0 || (len(row) == 1 && row[0] == "") {
			continue
		}

		if mapper == nil {
			mapper = utils.NewRecordMapper(row)
			continue
		}

		mappedRow := mapper.WithRow(row)

		next := db_model.CardNodeContent{
			Front: db_model.CardContentFace{
				Content: parseFrontContent(mappedRow),
			},
			Back: db_model.CardContentFace{
				Content: parseBackContent(mappedRow),
			},
		}

		if len(next.Front.Content) == 0 {
			fmt.Println("Empty row/card skipped")
			continue
		}

		cards = append(cards, next)
	}

	if len(cards) == 0 {
		fmt.Println("No cards have been imported")
		os.Exit(1)
	} else {
		fmt.Println("imported", len(cards), "cards")
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
