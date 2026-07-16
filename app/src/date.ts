
export const parseDateString = (value?: string | null): Date | null => {

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

export const fmtDateString = (value?: string | null) => {

	const date = parseDateString(value);
	if (!date) {
		return '-';
	}

	return date.toLocaleDateString('en-UK', {
		year: 'numeric',
		month: 'short',
		day: 'numeric',
	});
};

export const fmtTimeString = (value?: string | null) => {

	const date = parseDateString(value);
	if (!date) {
		return '-';
	}

	return date.toLocaleDateString('en-UK', {
		year: 'numeric',
		month: 'short',
		day: 'numeric',
		hour: 'numeric',
		minute: 'numeric',
		second: 'numeric',
	});
};
