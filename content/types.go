// Package content contains the type definitions for contents on the platform.
// It also defines the content server interface.
package content

import (
	"errors"
	"time"

	"github.com/polyglottis/platform/language"
	"github.com/polyglottis/platform/user"
)

// Sever is the interface a content server should comply to.
type Server interface {
	BasicServer

	// Read
	GetExtractId(slug string) (ExtractId, error)
}

type BasicServer interface {
	// Create
	NewExtract(author user.Name, e *Extract) error // sets e's ExtractId
	NewFlavor(author user.Name, f *Flavor) error   // sets f's FlavorId

	// Read
	ExtractList() ([]*Extract, error)
	GetExtract(id ExtractId) (*Extract, error)

	// Update
	UpdateExtract(author user.Name, e *Extract) error // doesn't update e.Flavors
	UpdateFlavor(author user.Name, f *Flavor) error   // doesn't update f.Units
	InsertOrUpdateUnits(author user.Name, units []*Unit) error

	// Delete (=invalidate)
}

var ErrNotFound = errors.New("Not found")
var ErrInvalidInput = errors.New("Invalid input")

type ExtractId string
type FlavorId int
type BlockId int
type UnitId int

type ExtractType string
type FlavorType string
type ContentType string

const (
	Article    ExtractType = "article"
	Dialog     ExtractType = "dialog"
	ShortStory ExtractType = "short_story"
	Poem       ExtractType = "poem"
	Song       ExtractType = "song"
	WordList   ExtractType = "word_list"
)

var AllExtractTypes = []ExtractType{
	Dialog,
	WordList,
	ShortStory,
	Article,
	Poem,
	Song,
}

const (
	Audio      FlavorType = "AUDIO"
	Text       FlavorType = "TEXT"
	Transcript FlavorType = "TRANSCRIPT"
)

const (
	TypeText ContentType = "text"
	TypeFile ContentType = "file"
)

type EditType string

const (
	EditNew    EditType = "new"
	EditUpdate EditType = "update"
	EditDelete EditType = "delete"
)

type Version struct {
	Author   user.Name
	Time     time.Time
	Number   int
	EditType EditType
}

type Metadata struct {
	SourceUrl      string        `json:",omitempty"`
	Previous       ExtractId     `json:",omitempty"`
	Next           ExtractId     `json:",omitempty"`
	TargetLanguage language.Code `json:",omitempty"`
}

// A Unit is the finest content unit on the platform.
// It corresponds to a line of text, an image, or an utterance.
type Unit struct {
	ExtractId   ExtractId
	FlavorId    FlavorId
	Language    language.Code
	FlavorType  FlavorType
	BlockId     BlockId
	Id          UnitId
	ContentType ContentType // text, table header, table row, file, ...
	Content     string
}

// A Flavor is a typically a short text in a certain language.
// The content Units are grouped into blocks.
type Flavor struct {
	ExtractId       ExtractId
	Id              FlavorId
	Summary         string
	Type            FlavorType // text, audio, ...
	Language        language.Code
	LanguageComment string
	Blocks          BlockSlice `json:",omitempty"`
}

type BlockSlice []UnitSlice
type UnitSlice []*Unit

// An Extract is the coarsest semantic unit on the platform.
// It corresponds to a group of Flavors describing the same content.
type Extract struct {
	Id       ExtractId
	Type     ExtractType
	UrlSlug  string
	Metadata *Metadata `json:",omitempty"`
	Flavors  FlavorMap `json:",omitempty"`
}

type FlavorMap map[language.Code]FlavorByType
type FlavorByType map[FlavorType][]*Flavor

type VersionedUnit struct {
	*Unit
	*Version
}

type VersionedFlavor struct {
	*Flavor
	*Version
}

type VersionedExtract struct {
	*Extract
	*Version
}

type VersionedUnitRef struct {
	ExtractId     ExtractId
	FlavorId      FlavorId
	BlockId       BlockId
	UnitId        UnitId
	VersionNumber int
}

type NoteType string

const (
	Vocabular     NoteType = "vocabular"
	Grammar       NoteType = "grammar"
	Pronunciation NoteType = "pronunciation"
	Culture       NoteType = "culture"
)

type NoteId int

type NotePointer struct {
	Unit  *VersionedUnitRef
	Start int
	End   int
	Type  NoteType
	Id    NoteId
	Notes []*Note
}

type Note struct {
	Unit     *VersionedUnitRef
	NoteId   NoteId
	Language language.Code
	Content  string
}

type ExtractShape []int
