// Implementation of binary trees in Go.
// 
// Based on binary-trees Rust #8 program
// from
// The Computer Language Benchmarks Game
//
// Technically this is cheating by the rules
// of the benchmark since we're pre-allocating
// memory.
// But lets see how fast it will be.

package main

import (
    "flag"
    "fmt"
    "strconv"
    "strings"
    "math"
)

type Tree struct {
    nodes []TreeNode
}

func (t *Tree) AddNode() uint32 {
    tn := TreeNode{ tree: t, l: MaxUint, r: MaxUint }
    t.nodes = append(t.nodes, tn)
    return uint32(len(t.nodes) - 1)
}

type TreeNode struct {
    tree *Tree
    l uint32
    r uint32
}

const MaxUint = ^uint32(0)

func (t TreeNode) Left() *TreeNode {
    if t.l == MaxUint {
        return nil
    } else {
        return &t.tree.nodes[t.l]
    }
}

func (t TreeNode) Right() *TreeNode {
    if t.r == MaxUint {
        return nil
    } else {
        return &t.tree.nodes[t.r]
    }
}

func itemCheck(tn *TreeNode) uint32 {
    l := tn.Left()
    r := tn.Right()

    if l != nil && r != nil {
        return uint32(1) + itemCheck(l) + itemCheck(r)
    }

    return 1
}

func bottomUpTree(depth uint32) *TreeNode {
    tree_size := uint32(math.Pow(2.0, float64(depth + 1)) - 1)
    tree := Tree{ nodes: make([]TreeNode, 0, tree_size) }
    bottomUpTreeInner(depth, &tree)

    //fmt.Printf("size of nodes: len %d cap %d\n", len(tree.nodes), cap(tree.nodes))

    return &tree.nodes[0]
}

func bottomUpTreeInner(depth uint32, tree *Tree) uint32 {
    result := tree.AddNode()
    if depth > uint32(0) {
        tree.nodes[result].l = bottomUpTreeInner(depth - 1, tree)
        tree.nodes[result].r = bottomUpTreeInner(depth - 1, tree)
    }

    return result
}

func inner(depth, iterations uint32) string {
    chk := innerImpl(depth, iterations)
    return fmt.Sprintf("%d\t trees of depth %d\t check: %d",
        iterations, depth, chk)
}

func innerImpl(depth, iterations uint32) uint32 {
    iter_chan := make([](chan uint32), 0, iterations)

    for i := uint32(0); i < iterations; i++ {
        ch := make(chan uint32, 1)
        iter_chan = append(iter_chan, ch)

        go func() {
            a := bottomUpTree(depth)
            chk := itemCheck(a)

            ch <- chk
            close(ch)
        }()
    }

    chk := uint32(0)

    for _, ch := range(iter_chan) {
        chk += <-ch
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
