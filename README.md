# 🛠️ ainvil

**Ainvil** is a command-line tool for exporting lifelogs from AI wearable pendants such as the [Limitless Pendant](https://www.limitless.ai/) and [Omi Pendant](https://omi.ai/). It normalizes and saves the data into a structured folder tree for easy access, analysis, and backup.

The name combines **AI** and **Anvil** — a nod to forging unified insight from fragmented sources.

---

## ✨ Features

- 📦 Supports multiple APIs (`omi`, `limitless`) with a modular architecture
- 🧠 Normalizes content into a single unified structure
- 🗃️ Exports entries into `./export/YYYY/MM/DD/entry.json`
- 🔁 Supports full or incremental export
- ⚙️ Configurable via flags or JSON config file
- 🧱 Written in Go, portable and dependency-free

---

## 🔧 Installation

### Option 1: From Compiled Release

_(To be added when releases are published)_

```bash
curl -LO https://github.com/sottey/ainvil/releases/latest/download/ainvil
chmod +x ainvil
./ainvil --help
```

### Option 2: From Source

```bash
git clone https://github.com/sottey/ainvil.git
cd ainvil
go build -o ainvil
./ainvil --help
```

Optional: use `build.sh` to include version info.

---

## 🚀 Usage

### Basic Export (Limitless)
```bash
./ainvil export --api limitless --token sk-xxxxx --start 2025-06-01 --end 2025-06-03
```

### Use Config File (`~/.ainvil.json`)
```json
{
  "api": "limitless",
  "token": "sk-xxxxx"
}
```

Then just run:
```bash
./ainvil export --start 2025-06-01 --end 2025-06-03
```

### Export Everything
```bash
./ainvil export --all
```

### Overwrite Existing Files
```bash
./ainvil export --all --full
```

### Check Version
```bash
./ainvil version
```

---

## 🛣️ Roadmap

- [x] Limitless export support (via API)
- [x] Date-based and full export options
- [x] Pluggable API client architecture
- [ ] Limitless export support (via API)
- [ ] Real-time streaming/log tailing
- [ ] Markdown/HTML rendering of lifelogs
- [ ] CSV and SQLite export
- [ ] Self-hosted dashboard to browse logs

---

## 🔐 Configuration

Ainvil supports configuration via:

1. `--config path/to/config.json`
2. `~/.ainvil.json` (default path)
3. Environment variables:
   - `AINVIL_API`
   - `AINVIL_TOKEN`

---

## 👨‍💻 Developer Guide

Clone and build:
```bash
git clone https://github.com/sottey/ainvil.git
cd ainvil
go build -o ainvil
```

Test:
```bash
./ainvil version
./ainvil export --start 2025-06-01 --end 2025-06-03
```

Use `build.sh` for versioned builds:
```bash
./build.sh
```

---

## 📄 License

MIT License — see [LICENSE](LICENSE)

---

## 🙏 Credits

- Developed by [sottey](https://github.com/sottey)
- Inspired by the growing ecosystem of AI wearables and the need for user-owned data

---

## 🧠 Contributing

Contributions welcome! Open issues, suggest features, or submit PRs.

---

## 🗂 Project Structure

```
ainvil/
├── cmd/              # Cobra CLI commands
├── config/           # Viper config loading
├── internal/storage/ # File writing logic
├── lib/
│   ├── api/          # API client logic
│   ├── model/        # Unified data model
│   └── util/         # Helpers
├── main.go
├── build.sh
└── .ainvil.json      # (optional) default config
```
