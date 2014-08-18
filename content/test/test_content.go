package test

import (
	"log"
	"testing"

	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/language"
	"github.com/polyglottis/platform/user"
)

var Author = user.Name("test_author")
var Extract = &content.Extract{
	Type:    content.Dialog,
	UrlSlug: "test_slug",
	Metadata: &content.Metadata{
		TargetLanguage: language.English.Code,
	},
	Flavors: content.FlavorMap{
		language.English.Code: content.FlavorByType{
			content.Text: []*content.Flavor{{
				Id:      1,
				Summary: "Test flavor",
				Blocks: content.BlockSlice{{{
					FlavorId:    1,
					BlockId:     1,
					Id:          1,
					ContentType: content.TypeText,
					Content:     "Title",
				}}, {{
					FlavorId:    1,
					BlockId:     2,
					Id:          1,
					ContentType: content.TypeText,
					Content:     "First line",
				}}},
			}},
		},
	},
}

type Tester struct {
	basic  content.BasicServer
	server content.Server
	*testing.T
}

func NewTester(server content.BasicServer, t *testing.T) *Tester {
	tester := &Tester{
		basic: server,
		T:     t,
	}
	if s, ok := server.(content.Server); ok {
		tester.server = s
	}
	return tester
}

func (t *Tester) All() {
	log.Print("Assert extract is missing")
	t.NotExist("")

	log.Println("Assert extract list is empty")
	t.ExtractList(nil)

	log.Print("New extract")
	t.New(Author, Extract)
	id := Extract.Id

	log.Println("Assert extract list contains one element")
	t.ExtractList([]*content.Extract{{
		Id:      Extract.Id,
		Type:    Extract.Type,
		UrlSlug: Extract.UrlSlug,
	}})

	log.Print("Get extract")
	t.Get(id, Extract)

	if t.server != nil {
		log.Print("Get extract id")
		t.GetExtractId(Extract.UrlSlug, Extract.Id)
		t.GetExtractIdFails("no such slug")
	}

	log.Print("New flavor")
	secondFlavor := &content.Flavor{
		ExtractId: Extract.Id,
		Type:      content.Audio,
		Language:  language.English.Code,
		Summary:   "Second flavor",
	}
	t.NewFlavor(Author, secondFlavor)
	Extract.Flavors[language.English.Code][content.Audio] = []*content.Flavor{secondFlavor}
	t.Get(Extract.Id, Extract)

	log.Print("Assert new flavor fails")
	t.NewFlavorFails(Author, &content.Flavor{ExtractId: Extract.Id, Language: language.English.Code})
	t.NewFlavorFails(Author, &content.Flavor{Language: language.English.Code, Type: content.Audio})
	t.NewFlavorFails(Author, &content.Flavor{ExtractId: Extract.Id, Type: content.Audio})

	log.Print("Update extract")
	Extract.Type = content.Article
	Extract.Metadata.Previous = "aaa"
	Extract.Metadata.TargetLanguage = ""
	slug := Extract.UrlSlug
	Extract.UrlSlug = "updated_slug"
	t.Update(Author, Extract)
	Extract.UrlSlug = slug // slug should not change
	t.Get(id, Extract)

	log.Print("Udpate flavor")
	f := Extract.Flavors[language.English.Code][content.Text][0]
	f.Summary = "This is a more elaborate test summary."
	f.LanguageComment = "Colloquial"
	t.UpdateFlavor(Author, f)
	t.Get(id, Extract)

	log.Print("Insert unit")
	thirdUnit := &content.Unit{
		ExtractId:   id,
		Language:    language.English.Code,
		FlavorType:  content.Text,
		FlavorId:    1,
		BlockId:     2,
		Id:          3,
		ContentType: content.TypeText,
		Content:     "Third line.",
	}
	blocks := Extract.Flavors[language.English.Code][content.Text][0].Blocks
	blocks[1] = append(blocks[1], thirdUnit)
	t.InsertOrUpdateUnits(Author, []*content.Unit{thirdUnit})
	t.Get(id, Extract)

	log.Print("Update unit")
	u := Extract.Flavors[language.English.Code][content.Text][0].Blocks[1][0]
	u.Content = "First line of text body."
	u.ContentType = content.TypeFile
	t.InsertOrUpdateUnits(Author, []*content.Unit{u})
	t.Get(id, Extract)

	log.Print("Insert and update units")
	title := Extract.Flavors[language.English.Code][content.Text][0].Blocks[0][0]
	title.Content = "This is the title."
	secondUnit := &content.Unit{
		ExtractId:   id,
		FlavorId:    1,
		Language:    language.English.Code,
		FlavorType:  content.Text,
		BlockId:     2,
		Id:          2,
		ContentType: content.TypeText,
		Content:     "Second line.",
	}
	blocks = Extract.Flavors[language.English.Code][content.Text][0].Blocks
	blocks[1] = []*content.Unit{blocks[1][0], secondUnit, blocks[1][1]}
	t.InsertOrUpdateUnits(Author, []*content.Unit{secondUnit, title})
	t.Get(id, Extract)

	log.Print("Test failures")
	t.UpdateFails(Author, &content.Extract{})
	t.UpdateFlavorFails(Author, &content.Flavor{ExtractId: id})
	t.InsertOrUpdateUnitsFails(Author, []*content.Unit{{ExtractId: id}})
	t.InsertOrUpdateUnitsFails(Author, []*content.Unit{{FlavorId: 1}})

	log.Print("Insert illegal units")
	t.InsertOrUpdateUnitsFails(Author, []*content.Unit{{ExtractId: id, FlavorId: 1, BlockId: 1}})
	t.InsertOrUpdateUnitsFails(Author, []*content.Unit{{ExtractId: id, FlavorId: 1, Id: 3}})
	t.InsertOrUpdateUnitsFails(Author, []*content.Unit{{ExtractId: id, FlavorId: 1, BlockId: 1, Id: 2}})
}

func (t *Tester) NotExist(id content.ExtractId) {
	_, err := t.basic.GetExtract("")
	if err != content.ErrNotFound {
		t.Fatalf("Extract should not exist, but got error '%v'", err)
	}
}

func (t *Tester) New(author user.Name, e *content.Extract) {
	oldId := e.Id
	err := t.basic.NewExtract(author, e)
	if err != nil {
		t.Fatal(err)
	}
	if e.Id == "" {
		t.Fatal("Extract id should not be empty.")
	}
	if e.Id == oldId {
		t.Fatalf("Extract id should have been set (%v vs %v)", e.Id, oldId)
	}
	for lang, fByType := range e.Flavors {
		for fType, flavors := range fByType {
			for _, f := range flavors {
				if f.ExtractId != e.Id {
					t.Fatal("Extract id should have been set on flavor.")
				}
				if f.Language != lang {
					t.Error("Flavor language should match map key.")
				}
				if f.Type != fType {
					t.Error("Flavor type should match map key.")
				}
				for _, block := range f.Blocks {
					for _, unit := range block {
						if unit.ExtractId != e.Id {
							t.Fatal("Extract id should have been set on unit.")
						}
						if unit.Language != lang {
							t.Error("Unit language should have been set.")
						}
						if unit.FlavorType != fType {
							t.Error("FlavorType should have been set on unit.")
						}
					}
				}
			}
		}
	}
}

func (t *Tester) NewFlavor(author user.Name, f *content.Flavor) {
	oldId := f.Id
	err := t.basic.NewFlavor(author, f)
	if err != nil {
		t.Fatal(err)
	}
	if f.Id == 0 {
		t.Fatal("Flavor id should not be zero.")
	}
	if f.Id == oldId {
		t.Fatalf("Flavor id should have been set (%v vs %v)", f.Id, oldId)
	}
	for _, block := range f.Blocks {
		for _, unit := range block {
			if unit.FlavorId != f.Id {
				t.Fatal("Flavor id should have been set on unit.")
			}
		}
	}
}

func (t *Tester) NewFlavorFails(author user.Name, f *content.Flavor) {
	err := t.basic.NewFlavor(author, f)
	if err == nil {
		t.Errorf("NewFlavor should fail for %v", f)
	}
}

func (t *Tester) Get(id content.ExtractId, check *content.Extract) {
	e, err := t.basic.GetExtract(id)
	if err != nil {
		t.Fatal(err)
	}
	if !check.Equals(e) {
		t.Fatalf("These extracts should coincide: %+v != %+v", check, e)
	}
}

func (t *Tester) ExtractList(expected []*content.Extract) {
	list, err := t.basic.ExtractList()
	if err != nil {
		t.Fatal(err)
	}
	if len(expected) != len(list) {
		t.Fatal("Expected list of length %d but got %d.", len(expected), len(list))
	}
	for i, e := range expected {
		other := list[i]
		if e.UrlSlug != other.UrlSlug || e.Type != other.Type || e.Id != other.Id {
			t.Errorf("Id, type and slug should coincide: %v != %v", e, other)
		}
	}
}

func (t *Tester) GetExtractId(slug string, expectedId content.ExtractId) {
	id, err := t.server.GetExtractId(slug)
	if err != nil {
		t.Fatal(err)
	}
	if expectedId != id {
		t.Fatalf("Wrong extract id: want %v but got %v", expectedId, id)
	}
}

func (t *Tester) GetExtractIdFails(nonslug string) {
	_, err := t.server.GetExtractId(nonslug)
	if err != content.ErrNotFound {
		t.Fatal("Expecting a not found error.")
	}
}

func (t *Tester) Update(author user.Name, e *content.Extract) {
	err := t.basic.UpdateExtract(author, e)
	if err != nil {
		t.Error(err)
	}
}

func (t *Tester) UpdateFails(author user.Name, e *content.Extract) {
	err := t.basic.UpdateExtract(author, e)
	if err == nil {
		t.Errorf("UpdateExtract should fail here: %v", e)
	}
}

func (t *Tester) UpdateFlavor(author user.Name, f *content.Flavor) {
	err := t.basic.UpdateFlavor(author, f)
	if err != nil {
		t.Error(err)
	}
}

func (t *Tester) UpdateFlavorFails(author user.Name, f *content.Flavor) {
	err := t.basic.UpdateFlavor(author, f)
	if err == nil {
		t.Errorf("UpdateFlavor should fail here: %v", f)
	}
}

func (t *Tester) InsertOrUpdateUnits(author user.Name, units []*content.Unit) {
	err := t.basic.InsertOrUpdateUnits(author, units)
	if err != nil {
		t.Error(err)
	}
}

func (t *Tester) InsertOrUpdateUnitsFails(author user.Name, units []*content.Unit) {
	err := t.basic.InsertOrUpdateUnits(author, units)
	if err == nil {
		t.Errorf("InsertOrUpdateUnits should fail here: %v", units)
	}
}
