package temp_leet

import (
	"sort"
	"testing"
)

func TestLongest(t *testing.T) {
	//longestWord([]string{"cat","banana","dog","nana","walk","walker","dogwalker"})
	longestWord([]string{"cat","banana","dog","nana","walk","walker","dogwalker"})
}


func longestWord(words []string) string {
	t := &Trie{}
	for _, word := range words{
		t.Insert(word)
	}

	var res = make([]byte,0)
	for i := 0; i < len(words); i++ {
		d := dfsHelper{lword:[]byte{}, cur:[]byte{}, finish: false, valid: false, root:t}
		d.dfs(t, words[i])
		if d.valid {
			if len(res) == 0 {
				res = make([]byte, len(d.lword))
				copy(res, d.lword)
			} else {
				if len(res) < len(d.lword) {
					res = make([]byte, len(d.lword))
					copy(res, d.lword)
				} else if len(res) == len(d.lword) {
					temp := []string{string(res), string(d.lword)}
					sort.Strings(temp)
					res = []byte(temp[0])
				}
			}
		}
	}
	return string(res)
}

type Trie struct{
	next [26]*Trie
	isEnd bool
}

func (this *Trie) Insert(word string) {
	cur := this
	for i := 0; i < len(word); i++ {
		pos := int(word[i]-'a')
		if cur.next[pos] == nil {
			cur.next[pos] = &Trie{}
		}
		cur = cur.next[pos]
	}
	cur.isEnd = true
}

type dfsHelper struct{
	lword []byte
	cur []byte
	root *Trie
	finish bool
	valid bool
}

func (this *dfsHelper)dfs(root *Trie, word string) {
	if root == nil {
		return
	}
	if len(word) == 0 {
		return
	}
	pos := int(word[0]) - int('a')
	if root.next[pos] == nil {
		return
	}
	t := root.next[pos]
	this.cur = append(this.cur, word[0])
	if !t.isEnd {
		this.dfs(t, word[1:])
	} else {
		if !this.finish {
			this.finish = true
			this.lword = make([]byte, len(this.cur))
			copy(this.lword, this.cur)
			this.cur = []byte{}
			this.dfs(this.root, word[1:])
		} else {
			if len(word) == 1 {
				this.valid = true
				this.lword = append(this.lword, this.cur...)
			} else {
				return
			}
		}
	}

}
