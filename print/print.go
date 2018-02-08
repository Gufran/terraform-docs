package print

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Gufran/terraform-docs/doc"
)

// Markdown prints the given doc as markdown.
func Markdown(d *doc.Doc, printRequired bool) (string, error) {
	var buf bytes.Buffer

	if len(d.Intro) > 0 {
		for _, v := range d.Intro {
			buf.WriteString(v + "\n\n")
		}
	}

	if len(d.Resources) > 0 || len(d.ResourcesIntro) > 0 {
		buf.WriteString("\n### Resources\n\n")
	}

	if len(d.ResourcesIntro) > 0 {
		for _, i := range d.ResourcesIntro {
			buf.WriteString(i + "\n\n")
		}
	}

	if len(d.Resources) > 0 {
		buf.WriteString("| Type | Name |\n")
		buf.WriteString("|------|------|\n")
		for _, v := range d.Resources {
			buf.WriteString(fmt.Sprintf("| <div id='resource.%s'></div> `%s` | %s |\n",
				fmt.Sprintf("%s.%s", v.Name, v.Label),
				v.Name,
				v.Label))

			desc := "No Description"
			if v.Description != "" {
				desc = normalizeMarkdownDesc(v.Description)
			}
			buf.WriteString("| <div class='description'></div> " + desc + " | <div class='empty'></div> |\n")
		}

		buf.WriteString("\n")
	}

	buf.WriteString("\n")

	if len(d.DataProviders) > 0 || len(d.DataProvidersIntro) > 0 {
		buf.WriteString("\n### Data Providers\n\n")
	}

	if len(d.DataProvidersIntro) > 0 {
		for _, i := range d.DataProvidersIntro {
			buf.WriteString(i + "\n\n")
		}
	}

	if len(d.DataProviders) > 0 {
		buf.WriteString("| Type | Name |\n")
		buf.WriteString("|------|------|\n")
		for _, v := range d.Resources {
			buf.WriteString(fmt.Sprintf("| <div id='data.%s'></div> `%s` | %s |\n",
				fmt.Sprintf("%s.%s", v.Name, v.Label),
				v.Name,
				v.Label))

			desc := "No Description"
			if v.Description != "" {
				desc = normalizeMarkdownDesc(v.Description)
			}
			buf.WriteString("| <div class='description'></div> " + desc + " | <div class='empty'></div> |\n")
		}

		buf.WriteString("\n")
	}

	buf.WriteString("\n")

	if len(d.IamPolicies) > 0 || len(d.IamPoliciesIntro) > 0 {
		buf.WriteString("\n### IAM Policies\n\n")
	}

	if len(d.IamPoliciesIntro) > 0 {
		for _, i := range d.IamPoliciesIntro {
			buf.WriteString(i + "\n\n")
		}
	}

	if len(d.IamPolicies) > 0 {
		for _, v := range d.IamPolicies {
			buf.WriteString("##### " + v.Name + "\n\n")
			buf.WriteString("!!! quote \"" + v.Name + " policy document\"\n")
			for _, l := range strings.Split(v.Description, "\n") {
				buf.WriteString("    " + l + "\n")
			}

			buf.WriteString("    ``` json\n")
			for _, l := range strings.Split(v.Policy, "\n") {
				buf.WriteString("    " + l + "\n")
			}
			buf.WriteString("    ```\n")
		}
	}

	if len(d.Inputs) > 0 || len(d.InputsIntro) > 0 {
		buf.WriteString("\n### Inputs\n\n")
	}

	if len(d.InputsIntro) > 0 {
		for _, i := range d.InputsIntro {
			buf.WriteString(i + "\n\n")
		}
	}

	if len(d.Inputs) > 0 {
		buf.WriteString("| Name | Description | Type |\n")
		buf.WriteString("|------|-------------|:----:|\n")
	}

	for _, v := range d.Inputs {
		def := v.Value()

		if def == "required" {
			def = "-"
		} else {
			def = fmt.Sprintf("`%s`", def)
		}

		buf.WriteString(fmt.Sprintf("| <div id='var.%s'></div> %s | %s | %s |\n",
			v.Name,
			v.Name,
			normalizeMarkdownDesc(v.Description),
			v.Type))
	}

	if len(d.Outputs) > 0 || len(d.OutputsInto) > 0 {
		buf.WriteString("\n### Outputs\n\n")
	}

	if len(d.OutputsInto) > 0 {
		for _, i := range d.OutputsInto {
			buf.WriteString(i + "\n\n")
		}
	}

	if len(d.Outputs) > 0 {
		buf.WriteString("| Name | Description |\n")
		buf.WriteString("|------|-------------|\n")
	}

	for _, v := range d.Outputs {
		buf.WriteString(fmt.Sprintf("| <div id='output.%s'></div> %s | %s |\n",
			v.Name,
			v.Name,
			normalizeMarkdownDesc(v.Description)))
	}

	return buf.String(), nil
}

// JSON prints the given doc as json.
func JSON(d *doc.Doc) (string, error) {
	s, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return "", err
	}

	return string(s), nil
}

// Humanize the given `v`.
func humanize(def *doc.Value) string {
	if def == nil {
		return "yes"
	}

	return "no"
}

// normalizeMarkdownDesc fixes line breaks in descriptions for markdown:
//
//  * Double newlines are converted to <br><br>
//  * A second pass replaces all other newlines with spaces
func normalizeMarkdownDesc(s string) string {
	return strings.Replace(strings.Replace(strings.TrimSpace(s), "\n\n", "<br><br>", -1), "\n", " ", -1)
}
