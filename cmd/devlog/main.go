package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/journeycodesayush/devlog/internal/store"
)

func printHelp() {
	fmt.Println(`devlog - a terminal-based devlog tool

Usage:
  devlog <command> [arguments]

Commands:
  add [-t tag] "message"      Add a new entry
  list [-n count]             List entries (latest first; all if -n omitted)
  search "keyword"            Search entries by keyword
  delete <id>                 Delete an entry by id

Run 'devlog help' to see this message.`)
}

func openStore() (*store.SQLiteStore, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("finding home directory: %w", err)
	}

	dbDir := filepath.Join(homeDir, ".devlog")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("creating %s: %w", dbDir, err)
	}

	return store.New(filepath.Join(dbDir, "devlog.db"))
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add", "list", "search", "delete":
	case "help":
		printHelp()
		os.Exit(0)
	default:
		printHelp()
		os.Exit(1)
	}

	st, err := openStore()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		tag := addCmd.String("t", "", "tag for the entry")
		addCmd.Parse(os.Args[2:])

		if addCmd.NArg() != 1 {
			fmt.Println("Usage: devlog add [-t tag] \"message\"")
			os.Exit(1)
		}

		message := addCmd.Arg(0)
		if err := st.Add(message, *tag); err != nil {
			fmt.Println("Error adding entry: ", err)
			os.Exit(1)
		}
		fmt.Println("Entry added successfully")
	case "list":
		listCmd := flag.NewFlagSet("list", flag.ExitOnError)
		n := listCmd.Int("n", 0, "show only the latest n entries (0 = all)")

		listCmd.Parse(os.Args[2:])

		if listCmd.NArg() > 0 || *n < 0 {
			fmt.Println("Usage: devlog list [-n count]")
			os.Exit(1)
		}

		list, err := st.List(*n)
		if err != nil {
			fmt.Println("Error listing entries:", err)
			os.Exit(1)
		}

		if len(list) == 0 {
			fmt.Println("No entries.")
			break
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTIME\tTAG\tPROJECT\tMESSAGE")
		for _, ent := range list {
			fmt.Fprintln(w, strings.Join(ent.Fields(), "\t"))
		}
		w.Flush()
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
	case "search":
		searchCmd := flag.NewFlagSet("search", flag.ExitOnError)

		searchCmd.Parse(os.Args[2:])

		if searchCmd.NArg() != 1 {
			fmt.Println("Usage: devlog search \"keyword\"")
			os.Exit(1)
		}

		keyword := searchCmd.Arg(0)
		matches, err := st.Search(keyword)
		if err != nil {
			fmt.Println("Error searching entries:", err)
			os.Exit(1)
		}

		if len(matches) == 0 {
			fmt.Println("No matches.")
			break
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTIME\tTAG\tPROJECT\tMESSAGE")
		for _, ent := range matches {
			fmt.Fprintln(w, strings.Join(ent.Fields(), "\t"))
		}
		w.Flush()
	case "help":
		printHelp()
	default:
		printHelp()
		os.Exit(1)
	}
}
