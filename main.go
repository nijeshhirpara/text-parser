package main

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
)

type stack struct {
	lock sync.Mutex
	s    []*node
}

func NewStack() *stack {
	return &stack{sync.Mutex{}, make([]*node, 0)}
}

func (s *stack) Push(v *node) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.s = append(s.s, v)
}

func (s *stack) Pop() (*node, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.s)
	if l == 0 {
		return nil, errors.New("Empty Stack")
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res, nil
}

func parse(v string) (*node, error) {
	// root node
	root := &node{}

	// stack to store tree nodes
	s := NewStack()

	// temporary variables
	tempStr := ""
	parent := root
	child := &node{}

	for _, elm := range v {
		elmStr := string(elm)
		if elmStr == "[" {
			c := parent.addChild(tempStr)
			if c == nil {
				continue
			}

			// clear temp string, once node is prepared
			tempStr = ""

			s.Push(parent)

			child = c
			parent = child
		} else if elmStr == "]" {
			if c := parent.addChild(tempStr); c == nil {
				continue
			}

			// clear temp string, once node is prepared
			tempStr = ""

			child = parent

			// check if anything stacked
			if p, _ := s.Pop(); p != nil {
				parent = p
			}
		} else if elmStr == "," {
			c := parent.addChild(tempStr)
			if c == nil {
				continue
			}

			// clear temp string, once node is prepared
			tempStr = ""

			child = c
		} else {
			// concatenate characters until delimiters detected
			tempStr += elmStr
		}
	}

	return root, nil
}

// creates a child and append to the parent
func (parent *node) addChild(v string) *node {
	if v == "" {
		return nil
	}
	child := &node{Name: v}
	parent.Children = append(parent.Children, child)
	return child
}

type node struct {
	Name     string  `json:"name"`
	Children []*node `json:"children,omitempty"`
}

var examples = []string{
	"[a,b,c]",
	"[a[aa[aaa],ab,ac],b,c[ca,cb,cc[cca]]]",
}

func main() {
	for i, example := range examples {
		result, err := parse(example)
		if err != nil {
			panic(err)
		}
		j, err := json.MarshalIndent(result, " ", " ")
		if err != nil {
			panic(err)
		}
		log.Printf("Example %d: %s - %s", i, example, string(j))
	}
}
