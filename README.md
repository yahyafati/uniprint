# Remote File Printer (Go)

This is a simple command-line tool written in Go that lets you print local files on a remote machine using SSH and SFTP. It asks for your email and password, uploads the files, and prints them using `lpr` on the remote server.

It includes:

- Password-protected login
- Basic printer selection
- Optional deletion of files after printing

> **NOTE:** The code is provided **as-is**. Use it at your own risk. No guarantees or support are provided.

**Note**: The code is provided **as-is**. Use it at your own risk. No guarantees or support are provided.

It is based on: https://fachschaft.tf/en/studium/pool_printing/

## Requirements

- Go 1.18 or newer

## Installation

### Option 1: Build the code yourself (Recommended)

> This option might require a few extra steps, but it's probably *safer*. Also, prebuilt binaries might not exist for your system architecture.


1. Clone the repository

```bash
git clone https://github.com/yahyafati/uniprint.git
cd uniprint
```

2. Install Dependencies

```bash
go mod tidy
```

3. Build the app

```bash
go build -o uniprint
```

### Option 2: Download the binary

If you don't want to build it yourself, you can download a precompiled binary (if available) from the [Releases](https://github.com/yahyafati/uniprint/releases) section.

> ⚠️ Make sure the binary matches your OS and architecture.

## How to print

```bash
uniprint # or ./uniprint
```

### Walkthrough

Example session (replace `tfusername` and `Dummy PDF.pdf` with your own):

```
Enter your email (e.g., user@host): tfusername
Enter your password: 
Hi, tfusername@login.informatik.uni-freiburg.de
Enter file path to print (leave empty to finish): Dummy PDF.pdf
Enter file path to print (leave empty to finish): 
Found 1 files to print.
Delete files after printing? (y/N): y
Files will be removed from the server once it has been printed.
Copying Dummy PDF.pdf to tfusername@login.informatik.uni-freiburg.de:~

Verifying remote login with 'whoami'...
Connected as: tfusername

Choose a printer:
1: tfppr1 (b/w)
2: tfppr2 (Color)
3: tfppr3 (b/w)
4: tfppr4 (Color)
Enter choice number: 3
Selected printer: tfppr3

(Now Printing) Command executed: lpr -P tfppr3 -r 'Dummy PDF.pdf'
Printed 'Dummy PDF.pdf' to tfppr3.
```

## Troubleshooting

- Make sure your username and password are correct.
- Ensure you can SSH into `login.informatik.uni-freiburg.de` manually.
- Check that `lpr` is installed and configured on the remote server.

## License

MIT License