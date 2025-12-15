// Package ui provides the Fyne-based user interface
package ui

import (
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/james-see/cleanpdfapp/pdf"
)

// Run starts the Clean PDF application
func Run() {
	a := app.NewWithID("us.jamescampbell.cleanpdf")
	a.Settings().SetTheme(theme.DefaultTheme())

	w := a.NewWindow("CLEAN PDF")
	w.Resize(fyne.NewSize(600, 250))
	w.SetMinSize(fyne.NewSize(450, 180))
	w.CenterOnScreen()

	// Create styled buttons
	analyzeBtn := widget.NewButton("ANALYZE MODE - Select PDF to view metadata", func() {
		showAnalyzeDialog(w)
	})
	analyzeBtn.Importance = widget.HighImportance

	cleanBtn := widget.NewButton("CLEAN MODE - Select PDF to wipe metadata", func() {
		showCleanDialog(w)
	})
	cleanBtn.Importance = widget.MediumImportance

	// Layout
	content := container.NewVBox(
		widget.NewLabel("Select a PDF to analyze or clean:"),
		analyzeBtn,
		cleanBtn,
	)

	w.SetContent(container.NewPadded(content))
	w.ShowAndRun()
}

func showAnalyzeDialog(w fyne.Window) {
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if reader == nil {
			return // User cancelled
		}
		defer reader.Close()

		filePath := reader.URI().Path()
		analyzeFile(w, filePath)
	}, w)

	fd.SetFilter(&pdfFilter{})
	fd.Resize(fyne.NewSize(800, 600))
	fd.Show()
}

func analyzeFile(w fyne.Window, filePath string) {
	meta, err := pdf.ReadMetadata(filePath)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to read PDF: %w", err), w)
		return
	}

	count := meta.Count()
	if count == 0 {
		dialog.ShowInformation("No Metadata", fmt.Sprintf("No metadata found in:\n\n%s", filepath.Base(filePath)), w)
		return
	}

	metaStr := meta.ToString()
	metaFile := pdf.GetMetadataFilename(filePath)

	confirmMsg := fmt.Sprintf("Found %d items in metadata.\n\n%s\nProceed to save metadata?", count, metaStr)

	dialog.ShowConfirm("Proceed", confirmMsg, func(save bool) {
		if save {
			err := pdf.SaveMetadataToFile(meta, metaFile)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			dialog.ShowInformation("Complete", fmt.Sprintf("Metadata saved as:\n\n%s", metaFile), w)
		} else {
			dialog.ShowInformation("Canceled", "Metadata shown but not saved.", w)
		}
	}, w)
}

func showCleanDialog(w fyne.Window) {
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if reader == nil {
			return // User cancelled
		}
		defer reader.Close()

		filePath := reader.URI().Path()
		cleanFile(w, filePath)
	}, w)

	fd.SetFilter(&pdfFilter{})
	fd.Resize(fyne.NewSize(800, 600))
	fd.Show()
}

func cleanFile(w fyne.Window, filePath string) {
	meta, err := pdf.ReadMetadata(filePath)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to read PDF: %w", err), w)
		return
	}

	count := meta.Count()
	cleanPath := pdf.GetCleanFilename(filePath)
	cleanName := filepath.Base(cleanPath)

	var confirmMsg string
	if count == 0 {
		confirmMsg = fmt.Sprintf("No metadata found in this PDF.\n\nProceed to create clean copy?\n\nFile will be saved as:\n%s", cleanName)
	} else {
		metaStr := meta.ToString()
		confirmMsg = fmt.Sprintf("Found %d items in metadata.\n\n%s\nProceed to wipe metadata?\n\nFile will be saved as:\n%s", count, metaStr, cleanName)
	}

	dialog.ShowConfirm("Proceed", confirmMsg, func(wipe bool) {
		if wipe {
			err := pdf.WipeMetadata(filePath, cleanPath)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			dialog.ShowInformation("Complete", fmt.Sprintf("File cleaned and saved as:\n\n%s", cleanName), w)
		} else {
			dialog.ShowInformation("Canceled", "Operation canceled. No clean file created.", w)
		}
	}, w)
}

// pdfFilter implements storage.FileFilter for PDF files
type pdfFilter struct{}

func (f *pdfFilter) Matches(uri fyne.URI) bool {
	ext := uri.Extension()
	return ext == ".pdf" || ext == ".PDF"
}

func (f *pdfFilter) Name() string {
	return "PDF Files (*.pdf)"
}
