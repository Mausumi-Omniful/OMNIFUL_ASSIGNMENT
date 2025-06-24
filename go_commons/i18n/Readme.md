# Package i18n

## Overview
The i18n package handles internationalization and localization, enabling the application to support multiple languages. It provides functions to load translation files and retrieve localized strings.

## Key Components
- Translation Functions: Obtain localized messages based on keys.
- Locale Management: Set and retrieve the current locale.
- Resource Loading: Load translations from files or other sources.

## Usage Example
~~~go
package main

import (
	"fmt"
	"github.com/omniful/go_commons/i18n"
)

func main() {
	i18n.LoadTranslations("en_US.json")
	message := i18n.Translate("welcome_message")
	fmt.Println(message)
}
~~~

## Notes
- Easily integrates with web applications for multilingual support.
