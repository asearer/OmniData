# OmniData CLI

**OmniData** is a forward-thinking, lightweight **universal data translator and inspector** written in Go.
It helps developers, data engineers, and analysts work with **heterogeneous structured data** formats (CSV, JSON, XML, SQL, and beyond) directly from the terminal.

---

## ✨ Features (MVP)

* 🔄 **Convert between formats**: CSV ↔ JSON (XML, Excel, and more planned).
* 👀 **Inspect data quickly**: `peek` command (coming soon) to view schema and top rows.
* 📊 **Schema detection & stats** (future roadmap).
* ⚡ **Fast & memory-efficient**: streams large files using Go’s native IO.
* 📦 **Portable**: single binary, no runtime dependencies.
* 🧩 **Extensible architecture**: add new formats with minimal code changes.

---

## 🚀 Usage

### Build

```bash
git clone https://github.com/asearer/omnidata.git
cd omnidata
go build -o omnidata
```

### Convert CSV → JSON

```bash
./omnidata convert -i data.csv -o data.json --from csv --to json
```

### Convert JSON → CSV

```bash
./omnidata convert -i data.json -o data.csv --from json --to csv
```

### Using STDIN/STDOUT

```bash
cat data.csv | ./omnidata convert - - --from csv --to json > data.json
```

### Dry-Run Mode

Preview what would happen without writing files:

```bash
./omnidata convert -i data.csv -o data.json --from csv --to json --dry-run
```

---

## 📂 Project Structure

```
OmniData/
├── main.go              # CLI entrypoint
├── cmd/                 # CLI command definitions
│   ├── root.go          # Base CLI with global flags
│   └── convert.go       # "convert" command (modular & extensible)
└── internal/            # Internal logic
    └── data/
        └── data.go      # Conversion implementation (plug-and-play converters)
```

---

## 💡 Why OmniData? (Real-World Use Cases)

Other tools like **jq** or **csvkit** are excellent, but they are specialized:

* `jq` → great for JSON but doesn’t handle CSV or XML.
* `csvkit` → powerful for CSV but doesn’t easily bridge to other formats.
* Many engineers end up stitching multiple tools together.

OmniData aims to be:

* 🌐 **Universal** → one CLI that understands multiple formats.
* 🧩 **Interoperable** → mix-and-match sources without friction.
* ⚡ **Efficient** → built in Go for speed and low memory usage.
* 🔮 **Forward-thinking** → roadmap includes XML, Excel, SQL, Parquet, and schema diffing.

---

### Example Scenarios

* **Data Engineer** → Convert CSV logs into JSON for ingestion into ELK or Kafka pipelines.
* **Developer** → Inspect JSON API responses and convert into CSV for spreadsheet analysis.
* **Analyst** → Merge multiple heterogeneous files (CSV + JSON + XML) into a unified dataset.
* **Ops** → Automate conversions in CI/CD pipelines with a single portable binary.

---

## 🛠️ Roadmap

* [ ] Add `peek` command → show first rows + column stats.
* [ ] Add SQL support → query database tables directly.
* [ ] Support XML, Excel, Parquet, and Avro.
* [ ] Schema diffing between files.
* [ ] Pluggable output formats → Markdown, HTML reports, JSON summaries.
* [ ] Streaming & memory optimizations for large datasets.



