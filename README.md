A Minimal Static Site Generator in Go
🧭 Overview
StaticGen is a fast, extensible, and minimal static site generator written in Go, inspired by tools like Hugo and Jekyll, designed for developers and content creators who prefer full control over their site architecture and build process.

The generator processes Markdown content files, parses frontmatter metadata, and renders HTML pages using customizable Go HTML templates.

🎯 Objectives
Build a fully functional static site generator in Go.

Support Markdown content with YAML frontmatter.

Provide a clean, extensible project architecture.

Enable easy theming and layout management via templates.

Generate clean, production-ready HTML output to a public/ folder.

Optional: CLI interface for building, previewing, and deploying sites.

Milestone | Description | ETA

✅ M1 | Project initialization, go.mod, directory structure | Day 1

✅ M2 | Parse .md files and render simple HTML pages | Day 2

✅ M3 | Integrate html/template for layout rendering | Day 3

🟡 M4 | Parse YAML frontmatter for metadata injection | Day 4

🔜 M5 | Implement global config (config.yaml) support | Day 5

🔜 M6 | Add clean CLI commands (build, serve) | Day 6–7

🔜 M7 | Add index page generation (index.html) | Day 8

🔜 M8 | Basic tagging system with tag index pages | Day 9–10

⏳ M9 | Live preview dev server | TBD

⏳ M10 | Deploy-ready GitHub repo + gh-pages instructions | Final day