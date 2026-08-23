package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/journeycodesayush/devlog/internal/store"
)

func printHelp() {
	fmt.Println(`devlog - a terminal-based devlog tool

Usage:
  devlog <command> [arguments]

Commands:
  add "message" [-t tag]      Add a new entry
  list [--today] [--since 7d] List entries
  search "keyword"            Search entries by keyword
  delete <id>                 Delete an entry by id

Run 'devlog help' to see this message.`)
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Unable to access home directory")
		os.Exit(1)
	}

	dbDir := filepath.Join(homeDir, ".devlog")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		fmt.Println("Unable to create devlog directory:", err)
		os.Exit(1)
	}

	dbPath := filepath.Join(dbDir, "devlog.db")

	st, err := store.New(dbPath)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		tag := addCmd.String("t", "", "tag for the entry")
		addCmd.Parse(os.Args[2:])

		if addCmd.NArg() < 1 {
			fmt.Println("Usage: devlog add \"message\" [-t tag]")
			os.Exit(1)
		}

		message := addCmd.Arg(0)
		_ = tag

		if err := st.Add(message); err != nil {
			fmt.Println("Error adding entry: ", err)
			os.Exit(1)
		}
		fmt.Println("Entry added successfully")
	case "list":
		listCmd := flag.NewFlagSet("list", flag.ExitOnError)
		listCmd.Bool("today", false, "logs saved today")

		listCmd.Parse(os.Args[2:])
		list, err := st.List()

		if err != nil {
			printHelp()
			os.Exit(1)
		}
		_ = list

		for _, ent := range list {
			fmt.Println(&ent)
		}

	case "delete":
		deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
		deleteCmd.Parse(os.Args[2:])

		if deleteCmd.NArg() < 1 {
			fmt.Println("Usage: devlog delete <id>")
			os.Exit(1)
		}

		id, err := strconv.Atoi(deleteCmd.Arg(0))
		if err != nil {
			fmt.Println("Invalid id:", deleteCmd.Arg(0))
			os.Exit(1)
		}

		err = st.Delete(id)
		if err != nil {
			fmt.Println("Error deleting message with id: ", id)
			os.Exit(1)
		}
	case "help":
		printHelp()
	default:
		printHelp()
		os.Exit(1)
	}
}
