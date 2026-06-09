import { useStorage } from "./storage/storage";

const store = useStorage();

export const defaultLang = 'en' as const;

export const useLanguage = (): string => {
	return store.preferences.language.load() || navigator.language.split('-')[0] || defaultLang;
};

export const intl = (lang: string, content: Record<string, string>): string => {
	return content[lang] ?? content[defaultLang] ?? Object.entries(content)?.[0]?.[1];
};
