// Implementation of binary trees in Go.
// 
// Based on binary-trees Rust #8 program
// from
// The Computer Language Benchmarks Game
//

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

func bottomUpTree(depth uint32) *Tree {
    if depth > uint32(0) {
        return &Tree{Left: bottomUpTree(depth - 1), Right: bottomUpTree(depth - 1)}
    } else {
        return &Tree{}
    }
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

    var messages = make([](chan string), 0)

    strech_chan := make(chan string, 1)
    messages = append(messages, strech_chan)
    
    
    go func() {
        tree := bottomUpTree(depth)
        strech_chan <- fmt.Sprintf("stretch tree of depth %d\t check: %d",
                depth, itemCheck(tree))
        close(strech_chan)
    }()

    longl_chan := make(chan string, 1)
    
    go func() {
        longLivedTree := bottomUpTree(maxDepth)
        longl_chan <-
            fmt.Sprintf("long lived tree of depth %d\t check: %d",
                maxDepth, itemCheck(longLivedTree))
        close(longl_chan)
    }()

    for halfDepth := minDepth / 2; halfDepth < maxDepth/2+1; halfDepth++ {
        depth := halfDepth * 2
        iterations := uint32(1 << (maxDepth - depth + minDepth))

        m_chan := make(chan string, 1)
        messages = append(messages, m_chan)

        func(d, i uint32) {
            go func() {
                m_chan <- inner(d, i)
                close(m_chan)
            }()
        }(depth, iterations)
    }

    messages = append(messages, longl_chan)

    var rb strings.Builder
    for _, m := range messages {
        rb.WriteString(<-m)
        rb.WriteString("\n")
    }
    return rb.String()
}
