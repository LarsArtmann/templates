package domain

import (
	"errors"
	"testing"
)

// MockFileRepository is a mock implementation of FileRepository for testing
type MockFileRepository struct {
	existsFunc   func(path string) (bool, error)
	generateFunc func(path, templatePath string, data interface{}) error
}

func (m *MockFileRepository) CheckExists(path string) (bool, error) {
	return m.existsFunc(path)
}

func (m *MockFileRepository) Generate(path, templatePath string, data interface{}) error {
	return m.generateFunc(path, templatePath, data)
}

func TestFileService_Validate(t *testing.T) {
	tests := []struct {
		name      string
		file      *File
		existsFunc func(path string) (bool, error)
		wantErr   bool
		wantExists bool
	}{
		{
			name: "file exists",
			file: &File{
				Path:     "test.txt",
				Required: true,
			},
			existsFunc: func(path string) (bool, error) {
				return true, nil
			},
			wantErr:   false,
			wantExists: true,
		},
		{
			name: "file does not exist but not required",
			file: &File{
				Path:     "test.txt",
				Required: false,
			},
			existsFunc: func(path string) (bool, error) {
				return false, nil
			},
			wantErr:   false,
			wantExists: false,
		},
		{
			name: "file does not exist and required",
			file: &File{
				Path:     "test.txt",
				Required: true,
			},
			existsFunc: func(path string) (bool, error) {
				return false, nil
			},
			wantErr:   true,
			wantExists: false,
		},
		{
			name: "error checking if file exists",
			file: &File{
				Path:     "test.txt",
				Required: true,
			},
			existsFunc: func(path string) (bool, error) {
				return false, errors.New("test error")
			},
			wantErr:   true,
			wantExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockFileRepository{
				existsFunc: tt.existsFunc,
			}
			
			service := NewFileService(repo)
			
			err := service.Validate(tt.file)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("FileService.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			
			if tt.file.Exists != tt.wantExists {
				t.Errorf("FileService.Validate() file.Exists = %v, wantExists %v", tt.file.Exists, tt.wantExists)
			}
		})
	}
}

func TestFileService_Generate(t *testing.T) {
	tests := []struct {
		name         string
		file         *File
		generateFunc func(path, templatePath string, data interface{}) error
		wantErr      bool
	}{
		{
			name: "generate file successfully",
			file: &File{
				Path:     "test.txt",
				Template: "test.tmpl",
			},
			generateFunc: func(path, templatePath string, data interface{}) error {
				return nil
			},
			wantErr: false,
		},
		{
			name: "no template",
			file: &File{
				Path:     "test.txt",
				Template: "",
			},
			generateFunc: func(path, templatePath string, data interface{}) error {
				return nil
			},
			wantErr: true,
		},
		{
			name: "error generating file",
			file: &File{
				Path:     "test.txt",
				Template: "test.tmpl",
			},
			generateFunc: func(path, templatePath string, data interface{}) error {
				return errors.New("test error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockFileRepository{
				generateFunc: tt.generateFunc,
			}
			
			service := NewFileService(repo)
			
			err := service.Generate(tt.file, nil)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("FileService.Generate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
