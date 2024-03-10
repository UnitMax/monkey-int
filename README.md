<p align=center><img src="logo.webp" alt="Logo courtesy of ChatGPT" width="200" /><br>
(Logo courtesy of ChatGPT... The future is now, I guess)</p>

# monkey-int
Monkey language interpreter based on Thorsten Ball's [book](https://interpreterbook.com). Including Lexer, Parser (Operator-Precedence/Pratt Parser), an object system and a very simple tree-walking interpreter/evaluator.

The interpreter is currently being expanded by a compiler and corresponding virtual machine (see subsequent book).

## Supported features

- 64bit integers
- booleans
- basic arithmetic
- function objects
- closures 😎
- simple tree walking interpreter
- strings
- arrays
- hashmaps
- printing to stdout
- reading and writing to the filesystem using `readfile` and `writefile`

## Progress

- [x] Lexer 
- [x] Parser
- [x] Eval
- [x] Extensions:
    - [x] Strings
    - [x] `len`
    - [x] Hashmaps
    - [x] Arrays
    - [x] File I/O
- [ ] Compiler
- [ ] MVM (Monkey Virtual Machine)