import { useStorage } from "./storage/storage";

const store = useStorage();

export const useLanguage = (): string => {
	return store.preferences.language.load() || navigator.language.split('-')[0] || 'en';
};

export const intl = (lang: string, strings: Record<string, string>) => {
	return strings[lang] ?? Object.entries(strings)?.[0]?.[1];
};
