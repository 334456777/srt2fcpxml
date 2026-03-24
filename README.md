# srt2fcpxml

Convert SRT subtitle files to Final Cut Pro FCPXML subtitle format.

This software uses Final Cut Pro X 10.4.6 version FCPXML files as a template. If you encounter any issues, please upgrade to the corresponding version.

## 🚀 Fork Improvements

This fork includes the following improvements over the original project:

### ✨ New Features
- **Auto-find SRT files** - No need to manually specify file paths
- **Simplified CLI** - More intuitive parameter settings
- **Version query support** - Supports `--version` parameter
- **Enhanced resolution settings** - Supports resolution settings with default frame rate

### 🤖 CI/CD Automation
- **Multi-platform builds** - Supports Windows, macOS, Linux
- **Native architecture support** - Supports AMD64 and ARM64
- **Automated release workflow** - Auto-release via GitHub Actions
- **Quality assurance** - Automated testing and validation

### 🔧 Technical Improvements
- **Dependency updates** - Modern Go module management
- **Optimized output format** - Clearer user feedback
- **Cross-platform compatibility** - Enhanced platform compatibility
- **Project structure optimization** - Standardized configuration and structure

## Compile

First, ensure you have the Go language development environment installed.

Then execute the `make` command in the project directory to generate the `srt2fcpxml` executable file in the `build` directory.

## Download

Users who prefer not to compile can download the [executable files](https://github.com/334456777/srt2fcpxml/releases) directly.

### Supported Platforms
- Windows (AMD64, ARM64)
- macOS (Intel, Apple Silicon)
- Linux (AMD64, ARM64)

## Usage

First, grant the program execute permission:
```bash
chmod +x ./srt2fcpxml
```

The program will automatically find SRT files in the current directory and convert them.

### Usage Patterns

```bash
# Auto find SRT file with default settings (1920x1080@30fps)
./srt2fcpxml

# Auto find SRT file with specified frame rate (1920x1080@60fps)
./srt2fcpxml 60

# Auto find SRT file with custom resolution and default frame rate (30fps)
./srt2fcpxml 1920 1080

# Auto find SRT file with custom resolution and frame rate
./srt2fcpxml 1920 1080 29.97
```

### Supported Frame Rates
23.98, 24, 25, 29.97, 30, 50, 59.94, 60

## Execution Examples

```bash
# Convert with default settings
./srt2fcpxml

# Convert with 60fps
./srt2fcpxml 60

# Convert with custom resolution and default frame rate
./srt2fcpxml 1920 1080

# Convert with custom settings
./srt2fcpxml 1920 1080 29.97
```

The `fcpxml` file will be automatically generated in the same directory as the SRT file, named after the SRT file.

## Development

### Local Build
```bash
# Build for current platform
go build -o srt2fcpxml cmd/main.go

# Build for all platforms using the build script
./build.sh v1.0.0
```

### Release Process
1. Create and push a version tag:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. GitHub Actions will automatically:
   - ✅ Run tests on Ubuntu
   - 🔨 Build binaries using native runners for better performance:
     - **Windows** builds on `windows-latest`
     - **macOS** builds on `macos-latest`
     - **Linux** builds on `ubuntu-latest`
   - 🎯 Support for AMD64 and ARM64 architectures
   - 📦 Create GitHub release with all binaries
   - 🔐 Generate SHA256 checksums for security
   - ✨ Test native binaries for quality assurance

### Workflow Triggers
- **Tag push** (`v*`) → Full build + release
- **Branch push** (main/master/develop) → Build only
- **Pull Request** → Test only
- **Manual trigger** → Custom version build

### Manual Trigger
You can manually trigger builds through the GitHub Actions interface with custom version numbers.

## License

MIT
