package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

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
	Theme   *CardFaceTheme    `json:"theme,omitempty"`
	Content []CardContentNode `json:"content"`
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

type BaseCardContentElement interface {
	ContentElementType() string
}

type CardContentNode struct {
	Element BaseCardContentElement
}

func (node *CardContentNode) Type() string {
	if node.Element == nil {
		return ""
	}
	return node.Element.ContentElementType()
}

func (node CardContentNode) MarshalJSON() ([]byte, error) {

	if node.Element == nil {
		return nil, nil
	}

	fields, err := utils.ExtractStructJSONFields(node.Element)
	if err != nil {
		return nil, fmt.Errorf("extract struct fields: %v", err)
	}

	fields["type"] = node.Element.ContentElementType()

	return json.Marshal(fields)
}

func (node *CardContentNode) UnmarshalJSON(data []byte) (err error) {

	nodeType, err := utils.ExtractJSONField[string](data, "type")
	if err != nil {
		return err
	}

	switch nodeType {
	case "title":
		node.Element, err = utils.DecodeGenericJSON[CardTitleElement](data)
	case "image":
		node.Element, err = utils.DecodeGenericJSON[CardImageElement](data)
	case "textbox":
		node.Element, err = utils.DecodeGenericJSON[CardTextBoxElement](data)
	case "poll":
		node.Element, err = utils.DecodeGenericJSON[CardPollElement](data)
	default:
		return fmt.Errorf("unsupported element type: '%v'", nodeType)
	}

	return
}

func NewCardTitleNode(title string) CardContentNode {
	return CardContentNode{
		Element: CardTitleElement{
			Content: title,
		},
	}
}

type CardTitleElement struct {
	Content string `json:"content"`
}

func (title CardTitleElement) ContentElementType() string {
	return "title"
}

func NewCardTextNode(text string) CardContentNode {

	var spans []CardTextboxElement

	for idx, content := range strings.Split(text, "\n") {

		spans = append(spans, CardTextboxElement{
			Element: CardTextboxElementTextNode{
				Content: content,
			},
		})

		if idx > 0 {
			spans = append(spans, CardTextboxElement{
				Element: CardTextboxElementNewlineNode{},
			})
		}
	}

	return CardContentNode{
		Element: CardTextBoxElement{
			Content: spans,
		},
	}
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

	fields, err := utils.ExtractStructJSONFields(element.Element)
	if err != nil {
		return nil, fmt.Errorf("extract struct fields: %v", err)
	}

	fields["type"] = element.Element.TextElementNodeType()

	return json.Marshal(fields)
}

func (element *CardTextboxElement) UnmarshalJSON(data []byte) (err error) {

	nodeType, err := utils.ExtractJSONField[string](data, "type")
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

type CardImageElement struct {
	MediaID string `json:"media_id,omitempty"`
}

func (textbox CardImageElement) ContentElementType() string {
	return "image"
}
