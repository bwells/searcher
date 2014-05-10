package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/bwells/trie"
	"github.com/bwells/uniq"
	"github.com/deckarep/golang-set"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
)

type Contact struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (c Contact) Key() int {
	return c.Id
}

func (c *Contact) String() string {
	return fmt.Sprintf("{%d %s %s}", c.Id, c.Name, c.Email)
}

func buildTrieFromDB() *trie.Trie {

	con, _ := sql.Open("mysql", "developer:w3bd3v2008@tcp(localhost:3306)/erp5?charset=utf8")
	defer con.Close()

	query := `select contacts.id as Id,
	   			concat_ws(' ', contacts.first_name, contacts.last_name) as Name,
	   			users.username as Email
			  from users join contacts on users.contact_id = contacts.id;`

	rows, _ := con.Query(query)
	defer rows.Close()

	t := trie.NewTrie()

	for rows.Next() {
		var id int
		var name string
		var email string
		_ = rows.Scan(&id, &name, &email)
		c := Contact{id, name, email}

		name = strings.ToLower(c.Name)
		email = strings.ToLower(c.Email)

		for _, f := range strings.Fields(name) {
			t.Add(f, &c)
		}

		for _, f := range strings.Split(email, "@") {
			t.Add(f, &c)
		}
	}

	return t
}

func multiMatch(t *trie.Trie, query string) []uniq.Interface {

	// strip whitespace
	query = strings.ToLower(strings.TrimSpace(query))

	// skip empty input
	if len(query) == 0 {
		return make([]uniq.Interface, 0)
	}

	// pull matches on the first term and create a Set
	terms := strings.Fields(query)
	m, _ := t.MatchPartial(terms[0])
	matches := mapset.NewSetFromSlice(m)

	// pull matches for remaining terms and intersect them with the first Set
	for _, term := range terms[1:] {
		m, _ = t.MatchPartial(term)
		matches = matches.Intersect(mapset.NewSetFromSlice(m))
	}

	// cast matches to []uniq.Interface
	results := make([]uniq.Interface, matches.Cardinality())
	i := 0
	for key := range matches.Iter() {
		results[i] = key.(*Contact)
		i++
	}

	// ensure the results are unique
	results = uniq.Uniq(results)

	return results
}

func main() {

	t := buildTrieFromDB()

	search_trie = t

	fmt.Println("Enter a search term:")

	serve()

	in := bufio.NewReader(os.Stdin)
	for {
		input, err := in.ReadString('\n')
		if err != nil {
			return
		}

		results := multiMatch(t, input)

		fmt.Println(results)
	}
}
