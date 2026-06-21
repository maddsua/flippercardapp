/// <reference types="vite/client" />

interface ImportMeta {
    readonly env: ImportMetaEnv
}

interface ImportMetaEnv {
    readonly VITE_APP_VERSION?: string
    readonly VITE_APP_BUILD_TS?: string
    readonly VITE_APP_PLATFORM?: string;
    readonly VITE_APP_SOURCE_VCS?: string;
    readonly VITE_APP_SOURCE_REPO?: string;
}
