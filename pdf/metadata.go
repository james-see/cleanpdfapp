// Package pdf provides PDF metadata operations
package pdf

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// Metadata represents PDF document metadata
type Metadata struct {
	Title        string
	Author       string
	Subject      string
	Keywords     string
	Creator      string
	Producer     string
	CreationDate string
	ModDate      string
}

// ReadMetadata extracts metadata from a PDF file
func ReadMetadata(filepath string) (*Metadata, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	info, err := api.PDFInfo(f, filepath, nil, false, model.NewDefaultConfiguration())
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF info: %w", err)
	}

	meta := &Metadata{
		Title:        info.Title,
		Author:       info.Author,
		Subject:      info.Subject,
		Keywords:     strings.Join(info.Keywords, ", "),
		Creator:      info.Creator,
		Producer:     info.Producer,
		CreationDate: info.CreationDate,
		ModDate:      info.ModificationDate,
	}

	return meta, nil
}

// ToMap converts metadata to a map for display
func (m *Metadata) ToMap() map[string]string {
	result := make(map[string]string)
	if m.Title != "" {
		result["Title"] = m.Title
	}
	if m.Author != "" {
		result["Author"] = m.Author
	}
	if m.Subject != "" {
		result["Subject"] = m.Subject
	}
	if m.Keywords != "" {
		result["Keywords"] = m.Keywords
	}
	if m.Creator != "" {
		result["Creator"] = m.Creator
	}
	if m.Producer != "" {
		result["Producer"] = m.Producer
	}
	if m.CreationDate != "" {
		result["CreationDate"] = m.CreationDate
	}
	if m.ModDate != "" {
		result["ModDate"] = m.ModDate
	}
	return result
}

// ToString formats metadata as a readable string
func (m *Metadata) ToString() string {
	var sb strings.Builder
	metaMap := m.ToMap()

	if len(metaMap) == 0 {
		return "No metadata found"
	}

	for key, value := range metaMap {
		sb.WriteString(fmt.Sprintf("%s: %s\n", key, value))
	}
	return sb.String()
}

// Count returns the number of non-empty metadata fields
func (m *Metadata) Count() int {
	return len(m.ToMap())
}

// WipeMetadata creates a new PDF with all metadata removed
func WipeMetadata(srcPath, destPath string) error {
	f, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer f.Close()

	conf := model.NewDefaultConfiguration()
	ctx, err := api.ReadContext(f, conf)
	if err != nil {
		return fmt.Errorf("failed to read PDF: %w", err)
	}

	// Clear all metadata fields on the XRefTable
	if ctx.XRefTable != nil {
		ctx.XRefTable.Title = ""
		ctx.XRefTable.Author = ""
		ctx.XRefTable.Subject = ""
		ctx.XRefTable.Creator = ""
		ctx.XRefTable.Producer = ""
		ctx.XRefTable.CreationDate = ""
		ctx.XRefTable.ModDate = ""
		ctx.XRefTable.Keywords = ""
		ctx.XRefTable.KeywordList = nil
	}

	// Write the cleaned PDF
	outFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	if err := api.WriteContext(ctx, outFile); err != nil {
		return fmt.Errorf("failed to write PDF: %w", err)
	}

	return nil
}

// SaveMetadataToFile saves metadata to a text file
func SaveMetadataToFile(meta *Metadata, outputPath string) error {
	content := meta.ToString()

	err := os.WriteFile(outputPath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	return nil
}

// GetCleanFilename generates the clean PDF filename
func GetCleanFilename(originalPath string) string {
	dir := filepath.Dir(originalPath)
	base := filepath.Base(originalPath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	return filepath.Join(dir, name+"-clean"+ext)
}

// GetMetadataFilename generates the metadata text filename
func GetMetadataFilename(originalPath string) string {
	dir := filepath.Dir(originalPath)
	base := filepath.Base(originalPath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	return filepath.Join(dir, name+"-metadata.txt")
}
