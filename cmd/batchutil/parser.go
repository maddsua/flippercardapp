package main

import (
	"fmt"
	"strings"

	db_model "github.com/maddsua/flippercardapp/db/model"
	"github.com/maddsua/flippercardapp/utils"
)

func parseFrontContent(row *utils.RecordRow) []db_model.CardContentNode {

	var nodes []db_model.CardContentNode

	if title, _ := row.Get("front_title"); title != "" {
		nodes = append(nodes, &db_model.CardTitleNode{Content: title})
	}

	if text, _ := row.Get("front_textarea"); text != "" {
		nodes = append(nodes, cardTextNodeFromString(text))
	}

	if opts, _ := row.Get("front_poll"); opts != "" {
		if poll, err := cardPollNodeFromString(opts); err != nil {
			fmt.Println("WARN: Empty poll options; Skipped")
		} else {
			nodes = append(nodes, poll)
		}
	}

	if len(nodes) == 0 {
		fmt.Println("WARN: Front side has no content nodes")
	}

	return nodes
}

func parseBackContent(row *utils.RecordRow) []db_model.CardContentNode {

	var nodes []db_model.CardContentNode

	if title, _ := row.Get("back_title"); title != "" {
		nodes = append(nodes, &db_model.CardTitleNode{Content: title})
	}

	if text, _ := row.Get("back_textarea"); text != "" {
		nodes = append(nodes, cardTextNodeFromString(text))
	}

	if len(nodes) == 0 {
		fmt.Println("WARN: Back side has no content nodes")
	}

	return nodes
}

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
