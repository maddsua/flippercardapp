/// <reference types="vite/client" />

import type { ApiClient } from "./api";

declare global {
    interface Window {
        appAPIClient?: ApiClient;
        appUserDB?: IDBDatabase;
    }
}
