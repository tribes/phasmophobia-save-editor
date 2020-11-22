package translations

import (
	"encoding/json"
	"io/ioutil"

	"github.com/gobuffalo/packr"
)

type translationItem struct {
	Language    string
	Dictionnary map[string]string `json:"translations"`
}

var translationsPath string = "../../assets/translations"
var translations = []translationItem{}
var currentLanguage = "English"

// LoadTranslations load the translations from the disk
func LoadTranslations() error {
	box := packr.NewBox(translationsPath)
	for _, translationFilePath := range box.List() {
		translationFile, err := box.Open(translationFilePath)
		if err != nil {
			return err
		}
		defer translationFile.Close()

		byteValue, _ := ioutil.ReadAll(translationFile)

		var translation translationItem
		err = json.Unmarshal(byteValue, &translation)
		if err != nil {
			return err
		}
		translations = append(translations, translation)
	}

	return nil
}

// AvailableLanguages lists all available languages
func AvailableLanguages() []string {
	languages := []string{}
	for _, trad := range translations {
		languages = append(languages, trad.Language)
	}
	return languages
}

// SetLanguage changes the current language used for translations
func SetLanguage(language string) {
	currentLanguage = language
}

// Get return a translated version of the given key, if missing, it will return the key itself
func Get(key string) string {
	return getFromLanguage(currentLanguage, key)
}

func getFromLanguage(language, key string) string {
	for _, translation := range translations {
		if translation.Language != language {
			continue
		}

		if value, ok := translation.Dictionnary[key]; ok {
			return value
		}
	}
	return key
}
