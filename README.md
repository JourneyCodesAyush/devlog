#  devlog - A Terminal Devlog in Go

![Go](https://img.shields.io/badge/Go-1.27-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-green)
![PRs Welcome](https://img.shields.io/badge/PRs-welcome-lightgrey)
![CI](https://github.com/journeycodesayush/devlog/actions/workflows/ci.yml/badge.svg)

---

##  Table of Contents

- [devlog - A Terminal Devlog in Go](#devlog---a-terminal-devlog-in-go)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Quick Install](#quick-install)
    - [Linux / macOS](#linux--macos)
    - [Windows (Command Prompt / PowerShell)](#windows-command-prompt--powershell)
  - [Usage](#usage)
  - [Examples](#examples)
  - [Features](#features)
  - [Dependencies](#dependencies)
  - [Known Limitations](#known-limitations)
  - [Design Philosophy](#design-philosophy)
  - [Project Structure](#project-structure)
  - [Development Guide](#development-guide)
  - [Contributing](#contributing)
  - [LICENSE](#license)
  - [Author](#author)


---

## Overview

**devlog** is a small **terminal devlog** written in Go. Jot down what you did, tag it, list it, search it. Everything is stored in a single local SQLite file.

It started as a learning project: I wanted to find out whether I actually enjoy Go or just had FOMO.

> Designed to stay small: standard library for the CLI, one pure-Go dependency for storage.

---

## Quick Install

```bash
go install github.com/journeycodesayush/devlog@latest
```

Or build from source.

### Linux / macOS

```bash
git clone https://github.com/journeycodesayush/devlog.git
cd devlog
go build .
./devlog help
```

### Windows (Command Prompt / PowerShell)

```powershell
git clone https://github.com/journeycodesayush/devlog.git
cd devlog
go build .
.\devlog.exe help
```

> The SQLite driver is pure Go, so **no C compiler or CGO is needed**.

---

## Usage

```text
devlog add [-t tag] "message"   Add a new entry
devlog list [-n count]          List entries, newest first (all if -n omitted)
devlog search "keyword"         Search entry messages
devlog delete <id>              Delete an entry by id
devlog help                     Show help
```

| Command  | Arguments            | Notes                                               |
| -------- | -------------------- | --------------------------------------------------- |
| `add`    | `[-t tag] "message"` | Exactly one message; the tag is optional            |
| `list`   | `[-n count]`         | `-n 0` or omitted shows everything                  |
| `search` | `"keyword"`          | Case-insensitive for ASCII; matches the message     |
| `delete` | `<id>`               | Id is the number shown in the `ID` column of `list` |

> [!IMPORTANT]
> **Flags go before the message.** `devlog add -t work "message"` works; `devlog add "message" -t work` does not. Go's `flag` package stops parsing at the first non-flag argument.

Data is stored in `~/.devlog/devlog.db` (on Windows, `C:\Users\<you>\.devlog\devlog.db`). The directory and database are created on first use. Timestamps are stored in UTC and displayed in local time.

---

## Examples

```bash
devlog add "fixed the flag parsing bug"
devlog add -t v0.1.0 "tagged first release"
devlog list -n 5
devlog search "flag"
devlog delete 3
```

Output example:

```text
ID  TIME              TAG     PROJECT  MESSAGE
3   2026-10-17 14:02  v0.1.0  -        tagged first release
2   2026-10-17 13:40  -       -        fixed the flag parsing bug
1   2026-10-17 12:41  -       -        new message
```

When nothing is found:

```text
$ devlog search zzzz
No matches.
```

---

## Features

- Add, list, search and delete entries
- Optional tag per entry
- Newest-first listing with an optional limit (`-n`)
- Aligned table output via `text/tabwriter`
- Cross-platform (Windows, Linux, macOS)
- Single local SQLite file, no server, no config
- No CGO required

---

## Dependencies

| Library              | License      | Notes                                          |
| -------------------- | ------------ | ---------------------------------------------- |
| `modernc.org/sqlite` | BSD-3-Clause | Pure-Go SQLite driver, used via `database/sql` |

> Everything else is the Go standard library (`flag`, `database/sql`, `text/tabwriter`, ...).

---

## Known Limitations

- `-t` and other flags must come **before** the message (see Usage).
- A tag can only be set when adding an entry; there is no edit command.
- The `project` column exists in the schema and is shown in `list`, but there is no command to set it yet, so it always shows `-`.
- `search` uses SQL `LIKE`: `%` and `_` in the keyword act as wildcards.
- `delete` on an id that doesn't exist does nothing and reports no error.
- Timestamps have one-second resolution; entries created in the same second are ordered by id.
- Messages containing tabs or newlines can misalign the `list` table.

---

## Design Philosophy

- **Small on purpose.** No edit, export, sync, TUI or colors.
- **Standard library first.** Subcommands use `flag.NewFlagSet`, one per command, instead of a CLI framework like cobra.
- **Pure Go.** A CGO-free SQLite driver keeps builds simple on every platform.
- **Thin `main`.** `main.go` only dispatches commands and prints; data types and storage live in `internal/`.

---

## Project Structure

```bash
devlog/
├── internal/
│   ├── entry/
│   │   └── entry.go       # Entry struct and table fields
│   └── store/
│       └── sqlite.go      # SQLiteStore: Add, List, Search, Delete
├── main.go                # Command dispatch and output
├── go.mod
├── go.sum
├── README.md
└── LICENSE
```

---

## Development Guide

- **Entry**: `internal/entry` defines the `Entry` struct and `Fields()`, which returns the columns for one table row.
- **Store**: `internal/store` wraps `database/sql`; `New(path)` opens the database and creates the schema if missing.
- **CLI**: each subcommand in `main.go` builds its own `FlagSet`, validates arguments, calls the store, and prints.
- **Adding a command**: add it to the early `switch` (so `help` and unknown commands never open the database), then add its `case` below `openStore()`.

Build and run:

```bash
go build .
go vet ./...
gofmt -l .
```

---

## Contributing

> [!NOTE]
> devlog is a small learning project. Bug fixes, documentation and small improvements are welcome. Larger feature requests may not be accepted.

1. Fork the repo

```bash
git clone https://github.com/journeycodesayush/devlog.git
```

2. Create a branch: `feat/something` or `fix/bug`

```bash
git switch -c fix/bug
```

3. Make changes & test locally
4. Open a PR with a clear description

Commit messages follow **Conventional Commits**:

```text
feat: add search command
fix: reject extra positional args in add and list
refactor: extract scanEntries and printEntries helpers
docs: add README
chore: add modernc/sqlite as sqlite driver
```

---

## LICENSE

MIT. See [LICENSE](LICENSE).

---

## Author

Made with ❤️ by [JourneyCodesAyush](https://github.com/JourneyCodesAyush)
