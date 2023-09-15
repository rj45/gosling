# Gosling

Gosling is a tiny go-like language.

I am following along with [ChibiCC's commits](https://github.com/rui314/chibicc/commits/main?after=90d1f7f199cc55b13c7fdb5839d1409806633fdb+300&branch=main) but in Go, and using my own ideas for how to structure things. I am mainly just taking the theme of each commit and implementing that my own way.

Initially I am generating AArch64 (64-bit ARM) assembly code, since I am on an ARM Mac. Eventually this will be able to cross compile to other architectures. If you are on x86, you might need an ARM emulator (or raspberry pi) to run the code / tests.

Follow the git commit history to see how I built this. Some commits have bugs I didn't notice, I didn't go back and rewrite history to fix them.

# License

[MIT](./LICENSE)
