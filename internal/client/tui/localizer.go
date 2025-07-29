// Package tui provides a Terminal User Interface (TUI) for the GophKeeper client.
// It implements an interactive command-line interface using the tview library
// for user authentication, data management, and server communication.
//
// The package provides a complete user interface with forms, menus, and modals
// for managing encrypted data through a secure gRPC connection to the server.
package tui

import (
	"os"
	"path/filepath"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// Localizer manages translations for the TUI interface.
// It provides a centralized way to handle internationalization
// and localization of user interface text.
type Localizer struct {
	bundle    *i18n.Bundle
	localizer *i18n.Localizer
}

// NewLocalizer creates a new Localizer instance.
// It loads translation files from the locales directory and sets up
// the localization bundle with support for Russian and English.
func NewLocalizer() (*Localizer, error) {
	bundle := i18n.NewBundle(language.Russian)

	// Load translation files
	localesDir := "locales"
	files, err := os.ReadDir(localesDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			_, err := bundle.LoadMessageFile(filepath.Join(localesDir, file.Name()))
			if err != nil {
				return nil, err
			}
		}
	}

	// Determine language from environment or default to Russian
	lang := os.Getenv("LANG")
	if lang == "" {
		lang = "ru"
	}

	// Map language codes to supported languages
	switch lang[:2] {
	case "en":
		lang = "en"
	default:
		lang = "ru"
	}

	localizer := i18n.NewLocalizer(bundle, lang)

	return &Localizer{
		bundle:    bundle,
		localizer: localizer,
	}, nil
}

// Translate returns the translated string for the given key.
// It provides a simple interface for getting localized text.
func (l *Localizer) Translate(key string) string {
	if l.localizer == nil {
		return key
	}
	msg, err := l.localizer.Localize(&i18n.LocalizeConfig{
		MessageID: key,
	})
	if err != nil {
		// Return the key if translation is not found
		return key
	}
	return msg
}

// TranslateWithParams returns the translated string with parameters.
// It supports parameter substitution in translation strings.
func (l *Localizer) TranslateWithParams(key string, params map[string]interface{}) string {
	if l.localizer == nil {
		return key
	}
	msg, err := l.localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: params,
	})
	if err != nil {
		// Return the key if translation is not found
		return key
	}
	return msg
}
