package ac

import (
	"bufio"
	"os"
)

type AcNode struct {
	fail      *AcNode
	isPattern bool
	next      map[rune]*AcNode
}

func newAcNode() *AcNode {
	return &AcNode{
		fail:      nil,
		isPattern: false,
		next:      map[rune]*AcNode{},
	}
}

type AcAutoMachine struct {
	root *AcNode
}

func NewAcAutoMachine() *AcAutoMachine {
	return &AcAutoMachine{
		root: newAcNode(),
	}
}

func (ac *AcAutoMachine) addPattern(pattern string) {
	chars := []rune(pattern)
	iter := ac.root
	for _, c := range chars {
		if _, ok := iter.next[c]; !ok {
			iter.next[c] = newAcNode()
		}
		iter = iter.next[c]
	}
	iter.isPattern = true
}

func (ac *AcAutoMachine) build() {
	queue := []*AcNode{}
	queue = append(queue, ac.root)
	for len(queue) != 0 {
		parent := queue[0]
		queue = queue[1:]

		for char, child := range parent.next {
			child.fail = ac.root
			failAcNode := parent.fail
			for failAcNode != nil {
				if _, ok := failAcNode.next[char]; ok {
					child.fail = failAcNode.next[char]
					break
				}
				failAcNode = failAcNode.fail
			}
			if failAcNode == nil {
				child.fail = ac.root
			}
			queue = append(queue, child)
		}
	}
}

func (ac *AcAutoMachine) Filter(content string) string {
	chars := []rune(content)
	iter := ac.root
	var start, end int
	for i, c := range chars {
		_, ok := iter.next[c]
		for !ok && iter != ac.root {
			iter = iter.fail
			_, ok = iter.next[c]
		}
		if _, ok = iter.next[c]; ok {
			if iter == ac.root { // this is the first match, record the start position
				start = i
			}
			iter = iter.next[c]
			if iter.isPattern {
				end = i // this is the end match, record one result
				for i := start; i <= end; i++ {
					chars[i] = '*'
				}
			}
		}
	}
	return string(chars)
}

func New(fileName string) *AcAutoMachine {
	ac := NewAcAutoMachine()
	f, err := os.OpenFile(fileName, os.O_RDONLY, 0660)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		l, _, err := r.ReadLine()
		if err != nil {
			break
		}
		ac.addPattern(string(l))
	}
	ac.build()
	return ac
}
