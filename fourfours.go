// Solve four fours puzzle with brute force.
// Only tries basic arithmetic operations.
package main

type node struct {
  kind rune  // '4', or rune representing a binary operation.
  left, right *node
}

func forall_tree(pp **node, n int, fun func()) {
  var t node
  *pp = &t
  if (n == 1) {
    t.kind = '4'
    fun()
    return
  }
  for _, op := range "+-/*" {
    t.kind = op
    for k := 1; k < n; k++ {
      forall_tree(&t.left, k, func() {
        forall_tree(&t.right, n - k, fun)
      })
    }
  }
}

func tree_print(p *node) {
  if p.kind == '4' {
    print("4")
    return
  }
  print("(")
  tree_print(p.left)
  print(string(p.kind))
  tree_print(p.right)
  print(")")
}

func tree_eval(p *node) int {
  switch p.kind {
    case '+': return tree_eval(p.left) + tree_eval(p.right)
    case '-': return tree_eval(p.left) - tree_eval(p.right)
    case '/': return tree_eval(p.left) / tree_eval(p.right)
    case '*': return tree_eval(p.left) * tree_eval(p.right)
  }
  return 4
}

func main() {
  var root *node
  forall_tree(&root, 4, func() {
    func() {
      defer func() {
        if r := recover(); r != nil {
          print("(undefined)")
        }
      }()
      print(tree_eval(root))
    }()
    print(" = ")
    tree_print(root)
    println()
  })
}
