package doc

import (
	"bytes"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/printer"
)

// Input represents a terraform input variable.
type Input struct {
	Name        string
	Description string
	Default     *Value
	Type        string
}

// Value returns the default value as a string.
func (i *Input) Value() string {
	if i.Default != nil {
		switch i.Default.Type {
		case "string":
			return i.Default.Literal
		case "map":
			return "<map>"
		case "list":
			return "<list>"
		}
	}

	return "required"
}

// Value represents a terraform value.
type Value struct {
	Type    string
	Literal string
}

// Output represents a terraform output.
type Output struct {
	Name        string
	Description string
}

// Resource represents a terraform resource.
type Resource struct {
	Name        string
	Label       string
	Description string
}

// Resource respresents a terraform data provider.
type DataProvider struct {
	Name        string
	Label       string
	Description string
}

// Resource represents a terraform aws iam policy document.
type IamPolicy struct {
	Name        string
	Policy      string
	Description string
}

// Doc represents a terraform module doc.
type Doc struct {
	Intro              []string
	InputsIntro        []string
	Inputs             []Input
	OutputsInto        []string
	Outputs            []Output
	ResourcesIntro     []string
	Resources          []Resource
	DataProvidersIntro []string
	DataProviders      []DataProvider
	IamPoliciesIntro   []string
	IamPolicies        []IamPolicy
}

type inputsByName []Input

func (a inputsByName) Len() int           { return len(a) }
func (a inputsByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a inputsByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type outputsByName []Output

func (a outputsByName) Len() int           { return len(a) }
func (a outputsByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a outputsByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type resourceByName []Resource

func (a resourceByName) Len() int           { return len(a) }
func (a resourceByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a resourceByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type dataProviderByName []DataProvider

func (a dataProviderByName) Len() int           { return len(a) }
func (a dataProviderByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a dataProviderByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type iamPoliciesByName []IamPolicy

func (a iamPoliciesByName) Len() int           { return len(a) }
func (a iamPoliciesByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a iamPoliciesByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

// Create creates a new *Doc from the supplied map
// of filenames and *ast.File.
func Create(files map[string]*ast.File) *Doc {
	doc := new(Doc)

	for _, f := range files {
		list := f.Node.(*ast.ObjectList)
		doc.Intro = append(doc.Intro, intro(f.Comments, "")...)
		doc.Inputs = append(doc.Inputs, inputs(list)...)
		doc.InputsIntro = append(doc.InputsIntro, intro(f.Comments, "input")...)
		doc.Outputs = append(doc.Outputs, outputs(list)...)
		doc.OutputsInto = append(doc.OutputsInto, intro(f.Comments, "output")...)
		doc.Resources = append(doc.Resources, resources(list)...)
		doc.ResourcesIntro = append(doc.ResourcesIntro, intro(f.Comments, "resource")...)
		doc.DataProviders = append(doc.DataProviders, dataProviders(list)...)
		doc.DataProvidersIntro = append(doc.DataProvidersIntro, intro(f.Comments, "data-provider")...)
		doc.IamPolicies = append(doc.IamPolicies, iamPolicies(list)...)
		doc.IamPoliciesIntro = append(doc.IamPoliciesIntro, intro(f.Comments, "iam-policy")...)
	}

	sort.Sort(inputsByName(doc.Inputs))
	sort.Sort(outputsByName(doc.Outputs))
	sort.Sort(resourceByName(doc.Resources))
	sort.Sort(dataProviderByName(doc.DataProviders))
	sort.Sort(iamPoliciesByName(doc.IamPolicies))
	return doc
}

func intro(comments []*ast.CommentGroup, category string) []string {
	ret := []string{}
	for _, c := range comments {
		data := comment(c.List)
		lines := strings.SplitN(data, "\n", 2)
		if len(lines) < 2 {
			continue
		}

		if strings.HasPrefix(lines[0], "@doc("+category+")") {
			ret = append(ret, lines[1])
		}
	}

	return ret
}

func iamPolicies(list *ast.ObjectList) []IamPolicy {
	var ret []IamPolicy

	for _, item := range list.Items {
		if !is(item, "data") {
			continue
		}

		if unquote(item.Keys[1].Token.Text) == "aws_iam_policy_document" {
			policyBuf := bytes.NewBufferString("")
			name, _ := strconv.Unquote(item.Keys[2].Token.Text)
			printer.Fprint(policyBuf, item.Val)

			comm := ""
			if item.LeadComment != nil {
				comm = comment(item.LeadComment.List)
			}

			ret = append(ret, IamPolicy{
				Name:        name,
				Description: comm,
				Policy:      policyBuf.String(),
			})
		}
	}

	return ret
}

func resources(list *ast.ObjectList) []Resource {
	var ret []Resource

	for _, item := range list.Items {
		if is(item, "resource") {
			name, _ := strconv.Unquote(item.Keys[1].Token.Text)
			label, _ := strconv.Unquote(item.Keys[2].Token.Text)
			comm := ""

			if item.LeadComment != nil {
				comm = comment(item.LeadComment.List)
			}

			ret = append(ret, Resource{
				Name:        name,
				Label:       label,
				Description: comm,
			})
		}
	}

	return ret
}

func dataProviders(list *ast.ObjectList) []DataProvider {
	var ret []DataProvider

	for _, item := range list.Items {
		if is(item, "data") {
			name, _ := strconv.Unquote(item.Keys[1].Token.Text)
			label, _ := strconv.Unquote(item.Keys[2].Token.Text)
			comm := ""

			if item.LeadComment != nil {
				comm = comment(item.LeadComment.List)
			}

			ret = append(ret, DataProvider{
				Name:        name,
				Label:       label,
				Description: comm,
			})
		}
	}

	return ret
}

// Inputs returns all variables from `list`.
func inputs(list *ast.ObjectList) []Input {
	var ret []Input

	for _, item := range list.Items {
		if is(item, "variable") {
			name, _ := strconv.Unquote(item.Keys[1].Token.Text)
			if name == "" {
				name = item.Keys[1].Token.Text
			}
			items := item.Val.(*ast.ObjectType).List.Items
			var desc string
			switch {
			case description(items) != "":
				desc = description(items)
			case item.LeadComment != nil:
				desc = comment(item.LeadComment.List)
			}

			var itemsType = get(items, "type")
			var itemType string

			if itemsType == nil || itemsType.Literal == "" {
				itemType = "string"
			} else {
				itemType = itemsType.Literal
			}

			def := get(items, "default")
			ret = append(ret, Input{
				Name:        name,
				Description: desc,
				Default:     def,
				Type:        itemType,
			})
		}
	}

	return ret
}

// Outputs returns all outputs from `list`.
func outputs(list *ast.ObjectList) []Output {
	var ret []Output

	for _, item := range list.Items {
		if is(item, "output") {
			name, _ := strconv.Unquote(item.Keys[1].Token.Text)
			if name == "" {
				name = item.Keys[1].Token.Text
			}
			items := item.Val.(*ast.ObjectType).List.Items
			var desc string
			switch {
			case description(items) != "":
				desc = description(items)
			case item.LeadComment != nil:
				desc = comment(item.LeadComment.List)
			}

			ret = append(ret, Output{
				Name:        name,
				Description: strings.TrimSpace(desc),
			})
		}
	}

	return ret
}

// Get `key` from the list of object `items`.
func get(items []*ast.ObjectItem, key string) *Value {
	for _, item := range items {
		if is(item, key) {
			v := new(Value)

			if lit, ok := item.Val.(*ast.LiteralType); ok {
				if value, ok := lit.Token.Value().(string); ok {
					v.Literal = value
				} else {
					v.Literal = lit.Token.Text
				}
				v.Type = "string"
				return v
			}

			if _, ok := item.Val.(*ast.ObjectType); ok {
				v.Type = "map"
				return v
			}

			if _, ok := item.Val.(*ast.ListType); ok {
				v.Type = "list"
				return v
			}

			return nil
		}
	}

	return nil
}

// description returns a description from items or an empty string.
func description(items []*ast.ObjectItem) string {
	if v := get(items, "description"); v != nil {
		return v.Literal
	}

	return ""
}

// Is returns true if `item` is of `kind`.
func is(item *ast.ObjectItem, kind string) bool {
	if len(item.Keys) > 0 {
		return item.Keys[0].Token.Text == kind
	}

	return false
}

// Unquote the given string.
func unquote(s string) string {
	s, _ = strconv.Unquote(s)
	return s
}

// Comment cleans and returns a comment.
func comment(l []*ast.Comment) string {
	var line string
	var ret string

	for _, t := range l {
		line = strings.TrimSpace(t.Text)
		line = strings.TrimPrefix(line, "#")
		line = strings.TrimPrefix(line, "//")
		line = strings.TrimPrefix(line, "/**")
		line = strings.TrimPrefix(line, "/*")
		line = strings.TrimSuffix(line, "**/")
		line = strings.TrimSuffix(line, "*/")
		ret += strings.TrimSpace(line) + "\n"
	}

	return ret
}

// Header returns the header comment from the list
// or an empty comment. The head comment must start
// at line 1 and start with `/**`.
func header(c *ast.CommentGroup) (comment string) {
	if len(c.List) == 0 {
		return comment
	}

	if c.Pos().Line != 1 {
		return comment
	}

	cm := strings.TrimSpace(c.List[0].Text)

	if strings.HasPrefix(cm, "/**") {
		lines := strings.Split(cm, "\n")

		if len(lines) < 2 {
			return comment
		}

		lines = lines[1 : len(lines)-1]
		for _, l := range lines {
			l = strings.TrimSpace(l)
			switch {
			case strings.TrimPrefix(l, "* ") != l:
				l = strings.TrimPrefix(l, "* ")
			default:
				l = strings.TrimPrefix(l, "*")
			}
			comment += l + "\n"
		}
	}

	return comment
}
