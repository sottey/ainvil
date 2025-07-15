# ainvil

**ainvil** is a Go-based command-line tool for standardizing, processing, and exporting data from multiple pendant recording sources.  

It supports merging different export formats into a consistent JSON structure for easy archiving, search, and downstream processing.

---

## ✨ Features

- Process text exports from **Omi** and **Bee** pendant formats (more to come!)
- Fetch and export lifelogs from the **Limitless** API
- Unified, consistent JSON schema
- Organized date-based output folders
- Incremental sync: auto-detects last saved lifelog and fetches only new entries
- Extensible design for new source types

---

## 📦 Folder Structure

```
ainvil/
  main.go
  cmd/
    root.go
    omi.go
    bee.go
    limitless.go
  common/
    processor.go
  out/
    (generated JSON output)
```

---

## ⚡️ Installation

**Clone the repo:**

```bash
git clone https://github.com/YOUR_USERNAME/ainvil.git
cd ainvil
```

**Build:**

```bash
go build -o ainvil
```

Or install directly:

```bash
go install github.com/YOUR_USERNAME/ainvil@latest
```

---

## 🚀 Usage

After building, you can run:

```
./ainvil [command] [flags]
```

### 📌 Global structure

Each command writes standardized JSON files under:

```
out/
  YYYY/
    MM/
      DD/
        [sourceType]_[ID].json
```

---

### 🎯 Commands

#### 1️⃣ omi

Process Omi pendant `.txt` files into standardized JSON.

```bash
ainvil omi --source ./omi_exports --out ./out
```

**Flags:**

- `--source` *(required)*: Directory containing `.txt` files.
- `--out`: Output root directory (default `./out`).

---

#### 2️⃣ bee

Process Bee pendant `.txt` files into standardized JSON.

```bash
ainvil bee --source ./bee_exports --out ./out
```

**Flags:**

- `--source` *(required)*: Directory containing `.txt` files.
- `--out`: Output root directory (default `./out`).

---

#### 3️⃣ limitless

Fetch lifelogs from the Limitless API.

```bash
ainvil limitless --token YOUR_API_KEY --url API_ENDPOINT --out ./out
```

**Flags:**

- `--token` *(required)*: Your Limitless API key.
- `--url` *(required)*: Base API URL.
- `--start`: Start date (RFC3339). Optional. If omitted, tool auto-detects last saved date from existing JSON.
- `--end`: End date (RFC3339). Optional.
- `--out`: Output root directory (default `./out`).

**Example with incremental sync:**

```bash
ainvil limitless --token YOUR_API_KEY --url https://api.limitless.ai/v1/logs
```

If you omit `--start`, it will scan existing output for the most recent saved lifelog and fetch only newer ones.

---

## 🗂 Output Example

```
out/
  2025/
    05/
      30/
        omi_abc123.json
  2025/
    05/
      30/
        bee_xyz456.json
  2025/
    05/
      30/
        limitless_789def.json
```

Each JSON contains a standardized structure:

```json
{
  "id": "...",
  "sourceType": "...",
  "startTime": "...",
  "endTime": "...",
  "title": "...",
  "overview": "...",
  "transcript": "...",
  "contents": [
    { "type": "heading1", "content": "..." },
    ...
  ],
  "exportDate": "...",
  "exportVersion": "ainvil 1.0.0",
  "sourceFile": "..."
}
```

---

## 💡 Planned / Suggested Features

- Config file support (Viper)
- --dry-run mode for previews
- Logging levels (--verbose / --quiet)
- Index file generation
- List / Show subcommands
- Unit tests for parsing

---

## 🛠 Development

**Run locally:**

```bash
go run main.go [command] [flags]
```

**Format code:**

```bash
go fmt ./...
```
