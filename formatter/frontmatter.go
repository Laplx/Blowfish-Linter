package formatter

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/Laplx/Blowfish-Linter/config"

	"gopkg.in/yaml.v3"
)

type FrontMatter struct {
	Title    string   `yaml:"title"`
	Date     string   `yaml:"date"`
	Tags     []string `yaml:"tags,omitempty,flow"`
	Authors  []string `yaml:"authors,omitempty,flow"`
	Category string   `yaml:"category,omitempty"`
	Draft    bool     `yaml:"draft"`
}

func ProcessFile(path string, within bool, force bool, inspect bool) {
	b, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("读取失败: %s\n", path)
		return
	}
	content := string(b)
	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		if within {
			fmt.Printf("无 Front Matter，跳过: %s\n", path)
			return
		}
		emptyFm := FrontMatter{}
		if config.Cfg.TitleMode == "FileName" {
			emptyFm.Title = inferTitle(path, inspect)
		} else if config.Cfg.TitleMode == "H1InMd" {
			h1, newContent := extractH1(content)
			if h1 != "" {
				emptyFm.Title = h1
				content = newContent
			} else {
				emptyFm.Title = inferTitle(path, inspect)
			}
		} else {
			emptyFm.Title = config.Cfg.TitleMode
			if emptyFm.Title == "" {
				emptyFm.Title = inferTitle(path, inspect)
			}
		}
		if config.Cfg.DateMode == "FileName" {
			emptyFm.Date = inferDate(path, inspect)
		} else if config.Cfg.DateMode == "Today" {
			emptyFm.Date = time.Now().Format("2006-01-02")
		} else {
			emptyFm.Date = config.Cfg.DateMode
			if emptyFm.Date == "" {
				emptyFm.Date = time.Now().Format("2006-01-02")
			}
		}
		if config.Cfg.TagsMode != nil {
			emptyFm.Tags = config.Cfg.TagsMode
		}
		if config.Cfg.Authors != nil {
			emptyFm.Authors = config.Cfg.Authors
		}
		if config.Cfg.Category != "" {
			emptyFm.Category = config.Cfg.Category
		}
		if config.Cfg.DraftMode != "" {
			emptyFm.Draft = config.Cfg.DraftMode == "true"
		} else {
			emptyFm.Draft = true
		}

		yamlStr, _ := yaml.Marshal(&emptyFm)
		yamlStr = []byte(forceQuotedStrings(string(yamlStr)))
		newContent := fmt.Sprintf("---\n%s---\n", yamlStr)
		if !strings.Contains(content, "{{< katex >}}") {
			newContent += "{{< katex >}}\n"
		}
		newContent += FixLatex(content, config.Cfg.InlineKatex, config.Cfg.MultiKatex)
		os.WriteFile(path, []byte(newContent), 0644)
		fmt.Printf("添加 Front Matter: %s\n", path)
		return
	}
	if parts[0] != "" {
		parts[2] = fmt.Sprintf("%s---\n%s---\n%s", parts[0], parts[1], parts[2])
	}
	fm := FrontMatter{}
	yaml.Unmarshal([]byte(parts[1]), &fm)
	changed := false

	if (fm.Title == "" || force) && config.Cfg.TitleMode != "" {
		if config.Cfg.TitleMode == "FileName" {
			fm.Title = inferTitle(path, inspect)
		} else if config.Cfg.TitleMode == "H1InMd" {
			h1, newContent := extractH1(content)
			if h1 != "" {
				fm.Title = h1
				content = newContent
			} else {
				fm.Title = inferTitle(path, inspect)
			}
		} else {
			fm.Title = config.Cfg.TitleMode
		}
		changed = true
	}
	if (fm.Date == "" || force) && config.Cfg.DateMode != "" {
		if config.Cfg.DateMode == "FileName" {
			fm.Date = inferDate(path, inspect)
		} else if config.Cfg.DateMode == "Today" {
			fm.Date = time.Now().Format("2006-01-02")
		} else {
			fm.Date = config.Cfg.DateMode
		}
		changed = true
	}
	if (len(fm.Tags) == 0 || force) && config.Cfg.TagsMode != nil {
		fm.Tags = config.Cfg.TagsMode
		changed = true
	}
	if (len(fm.Authors) == 0 || force) && config.Cfg.Authors != nil {
		fm.Authors = config.Cfg.Authors
		changed = true
	}
	if (fm.Category == "" || force) && config.Cfg.Category != "" {
		fm.Category = config.Cfg.Category
		changed = true
	}
	if (!fm.Draft || force) && config.Cfg.DraftMode != "" {
		fm.Draft = config.Cfg.DraftMode == "true"
		changed = true
	}
	if changed {
		yamlStr, _ := yaml.Marshal(&fm)
		yamlStr = []byte(forceQuotedStrings(string(yamlStr)))
		fixed := fmt.Sprintf("---\n%s---\n", yamlStr)
		parts := strings.SplitN(content, "---", 3)
		if !strings.Contains(parts[2], "{{< katex >}}") {
			fixed += "{{< katex >}}\n"
		}
		fixed += FixLatex(parts[2], config.Cfg.InlineKatex, config.Cfg.MultiKatex)
		os.WriteFile(path, []byte(fixed), 0644)
		fmt.Printf("已修正: %s\n", path)
	}
}

func forceQuotedStrings(yamlStr string) string {
	re := regexp.MustCompile(`(?m)^(authors|tags):\s*\[([^\]]*)\]`)
	return re.ReplaceAllStringFunc(yamlStr, func(match string) string {
		parts := re.FindStringSubmatch(match)
		key := parts[1]
		rawList := parts[2]

		items := strings.Split(rawList, ",")
		quoted := make([]string, 0, len(items))
		for _, item := range items {
			item = strings.TrimSpace(item)
			if item != "" {
				if !(strings.HasPrefix(item, `"`) && strings.HasSuffix(item, `"`)) {
					item = fmt.Sprintf(`"%s"`, item)
				}
				quoted = append(quoted, item)
			}
		}
		return fmt.Sprintf("%s: [%s]", key, strings.Join(quoted, ", "))
	})
}

func inferTitle(path string, inspect bool) string {
	parts := strings.Split(path, string(os.PathSeparator))
	filename := strings.TrimSuffix(parts[len(parts)-1], ".md")
	if !inspect {
		return filename
	}
	if len(parts) > 1 {
		return parts[len(parts)-2]
	}
	return filename
}

func inferDate(path string, inspect bool) string {
	filename := inferTitle(path, inspect)
	// YYYY-MM-DD
	re := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})`)
	if match := re.FindStringSubmatch(filename); match != nil {
		return match[1]
	}
	// YY.MM.DD
	re = regexp.MustCompile(`^(\d{2})\.(\d{2})\.(\d{2})`)
	if match := re.FindStringSubmatch(filename); match != nil {
		year := "20" + match[1]
		return fmt.Sprintf("%s-%s-%s", year, match[2], match[3])
	}
	return time.Now().Format("2006-01-02")
}

func extractH1(content string) (string, string) {
	lines := strings.Split(content, "\n")
	newLines := []string{}
	h1Found := ""
	h1Skipped := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !h1Skipped && strings.HasPrefix(trimmed, "# ") {
			h1Found = strings.TrimSpace(trimmed[2:])
			h1Skipped = true
			continue
		}
		newLines = append(newLines, line)
	}

	return h1Found, strings.Join(newLines, "\n")
}
