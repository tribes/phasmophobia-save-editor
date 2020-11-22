package phasmophobia

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
	"github.com/tribes/phasmophobia-save-editor/internal/translations"
)

var selectValues = map[string][]string{
	"GhostType": {"Phantom", "Banshee", "Jinn", "Revenant", "Shade", "Oni", "Wraith", "Mare", "Demon", "Yurei", "Poltergeist", "Spirit"},
}
var regexpFloat = regexp.MustCompile(`^\d+(\.\d+)?$`)
var regexpInt = regexp.MustCompile(`^\d+$`)
var maxIntValues = map[string]int{
	"PlayersMoney":                 249999,
	"EMFReaderInventory":           99,
	"FlashlightInventory":          99,
	"CameraInventory":              99,
	"LighterInventory":             99,
	"CandleInventory":              99,
	"UVFlashlightInventory":        99,
	"CrucifixInventory":            99,
	"DSLRCameraInventory":          99,
	"EVPRecorderInventory":         99,
	"SaltInventory":                99,
	"SageInventory":                99,
	"TripodInventory":              99,
	"StrongFlashlightInventory":    99,
	"MotionSensorInventory":        99,
	"SoundSensorInventory":         99,
	"SanityPillsInventory":         99,
	"ThermometerInventory":         99,
	"GhostWritingBookInventory":    99,
	"IRLightSensorInventory":       99,
	"ParabolicMicrophoneInventory": 99,
	"GlowstickInventory":           99,
	"HeadMountedCameraInventory":   99,
}

// GenerateWidget return the associated widget with all the necessary event listeners
func (s Save) GenerateWidget(w fyne.Window) *widget.Form {

	// Build form with properties
	formItems := []*widget.FormItem{}
	// Handle strings
	for _, item := range s.StringData {
		// Handle predefined values
		if values, ok := selectValues[item.Key]; ok {
			entry := widget.NewSelectEntry(values)
			entry.Text = item.Value
			entry.OnChanged = func(item *StringSaveEntry) func(s string) { return func(s string) { item.Value = s } }(item)
			formItems = append(formItems, widget.NewFormItem(translations.Get(item.Key), entry))
			continue
		}

		// Handle basic string
		entry := widget.NewEntry()
		entry.Validator = func(a string) error { return nil }
		entry.Text = item.Value
		entry.OnChanged = func(item *StringSaveEntry) func(s string) { return func(s string) { item.Value = s } }(item)
		formItems = append(formItems, widget.NewFormItem(translations.Get(item.Key), entry))
	}

	// Handle floats
	for _, item := range s.FloatData {
		entry := widget.NewEntry()
		entry.Validator = func(a string) error {
			if regexpFloat.MatchString(a) {
				return nil
			}
			return errors.New("invalid entry")
		}
		entry.Text = fmt.Sprintf("%f", item.Value)
		entry.OnChanged = func(item *FloatSaveEntry) func(s string) {
			return func(s string) {
				i, err := strconv.ParseFloat(s, 64)
				if err != nil {
					return
				}
				item.Value = i
			}
		}(item)
		formItems = append(formItems, widget.NewFormItem(translations.Get(item.Key), entry))
	}

	// Handle integers
	for _, item := range s.IntData {
		entry := widget.NewEntry()
		entry.Validator = func(s string) error {
			if !regexpInt.MatchString(s) {
				return errors.New("invalid entry")
			}

			if max, ok := maxIntValues[s]; ok {
				if i, _ := strconv.Atoi(s); i > max {
					return errors.New("value is too big")
				}
			}

			return nil
		}
		entry.Text = strconv.FormatInt(item.Value, 10)
		entry.OnChanged = func(item *IntSaveEntry) func(s string) {
			return func(s string) {
				i, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return
				}
				item.Value = i
			}
		}(item)
		formItems = append(formItems, widget.NewFormItem(translations.Get(item.Key), entry))
	}

	// Handle booleans
	for _, item := range s.BoolData {
		check := widget.NewCheck("", func(newValue bool) { item.Value = newValue })
		check.SetChecked(item.Value)
		formItems = append(formItems, widget.NewFormItem(translations.Get(item.Key), check))
	}

	form := widget.NewForm(formItems...)
	form.SubmitText = translations.Get("SubmitText")
	form.OnSubmit = func() {
		// Apply changes and save configuration
		err := s.Save()
		if err != nil {
			dialog.ShowInformation(translations.Get("ErrorWindowTitle"), fmt.Sprintf(translations.Get("ErrorOnSaveFile"), err), w)
		}

		return
	}

	return form
}
