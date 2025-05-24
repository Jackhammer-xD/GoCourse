package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// FileServer represents a simple file sharing server
type FileServer struct {
	uploadDir string
}

// NewFileServer creates a new instance of FileServer
func NewFileServer(uploadDir string) *FileServer {
	return &FileServer{
		uploadDir: uploadDir,
	}
}

// handleUpload handles file upload requests
func (fs *FileServer) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(fs.uploadDir, 0755); err != nil {
		http.Error(w, "Error creating upload directory", http.StatusInternalServerError)
		return
	}

	// Create destination file
	dst, err := os.Create(filepath.Join(fs.uploadDir, header.Filename))
	if err != nil {
		http.Error(w, "Error creating destination file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy file contents
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File %s uploaded successfully", header.Filename)
}

// handleDownload handles file download requests
func (fs *FileServer) handleDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filename := r.URL.Query().Get("file")
	if filename == "" {
		http.Error(w, "No file specified", http.StatusBadRequest)
		return
	}

	filepath := filepath.Join(fs.uploadDir, filename)
	file, err := os.Open(filepath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Set headers for file download
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Type", "application/octet-stream")

	// Stream file to client
	if _, err := io.Copy(w, file); err != nil {
		http.Error(w, "Error sending file", http.StatusInternalServerError)
		return
	}
}

// handleList handles requests to list available files
func (fs *FileServer) handleList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	files, err := os.ReadDir(fs.uploadDir)
	if err != nil {
		http.Error(w, "Error reading directory", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Available files:")
	for _, file := range files {
		if !file.IsDir() {
			fmt.Fprintf(w, "- %s\n", file.Name())
		}
	}
}

func main() {
	uploadDir := "./uploads"
	server := NewFileServer(uploadDir)

	// Register handlers
	http.HandleFunc("/upload", server.handleUpload)
	http.HandleFunc("/download", server.handleDownload)
	http.HandleFunc("/list", server.handleList)

	// Start server
	fmt.Println("File sharing server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
