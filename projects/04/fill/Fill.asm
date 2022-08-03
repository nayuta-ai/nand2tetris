// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.

// R0: 
// R1:

// Put your code here.
    // define screen size
    @8192
    D=A
    @num_pixel
    M=D     // num_pixel = 8192

    @R0
    M=0     // R0 = 0
(LOOP)
    @KBD
    D=M     // key pressed code
    @KEY_PUSH
    D;JNE   // if KBD != 0 go to KEY_PUSH
(KEY_NOT_PUSH)
    @R1
    M=0     // R1 = 0
    @STATE_CHANGE
    0;JMP   // Go to STATE_CHANGE
(KEY_PUSH)
    @R1
    M=1     // R1 = 1
(STATE_CHANGE)
    @R0
    D=M     // D = R0
    @R1
    D=D-M   // D = R0 - R1
    @LOOP
    D;JEQ   // R0 == R1 (D == 0) go to LOOP
    @i
    M=0     // i = 0 (initialize)
    @R1
    D=M     // D = R1
    @R0
    M=D     // R0 = D thus R0 = R1
    @CLEAR
    D;JEQ   // if R0 == R1 (D == 0) go to CLEAR
(FILL)
    @i
    D=M     // D = i
    @num_pixel
    D=M-D   // D = num_pixel - i
    @LOOP
    D;JLT   // if D < 0 go to LOOP
    @SCREEN // 17536 (initial state)
    A=A+D   // A = A + D (D < 0)
    M=-1    // M[A] = -1
    @i
    MD=M+1  // D = i + 1, i++ 
    @FILL
    0;JMP
(CLEAR)
    @i
    D=M     // D = i
    @num_pixel
    D=M-D   // D = num_pixel - i
    @LOOP
    D;JLT   // if D < 0 go to LOOP
    @SCREEN // 17536
    A=A+D   // A = A + D (D < 0)
    M = 0   // M[A] = 0
    @i
    MD=M+1  // D = i + 1, i++ 
    @CLEAR
    0;JMP   // go to CLEAR
