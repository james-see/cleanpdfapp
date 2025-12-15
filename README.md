# Clean PDF

A fast, native desktop application to view and wipe metadata from PDF files. Built with Go and [Fyne](https://fyne.io/).

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.22+-00ADD8.svg)
![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Windows%20%7C%20Linux-lightgrey.svg)

## Features

- **Analyze Mode** - View all metadata hidden in your PDFs (author, creation date, software used, etc.) and export it to a text file
- **Clean Mode** - Wipe all metadata from a PDF and save a pristine copy
- **Cross-platform** - Native builds for macOS, Windows, and Linux
- **Privacy-focused** - No telemetry, no network access, fully offline

## Installation

### macOS (Homebrew)

```bash
brew install james-see/tap/cleanpdf
```

### Manual Download

Download the latest release for your platform from [GitHub Releases](https://github.com/james-see/cleanpdfapp/releases/latest):

- **macOS**: `CleanPDF-macos-vX.X.X.zip`
- **Windows**: `CleanPDF-windows-vX.X.X.zip`
- **Linux**: `CleanPDF-linux-vX.X.X.tar.gz`

## Usage

1. Launch Clean PDF
2. Choose your mode:
   - **ANALYZE MODE** - Select a PDF to view its metadata. Optionally save the metadata to a text file.
   - **CLEAN MODE** - Select a PDF to create a clean copy with all metadata removed.
3. Clean PDFs are saved with `-clean.pdf` suffix in the same directory as the original.

## Building from Source

### Prerequisites

- Go 1.22 or later
- For GUI: Platform-specific dependencies for [Fyne](https://developer.fyne.io/started/)

### Build

```bash
# Clone the repository
git clone https://github.com/james-see/cleanpdfapp.git
cd cleanpdfapp

# Install dependencies
go mod download

# Build
go build -o cleanpdf .

# Or use make
make build
```

### Cross-platform Builds

```bash
# Build macOS app bundle
make package-macos

# Build for all platforms (requires fyne-cross or platform toolchains)
make build-macos
make build-windows
make build-linux
```

## Project Structure

```
cleanpdfapp/
├── main.go              # Application entry point
├── ui/
│   └── window.go        # Fyne GUI implementation
├── pdf/
│   └── metadata.go      # PDF metadata operations
├── go.mod               # Go module definition
├── Makefile             # Build automation
├── FyneApp.toml         # Fyne packaging config
├── Icon.png             # Application icon
├── docs/                # GitHub Pages site
└── homebrew/            # Homebrew cask formula
```

## Dependencies

- [Fyne](https://fyne.io/) - Cross-platform GUI toolkit for Go
- [pdfcpu](https://github.com/pdfcpu/pdfcpu) - PDF processing library (pure Go)

## Version History

- **v2.0.0** - Complete rewrite in Go/Fyne, cross-platform support
- **v1.1.0** - Python/Tkinter version (legacy)

## License

MIT License - see [LICENSE](LICENSE) for details.

## Author

[James Campbell](https://github.com/james-see)

---

[Website](https://james-see.github.io/cleanpdfapp/) • [Report Issue](https://github.com/james-see/cleanpdfapp/issues)
