package trietree

type alphbt interface {
	// ToRune convert index to rune
	ToRune(int) rune
	// ToIndex convert rune to index
	ToIndex(rune) int
	// R is the size of this Alphabet
	R() int
}
