package filesys

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/heroku/self/MetaManager/internal/file"
	contextrepo "github.com/heroku/self/MetaManager/internal/repository/filesys"
	contextmocks "github.com/heroku/self/MetaManager/mocks/repository/filesys"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", "/"},
		{"dot", ".", "/"},
		{"root", "/", "/"},
		{"single folder", "Folder", "/Folder"},
		{"single folder with leading slash", "/Folder", "/Folder"},
		{"single folder with trailing slash", "Folder/", "/Folder"},
		{"nested path", "/Folder/Sub", "/Folder/Sub"},
		{"nested path no leading slash", "Folder/Sub", "/Folder/Sub"},
		{"nested path with trailing slash", "/Folder/Sub/", "/Folder/Sub"},
		{"whitespace", "  /Folder/Sub  ", "/Folder/Sub"},
		{"multiple slashes", "///Folder///Sub///", "///Folder///Sub//"},
		{"just slashes", "///", "//"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizePath(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestResolvePath(t *testing.T) {
	tests := []struct {
		name     string
		cwd      string
		target   string
		expected string
	}{
		// Empty cwd defaults to root
		{"empty cwd, empty target", "", "", "/"},
		{"empty cwd, dot target", "", ".", "/"},
		{"empty cwd, root target", "", "/", "/"},
		{"empty cwd, absolute path", "", "/Folder", "/Folder"},
		{"empty cwd, relative path", "", "Folder", "/Folder"},

		// Root cwd
		{"root cwd, empty target", "/", "", "/"},
		{"root cwd, dot target", "/", ".", "/"},
		{"root cwd, root target", "/", "/", "/"},
		{"root cwd, absolute path", "/", "/Folder", "/Folder"},
		{"root cwd, relative path", "/", "Folder", "/Folder"},

		// Nested cwd
		{"nested cwd, empty target", "/Folder", "", "/Folder"},
		{"nested cwd, dot target", "/Folder", ".", "/Folder"},
		{"nested cwd, absolute path", "/Folder", "/Sub", "/Sub"},
		{"nested cwd, relative path", "/Folder", "Sub", "/Folder/Sub"},
		{"nested cwd, parent relative", "/Folder/Sub", "..", "/Folder"},
		{"nested cwd, parent parent relative", "/Folder/Sub", "../..", "/"},
		{"nested cwd, complex relative", "/Folder/Sub", "../Other/File", "/Folder/Other/File"},

		// Whitespace handling
		{"whitespace cwd", "  /Folder  ", "Sub", "/Folder  /Sub"},
		{"whitespace target", "/Folder", "  Sub  ", "/Folder/Sub"},

		// Edge cases
		{"trailing slash cwd", "/Folder/", "Sub", "/Folder/Sub"},
		{"trailing slash target", "/Folder", "Sub/", "/Folder/Sub"},
		{"both trailing slashes", "/Folder/", "Sub/", "/Folder/Sub"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ResolvePath(tt.cwd, tt.target)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBasicResolver_Resolve_Local(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		setupMock func(*contextmocks.MockContextRepository)
		checkPath func(*testing.T, string)
		wantError bool
	}{
		{
			name: "resolve relative path",
			path: "testfile.txt",
			setupMock: func(m *contextmocks.MockContextRepository) {
				m.On("GetContext").Return("localctx", nil)
				m.On("GetContextType", "localctx").Return(contextrepo.TypeLocal, nil)
			},
			checkPath: func(t *testing.T, resolved string) {
				abs, _ := filepath.Abs("testfile.txt")
				assert.Equal(t, abs, resolved)
			},
			wantError: false,
		},
		{
			name: "resolve current directory",
			path: ".",
			setupMock: func(m *contextmocks.MockContextRepository) {
				m.On("GetContext").Return("localctx", nil)
				m.On("GetContextType", "localctx").Return(contextrepo.TypeLocal, nil)
			},
			checkPath: func(t *testing.T, resolved string) {
				abs, _ := filepath.Abs(".")
				assert.Equal(t, abs, resolved)
			},
			wantError: false,
		},
		{
			name: "get context error",
			path: "testfile.txt",
			setupMock: func(m *contextmocks.MockContextRepository) {
				m.On("GetContext").Return("", assert.AnError)
			},
			wantError: true,
		},
		{
			name: "get context type error",
			path: "testfile.txt",
			setupMock: func(m *contextmocks.MockContextRepository) {
				m.On("GetContext").Return("localctx", nil)
				m.On("GetContextType", "localctx").Return("", assert.AnError)
			},
			wantError: true,
		},
		{
			name: "invalid context type",
			path: "testfile.txt",
			setupMock: func(m *contextmocks.MockContextRepository) {
				m.On("GetContext").Return("invalidctx", nil)
				m.On("GetContextType", "invalidctx").Return("invalid", nil)
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := contextmocks.NewMockContextRepository(t)
			tt.setupMock(mockRepo)

			resolver := NewBasicResolver(mockRepo)
			result, err := resolver.Resolve(tt.path)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.checkPath != nil {
					tt.checkPath(t, result)
				}
			}
		})
	}
}

func TestBasicResolver_Resolve_GDrive(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		cwd       string
		setupMock func(*contextmocks.MockContextRepository, string)
		expected  string
		wantError bool
	}{
		{
			name: "resolve root",
			path: "/",
			cwd:  "/",
			setupMock: func(m *contextmocks.MockContextRepository, cwd string) {
				m.On("GetContext").Return("gdrivectx", nil)
				m.On("GetContextType", "gdrivectx").Return(contextrepo.TypeGDrive, nil)
				m.On("GetGDriveCwd").Return(cwd, nil)
			},
			expected:  file.GDrivePathPrefix,
			wantError: false,
		},
		{
			name: "resolve absolute path",
			path: "/Folder",
			cwd:  "/",
			setupMock: func(m *contextmocks.MockContextRepository, cwd string) {
				m.On("GetContext").Return("gdrivectx", nil)
				m.On("GetContextType", "gdrivectx").Return(contextrepo.TypeGDrive, nil)
				m.On("GetGDriveCwd").Return(cwd, nil)
			},
			expected:  file.GDrivePathRoot + "/Folder",
			wantError: false,
		},
		{
			name: "resolve relative path from root",
			path: "Folder",
			cwd:  "/",
			setupMock: func(m *contextmocks.MockContextRepository, cwd string) {
				m.On("GetContext").Return("gdrivectx", nil)
				m.On("GetContextType", "gdrivectx").Return(contextrepo.TypeGDrive, nil)
				m.On("GetGDriveCwd").Return(cwd, nil)
			},
			expected:  file.GDrivePathRoot + "/Folder",
			wantError: false,
		},
		{
			name: "resolve relative path from nested cwd",
			path: "Sub",
			cwd:  "/Folder",
			setupMock: func(m *contextmocks.MockContextRepository, cwd string) {
				m.On("GetContext").Return("gdrivectx", nil)
				m.On("GetContextType", "gdrivectx").Return(contextrepo.TypeGDrive, nil)
				m.On("GetGDriveCwd").Return(cwd, nil)
			},
			expected:  file.GDrivePathRoot + "/Folder/Sub",
			wantError: false,
		},
		{
			name: "resolve current directory",
			path: ".",
			cwd:  "/Folder",
			setupMock: func(m *contextmocks.MockContextRepository, cwd string) {
				m.On("GetContext").Return("gdrivectx", nil)
				m.On("GetContextType", "gdrivectx").Return(contextrepo.TypeGDrive, nil)
				m.On("GetGDriveCwd").Return(cwd, nil)
			},
			expected:  file.GDrivePathRoot + "/Folder",
			wantError: false,
		},
		{
			name: "resolve parent directory",
			path: "..",
			cwd:  "/Folder/Sub",
			setupMock: func(m *contextmocks.MockContextRepository, cwd string) {
				m.On("GetContext").Return("gdrivectx", nil)
				m.On("GetContextType", "gdrivectx").Return(contextrepo.TypeGDrive, nil)
				m.On("GetGDriveCwd").Return(cwd, nil)
			},
			expected:  file.GDrivePathRoot + "/Folder",
			wantError: false,
		},
		{
			name: "resolve nested absolute path",
			path: "/Folder/Sub/File",
			cwd:  "/Other",
			setupMock: func(m *contextmocks.MockContextRepository, cwd string) {
				m.On("GetContext").Return("gdrivectx", nil)
				m.On("GetContextType", "gdrivectx").Return(contextrepo.TypeGDrive, nil)
				m.On("GetGDriveCwd").Return(cwd, nil)
			},
			expected:  file.GDrivePathRoot + "/Folder/Sub/File",
			wantError: false,
		},
		{
			name: "get gdrive cwd error",
			path: "Folder",
			cwd:  "/",
			setupMock: func(m *contextmocks.MockContextRepository, cwd string) {
				m.On("GetContext").Return("gdrivectx", nil)
				m.On("GetContextType", "gdrivectx").Return(contextrepo.TypeGDrive, nil)
				m.On("GetGDriveCwd").Return("", assert.AnError)
			},
			wantError: true,
		},
		{
			name: "empty cwd defaults to root",
			path: "Folder",
			cwd:  "",
			setupMock: func(m *contextmocks.MockContextRepository, cwd string) {
				m.On("GetContext").Return("gdrivectx", nil)
				m.On("GetContextType", "gdrivectx").Return(contextrepo.TypeGDrive, nil)
				m.On("GetGDriveCwd").Return(cwd, nil)
			},
			expected:  file.GDrivePathRoot + "/Folder",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := contextmocks.NewMockContextRepository(t)
			tt.setupMock(mockRepo, tt.cwd)

			resolver := NewBasicResolver(mockRepo)
			result, err := resolver.Resolve(tt.path)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestBasicResolver_resolveLocal(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "testfile.txt")
	require.NoError(t, os.WriteFile(testFile, []byte("test"), 0644))

	mockRepo := contextmocks.NewMockContextRepository(t)
	resolver := NewBasicResolver(mockRepo)

	tests := []struct {
		name     string
		path     string
		expected string
		wantErr  bool
	}{
		{
			name:     "absolute path",
			path:     testFile,
			expected: testFile,
			wantErr:  false,
		},
		{
			name:    "relative path",
			path:    "testfile.txt",
			wantErr: false,
			// Expected will be current working dir + testfile.txt
		},
		{
			name:    "current directory",
			path:    ".",
			wantErr: false,
			// Expected will be current working directory
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Change to tmpDir for relative path tests
			if tt.path == "testfile.txt" || tt.path == "." {
				oldWd, _ := os.Getwd()
				require.NoError(t, os.Chdir(tmpDir))
				defer os.Chdir(oldWd)
			}

			result, err := resolver.resolveLocal(tt.path)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.expected != "" {
					assert.Equal(t, tt.expected, result)
				} else {
					// For relative paths, just check it's an absolute path
					assert.True(t, filepath.IsAbs(result))
				}
			}
		})
	}
}

func TestBasicResolver_resolveGDrive(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		cwd       string
		setupMock func(*contextmocks.MockContextRepository, string)
		expected  string
		wantError bool
	}{
		{
			name: "root path",
			path: "/",
			cwd:  "/",
			setupMock: func(m *contextmocks.MockContextRepository, cwd string) {
				m.On("GetGDriveCwd").Return(cwd, nil)
			},
			expected:  file.GDrivePathPrefix,
			wantError: false,
		},
		{
			name: "single folder",
			path: "/Folder",
			cwd:  "/",
			setupMock: func(m *contextmocks.MockContextRepository, cwd string) {
				m.On("GetGDriveCwd").Return(cwd, nil)
			},
			expected:  file.GDrivePathRoot + "/Folder",
			wantError: false,
		},
		{
			name: "nested path",
			path: "/Folder/Sub",
			cwd:  "/",
			setupMock: func(m *contextmocks.MockContextRepository, cwd string) {
				m.On("GetGDriveCwd").Return(cwd, nil)
			},
			expected:  file.GDrivePathRoot + "/Folder/Sub",
			wantError: false,
		},
		{
			name: "relative path",
			path: "Sub",
			cwd:  "/Folder",
			setupMock: func(m *contextmocks.MockContextRepository, cwd string) {
				m.On("GetGDriveCwd").Return(cwd, nil)
			},
			expected:  file.GDrivePathRoot + "/Folder/Sub",
			wantError: false,
		},
		{
			name: "get cwd error",
			path: "/Folder",
			cwd:  "/",
			setupMock: func(m *contextmocks.MockContextRepository, cwd string) {
				m.On("GetGDriveCwd").Return("", assert.AnError)
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := contextmocks.NewMockContextRepository(t)
			tt.setupMock(mockRepo, tt.cwd)

			resolver := NewBasicResolver(mockRepo)
			result, err := resolver.resolveGDrive(tt.path)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestNewBasicResolver(t *testing.T) {
	mockRepo := contextmocks.NewMockContextRepository(t)
	resolver := NewBasicResolver(mockRepo)

	assert.NotNil(t, resolver)
	assert.Equal(t, mockRepo, resolver.ctxRepo)
}
