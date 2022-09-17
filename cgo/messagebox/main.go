package main

// go调用c,会弹出一个框

/*
#include <stdio.h>
#include <windows.h>

void box(){
	MessageBox(0,"ddd","aa",MB_OK);
}

*/
import "C"

func main() {
	C.box()
}
