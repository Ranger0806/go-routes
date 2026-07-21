# Go Routes

Go Routes is a command-line utility that collects routes from multiple Windows BAT files, safely optimizes them, and generates a single BAT file suitable for import into a Keenetic router.

The utility was created to process route lists such as those found in the [RockBlack-VPN/ip-address](https://github.com/RockBlack-VPN/ip-address) repository.

## Features

- Recursively searches for `.bat` files
- Parses IPv4 `route add` commands
- Uses source filenames as route descriptions
- Removes exact duplicate routes
- Removes networks already covered by broader routes
- Safely merges sibling CIDR networks
- Preserves and combines route descriptions
- Generates a single BAT file with Windows CRLF line endings
- Produces stable and sorted output
- Ignores a previously generated output file during repeated runs

## Example

Input files:

```bat
route add 10.0.0.0 mask 255.255.255.128 0.0.0.0
route add 10.0.0.128 mask 255.255.255.128 0.0.0.0
```

Go Routes safely combines them into:

```bat
route add 10.0.0.0 mask 255.255.255.0 0.0.0.0 & rem example
```

The optimization does not add new IP addresses. Two networks are merged only when they are complete sibling halves of the same parent network.

## Requirements

- Go 1.22 or newer

Check your installed Go version:

```bash
go version
```

## Build

Clone the repository:

```bash
git clone https://github.com/Ranger0806/go-routes.git
cd go-routes
```

Run tests:

```bash
go test ./...
go vet ./...
```

Build the application:

```bash
go build -o goroutes ./cmd/goroutes
```

## Usage

```bash
./goroutes \
  --input /path/to/source-directory \
  --output ./keenetic-routes.bat
```

Example on macOS:

```bash
./goroutes \
  --input "$HOME/Downloads/ip-address" \
  --output "./keenetic-routes.bat"
```

Example on Windows:

```powershell
.\goroutes.exe `
  --input "C:\Users\User\Downloads\ip-address" `
  --output ".\keenetic-routes.bat"
```

## Command-line options

```text
--input
    Directory containing source BAT files.
    This option is required.

--output
    Path and filename of the generated BAT file.
    Default: routes.bat

--version
    Print the application version.
```

Show help:

```bash
./goroutes --help
```

Show version:

```bash
./goroutes --version
```

## Output

After a successful run, the application prints statistics:

```text
Go Routes

Parsed routes:    1247
Optimized routes: 727
Removed routes:   520

Created: /path/to/keenetic-routes.bat
```

Generated commands look like this:

```bat
route add 142.250.0.0 mask 255.255.0.0 0.0.0.0 & rem google, youtube
```

## Import into Keenetic

1. Generate the resulting BAT file.
2. Open the Keenetic web interface.
3. Import or execute the generated route commands using the appropriate Keenetic configuration workflow.
4. Select the required VPN interface for the imported routes.

Always review generated routes before applying them to a production router.

## Optimization rules

Go Routes performs only lossless route optimizations.

### Exact duplicates

```text
10.0.0.0/24
10.0.0.0/24
```

Become:

```text
10.0.0.0/24
```

Descriptions from both routes are preserved.

### Covered networks

```text
10.0.0.0/16
10.0.20.0/24
```

Become:

```text
10.0.0.0/16
```

The `/24` route is already fully covered by the `/16` route.

### Sibling networks

```text
10.0.0.0/25
10.0.0.128/25
```

Become:

```text
10.0.0.0/24
```

Sibling networks are merged only when together they cover the parent network completely.

## Project structure

```text
.
├── cmd/
│   └── goroutes/
│       └── main.go
├── internal/
│   ├── app/
│   ├── generator/
│   ├── optimizer/
│   ├── parser/
│   └── route/
├── go.mod
├── go.sum
├── .gitignore
└── README.md
```

## Development

Format the project:

```bash
gofmt -w .
```

Run all tests:

```bash
go test ./...
```

Run static analysis:

```bash
go vet ./...
```

Run everything before committing:

```bash
gofmt -w .
go test ./...
go vet ./...
go build ./cmd/goroutes
```

## License

This project is licensed under the MIT License.