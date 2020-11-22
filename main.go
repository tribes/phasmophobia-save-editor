package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
	"github.com/tribes/phasmophobia-save-editor/internal/phasmophobia"
	"github.com/tribes/phasmophobia-save-editor/internal/translations"
)

var saveLocations = []string{
	"saveData.txt",
}

func main() {
	// %USERPROFILE% is not recognized directly in paths so we add it here
	home, _ := os.UserHomeDir()
	saveLocations = append(saveLocations, home+`\AppData\LocalLow\Kinetic Games\Phasmophobia\saveData.txt`)

	// Boot app
	a := app.New()
	w := a.NewWindow("")

	// Load translations
	if err := translations.LoadTranslations(); err != nil {
		log.Fatalf("Unable to load translations: %s", err)
	}
	w.SetTitle(translations.Get("WindowTitle"))

	// Load savegame
	var save *phasmophobia.Save
	for _, saveLoc := range saveLocations {
		stat, err := os.Stat(saveLoc)
		if err != nil || stat.IsDir() {
			continue
		}

		save, err = phasmophobia.ReadSave(saveLoc)
		if err != nil {
			//displayError(w, "ErrorOnSaveLoad", err)
			continue
		}

		break
	}

	// Ask user to indicate save location if we couldn't find it
	if save == nil {
		getSaveFromFiler(a, w)
	}

	// Build languages menu
	languagesMenuItems := []*fyne.MenuItem{}
	for _, language := range translations.AvailableLanguages() {
		languagesMenuItems = append(languagesMenuItems, fyne.NewMenuItem(language, func(newLanguage string) func() {
			return func() {
				translations.SetLanguage(newLanguage)
				if save != nil {
					w.SetContent(save.GenerateWidget(w))
				}
			}
		}(language)))
	}
	languagesMenu := fyne.NewMenu(translations.Get("MenuItemLanguages"), languagesMenuItems...)

	// Add some menu interactions
	mainMenu := fyne.NewMainMenu(
		// a quit item will be appended to our first menu
		fyne.NewMenu(translations.Get("MenuItemFile"),
			fyne.NewMenuItem(translations.Get("MenuItemLoad"), func() {
				getSaveFromFiler(a, w)
			}),
		),
		languagesMenu,
	)
	w.SetMainMenu(mainMenu)

	// Generate our main widget
	if save != nil {
		w.SetContent(save.GenerateWidget(w))
	} else {
		w.SetContent(widget.NewButton(translations.Get("ButtonOpenSave"), func() { getSaveFromFiler(a, w) }))
	}

	// Show our main window
	w.SetMaster()
	w.ShowAndRun()
}

func getSaveFromFiler(a fyne.App, w fyne.Window) (save *phasmophobia.Save, err error) {
	w2 := a.NewWindow(translations.Get("BrowserTitle"))
	w.Hide()

	fd := dialog.NewFileOpen(func(r fyne.URIReadCloser, err error) {
		defer w2.Close()
		defer w.Show()
		if err != nil {
			displayError(w, "ErrorOnSaveLoad", err)
			return
		}
		if r == nil {
			// Cancel doesn't return an error
			return
		}
		defer r.Close()

		// Check mime type
		if r.URI().MimeType() != "text/plain" {
			displayError(w, "%s", errors.New(translations.Get("ErrorBadFileFormat")))
			return
		}

		// Check file extension
		if r.URI().Extension() != ".txt" {
			displayError(w, "%s", errors.New(translations.Get("ErrorBadFileExtension")))
			return
		}

		// Try to load the selected file
		save, err = phasmophobia.ReadSave(strings.TrimPrefix(r.URI().String(), "file://"))
		if err != nil {
			displayError(w, "ErrorOnSaveLoad", err)
			return
		}
		w.SetContent(save.GenerateWidget(w))

	}, w2)

	// Full screen dialog
	w2.Resize(fyne.NewSize(1280, 720))
	w2.SetFixedSize(true)
	w2.Show()
	fd.Show()

	// Windows won't resize right at start, so, little hacky dockey ...
	go func(fileDialog *dialog.FileDialog) {
		i := 1
		for i < 10 {
			time.Sleep(time.Millisecond * time.Duration(100))
			fileDialog.Resize(fyne.NewSize(1280, 720))
			i++
		}
		w2.Show()
	}(fd)

	return
}

func displayError(w fyne.Window, label string, err error) {
	dialog.ShowInformation(translations.Get("ErrorWindowTitle"), fmt.Sprintf(translations.Get(label), err), w)
}
