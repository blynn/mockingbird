// Solves Chapter 9, Problem 29 from Smullyan's "To Mock a Mockingbird" with
// brute force.
//
// Takes under a minute to process ~80000 trees.
//
// Originally my program had a lot more nested functions and seemed a few
// seconds faster. I'm willing to sacrifice a few seconds so this code looks
// more conventional.
#include <stdio.h>
#include <stdlib.h>
#include "blt.h"

struct node_s {
  int type;  // 0 = leaf, 1 = internal.
  struct node_s *left, *right;
};
typedef struct node_s *node_t;

// For each full binary tree with n leaf nodes, assigns the tree to '*p' and
// calls the given function. 
void forall_tree(node_t *p, int n, void (*fun)()) {
  struct node_s t[1];
  *p = t;
  if (n == 1) {
    t->type = 0; // *t = (node_s) {0, 0, 0};
    fun();
  } else {
    t->type = 1;
    for (int k = 1; k < n; k++) {
      void g() {
        forall_tree(&t->right, n - k, fun);
      }
      forall_tree(&t->left, k, g);
    }
  }
}

// Prints a tree as a parenthesized expression.
void tree_print(node_t p) {
  if (p->type) {
    putchar('(');
    tree_print(p->left);
    tree_print(p->right);
    putchar(')');
    return;
  }
  putchar('L');
}

// Returns a new leaf node.
node_t node_new_leaf() {
  node_t r = malloc(sizeof(*r));
  *r = (struct node_s) {0, 0, 0};
  return r;
}

// Returns a new internal node.
node_t node_new_branch(node_t left, node_t right) {
  node_t r = malloc(sizeof(*r));
  *r = (struct node_s) {1, left, right};
  return r;
}

// Copies a tree.
node_t tree_dup(node_t p) {
  return p->type ? node_new_branch(tree_dup(p->left), tree_dup(p->right)) : node_new_leaf();
}

// Deletes a tree.
void tree_clear(node_t p) {
  if (p->type) {
    tree_clear(p->left);
    tree_clear(p->right);
  }
  free(p);
}

// Compares trees, returning 0 on equality, and nonzero otherwise.
int tree_cmp(node_t p, node_t q) {
  int k = q->type - p->type;
  if (k) return k;
  if (!p->type) return 0;
  return tree_cmp(p->left, q->left) || tree_cmp(p->right, q->right);
}

// Copies a tree, but substitutes any 'target' node with 'repl'.
node_t tree_sub(node_t p, node_t target, node_t repl) {
  if (p == target) return repl;
  return p->type ? node_new_branch(tree_sub(p->left, target, repl), tree_sub(p->right, target, repl)) : node_new_leaf();
}

int already_tried(BLT* blt, node_t t) {
  // Remember each tree to avoid trying it more than once.
  char buf[128], *s = buf;
  void f(node_t p) {
    // For internal use only; prefix notation will do. Efficient enough.
    if (p->type) {
      *s++ = '(';
      f(p->left);
      f(p->right);
    }
    *s++ = 'L';
  }
  f(t);
  *s = 0;
  return blt_put_if_absent(blt, buf, 0);
}

int main() {
  // For all trees with up to 12 leaf nodes, we "double it up" by construct a
  // tree where both the left and right child nodes of the root are copies of
  // the original. Then we try every possible replacement of subtrees of the
  // form x(yy) by (Lx)y, and see if any of them are equal to the original
  // tree. This must terminate, since (Lx)y is either strictly smaller than
  // x(yy), or cannot be reduced further.
  node_t orig;
  int n = 0;
  void f() {
    BLT *blt = blt_new();
    if (!(++n % 10000)) printf("%d trees\n", n);
    // Prints the tree if it equals the original tree.
    void check(node_t p) {
      if (!tree_cmp(p, orig)) tree_print(p), putchar('\n');
    }
    // Traverses the tree, and along the way, if we encounter a sub-tree of the
    // form x(yy), make a new tree where this sub-tree has been replaced with
    // (Lx)y and spawn a recursive call on this new tree.
    // After reaching the last node, calls the given callback.
    void try(node_t r, node_t p, void (*fun)(node_t)) {
      if (!p->type) return fun(r);
      if (p->right->type && !tree_cmp(p->right->left, p->right->right)) {
        // x(yy) --> (Lx)y
        node_t t = tree_sub(r, p, node_new_branch(node_new_branch(node_new_leaf(), tree_dup(p->left)), tree_dup(p->right->left)));
        if (!already_tried(blt, t)) try(t, t, check);
        tree_clear(t);
      }
      void g() {
        try(r, p->right, fun);
      }
      try(r, p->left, g);
    }
    node_t dbl = node_new_branch(tree_dup(orig), tree_dup(orig));
    try(dbl, dbl, check);
    tree_clear(dbl);
    blt_clear(blt);
  }
  for(int k = 1; k <= 12; k++) {
    printf("leaves: %d\n", k);
    forall_tree(&orig, k, f);
  }
  return 0;
}
