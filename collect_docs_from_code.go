package collect_docs_from_code

import (
	"fmt"
	. "github.com/phannam1412/go-pattern-matching/parser"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Item struct {
	Number   string
	Header   string
	Body     string
	Children []Item
}

func (this *Item) Compare(other *Item) int {
	pieces := strings.Split(this.Number, ".")
	otherPieces := strings.Split(other.Number, ".")
	pieces = pieces[:len(pieces)-1]
	otherPieces = otherPieces[:len(otherPieces)-1]
	max := 0
	if len(pieces) > len(otherPieces) {
		max = len(otherPieces)
	} else {
		max = len(pieces)
	}
	for a := 0; a < max; a++ {
		num1, err := strconv.Atoi(pieces[a])
		panicOnError(err)
		num2, err := strconv.Atoi(otherPieces[a])
		panicOnError(err)
		if num1 > num2 {
			return 1
		}
		if num1 < num2 {
			return -1
		}
	}
	if len(pieces) > len(otherPieces) {
		return 1
	}
	if len(pieces) < len(otherPieces) {
		return -1
	}
	return 0
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func getText(item *Item, level int, forIntro bool) string {
	text := ""
	if forIntro {
		text += strings.Repeat("    ", level) + " " + item.Header + "\n\n"
	} else {
		text += strings.Repeat("#", level+1) + " " + item.Header + "\n\n"
		if item.Body != "" {
			text += item.Body + "\n\n"
		}
	}
	for _, child := range item.Children {
		text += getText(&child, level+1, forIntro)
	}
	return text
}

func Run(paths []string, findExtensions []string, writeResultTo string) {
	parserFactory := func() func(text string) *Res {
		numberWithDot := Combine(Number, Dot)
		somethingBeforeComment := Any(Or(Whitespace, Tab))
		startOfComment := Or(Sharp, Combine(Backsplash, Backsplash))
		header := Combine(somethingBeforeComment, startOfComment, startOfComment, Whitespace, Combine(Some(numberWithDot), Whitespace, TextUntilLineEnd))
		body := AndBut(Combine(somethingBeforeComment, startOfComment, TextUntilLineEnd), header)
		multiline := Combine(header, Any(body))
		formula := FullSearch(multiline, 0)
		return func(text string) *Res {
			tokens := Tokenize(text)
			return formula(tokens, 0)
		}
	}
	var items []*Item
	resolver := func(res *Res) {
		if res == nil {
			println("no result")
			return
		}
		for _, v := range res.Children {
			number := v.Children[0].Children[4].Children[0].Value
			header := strings.Trim(v.Children[0].Children[4].Value, " \n")
			body := ""
			for _, v2 := range v.Children[1].Children {
				body += v2.Children[2].Value
			}
			body = strings.Trim(body, " \n")
			items = append(items, &Item{
				Number:   number,
				Header:   header,
				Body:     body,
				Children: nil,
			})
		}
	}
	parser := parserFactory()

	for _, v := range paths {
		panicOnError(filepath.Walk(v, func(path string, info os.FileInfo, err error) error {
			// ignore if this is directory, we only manipulate on file
			if info.IsDir() {
				return nil
			}
			name := info.Name()
			isInvalidExtension := false
			for _, findForExtension := range findExtensions {
				if len(name) < len("."+findForExtension) {
					continue
				}
				// cut last characters to see if it match with our expected extension
				extension := name[len(name)-len(findForExtension)-1:]
				if extension != "."+findForExtension {
					continue
				}
				isInvalidExtension = true
				break
			}
			if isInvalidExtension == false {
				return nil
			}
			fmt.Printf("===> extracting docs from %s...\n", path)
			b, err := ioutil.ReadFile(path)
			panicOnError(err)
			res := parser(string(b))
			resolver(res)
			return nil
		}))
	}

	// initialize a hash map for quick lookup Item by number 1.1., 2.2. etc...
	fromNumberToItem := map[string]*Item{}
	for _, v := range items {
		fromNumberToItem[v.Number] = v
	}

	// distribute children to corresponding parent
	for _, v := range items {
		number := strings.Split(v.Number, ".")
		// ignore if this item is at level 1
		if len(number) < 3 {
			continue
		}
		// get parent of this item, e.g. 4.2.1 has parent of 4.2.
		parent := strings.Join(number[:len(number)-2], ".") + "."
		if _, ok := fromNumberToItem[parent]; !ok {
			panic(fmt.Errorf("cannot find key %s", parent))
		}
		mparent := fromNumberToItem[parent]
		mparent.Children = append(mparent.Children, *v)
	}

	// sort topics
	sort.Slice(items, func(i, j int) bool {
		return items[i].Compare(items[j]) < 0
	})

	text := ""
	intro := ""
	// only start with items of level 1 and then recursively travel down to children
	for _, v := range items {
		number := strings.Split(v.Number, ".")
		if len(number) >= 3 {
			continue
		}
		intro += getText(v, 0, true)
		text += getText(v, 0, false)
	}
	text = "**WARNING: GENERATED, DO NOT EDIT !!!**\n\n" +
		intro +
		"\n\n" +
		text
	panicOnError(ioutil.WriteFile(writeResultTo, []byte(text), 0777))
}
