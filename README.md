# OmniData CLI

> 🧠 A forward-thinking, lightweight **universal data translator and inspector** written in Go.
> Work with heterogeneous structured data formats (CSV, JSON, XML, SQL, Excel, Parquet, Avro, and beyond) directly from the terminal.

[![CI](https://github.com/asearer/OmniData/actions/workflows/ci.yml/badge.svg)](https://github.com/asearer/OmniData/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/asearer/OmniData)](https://goreportcard.com/report/github.com/asearer/OmniData)
[![Go Version](https://img.shields.io/github/go-mod/go-version/asearer/OmniData)](go.mod)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

---

## ✨ Features

* 🔄 **Convert between formats**: CSV ↔ JSON ↔ XML ↔ XLSX ↔ SQL ↔ Parquet ↔ Avro
* 👀 **Inspect data quickly**: `peek` command to view schema and top rows with statistics
* 📊 **Schema detection & stats**: Automatic type inference and column statistics
* 🔍 **Schema diffing**: Compare schemas between two files with `diff` command
* 🗄️ **SQL support**: Query database tables directly and convert results to any supported format
* 📄 **Multiple output formats**: Export schema and diffs as Markdown, HTML, or JSON
* ⚡ **Streaming mode**: Process large files efficiently with `--stream` flag
* ⚡ **Fast & memory-efficient**: Stream large files using Go’s native IO
* 📦 **Portable**: single binary, no runtime dependencies
* 🧩 **Extensible architecture**: Add new formats or outputs with minimal code changes

---

## 🚀 Installation & Build

### Using Go

```bash
go install github.com/asearer/OmniData@latest
```

### Clone and Build

```bash
git clone https://github.com/asearer/OmniData.git
cd OmniData
go build -o omnidata
```

### Docker

```bash
docker build -t omnidata .
docker run --rm -v "$(pwd)":/data omnidata convert /data/input.csv /data/output.json
```

---

## 🛠️ CI/CD

- Automated CI is powered by **GitHub Actions**! Each push or PR runs:
  - Full Go build/lint/test
  - Docker build (verifies Docker image compiles successfully)
- See the latest results and workflow run details above.

---

## 🧰 Usage

### Convert CSV → JSON

```bash
./omnidata convert -i data.csv -o data.json --from csv --to json
```

### Convert JSON → CSV

```bash
./omnidata convert -i data.json -o data.csv --from json --to csv
```

### Query SQL Database

```bash
./omnidata query -d "postgres://user:pass@localhost/db" -q "SELECT * FROM users" -o users.json --to json
```

### Using STDIN/STDOUT

```bash
cat data.csv | ./omnidata convert - - --from csv --to json > data.json
```

### Dry-Run Mode

```bash
./omnidata convert -i data.csv -o data.json --from csv --to json --dry-run
```

### Peek Command

```bash
./omnidata peek -i data.csv --format csv
./omnidata peek -i data.json --format json --rows 10 --stats
./omnidata peek -i data.csv --format csv --output-format markdown -o schema.md
```

### Diff Command

```bash
./omnidata diff -1 old.csv -2 new.csv --format1 csv --format2 csv
./omnidata diff -1 schema1.json -2 schema2.json --format1 json --format2 json --output-format html -o diff.html
```

### Streaming Mode

```bash
./omnidata convert -i large.csv -o large.json --from csv --to json --stream
```

---

## 📂 Project Structure

```
OmniData/
├── main.go
├── go.mod
├── go.sum
├── README.md
├── cmd/
│   ├── root.go
│   ├── convert.go
│   ├── diff.go
│   └── peek.go
├── internal/
│   ├── convert/
│   │   ├── registry.go
│   │   ├── runner.go
│   │   └── validator.go
│   ├── formats/
│   │   ├── avro.go
│   │   ├── csv.go
│   │   ├── json.go
│   │   ├── parquet.go
│   │   ├── sql.go
│   │   ├── xlsx.go
│   │   └── xml.go
│   ├── inspect/
│   │   ├── diff.go
│   │   ├── peek.go
│   │   └── schema.go
│   ├── output/
│   │   └── formatters.go
│   └── stream/
│       └── reader.go
├── tests/
│   ├── convert/
│   │   ├── registry_test.go
│   │   ├── runner_test.go
│   │   └── validator_test.go
│   └── formats/
│       ├── csv_test.go
│       ├── json_test.go
│       ├── xlsx_test.go
│       └── xml_test.go
```

---

## 📋 Supported Formats & Commands

| Format  | Convert | Peek | Diff | Query |
| ------- | :-----: | :--: | :--: | :---: |
| CSV     |    ✅    |   ✅  |   ✅  |   ❌   |
| JSON    |    ✅    |   ✅  |   ✅  |   ❌   |
| XML     |    ✅    |   ✅  |   ✅  |   ❌   |
| XLSX    |    ✅    |   ✅  |   ✅  |   ❌   |
| SQL     |    ✅    |   ❌  |   ❌  |   ✅   |
| Parquet |    ✅    |   ✅  |   ✅  |   ❌   |
| Avro    |    ✅    |   ✅  |   ✅  |   ❌   |

**Legend:** ✅ = supported, ❌ = not supported / not applicable

---

## ⚙️ Command vs Feature Matrix

| Command   | Input Formats                            | Output Formats                           | Flags & Options                                                                     | Notes                                                          |
| --------- | ---------------------------------------- | ---------------------------------------- | ----------------------------------------------------------------------------------- | -------------------------------------------------------------- |
| `convert` | CSV, JSON, XML, XLSX, SQL, Parquet, Avro | CSV, JSON, XML, XLSX, SQL, Parquet, Avro | `--from <format>` `--to <format>` `--stream` `--dry-run` `-i` `-o`                  | Supports streaming for large datasets; dry-run previews output |
| `peek`    | CSV, JSON, XML, XLSX, Parquet, Avro      | Markdown, HTML, JSON                     | `--rows <n>` `--stats` `--output-format <format>` `-i` `-o`                         | Preview schema + top rows; includes column stats               |
| `diff`    | CSV, JSON, XML, XLSX, Parquet, Avro      | Markdown, HTML, JSON                     | `-1` `-2` `--format1 <format>` `--format2 <format>` `--output-format <format>` `-o` | Compares schemas between two files                             |
| `query`   | SQL databases                            | CSV, JSON, XML, XLSX, Parquet, Avro      | `-d <db-connection>` `-q <query>` `--to <format>` `-o`                              | Execute SQL queries and convert results to supported formats   |

---

## 💡 Why OmniData? (Real-World Use Cases)

OmniData is a universal CLI alternative to **jq**, **csvkit**, **xlsx2csv**, and other specialized tools:

* 🌐 **Universal** → one CLI for multiple formats
* 🧩 **Interoperable** → mix-and-match sources easily
* ⚡ **Efficient** → Go-powered speed & memory efficiency
* 🔮 **Forward-thinking** → supports SQL, XML, Excel, Parquet, Avro, and schema diffing

### Example Scenarios

* **Data Engineer** → Convert CSV logs to JSON for ELK/Kafka pipelines
* **Developer** → Convert API JSON responses into CSV or Excel
* **Analyst** → Merge multiple file types into unified datasets
* **Ops** → Automate conversions and database queries in CI/CD pipelines

---

## ⚙️ Architecture

OmniData is modular and extensible, built on Go’s Cobra CLI framework.

```
cmd/       → CLI commands (root, convert, query, inspect)
internal/  → Core logic & converters
pkg/       → Optional reusable modules
main.go    → Program entrypoint
```

**Converter Interface Example:**

```go
type Converter interface {
    Read(input io.Reader) (DataSet, error)
    Write(output io.Writer, data DataSet) error
}
```

Adding new formats? Implement `Converter` and register in `registry`.

---

## 🧪 Testing

Run tests:

```bash
go test ./... -v
```

---

## 📌 Quick Reference CLI Examples


### Convert CSV → JSON
----
```
./omnidata convert -i data.csv -o data.json --from csv --to json
```
### Convert JSON → CSV
----
```
./omnidata convert -i data.json -o data.csv --from json --to csv
```
### Query SQL database
----
```
./omnidata query -d "postgres://user:pass@localhost/db" -q "SELECT * FROM users" -o users.json --to json
```
### Peek first 10 rows and stats
----
```
./omnidata peek -i data.csv --format csv --rows 10 --stats
```
### Compare schemas between two files
----
```
./omnidata diff -1 old.csv -2 new.csv --format1 csv --format2 csv --output-format html -o diff.html
```
### Streaming conversion for large files
----
```
./omnidata convert -i large.csv -o large.json --from csv --to json --stream
```

### Dry-run conversion
----
```
./omnidata convert -i data.csv -o data.json --from csv --to json --dry-run
```

---

💡 Tips & Best Practices

Use --stream for large files to avoid memory issues; works for CSV, JSON, Parquet, and Avro

Peek before converting to verify schema and data types using ./omnidata peek -i file

SQL queries: Use parameterized queries for large datasets and export directly to a supported format (--to csv/json/xlsx)

Schema diffing: Use diff to check changes before merging new datasets

Dry-run mode: Use --dry-run to preview conversions without writing output

Pipeline-friendly: OmniData works with STDIN/STDOUT, making it easy to chain commands in scripts