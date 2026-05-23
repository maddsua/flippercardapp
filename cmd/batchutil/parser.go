package main

import (
	"fmt"
	"strings"

	db_model "github.com/maddsua/flippercardapp/db/model"
)

func cardTextNodeFromString(text string) db_model.CardContentNode {

	var spans []db_model.CardTextboxNode

	for idx, content := range strings.Split(text, "\n") {

		spans = append(spans, &db_model.CardTextboxTextNode{
			Content: content,
		})

		if idx > 0 {
			spans = append(spans, &db_model.CardTextboxNewlineNode{})
		}
	}

	return &db_model.CardTextBoxNode{
		Content: spans,
	}
}

func cardPollNodeFromString(text string) (db_model.CardContentNode, error) {

	var poll db_model.CardPollNode

	for idx, option := range strings.Split(text, ",") {

		if option = strings.TrimSpace(option); option == "" {
			continue
		}

		poll.Content = append(poll.Content, db_model.CardPollNodeOption{
			Value:    option,
			IsAnswer: idx == 0,
		})
	}

	if len(poll.Content) < 2 {
		return nil, fmt.Errorf("invalid poll answers: only %d node(s) parsed", len(poll.Content))
	}

	return &poll, nil
}
