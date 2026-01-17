package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFile(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T, tmpDir string) (src, dst string)
		validate func(t *testing.T, src, dst string)
		wantErr  bool
	}{
		{
			name: "Success",
			setup: func(t *testing.T, tmpDir string) (string, string) {
				src := filepath.Join(tmpDir, "src.txt")
				dst := filepath.Join(tmpDir, "dst.txt")
				content := "hello world"
				if err := os.WriteFile(src, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to create src file: %v", err)
				}
				return src, dst
			},
			validate: func(t *testing.T, src, dst string) {
				got, err := os.ReadFile(dst)
				if err != nil {
					t.Fatalf("Failed to read dst file: %v", err)
				}
				if string(got) != "hello world" {
					t.Errorf("Expected content 'hello world', got %q", string(got))
				}

				srcInfo, _ := os.Stat(src)
				dstInfo, _ := os.Stat(dst)
				if srcInfo.Mode() != dstInfo.Mode() {
					t.Errorf("Expected mode %v, got %v", srcInfo.Mode(), dstInfo.Mode())
				}
			},
			wantErr: false,
		},
		{
			name: "Source Not Found",
			setup: func(t *testing.T, tmpDir string) (string, string) {
				return filepath.Join(tmpDir, "nonexistent.txt"), filepath.Join(tmpDir, "out.txt")
			},
			wantErr: true,
		},
		{
			name: "Destination Error",
			setup: func(t *testing.T, tmpDir string) (string, string) {
				src := filepath.Join(tmpDir, "src.txt")
				if err := os.WriteFile(src, []byte("content"), 0644); err != nil {
					t.Fatal(err)
				}
				// Trying to write to a path where a directory already exists (or parent dir invalid)
				invalidDst := filepath.Join(tmpDir, "some-dir")
				os.Mkdir(invalidDst, 0755)
				// trying to write to invalidDst/another-dir/file.txt but invalidDst/another-dir doesn't exist
				return src, filepath.Join(invalidDst, "another-dir", "file.txt")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			src, dst := tt.setup(t, tmpDir)

			err := CopyFile(src, dst)
			if (err != nil) != tt.wantErr {
				t.Errorf("CopyFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, src, dst)
			}
		})
	}
}

func TestCopyDir(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T, tmpDir string) (src, dst string)
		validate func(t *testing.T, dst string)
		wantErr  bool
	}{
		{
			name: "Success",
			setup: func(t *testing.T, tmpDir string) (string, string) {
				srcDir := filepath.Join(tmpDir, "src")
				dstDir := filepath.Join(tmpDir, "dst")

				// Create source structure
				if err := os.MkdirAll(filepath.Join(srcDir, "sub"), 0755); err != nil {
					t.Fatal(err)
				}
				os.WriteFile(filepath.Join(srcDir, "file1.txt"), []byte("content1"), 0644)
				os.WriteFile(filepath.Join(srcDir, "sub", "file2.txt"), []byte("content2"), 0644)

				return srcDir, dstDir
			},
			validate: func(t *testing.T, dstDir string) {
				// Verify file1
				got1, err := os.ReadFile(filepath.Join(dstDir, "file1.txt"))
				if err != nil {
					t.Errorf("file1.txt missing in dst: %v", err)
				} else if string(got1) != "content1" {
					t.Errorf("file1 content mismatch. Got %q", string(got1))
				}

				// Verify file2
				got2, err := os.ReadFile(filepath.Join(dstDir, "sub", "file2.txt"))
				if err != nil {
					t.Errorf("sub/file2.txt missing in dst: %v", err)
				} else if string(got2) != "content2" {
					t.Errorf("file2 content mismatch. Got %q", string(got2))
				}
			},
			wantErr: false,
		},
		{
			name: "Source Not Found",
			setup: func(t *testing.T, tmpDir string) (string, string) {
				return filepath.Join(tmpDir, "nonexistent-dir"), filepath.Join(tmpDir, "out-dir")
			},
			wantErr: true,
		},
		{
			name: "Destination Mkdir Error",
			setup: func(t *testing.T, tmpDir string) (string, string) {
				srcDir := filepath.Join(tmpDir, "src")
				os.MkdirAll(srcDir, 0755)
				// Create a file where a directory should be created
				blockedDst := filepath.Join(tmpDir, "blocked")
				os.WriteFile(blockedDst, []byte("i am a file"), 0644)
				return srcDir, filepath.Join(blockedDst, "sub")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			src, dst := tt.setup(t, tmpDir)

			err := CopyDir(src, dst)
			if (err != nil) != tt.wantErr {
				t.Errorf("CopyDir() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, dst)
			}
		})
	}
}
