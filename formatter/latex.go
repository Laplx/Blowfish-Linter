package formatter

import (
	"regexp"
)

func FixLatex(content string, inline string, multiline string) string {
	if inline != "false" {
		re := regexp.MustCompile(`\$(.+?)\$`)
		if inline == "true" {
			content = re.ReplaceAllString(content, `\\($1\\)`)
		} else {
			delimit_re := regexp.MustCompile(`(.+?)to(.+?)`)
			if match := delimit_re.FindStringSubmatch(inline); match != nil {
				front, back := match[1], match[2]
				re = regexp.MustCompile(regexp.QuoteMeta(front) + `(.+?)` + regexp.QuoteMeta(back))
				content = re.ReplaceAllString(content, front+`${1}`+back)
			} else {
				content = re.ReplaceAllString(content, `$1`)
			}
		}
	}
	if multiline != "false" {
		re := regexp.MustCompile(`\$\$(.*?)\$\$`)
		content = re.ReplaceAllString(content, `$$$1$$`) // preserve delimiters $$
		re = regexp.MustCompile(`\\{2}`)
		if multiline == "true" {
			content = re.ReplaceAllString(content, `\\\`)
		} else {
			content = re.ReplaceAllString(content, multiline)
		}

	}
	return content
}
