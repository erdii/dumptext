# dumptext

dump an ELF's .text section in various formats and optionally verify that it does not contain any null bytes.

## Installation

`go install github.com/erdii/dumptext@main`


## Usage

```
Usage: dumptext path/to/elf/binary
Usage: cat path/to/elf/binary | dumptext
Description: Reads the .text section from an ELF binary and dumps (optionally formatted) bytes to stdout.
Flags:
        Envvar: FORMAT=escape(default)|dump|raw - specifies output format.
        Envvar: VALIDATE=1 - validates that there are no null bytes in the data.
```

## Formats

### Raw

Dumps raw bytes from the `.text` section in the host's native endianness.

### Escape

Dumps the bytes from the `.text` section as a hex-escaped string. Handy for copy-pasting into terminals or code.

Example:

```bash
nasm -f elf64 misc/example-shellcode.asm -o /dev/stdout | FORMAT=escape dumptext
# =>
# \x6a\x31\x58\x99\xcd\x80\x89\xc3\x89\xc1\x6a\x46\x58\xcd\x80\xb0\x0b\x52\x68\x6e\x2f\x73\x68\x68\x2f\x2f\x62\x69\x89\xe3\x89\xd1\xcd\x80
```

### Dump

Prints a hexdump of the bytes in the `.text section`.

Example:

```bash
nasm -f elf64 misc/example-shellcode.asm -o /dev/stdout | FORMAT=dump dumptext
# =>
# 00000000  6a 31 58 99 cd 80 89 c3  89 c1 6a 46 58 cd 80 b0  |j1X.......jFX...|
# 00000010  0b 52 68 6e 2f 73 68 68  2f 2f 62 69 89 e3 89 d1  |.Rhn/shh//bi....|
# 00000020  cd 80                                             |..|
```
