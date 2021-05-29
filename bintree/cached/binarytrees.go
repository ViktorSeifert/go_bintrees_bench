// Implementation of binary trees in Go.
// 
// Based on binary-trees Rust #8 program
// from
// The Computer Language Benchmarks Game
//

// This is actually against the rules of the Benchmarks Game.
// I added this to be able to check how "fast" tree traversal is.

package main

import (
    "flag"
    "fmt"
    "strconv"
    "strings"
)

type Tree struct {
    Left  *Tree
    Right *Tree
}

func itemCheck(tree *Tree) uint32 {
    if tree.Left != nil && tree.Right != nil {
        return uint32(1) + itemCheck(tree.Right) + itemCheck(tree.Left)
    }

    return 1
}

var tree_cache = make(map[uint32]*Tree)

func bottomUpTree(depth uint32) (t *Tree) {
    t, ok := tree_cache[depth]
    if ok {
        return
    }

    if depth > uint32(0) {
        t = &Tree{Left: bottomUpTree(depth - 1), Right: bottomUpTree(depth - 1)}
    } else {
        t = &Tree{}
    }

    tree_cache[depth] = t

    return
}

func inner(depth, iterations uint32) string {
    chk := innerImpl(depth, iterations)
    return fmt.Sprintf("%d\t trees of depth %d\t check: %d",
        iterations, depth, chk)
}

func innerImpl(depth, iterations uint32) uint32 {
    chk := uint32(0)
    for i := uint32(0); i < iterations; i++ {
        a := bottomUpTree(depth)
        chk += itemCheck(a)
    }
    return chk
}

const minDepth = uint32(4)

func main() {
    n := 0
    flag.Parse()
    if flag.NArg() > 0 {
        n, _ = strconv.Atoi(flag.Arg(0))
    }

    fmt.Println(run(uint32(n)))
}

func run(n uint32) string {
    maxDepth := n
    if minDepth+2 > n {
        maxDepth = minDepth + 2
    }

    depth := maxDepth + 1

    var messages = make([]string, 0)
    
    tree := bottomUpTree(depth)
    messages = append(messages, fmt.Sprintf("stretch tree of depth %d\t check: %d",
            depth, itemCheck(tree)))
    tree = nil
    
    longLivedTree := bottomUpTree(maxDepth)

    for halfDepth := minDepth / 2; halfDepth < maxDepth/2+1; halfDepth++ {
        depth := halfDepth * 2
        iterations := uint32(1 << (maxDepth - depth + minDepth))

        func(d, i uint32) {
            messages = append(messages, inner(d, i))
        }(depth, iterations)
    }

    messages = append(messages, 
        fmt.Sprintf("long lived tree of depth %d\t check: %d",
            maxDepth, itemCheck(longLivedTree)))

    var rb strings.Builder
    for _, m := range messages {
        rb.WriteString(m)
        rb.WriteString("\n")
    }
    return rb.String()
}
