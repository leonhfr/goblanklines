package funcdecl

import "fmt"

func multiLineA() {
	fmt.Println("a")
}
func multiLineB() { // want "missing newline between function declarations"
	fmt.Println("b")
}

func multiLineC() {
	fmt.Println("c")
}

func multiLineD() {
	fmt.Println("d")
}

func oneLinerA() { fmt.Println("a") }
func oneLinerB() { fmt.Println("b") }
func oneLinerC() { fmt.Println("c") }

func oneLinerD() { fmt.Println("d") }

func oneLinerE() { fmt.Println("e") }
func multiLineAfterOneLiner() { // want "missing newline between function declarations"
	fmt.Println("after")
}

func multiLineBeforeOneLiner() {
	fmt.Println("before")
}
func oneLinerF() { fmt.Println("f") } // want "missing newline between function declarations"

type T struct{}

func (t T) method1() {
	fmt.Println("1")
}
func (t T) method2() { // want "missing newline between function declarations"
	fmt.Println("2")
}

var x = 1

func afterUnrelatedDecl() {
	fmt.Println(x)
}

// Doc comment.
func withDocComment() {
	fmt.Println("doc")
}
func noDocComment() { // want "missing newline between function declarations"
	fmt.Println("no doc")
}

func lastFunc() {
	fmt.Println("last")
}
