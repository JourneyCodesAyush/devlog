package store

import (
	"path/filepath"
	"slices"
	"testing"

	"github.com/journeycodesayush/devlog/internal/entry"
)

func newTestStore(t *testing.T) *SQLiteStore {
	t.Helper()
	st, err := New(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	return st
}

func mustAdd(t *testing.T, st *SQLiteStore, message, tag string) {
	t.Helper()
	if err := st.Add(message, tag); err != nil {
		t.Fatalf("Add(%q, %q): %v", message, tag, err)
	}
}

func messages(entries []entry.Entry) []string {
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		out = append(out, e.Message)
	}
	return out
}

func TestAddStoresTagAsNullWhenEmpty(t *testing.T) {
	st := newTestStore(t)
	mustAdd(t, st, "with tag", "work")
	mustAdd(t, st, "without tag", "")

	got, err := st.List(0)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d entries, want 2", len(got))
	}

	// newest first: "without tag" is index 0
	if got[0].Tag.Valid {
		t.Errorf("entry without tag: Tag.Valid = true, want false")
	}
	if !got[1].Tag.Valid || got[1].Tag.String != "work" {
		t.Errorf("entry with tag: Tag = %+v, want valid \"work\"", got[1].Tag)
	}
}

func TestListOrderNewestFirst(t *testing.T) {
	st := newTestStore(t)
	for _, m := range []string{"one", "two", "three"} {
		mustAdd(t, st, m, "")
	}

	got, err := st.List(0)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	want := []string{"three", "two", "one"}
	if !slices.Equal(messages(got), want) {
		t.Errorf("order = %v, want %v", messages(got), want)
	}
}

func TestListLimit(t *testing.T) {
	st := newTestStore(t)
	for _, m := range []string{"a", "b", "c", "d", "e"} {
		mustAdd(t, st, m, "")
	}

	tests := []struct {
		name  string
		limit int
		want  []string
	}{
		{"limit 2 returns two newest", 2, []string{"e", "d"}},
		{"limit 0 returns all", 0, []string{"e", "d", "c", "b", "a"}},
		{"limit above count returns all", 10, []string{"e", "d", "c", "b", "a"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := st.List(tt.limit)
			if err != nil {
				t.Fatalf("List(%d): %v", tt.limit, err)
			}
			if !slices.Equal(messages(got), tt.want) {
				t.Errorf("List(%d) = %v, want %v", tt.limit, messages(got), tt.want)
			}
		})
	}
}

func TestListEmpty(t *testing.T) {
	st := newTestStore(t)
	got, err := st.List(0)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("got %d entries, want 0", len(got))
	}
}

func TestSearch(t *testing.T) {
	st := newTestStore(t)
	mustAdd(t, st, "new message", "")
	mustAdd(t, st, "fixed the flag bug", "")
	mustAdd(t, st, "100% done", "")
	mustAdd(t, st, "unrelated", "")

	tests := []struct {
		name    string
		keyword string
		want    []string
	}{
		{"substring match", "flag", []string{"fixed the flag bug"}},
		{"case insensitive", "MESSAGE", []string{"new message"}},
		{"no match", "zzzz", []string{}},
		// Documents current behavior: % is a LIKE wildcard, so "%" matches everything.
		// If you escape wildcards later, change this case on purpose.
		{"percent acts as wildcard", "%", []string{"unrelated", "100% done", "fixed the flag bug", "new message"}},
		{"empty keyword matches all", "", []string{"unrelated", "100% done", "fixed the flag bug", "new message"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := st.Search(tt.keyword)
			if err != nil {
				t.Fatalf("Search(%q): %v", tt.keyword, err)
			}
			if !slices.Equal(messages(got), tt.want) {
				t.Errorf("Search(%q) = %v, want %v", tt.keyword, messages(got), tt.want)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	st := newTestStore(t)
	mustAdd(t, st, "keep", "")
	mustAdd(t, st, "remove", "")

	all, err := st.List(0)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	var removeID int
	for _, e := range all {
		if e.Message == "remove" {
			removeID = e.Id
		}
	}

	if err := st.Delete(removeID); err != nil {
		t.Fatalf("Delete(%d): %v", removeID, err)
	}

	got, err := st.List(0)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if want := []string{"keep"}; !slices.Equal(messages(got), want) {
		t.Errorf("after delete = %v, want %v", messages(got), want)
	}
}

func TestDeleteMissingIDIsNoError(t *testing.T) {
	st := newTestStore(t)
	if err := st.Delete(999); err != nil {
		t.Errorf("Delete(999) on empty store: %v, want nil", err)
	}
}
