# Assembler
## Overview
My Hack assembler reads as input a text file named *.asm, containing a Hack assembly program, and produces as output a text file, containig the translated hack machine code.
I proposed an assembler based on four modules: a Parser module that parses the input, a Code module that provides the binary codes of all the assembly mnemonics, a SymbolTable module that handles symbols, and a main program that drives the entire translation process.
## To Do
- [x] Parser module
- [x] Code module
- [x] SymbolTable module
- [x] main program

### Parser module
The main function of the parser is to break each assembly command into its underlying components(fields and symbols).
