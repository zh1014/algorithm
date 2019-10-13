package regexp

import (
	"fmt"
	"github.com/zh1014/algorithm/graphs/digraph"
	"github.com/zh1014/algorithm/set"
	"github.com/zh1014/algorithm/stack"
)

func IsMatch(pattern, txt string) bool {
	symbolTbl := parsePattern(pattern)
	g := createNFA(symbolTbl)
	reachableStatus := getStartStatus(g)
	for _, r := range txt {
		statusAfterMatch := set.New()
		for !reachableStatus.IsEmpty() {
			status := reachableStatus.RemoveOne()
			if status < len(symbolTbl) && match(symbolTbl[status], r) {
				statusAfterMatch.Add(status+1)
			}
		}
		if statusAfterMatch.IsEmpty() {
			return false
		}
		getReachableStatus(g, statusAfterMatch, reachableStatus)
	}
	return reachableStatus.Contains(len(symbolTbl))
}

func match(s symbol, r rune) bool {
	return isPrimeDot(s) || (!s.isPrime && s.r == r)
}

type symbol struct {
	isPrime bool
	r       rune
}

func parsePattern(pattern string) []symbol {
	// this can only transfer \( \) \| \* \.
	pttrnRunes := []rune(pattern)
	numRunes := len(pttrnRunes)
	symbolTable := make([]symbol, 0, numRunes)
	for i := 0; i < numRunes; i++ {
		if pttrnRunes[i] == '\\' {
			i++
			if !isPrimeRune(pttrnRunes[i]) {
				panic(fmt.Sprintf("invalid transfer: \\%v", string(pttrnRunes[i])))
			}
			symbolTable = append(symbolTable, symbol{
				isPrime: false,
				r:       pttrnRunes[i],
			})
			continue
		}
		symbolTable = append(symbolTable, symbol{
			isPrime: isPrimeRune(pttrnRunes[i]),
			r:       pttrnRunes[i],
		})
	}
	return symbolTable
}

func isPrimeRune(r rune) bool {
	return r == '(' || r == ')' || r == '|' || r == '*' || r == '.'
}

func getReachableStatus(g digraph.Digraph, src, reachable set.Set) {
	marked := make([]bool, g.NumV())
	for !src.IsEmpty() {
		aSourceStatus := src.RemoveOne()
		dfs(g, aSourceStatus, marked)
	}
	for i, b := range marked {
		if b {
			reachable.Add(i)
		}
	}
}

func getStartStatus(g digraph.Digraph) set.Set {
	marked := make([]bool, g.NumV())
	dfs(g, 0, marked)
	reachable := set.New()
	for i, b := range marked {
		if b {
			reachable.Add(i)
		}
	}
	return reachable
}

func dfs(g digraph.Digraph, v int, marked []bool) {
	marked[v] = true
	adj := g.Adjacent(v)
	for _, w := range adj {
		if !marked[w] {
			dfs(g, w, marked)
		}
	}
}

func createNFA(symbolTable []symbol) digraph.Digraph {
	tblSize := len(symbolTable)
	g := digraph.NewDigraph(tblSize + 1)
	stck := stack.NewStackInt(tblSize)
	for i, symbl := range symbolTable {
		leftBracket := i
		if symbl.isPrime && (symbl.r == '(' || symbl.r == ')' || symbl.r == '*') {
			g.AddEdge(i, i+1)
		}
		if symbl.isPrime && (symbl.r == '(' || symbl.r == '|') {
			stck.Push(i)
		}
		if symbl.isPrime && symbl.r == ')' {
			or := stck.Pop()
			if symbolTable[or].r == '(' {
				leftBracket = or
			} else { // symbolTable[or].r == '|'
				leftBracket = stck.Pop()
				g.AddEdge(leftBracket, or+1)
				g.AddEdge(or, i)
			}
		}
		if i+1 < tblSize && isPrimeWildCard(symbolTable[i+1]) {
			g.AddEdge(leftBracket, i+1)
			g.AddEdge(i+1, leftBracket)
		}
	}
	return g
}

func isPrimeDot(s symbol) bool {
	return s.isPrime && s.r == '.'
}

func isPrimeWildCard(s symbol) bool {
	return s.isPrime && s.r == '*'
}
