// Solves Chapter 9, Problem 29 from Smullyan's "To Mock a Mockingbird" with
// brute force.
//
// On my laptop, this program about 10 minutes to run.
package main

type node struct {
  kind int  // 0 = leaf, 1 = branch.
  left, right *node
}

// For each full binary tree with n leaf nodes, assigns a pointer to the tree
// to '*pp' and calls the given function. 
func forall_tree(pp **node, n int, fun func()) {
  var t node
  *pp = &t
  if (n == 1) {
    t.kind = 0
    fun()
    return
  }
  t.kind = 1
  for k := 1; k < n; k++ {
    forall_tree(&t.left, k, func() {
      forall_tree(&t.right, n - k, fun)
    })
  }
}

func node_new_leaf() *node {
  return &node{0,nil,nil}
}

func node_new_branch(left, right *node) *node {
  return &node{1,left,right}
}

// Prints a tree.
func tree_print(p *node) {
  if p.kind == 1 {
    print("(")
    tree_print(p.left)
    tree_print(p.right)
    print(")")
    return
  }
  print("L")
}

// Prints a tree to a string.
func tree_sprint(p *node) string {
  if p.kind == 1 {
    return "(" + tree_sprint(p.left) + tree_sprint(p.right) + ")"
  }
  return "L"
}

// Returns true if given trees are equal.
func tree_eq(x, y *node) bool {
  if x.kind != y.kind {
    return false
  }
  if x.kind == 0 {
    return true
  }
  return tree_eq(x.left, y.left) && tree_eq(x.right, y.right)
}

// Copies a tree.
func tree_dup(p *node) *node {
  if p.kind == 0 {
    return node_new_leaf()
  }
  return node_new_branch(tree_dup(p.left), tree_dup(p.right))
}

// Copies a tree, but substitutes any 'target' node with 'repl'.
func tree_sub(p, target, repl *node) *node {
  if p == target {
    return repl
  }
  if p.kind == 0 {
    return node_new_leaf()
  }
  return node_new_branch(tree_sub(p.left, target, repl), tree_sub(p.right, target, repl))
}

func main() {
  // For all trees with up to 12 leaf nodes, we "double it up" by construct a
  // tree where both the left and right child nodes of the root are copies of
  // the original. Then we try every possible replacement of subtrees of the
  // form x(yy) by (Lx)y, and see if any of them are equal to the original
  // tree. This must terminate, since (Lx)y is either strictly smaller than
  // x(yy), or cannot be reduced further.
  n := 0
  for k := 1; k <= 12; k++ {
    println("leaves: ", k)
    var orig *node
    forall_tree(&orig, k, func() {
      n++
      if n % 10000 == 0 {
        println(n, " trees")
      }
      // Prints the tree if it equals the original tree.
      check := func(p *node) {
        if tree_eq(p, orig) {
          tree_print(p)
          println()
        }
      }
      // Traverses the tree, and along the way, if we encounter a sub-tree of
      // the form x(yy), make a new tree where this sub-tree has been replaced
      // with (Lx)y and spawn a recursive call on this new tree. After
      // reaching the last node, calls the given callback.
      m := make(map[string]bool)
      var try func(r, p *node, fun func(*node))
      try = func(r, p *node, fun func(*node)) {
        if p.kind == 0 {
          fun(r)
          return
        }
        if p.right.kind == 1 && tree_eq(p.right.left, p.right.right) {
          // x(yy) --> (Lx)y
          t := tree_sub(r, p, node_new_branch(node_new_branch(node_new_leaf(), tree_dup(p.left)), tree_dup(p.right.left)))
          s := tree_sprint(t)
          if _, ok := m[s]; !ok {
            m[s] = true
            try(t, t, check)
          }
        }
        try(r, p.left, func(*node) { try(r, p.right, fun) })
      }
      dbl := node_new_branch(tree_dup(orig), tree_dup(orig))
      try(dbl, dbl, check)
    })
  }
}
