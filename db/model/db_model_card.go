package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type CardNodeContent struct {
	Front CardContentFace `json:"front"`
	Back  CardContentFace `json:"back"`
}

func (perms CardNodeContent) Value() (driver.Value, error) {
	return json.Marshal(perms)
}

func (perms *CardNodeContent) Scan(src any) error {
	switch src := src.(type) {
	case []byte:
		return json.Unmarshal(src, perms)
	case string:
		return json.Unmarshal([]byte(src), perms)
	default:
		return fmt.Errorf("unable to scan %T into CardNodeContent", src)
	}
}

type CardContentFace struct {
	Theme   *CardFaceTheme      `json:"theme,omitempty"`
	Content CardContentNodeList `json:"content"`
}

type CardContentNode interface {
	CardContentNodeType() string
}

type CardContentNodeBase struct {
	Type string `json:"type"`
}

type CardContentNodeList []CardContentNode

func (list CardContentNodeList) MarshalJSON() ([]byte, error) {

	for _, node := range list {
		switch node := node.(type) {
		case *CardTitleNode:
			node.CardContentNodeBase.Type = node.CardContentNodeType()
		case *CardTextBoxNode:
			node.CardContentNodeBase.Type = node.CardContentNodeType()
		case *CardPollNode:
			node.CardContentNodeBase.Type = node.CardContentNodeType()
			node.IsQuiz = node.IsReallyQuiz()
		case *CardImageNode:
			node.CardContentNodeBase.Type = node.CardContentNodeType()
		}
	}

	return json.Marshal([]CardContentNode(list))
}

func (list *CardContentNodeList) UnmarshalJSON(data []byte) (err error) {

	var entries []json.RawMessage
	if err := json.Unmarshal(data, &entries); err != nil {
		return fmt.Errorf("decode slice: %v", err)
	}

	*list = make(CardContentNodeList, len(entries))

	for idx, nodeData := range entries {

		var base CardContentNodeBase
		if err := json.Unmarshal(nodeData, &base); err != nil {
			return fmt.Errorf("decode base node: %v", err)
		}

		var node CardContentNode

		switch base.Type {
		case "title":
			node = &CardTitleNode{}
		case "image":
			node = &CardImageNode{}
		case "textbox":
			node = &CardTextBoxNode{}
		case "poll":
			node = &CardPollNode{}
		default:
			return fmt.Errorf("unsupported element type: '%v'", base.Type)
		}

		if err := json.Unmarshal(nodeData, node); err != nil {
			return fmt.Errorf("decode node: %v", err)
		}

		(*list)[idx] = node
	}

	return
}

type CardFaceTheme struct {
	Card         *CardCanvasTheme         `json:"card,omitempty"`
	Interactives *CardContentElementTheme `json:"interactives,omitempty"`
}

type CardContentElementTheme struct {
	FillColor string `json:"fill_color,omitempty"`
	MaskColor string `json:"mask_color,omitempty"`
}

type CardCanvasTheme struct {
	CardContentElementTheme
	OutlineColor string `json:"outline_color,omitempty"`
}

type CardTitleNode struct {
	CardContentNodeBase
	Content string `json:"content"`
}

func (elem *CardTitleNode) CardContentNodeType() string {
	return "title"
}

type CardTextBoxNode struct {
	CardContentNodeBase
	Content CardTextboxNodeList `json:"content"`
}

type CardTextboxNode interface {
	CardTextboxNodeType() string
}

type CardTextboxNodeList []CardTextboxNode

func (list CardTextboxNodeList) MarshalJSON() ([]byte, error) {

	for _, node := range list {
		switch node := node.(type) {
		case *CardTextboxTextNode:
			node.CardContentNodeBase.Type = node.CardTextboxNodeType()
		case *CardTextboxNewlineNode:
			node.CardContentNodeBase.Type = node.CardTextboxNodeType()
		}
	}

	return json.Marshal([]CardTextboxNode(list))
}

func (list *CardTextboxNodeList) UnmarshalJSON(data []byte) (err error) {

	var entries []json.RawMessage
	if err := json.Unmarshal(data, &entries); err != nil {
		return fmt.Errorf("decode slice: %v", err)
	}

	*list = make(CardTextboxNodeList, len(entries))

	for idx, nodeData := range entries {

		var base CardContentNodeBase
		if err := json.Unmarshal(nodeData, &base); err != nil {
			return fmt.Errorf("decode base node: %v", err)
		}

		var node CardTextboxNode

		switch base.Type {
		case "text":
			node = &CardTextboxTextNode{}
		case "newline":
			node = &CardTextboxNewlineNode{}
		default:
			return fmt.Errorf("unsupported element type: '%v'", base.Type)
		}

		if err := json.Unmarshal(nodeData, node); err != nil {
			return fmt.Errorf("decode node: %v", err)
		}

		(*list)[idx] = node
	}

	return
}

func (elem *CardTextBoxNode) CardContentNodeType() string {
	return "textbox"
}

type CardTextboxTextNode struct {
	CardContentNodeBase
	Content string                   `json:"content"`
	Theme   *CardTextboxElementTheme `json:"theme,omitempty"`
}

func (text *CardTextboxTextNode) CardTextboxNodeType() string {
	return "text"
}

type CardTextboxElementTheme struct {
	Highlight  *CardTextboxElementTextHighlight `json:"highlight,omitempty"`
	Bold       bool                             `json:"bold,omitempty"`
	Italic     bool                             `json:"italic,omitempty"`
	Decoration CardTextboxElementDecoration     `json:"decoration,omitempty"`
}

type CardTextboxElementTextHighlight struct {
	TextColor string `json:"text_color"`
	FillColor string `json:"fill_color"`
}

type CardTextboxElementDecoration string

const (
	CardTextboxElementDecorationUnderline     = CardTextboxElementDecoration("underline")
	CardTextboxElementDecorationStrikethrough = CardTextboxElementDecoration("strikethrough")
)

type CardTextboxNewlineNode struct {
	CardContentNodeBase
}

func (newline *CardTextboxNewlineNode) CardTextboxNodeType() string {
	return "newline"
}

type CardPollNode struct {
	CardContentNodeBase
	IsQuiz  bool                 `json:"is_quiz,omitempty"`
	Content []CardPollNodeOption `json:"content"`
}

func (poll *CardPollNode) CardContentNodeType() string {
	return "poll"
}

func (poll *CardPollNode) IsReallyQuiz() bool {
	for _, opt := range poll.Content {
		if opt.IsAnswer {
			return true
		}
	}
	return false
}

type CardPollNodeOption struct {
	Value    string `json:"value"`
	IsAnswer bool   `json:"is_answer,omitempty"`
}

type CardImageNode struct {
	CardContentNodeBase
	MediaID string `json:"media_id,omitempty"`
}

func (image *CardImageNode) CardContentNodeType() string {
	return "image"
}
