
export const useLanguage = (): string => {
	const langCode = navigator.language;
	return langCode.split('-')[0] || 'en';
};

export const intl = (lang: string, strings: Record<string, string>) => {
	return strings[lang] ?? Object.entries(strings)?.[0]?.[1];
};
