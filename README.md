# Filesys Package

This package provides file system scanning and tracking functionality for both local file systems and Google Drive.

## Overview

The `filesys` package contains:

- **Scanner Interface**: Defines the interface for scanning directories and returning tree nodes
- **File System Scanner**: Implementation for local Unix file systems
- **Google Drive Scanner**: Implementation for Google Drive directories
- **Context-Aware Tracker**: Tracks files/directories based on the current context (local or gdrive)

## Mocking with Mockery

This package uses [mockery](https://github.com/vektra/mockery) to generate mocks for testing. The mocks are generated for interfaces to enable unit testing without external dependencies.

### Setup

1. **Install mockery**:
   ```bash
   go install github.com/vektra/mockery/v3@latest
   ```

2. **Install gomock dependency** (required for expecter-style mocks):
   ```bash
   go get go.uber.org/mock/gomock
   ```

3. **Generate mocks**:
   ```bash
   # From the project root
   mockery
   ```

   This will generate mocks based on the `.mockery.yaml` configuration file in the project root.
   
   **Note**: Mock generation requires network access to download Go dependencies if needed.

### Generated Mocks

The following mocks are generated:

- `internal/services/mocks/MockGDriveServiceInterface.go` - Mock for Google Drive service
- `internal/filesys/mocks/MockScanner.go` - Mock for Scanner interface
- `internal/filesys/mocks/MockScannableNode.go` - Mock for ScannableNode interface

### Using Mocks in Tests

#### Example: Testing with MockGDriveServiceInterface

```go
package cmd_test

import (
    "context"
    "testing"
    
    "github.com/heroku/self/MetaManager/internal/filesys"
    "github.com/heroku/self/MetaManager/internal/services"
    "github.com/heroku/self/MetaManager/internal/services/mocks"
    "github.com/stretchr/testify/require"
    "go.uber.org/mock/gomock"
)

func TestTrackGDrive(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    // Create a mock GDrive service
    mockSvc := mocks.NewMockGDriveServiceInterface(ctrl)
    
    // Set up expectations
    mockSvc.EXPECT().
        ResolvePath(gomock.Any(), "/Folder1").
        Return("folder1", nil)
    
    mockSvc.EXPECT().
        ListFolder(gomock.Any(), "folder1").
        Return([]services.RootEntry{
            {Id: "sub1", Name: "Sub", IsFolder: true, MimeType: services.DriveFolderMimeType},
            {Id: "file1", Name: "file1.txt", IsFolder: false, MimeType: "text/plain"},
        }, nil)
    
    // Create scanner with mock
    scanner := filesys.NewGDriveScanner(mockSvc)
    
    // Test the functionality
    tree, err := scanner.TrackGDrive(context.Background(), "/Folder1", false)
    require.NoError(t, err)
    require.NotNil(t, tree)
}
```

**Note**: The mockery-generated mocks use `go.uber.org/mock/gomock` for expecter-style mocking. Make sure to add this dependency:

```bash
go get go.uber.org/mock/gomock
```

#### Example: Testing with MockScanner

```go
func TestContextAwareTracker(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    // Create a mock scanner
    mockScanner := mocks.NewMockScanner(ctrl)
    
    // Set up expectations
    mockScanner.EXPECT().
        Scan("/some/path").
        Return(&ds.TreeNode{...}, nil)
    
    // Use the mock in your test
    // ...
}
```

### Regenerating Mocks

After modifying interfaces, regenerate the mocks:

```bash
mockery
```

Or regenerate specific interfaces using the interface name:

```bash
mockery --name GDriveServiceInterface --dir internal/services
```

**Note**: Mock generation requires network access if Go dependencies need to be downloaded.

### Mockery Configuration

The mockery configuration is stored in `.mockery.yaml` at the project root. Key settings:

- `template: testify` - Uses the testify mock template (supports expecter-style mocks)
- `template-data.with-expecter: true` - Generates expecter-style mocks (gomock compatible)
- `structname: "Mock{{.InterfaceName}}"` - Naming convention for generated mocks
- `dir: "mocks"` - Output directory for mocks (relative to interface package)
- `filename: "{{.MockName}}.go"` - Filename pattern for generated mocks

The configuration generates mocks for:
- `GDriveServiceInterface` in `internal/services/mocks/`
- `Scanner` and `ScannableNode` in `internal/filesys/mocks/`

### Best Practices

1. **Always use mocks in tests** - Don't rely on real external services (like Google Drive API) in unit tests
2. **Set up expectations** - Use `EXPECT()` to define what methods should be called and what they should return
3. **Verify calls** - Use `gomock.Controller` to ensure all expected calls were made
4. **Keep mocks in sync** - Regenerate mocks after interface changes
5. **Commit mocks** - Include generated mock files in version control for consistency

### Troubleshooting

**Issue**: Mocks not generating
- **Solution**: Ensure mockery is installed and `.mockery.yaml` is in the project root

**Issue**: Import errors in generated mocks
- **Solution**: Run `go mod tidy` after generating mocks

**Issue**: Mocks out of date
- **Solution**: Regenerate mocks with `mockery` command

## Package Structure

```
internal/filesys/
├── scanner.go              # Scanner interface and factory
├── file_scanner.go         # Local file system scanner
├── gdrive_scanner.go       # Google Drive scanner
├── track.go                # Context-aware tracker
├── mocks/                  # Generated mocks (git tracked)
│   ├── MockScanner.go
│   └── MockScannableNode.go
└── README.md               # This file
```

## Related Packages

- `internal/services` - Google Drive service implementation
- `internal/data` - Tree data structures and management
- `internal/file` - File node definitions
- `internal/repository/filesys` - Context repository for file system contexts
