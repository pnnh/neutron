package models

const LangEn = "en"
const LangZh = "zh"
const LangEs = "es"
const LangFr = "fr"
const LangJa = "ja"
const DefaultLanguage = LangEn

func IsValidLanguage(lang string) bool {
	return lang == LangEn || lang == LangZh || lang == LangEs || lang == LangFr || lang == LangJa
}
