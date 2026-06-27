import { ref } from "vue";

export interface AppInfo {
	readonly version: string | null;
	readonly buildTime: Date | null;
	readonly mode: 'pwa' | 'web' | 'web-limited' | null;
	readonly platform: string | null;
	readonly source: AppSource;
};

export interface AppSource {
	vcs?: string;
	repo?: string;
};

declare let window: Window & {
	appInfo?: AppInfo;
};

const parseDateString = (value?: string | null) => {

	if (!value) {
		return null;
	}

	try {
		const date = new Date(value);
		return date.toISOString().length ? date : null;
	} catch (_) {
		return null;
	}
};

const pwaDisplayModes = ['standalone', 'minimal-ui', 'fullscreen'] as const;

const detectAppMode = () => {

	for (const query of pwaDisplayModes) {
		try {
			if (window.matchMedia(`(display-mode: ${query})`).matches) {
				return 'pwa';
			}
		} catch (_) {
			return 'web-limited';
		}
	}

	return 'web';
};

export const getAppInfo = (): AppInfo => ({
	version: import.meta.env.VITE_APP_VERSION || null,
	buildTime: parseDateString(import.meta.env.VITE_APP_BUILD_TS),
	mode: detectAppMode(),
	platform: import.meta.env.VITE_APP_PLATFORM || null,
	source: {
		vcs: import.meta.env.VITE_APP_SOURCE_VCS,
		repo: import.meta.env.VITE_APP_SOURCE_REPO
	},
});

export const pwaInstallPrompt = ref<BeforeInstallPromptEvent | null>(null);

export const enablePwaInstall = () => {

	window.addEventListener('beforeinstallprompt', (event) => {
		event.preventDefault();
		pwaInstallPrompt.value = event as BeforeInstallPromptEvent;
	});

	window.addEventListener("appinstalled", () => pwaInstallPrompt.value = null);
};
