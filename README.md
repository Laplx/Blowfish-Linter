# Blowfish-Linter

<img src="./bflint_logo.png" alt="bflint_logo" style="zoom: 10%;" />

轻量、可配置、针对 [Blowfish Hugo 主题](https://blowfish.page/) 定制的 Markdown 批量格式规范与增强工具。

## Features

 ✅ 自动补充或修正页面文件的 Front Matter
 ✅ 支持行内与多行公式语法规范替换（自动添加 KaTeX 兼容）
 ✅ 支持子文件夹结构检查与整理

## Usage

选择对应版本的 bflint 可执行程序下载，并将其加入环境变量。在其同目录下编辑配置文件或置于 `$USERPROFILE/.config/bflint/`。

## Params

```bash
bflint [flags] <directory>
```

| 参数            | 作用                                            |
| --------------- | ----------------------------------------------- |
| `-w, --within`  | 仅有 Front Matter 文件头的文件                  |
| `-f, --force`   | 强制覆盖已有 Front Matter 字段                  |
| `-i, --inspect` | 检查并整理子文件夹结构，并自动规范为 `index.md` |
| `-d, --default` | 指定默认文件来源文件夹，配合 `--inspect` 使用   |
| `-h, --help`    | 查看帮助文档                                    |

示例：

```bash
bflint -w -f ./content
bflint -i -d ./defaults ./content
```

##  Configuration

```ini
# This is an example of config.yaml.
# Title, date, and draft are required in front matter.

Front Matter:
  title: H1InMd          # FileName | your_fix_title | (first from top)H1InMd | ""
  date: FileName         # FileName | your_fix_date | Today | ""
  tags: []               # your_fix_tags | [] | null
  authors: ["Laplx"]     # your_fix_authors | [] | null
  category: ""           # your_fix_category | ""
  draft: true            # true | false | ""

Katex:
  inline: true           # true(default: \\[formula\\]) | false | your_symbol(e.g."$$to$$", needs escaping)
  multiline: true        # true(default: $$formula\\\next_formula$$) | false | your_symbol
```

## Notes

- `--inspect` 跳过已为 `index.md` 的文件
- KaTeX 仅通过简单正则匹配 `$...$` 与 `\\`，需避免嵌套复杂场景
- `date` 使用文件名配置时需满足 `YY.MM.DD` 或 `YYYY-MM-DD` 格式
- `draft` 字段未强制覆盖时，若配置为 `true`，仍会强制补充该字段
- 若文档无 H1 标题，`title` 默认回落为文件名

## Further Usage

-  TypeIt 块避免 `>` 引用前缀
-  Markdown 语法一致性检查
-  详细统计报告输出
-  可选日志文件功能

## License

MIT License
