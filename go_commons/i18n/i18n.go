package i18n

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/omniful/go_commons/constants"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

// LanguageTag represents a base language and region combination e.g. "en", "tr", "pt-BR"
type LanguageTag = language.Tag

// List of language tags which will be sent by the clients
var (
	English LanguageTag = language.English
	Arabic  LanguageTag = language.Arabic
	Turkish LanguageTag = language.Turkish
	Urdu    LanguageTag = language.Urdu
)

type contextLangCodeType string

var contextLangCodeKey contextLangCodeType = "langCode"

type langValue struct {
	Message string `json:"message"`
}

type langContext struct {
	Context context.Context
}

type langConfig struct {
	rootPath      string
	mergeExisting bool
}

type option func(*langConfig)

var (
	langKeys    = make(map[string]map[string]langValue)
	initialized bool
	initOnce    sync.Once
)

// WithRootPath initializes localization directory for i18n package
func WithRootPath(rootPath string) option {
	return func(config *langConfig) {
		config.rootPath = rootPath
	}
}

// WithMergeExisting merges new translations with existing ones instead of replacing them
func WithMergeExisting(merge bool) option {
	return func(config *langConfig) {
		config.mergeExisting = merge
	}
}

// Initialize function initializes all the key value pairs for all the languages.
// Does not initialize any key-value pairs if the localization folder is not correctly set up.
// Returns error if any individual lang file is not a valid json file or does not follow the lang file creation standard.
// It initializes lang key-value pairs in memory for correct files.
func Initialize(opts ...option) error {
	var initErr error

	// Use sync.Once to ensure initialization happens only once
	initOnce.Do(func() {
		initErr = initializeInternal(opts...)
	})

	return initErr
}

// initializeInternal does the actual initialization work
func initializeInternal(opts ...option) error {
	config := setupConfig(opts...)

	// Initialize or prepare for merging
	if !config.mergeExisting || !initialized {
		langKeys = make(map[string]map[string]langValue)
	}

	// First, load service-specific localization files
	serviceErrStrings := loadServiceLocalizations(config)

	// Next, load common localization files as fallback
	loadCommonLocalizations()

	// Return any errors from loading service localizations
	if len(serviceErrStrings) > 0 {
		return errors.New(strings.Join(serviceErrStrings[:], "\n"))
	}

	initialized = true
	return nil
}

// loadServiceLocalizations loads localization files from the service's specified path
func loadServiceLocalizations(config *langConfig) []string {
	serviceFiles, err := findLocalizationFiles(config.rootPath)
	if err != nil {
		return []string{err.Error()}
	}

	return processLocalizationFiles(serviceFiles, config, false)
}

// loadCommonLocalizations loads localization files from the go_commons package
func loadCommonLocalizations() {
	// Find the path to the go_commons package using runtime.Caller
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		return
	}

	// We're currently in the i18n package, so go up one level to get go_commons root
	goCommonsRoot := filepath.Dir(filepath.Dir(thisFile))
	commonsLocalizationPath := filepath.Join(goCommonsRoot, "localization")

	// Check if the commons localization directory exists
	if _, err := os.Stat(commonsLocalizationPath); err != nil {
		return
	}

	// Process common localization files with merge=true to ensure fallback behavior
	commonFiles, err := findLocalizationFiles(commonsLocalizationPath)
	if err != nil || len(commonFiles) == 0 {
		return
	}

	// Create a config that always merges for common files
	commonConfig := &langConfig{
		rootPath:      commonsLocalizationPath,
		mergeExisting: true,
	}

	// Process the common files with isCommon=true for proper precedence
	_ = processLocalizationFiles(commonFiles, commonConfig, true)
}

// setupConfig creates and returns a configuration with provided options
func setupConfig(opts ...option) *langConfig {
	config := &langConfig{
		rootPath:      "./localization",
		mergeExisting: true, // Default to merging
	}

	for _, opt := range opts {
		opt(config)
	}

	return config
}

// findLocalizationFiles walks the directory tree to find localization files
func findLocalizationFiles(rootPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})

	return files, err
}

// processLocalizationFiles processes each localization file and loads translations
func processLocalizationFiles(files []string, config *langConfig, isCommon bool) []string {
	var errstrings []string
	regex := constants.LocalizationRegex

	for _, file := range files {
		matchedStrings := regex.FindStringSubmatch(file)

		if len(matchedStrings) < 2 {
			errstrings = append(errstrings, string("File with name "+file+" not correct"))
			continue
		}

		langCode := matchedStrings[1]

		// Validate language code
		if err := validateLanguageCode(langCode, file); err != nil {
			errstrings = append(errstrings, err.Error())
		}

		// Load and parse translations
		translations, errs := loadTranslations(file, langCode)
		if len(errs) > 0 {
			errstrings = append(errstrings, errs...)
			continue
		}

		// Merge translations
		mergeTranslations(langCode, translations, isCommon)
	}

	return errstrings
}

// validateLanguageCode validates that a language code is properly formatted
func validateLanguageCode(langCode, file string) error {
	_, err := language.Parse(langCode)
	if err != nil {
		return errors.New("Non standard lang code: " + langCode + " for file: " + file)
	}
	return nil
}

// loadTranslations loads and parses translations from a file
func loadTranslations(file, langCode string) (map[string]langValue, []string) {
	var errstrings []string
	var currentLangKeys map[string]langValue

	data, err := os.ReadFile(file)
	if err != nil {
		errstrings = append(errstrings, file+": "+err.Error())
		return nil, errstrings
	}

	err = json.Unmarshal(data, &currentLangKeys)
	if err != nil {
		errstrings = append(errstrings, file+": "+err.Error())
		return nil, errstrings
	}

	return currentLangKeys, errstrings
}

// mergeTranslations merges translations into the global map
// The isCommon parameter indicates if these are common translations (which have lower precedence)
func mergeTranslations(langCode string, translations map[string]langValue, isCommon bool) {
	// Create language map if it doesn't exist
	if _, exists := langKeys[langCode]; !exists {
		langKeys[langCode] = make(map[string]langValue)
	}

	// Merge the translations
	for key, value := range translations {
		if isCommon {
			// For common translations, only add if the key doesn't already exist
			// This ensures service-specific translations take precedence
			if _, exists := langKeys[langCode][key]; !exists {
				langKeys[langCode][key] = value
			}
		} else {
			// For service-specific translations, always override
			langKeys[langCode][key] = value
		}
	}
}

// Parse function takes lang code as string and returns LanguageTag
// returns error if not a valid lang code
func Parse(langCode string) (LanguageTag, error) {
	return language.Parse(langCode)
}

// ContextWithLanguage puts the lang code into the current context.
func ContextWithLanguage(ctx context.Context, lang LanguageTag) context.Context {
	switch ctx.(type) {
	case *gin.Context:
		gctx := ctx.(*gin.Context)
		gctx.Set(string(contextLangCodeKey), lang)
		return gctx
	default:
		return context.WithValue(ctx, contextLangCodeKey, lang)
	}
}

// WithContext initializes context from the package to read lang code from it.
func WithContext(ctx context.Context) langContext {
	return langContext{Context: ctx}
}

// WithLanguage overrides the lang in the context with specific lang
func WithLanguage(ctx context.Context, lang LanguageTag) langContext {
	context := context.WithValue(ctx, contextLangCodeKey, lang)
	return langContext{Context: context}
}

// LangCodeFromContext returns the language code in the provided context
func LangCodeFromContext(ctx context.Context) string {
	var lang interface{}

	switch ctx.(type) {
	case *gin.Context:
		gctx := ctx.(*gin.Context)
		lang, _ = gctx.Get(string(contextLangCodeKey))
	default:
		lang = ctx.Value(contextLangCodeKey)
	}

	if lang == nil {
		return English.String()
	}

	return lang.(LanguageTag).String()
}

// getLangKeyValue gets a translation for a key in a specific language
// Falls back to English if not found in the requested language
// Returns empty string if no translation found in any language
func getLangKeyValue(key string, lang string) string {
	// Ensure we're initialized
	if !initialized {
		Initialize()
	}

	// Try the requested language first
	if langMap, ok := langKeys[lang]; ok {
		if val, ok := langMap[key]; ok && val.Message != "" {
			return val.Message
		}
	}

	// Fall back to English if not English already
	if lang != English.String() {
		if langMap, ok := langKeys[English.String()]; ok {
			if val, ok := langMap[key]; ok && val.Message != "" {
				return val.Message
			}
		}
	}

	// Return empty string if no translation found
	return ""
}

// Translate returns value of a lang key based on the language in the context
func (lc langContext) Translate(key string) string {
	langCode := LangCodeFromContext(lc.Context)
	translatedString := getLangKeyValue(key, langCode)
	return translatedString
}

// Translate is a shorthand for i18n.WithContext(ctx).Translate(key)
// If no translation is found, returns the key itself
func Translate(ctx context.Context, key string) string {
	if key == "" {
		return ""
	}
	translation := WithContext(ctx).Translate(key)
	if translation == "" {
		return key
	}
	return translation
}
