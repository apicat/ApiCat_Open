package spec

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestParseSpec(t *testing.T) {
	raw, err := os.ReadFile("./testdata/spec.json")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	spec, err := ParseJSON(raw)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	w, err := spec.ToJSON(JSONOption{Indent: "  "})
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	fmt.Println(string(w))
}

func TestDereferenceSchema(t *testing.T) {

	ab, _ := os.ReadFile("./testdata/self_to_self.json")

	source, err := ParseJSON(ab)
	if err != nil {
		fmt.Println(err)
	}

	parent := source.Definitions.Schemas.LookupID(2068)
	sub := source.Definitions.Schemas.LookupID(2332)

	parent.DereferenceSchema(sub)

	bs, _ := json.MarshalIndent(parent, "", " ")

	fmt.Println(string(bs))
}

func TestRemoveSchema(t *testing.T) {
	ab, _ := os.ReadFile("./testdata/self_to_self.json")

	source, err := ParseJSON(ab)
	if err != nil {
		fmt.Println(err)
	}

	parent := source.Definitions.Schemas.LookupID(2068)
	// sub := source.Definitions.Schemas.LookupID(2332)

	// parent.RemoveSchema(sub)

	// bs, _ := json.MarshalIndent(parent, "", " ")

	// fmt.Println(string(bs))

	for _, c := range source.Collections {
		c.RemoveSchema(parent)
	}

	bs, _ := json.MarshalIndent(source, "", " ")

	fmt.Println(string(bs))

}

func TestDereferenceSelf(t *testing.T) {

	ab, _ := os.ReadFile("./testdata/self_to_self.json")

	source, err := ParseJSON(ab)
	if err != nil {
		fmt.Println(err)
	}

	onlySelf := source.Definitions.Schemas.LookupID(2332)

	onlySelf.dereferenceSelf()

	bs, _ := json.MarshalIndent(onlySelf, "", " ")

	fmt.Println(string(bs))

}

func TestResponseRef(t *testing.T) {

	ab, err := os.ReadFile("./testdata/response_ref.json")
	if err != nil {
		fmt.Println(err)
	}
	source, err := ParseJSON(ab)
	if err != nil {
		fmt.Println(err)
	}

	resp := source.Definitions.Responses.LookupID(378)
	for _, c := range source.Collections {
		c.DereferenceResponse(resp)
	}

	bs, _ := json.MarshalIndent(source, "", " ")

	fmt.Println(string(bs))
}

func TestResponseRemove(t *testing.T) {
	ab, err := os.ReadFile("./testdata/response_ref.json")
	if err != nil {
		fmt.Println(err)
	}
	source, err := ParseJSON(ab)
	if err != nil {
		fmt.Println(err)
	}

	resp := source.Definitions.Responses.LookupID(378)

	for _, c := range source.Collections {
		c.RemoveResponse(resp)
	}

	bs, _ := json.MarshalIndent(source, "", " ")

	fmt.Println(string(bs))
}

func TestDereferenceGlobalParameters(t *testing.T) {

	ab, err := os.ReadFile("./testdata/global_excepts.json")
	if err != nil {
		fmt.Println(err)
	}
	source, err := ParseJSON(ab)
	if err != nil {
		fmt.Println(err)
	}

	for _, c := range source.Collections {
		c.DereferenceGlobalParameters("header", source.Globals.Parameters.Header.LookupID(31))
	}

	bs, _ := json.MarshalIndent(source, "", " ")

	fmt.Println(string(bs))

}

func TestAddGlobalParameters(t *testing.T) {

	ab, err := os.ReadFile("./testdata/global_excepts.json")
	if err != nil {
		fmt.Println(err)
	}
	source, err := ParseJSON(ab)
	if err != nil {
		fmt.Println(err)
	}

	for _, c := range source.Collections {
		c.AddParameters("header", source.Globals.Parameters.Header.LookupID(31))
	}

	bs, _ := json.MarshalIndent(source, "", " ")

	fmt.Println(string(bs))

}

func TestOpenGlobalParameters(t *testing.T) {

	ab, err := os.ReadFile("./testdata/global_excepts.json")
	if err != nil {
		fmt.Println(err)
	}
	source, err := ParseJSON(ab)
	if err != nil {
		fmt.Println(err)
	}

	for _, c := range source.Collections {
		c.OpenGlobalParameters("path", source.Globals.Parameters.Path.LookupID(34))
	}

	bs, _ := json.MarshalIndent(source, "", " ")

	fmt.Println(string(bs))

}
