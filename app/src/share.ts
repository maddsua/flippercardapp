import { getAppInfo } from "./app";

export const appCanShareData = () => getAppInfo().mode === 'PWA' && 'share' in navigator;

export const appShareData = async (data?: ShareData | null) => {

	if (!data?.url) {
		return;
	}

	if (await shareDataProperly(data)) {
		return;
	}

	await shareViaClipboard(data);
};

const shareDataProperly = async (data: ShareData) => {

	if (!('share' in navigator)) {
		return false;
	}

	return await navigator.share(data).then(() => true).catch(() => false);
};

const shareViaClipboard = async (data: ShareData) => {

	if (!('clipboard' in navigator) || !data.url) {
		return false;
	}

	return await navigator.clipboard.writeText(data.url).then(() => true).catch(() => false);
};
