package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/flytam/filenamify"
	"gopkg.in/yaml.v2"
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

	type Rc struct {
		Id int `xml:"id,attr"`
	}

	type Rcs struct {
		Rc []Rc `xml:"rc"`
	}

	type Note struct {
		Value string `xml:",chardata"`
	}

	type Recipe struct {
		Id int `xml:"id,attr"`

		CookTime    int    `xml:"cooktime,attr" yaml:"cookTime"`
		PrepTime    int    `xml:"preptime,attr"`
		Name        string `xml:"name,attr"`
		Source      string `xml:"source,attr"`
		Category    int    `xml:"category,attr"`
		PreHeat     string `xml:"preheat,attr"`
		CreateDate  string `xml:"createdate,attr"`
		SubmittedBy int    `xml:"submittedby,attr"`

		IngredientSets []IngredientSet `xml:"ingredientset"`
		Steps          Steps           `xml:"steps"`
		Notes          []Note          `xml:"note"`

		Rcs Rcs `xml:"rcs"`
	}

	type Recipes struct {
		Recipes []Recipe `xml:"recipe"`
	}

	type Rbook struct {
		Categories Categories `xml:"categories"`
		Users      Users      `xml:"users"`
		Recipes    Recipes    `xml:"recipes"`
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
	fmt.Println()
	fmt.Println("-----------------------")
	fmt.Println()

	// 1: Generate Categories map
	categories := make(map[int]string)
	for _, cat := range rbook.Categories.Categories {
		categories[cat.Id] = cat.Name
	}

	fmt.Println("Categories map:")
	fmt.Println(categories)

	// 2: Generate Users map
	users := make(map[int]string)
	for _, user := range rbook.Users.Users {
		users[user.Id] = user.Name
	}

	fmt.Println("Users map:")
	fmt.Println(users)

	fmt.Println()
	fmt.Println("-----------------------")
	fmt.Println()

	// 3: For every recipe:
	//	- generate a corresponding export-friendly recipe
	//	- generate markdown with front matter
	//	- save to file

	type ExportIngredientSet struct {
		Ingredients []map[string]string `yaml:""`
	}

	type ExportRecipe struct {
		Title      string   `yaml:"title"`
		Author     string   `yaml:"author"`
		Source     string   `yaml:"source"`
		Categories []string `yaml:"categories,flow"`

		CookTime int    `yaml:"cookTime"` //minutes
		PrepTime int    `yaml:"prepTime"`
		PreHeat  string `yaml:"preheat"`
		Date     string `yaml:"date"`

		ExportIngredientSets    []map[string][]map[string]string `yaml:"ingredientsA,omitempty"`
		ExportSimpleIngredients []map[string]string              `yaml:"ingredientsB,omitempty"`
		Steps                   []string                         `yaml:"directions"`
		Notes                   []string                         `yaml:"notes"`
	}

	for _, recipe := range rbook.Recipes.Recipes {
		var exportRecipe ExportRecipe

		exportRecipe.Title = recipe.Name
		exportRecipe.Author = users[recipe.SubmittedBy]
		exportRecipe.Source = recipe.Source

		for _, rc := range recipe.Rcs.Rc {
			exportRecipe.Categories = append(exportRecipe.Categories, categories[rc.Id])
		}

		exportRecipe.CookTime = recipe.CookTime
		exportRecipe.PrepTime = recipe.PrepTime
		exportRecipe.PreHeat = recipe.PreHeat
		exportRecipe.Date = recipe.CreateDate

		if len(recipe.IngredientSets) > 1 {
			exportRecipe.ExportIngredientSets = make([]map[string][]map[string]string, len(recipe.IngredientSets))
			for isetIndex, iset := range recipe.IngredientSets {
				exportRecipe.ExportIngredientSets[isetIndex] = make(map[string][]map[string]string)

				exportSet := make([]map[string]string, 0)

				ing_count := 0
				for _, ingredient := range iset.Ingredients {
					if ingredient.Amount != "" || ingredient.Description != "" {
						exportSet = append(exportSet, make(map[string]string))
						exportSet[ing_count][ingredient.Description] = ingredient.Amount
						ing_count++
					}
				}

				exportRecipe.ExportIngredientSets[isetIndex][iset.Name] = exportSet
			}
		} else if len(recipe.IngredientSets) > 0 {
			exportRecipe.ExportSimpleIngredients = make([]map[string]string, 0)
			ing_count := 0
			for _, ingredient := range recipe.IngredientSets[0].Ingredients {
				if ingredient.Amount != "" || ingredient.Description != "" {
					exportRecipe.ExportSimpleIngredients = append(exportRecipe.ExportSimpleIngredients, make(map[string]string))
					exportRecipe.ExportSimpleIngredients[ing_count][ingredient.Description] = ingredient.Amount

					ing_count++
				}
			}
		}

		exportRecipe.Steps = make([]string, len(recipe.Steps.Steps))
		for i, step := range recipe.Steps.Steps {
			exportRecipe.Steps[i] = step.Value
		}

		exportRecipe.Notes = make([]string, len(recipe.Notes))
		for i, note := range recipe.Notes {
			exportRecipe.Notes[i] = note.Value
		}

		r, err := yaml.Marshal(&exportRecipe)
		if err != nil {
			fmt.Println("error:", err)
		}

		exportFrontMatter := string(r)
		exportFrontMatter = strings.ReplaceAll(exportFrontMatter, "ingredientsA", "ingredients")
		exportFrontMatter = strings.ReplaceAll(exportFrontMatter, "ingredientsB", "ingredients")
		//		fmt.Printf(string(exportFrontMatter))

		exportString := fmt.Sprintf("---\n%s---\n\n", exportFrontMatter)

		os.Mkdir("./export", 0777)

		filename, err := filenamify.Filenamify(exportRecipe.Title, filenamify.Options{})
		filename = strings.ReplaceAll(filename, " ", "-")
		filename = strings.ReplaceAll(filename, ",", "")
		filename = strings.ReplaceAll(filename, "'", "")
		filename = strings.ReplaceAll(filename, "&", "n")
		filename = strings.ReplaceAll(filename, "#", "")

		f, err := os.Create(fmt.Sprintf("./export/%s.md", filename))
		if err != nil {
			panic(err)
		}
		defer f.Close()

		_, err = f.WriteString(exportString)
		if err != nil {
			panic(err)
		}
	}

	// You're done, that was everything.
}
