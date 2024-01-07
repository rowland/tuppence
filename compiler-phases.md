# Compiler Phases

1.  Lexing

    Speed syntax checking, simplify parsing and improve error messages by lexing first.

2.  Parsing

    Build AST, preserving operators and other syntax sugar.

3.  Type Checking

    Build CST, transforming operators to function calls, where appropriate.

4.  Optimization

    Constant propagation, constant folding, sparse conditional constant propagation, common subexpression elimination, dead code elimination.

4. Code Generation
