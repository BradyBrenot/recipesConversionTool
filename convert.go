package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {

	type Category struct {
		Id   int    `xml:"id,attr"`
		Name string `xml:"name,attr"`
	}

	type Categories struct {
		Categories []Category `xml:"category"`
	}

	type User struct {
		Id       int    `xml:"id,attr"`
		Email    string `xml:"email,attr"`
		Username string `xml:"username,attr"`
		Name     string `xml:"name,attr"`
	}

	type Users struct {
		Users []User `xml:"user"`
	}

	type Ingredient struct {
		Amount      string `xml:"amount,attr"`
		Description string `xml:"description,attr"`
	}

	type IngredientSet struct {
		Id          int          `xml:"id,attr"`
		Name        string       `xml:"name,attr"`
		Ingredients []Ingredient `xml:"ingredient"`
	}

	type Step struct {
		Value string `xml:",chardata"`
	}

	type Steps struct {
		Steps []Step `xml:"step"`
	}

	type Note struct {
		Value string `xml:",chardata"`
	}

	type Recipe struct {
		Id int `xml:"id,attr"`

		CookTime    int    `xml:"cooktime,attr"`
		PrepTime    int    `xml:"preptime,attr"`
		Name        string `xml:"name,attr"`
		Source      string `xml:"source,attr"`
		Category    string `xml:"category,attr"`
		PreHeat     string `xml:"preheat,attr"`
		CreateDate  string `xml:"createdate,attr"`
		SubmittedBy int    `xml:"submittedby,attr"`

		IngredientSets []IngredientSet `xml:"ingredientset"`
		Steps          Steps           `xml:"steps"`
		Notes          []Note          `xml:"note"`
	}

	type Recipes struct {
		Recipes []Recipe `xml:"recipe"`
	}

	type Rbook struct {
		Categories []Categories `xml:"categories"`
		Users      []Users      `xml:"users"`
		Recipes    []Recipes    `xml:"recipes"`
	}

	// 0: Parse the rbook XML doc
	xmlFile, err := os.Open("export-1427076510.xml")
	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var rbook Rbook
	err = xml.Unmarshal(byteValue, &rbook)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(rbook.Recipes)

	// 1: Generate Categories map

	// 2: Generate Users map

	// 3: For every recipe, generate a corresponding markdown file

	// You're done, that was everything.
}
