package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"kvstore/internal/config"
	"kvstore/internal/logger"
)

var GlobalTranslator *Translator

// Translator holds the translations for different languages.
type Translator struct {
	mu           sync.RWMutex
	translations map[string]map[string]string
	defaultLang  string
}

// NewTranslator creates a new Translator instance.
func NewTranslator() (*Translator, error) {
	t := &Translator{
		translations: make(map[string]map[string]string),
		defaultLang:  config.Cfg.DefaultLanguage,
	}

	if err := t.loadTranslations(); err != nil {
		return nil, fmt.Errorf("failed to load translations: %v", err)
	}

	return t, nil
}

// loadTranslations loads translation files from the translations directory.
func (t *Translator) loadTranslations() error {
	dir := config.Cfg.TranslationsDir
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read translations directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() || len(file.Name()) < 6 || file.Name()[len(file.Name())-5:] != ".json" {
			continue
		}

		langCode := file.Name()[:len(file.Name())-5] // Remove ".json"
		filePath := fmt.Sprintf("%s/%s", dir, file.Name())

		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read translation file %s: %v", filePath, err)
		}

		var translations map[string]string
		if err := json.Unmarshal(data, &translations); err != nil {
			return fmt.Errorf("failed to parse translation file %s: %v", filePath, err)
		}

		t.translations[langCode] = translations
	}

	if _, ok := t.translations[t.defaultLang]; !ok {
		return fmt.Errorf("default language %s not found in translations", t.defaultLang)
	}

	return nil
}

// Translate translates a uint32 error code to the current language.
func (t *Translator) Translate(code uint32) string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	keyStr := fmt.Sprintf("0x%08X", code)

	if val, ok := t.translations[t.defaultLang][keyStr]; ok {
		return val
	}

	logger.Log.Warnf("Translation for error code '%s' not found in language %s", keyStr, t.defaultLang)
	return keyStr
}

// TranslateForLang translates an error code for a specific language.
func (t *Translator) TranslateForLang(code uint32, lang string) string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	translations, ok := t.translations[lang]
	if !ok {
		logger.Log.Warnf("Requested language %s not found", lang)
		return fmt.Sprintf("0x%08X", code)
	}

	keyStr := fmt.Sprintf("0x%08X", code)
	if val, ok := translations[keyStr]; ok {
		return val
	}

	logger.Log.Warnf("Translation for error code '%s' not found in language %s", keyStr, lang)
	return keyStr
}

// SetLanguage changes the language for translations.
func (t *Translator) SetLanguage(lang string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.translations[lang]; ok {
		t.defaultLang = lang
	} else {
		logger.Log.Warnf("Language %s not found, using previous default language %s", lang, t.defaultLang)
	}
}

// GetLanguage returns the currently set default language.
func (t *Translator) GetLanguage() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.defaultLang
}

// InitializeTranslator initializes the global translator.
func InitializeTranslator() error {
	translator, err := NewTranslator()
	if err != nil {
		return err
	}
	GlobalTranslator = translator
	return nil
}

// Tr is a shorthand for translating a uint32 error code using the global translator.
func Tr(code uint32) string {
	if GlobalTranslator == nil {
		return fmt.Sprintf("0x%08X", code)
	}
	return GlobalTranslator.Translate(code)
}
