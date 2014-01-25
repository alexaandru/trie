package utils

import (
    "fmt"
    u "github.com/alexaandru/utils"
    "strings"
)

// FIXME: Ugly and most probably not goroutine friendly either. Gets the job done though.
var nextID = trieID()

// Trie models a trie data structure.
type Trie struct {
    Label    string
    ID       int
    Children map[string]*Trie
}

// TrieNew initializes and a new Trie returns a pointer to it.
func TrieNew(labels ...string) (t *Trie) {
    t = new(Trie)
    if len(labels) > 0 {
        t.Label = labels[0]
    }
    t.ID = nextID()
    t.Children = map[string]*Trie{}

    return
}

// LoadTrieFromFile loads a trie definition from file.
func LoadTrieFromFile(fname string, opts ...bool) (t *Trie, pattern string) {
    t = TrieNew()
    firstStringAsPattern, lines := false, strings.Split(u.LoadFile(fname), "\n")
    if len(opts) > 0 {
        firstStringAsPattern = opts[0]
    }

    if firstStringAsPattern {
        pattern, lines = lines[0], lines[1:]
    }

    for _, line := range lines {
        t.Add(line)
    }

    return
}

// Add adds one string to a trie.
func (t *Trie) Add(s string) {
    ch := string(s[0])
    if t.Children[ch] == nil {
        t.Children[ch] = TrieNew(ch)
    }
    if len(s) > 1 {
        t.Children[ch].Add(s[1:])
    }
}

// IsLeaf determines if a trie is a leaf.
func (t *Trie) IsLeaf() bool {
    return len(t.Children) == 0
}

func (t *Trie) String() (out string) {
    for _, v := range t.Children {
        out += fmt.Sprintf("%d %d %s\n", t.ID, v.ID, v.Label)
    }
    for _, v := range t.Children {
        out += v.String()
    }

    return
}

// Match tries to match the t trie against the str string.
func (t *Trie) Match(str string) (match string, found bool) {
    if str == "" {
        return
    }

    str += "0"

    i := 0
    s, v := string(str[i]), t
    for {
        if v.Children[s] != nil {
            i++
            match += s
            s, v = string(str[i]), v.Children[s]
        } else if v.IsLeaf() {
            found = true
            return
        } else {
            match = ""
            return
        }
    }

    return
}

// Matches tries to find all the t trie matches agains the str string.
func (t *Trie) Matches(str string) (pos map[string]([]int)) {
    pos = map[string]([]int){}
    for i := 0; i < len(str); i++ {
        if match, found := t.Match(str[i:]); found {
            pos[match] = append(pos[match], i)
        }
    }

    return
}

// trieID returns the current id for a trie. Used for labelling tries during build.
func trieID() func() int {
    lastID := 0

    return func() int {
        lastID++
        return lastID
    }
}
