package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/maddsua/flippercardapp/utils"
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
	Theme   *CardFaceTheme       `json:"theme,omitempty"`
	Content []CardContentElement `json:"content"`
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

func peekCardContentNodeType(data []byte) (string, error) {

	var ts struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal(data, &ts); err != nil {
		return "", err
	}

	return ts.Type, nil
}

func marshalCardContentNode(node any) ([]byte, error) {

	var typeLabel string

	switch node := node.(type) {
	case BaseCardContentElement:
		typeLabel = node.ContentElementType()
	case BaseCardTextboxElement:
		typeLabel = node.TextElementNodeType()
	default:
		return nil, fmt.Errorf("%T doesn't implement content node", node)
	}

	fields := map[string]json.RawMessage{}

	if data, err := json.Marshal(node); err != nil {
		return nil, fmt.Errorf("marshal content node field: %v", err)
	} else if err = json.Unmarshal(data, &fields); err != nil {
		return nil, fmt.Errorf("marshal content node field: %v", err)
	}

	if typeData, err := json.Marshal(typeLabel); err != nil {
		return nil, fmt.Errorf("marshal content node type: %v", err)
	} else {
		fields["type"] = typeData
	}

	return json.Marshal(fields)
}

type BaseCardContentElement interface {
	ContentElementType() string
}

type CardContentElement struct {
	Element BaseCardContentElement
}

func (element *CardContentElement) Type() string {
	if element.Element == nil {
		return ""
	}
	return element.Element.ContentElementType()
}

func (element CardContentElement) MarshalJSON() ([]byte, error) {
	if element.Element == nil {
		return nil, nil
	}
	return marshalCardContentNode(element.Element)
}

func (element *CardContentElement) UnmarshalJSON(data []byte) (err error) {

	nodeType, err := peekCardContentNodeType(data)
	if err != nil {
		return err
	}

	switch nodeType {
	case "title":
		element.Element, err = utils.DecodeGenericJSON[CardTitleElement](data)
	case "textbox":
		element.Element, err = utils.DecodeGenericJSON[CardTextBoxElement](data)
	case "poll":
		element.Element, err = utils.DecodeGenericJSON[CardPollElement](data)
	default:
		return fmt.Errorf("unsupported element type: '%v'", nodeType)
	}

	return
}

type CardTitleElement struct {
	Content string `json:"content"`
}

func (title CardTitleElement) ContentElementType() string {
	return "title"
}

type CardTextBoxElement struct {
	Content []CardTextboxElement `json:"content"`
}

func (textbox CardTextBoxElement) ContentElementType() string {
	return "textbox"
}

type BaseCardTextboxElement interface {
	TextElementNodeType() string
}

type CardTextboxElement struct {
	Element BaseCardTextboxElement
}

func (element *CardTextboxElement) Type() string {
	if element.Element == nil {
		return ""
	}
	return element.Element.TextElementNodeType()
}

func (element CardTextboxElement) MarshalJSON() ([]byte, error) {
	if element.Element == nil {
		return nil, nil
	}
	return marshalCardContentNode(element.Element)
}

func (element *CardTextboxElement) UnmarshalJSON(data []byte) (err error) {

	nodeType, err := peekCardContentNodeType(data)
	if err != nil {
		return err
	}

	switch nodeType {
	case "text":
		element.Element, err = utils.DecodeGenericJSON[CardTextboxElementTextNode](data)
	case "newline":
		element.Element, err = utils.DecodeGenericJSON[CardTextboxElementNewlineNode](data)
	default:
		return fmt.Errorf("unsupported element type: '%v'", nodeType)
	}

	return
}

type CardTextboxElementTextNode struct {
	Content string                           `json:"content"`
	Theme   *CardTextboxElementTextNodeTheme `json:"theme,omitempty"`
}

func (text CardTextboxElementTextNode) TextElementNodeType() string {
	return "text"
}

type CardTextboxElementTextNodeTheme struct {
	Highlight  *CardTextboxElementTextNodeHighlight `json:"highlight,omitempty"`
	Bold       bool                                 `json:"bold,omitempty"`
	Italic     bool                                 `json:"italic,omitempty"`
	Decoration CardTextboxElementTextDecoration     `json:"decoration,omitempty"`
}

type CardTextboxElementTextNodeHighlight struct {
	TextColor string `json:"text_color"`
	FillColor string `json:"fill_color"`
}

type CardTextboxElementTextDecoration string

const (
	CardTextboxElementTextDecorationUnderline     = CardTextboxElementTextDecoration("underline")
	CardTextboxElementTextDecorationStrikethrough = CardTextboxElementTextDecoration("strikethrough")
)

type CardTextboxElementNewlineNode struct{}

func (newline CardTextboxElementNewlineNode) TextElementNodeType() string {
	return "newline"
}

type CardPollElement struct {
	IsQuiz  bool                        `json:"is_quiz,omitempty"`
	Content []CardPollElementOptionNode `json:"content"`
}

func (poll CardPollElement) ContentElementType() string {
	return "poll"
}

type CardPollElementOptionNode struct {
	Value    string `json:"value"`
	IsAnswer bool   `json:"is_answer,omitempty"`
}
