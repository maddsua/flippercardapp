
export const escapeFileName = (val: string) => val.replace(/[^a-z0-9]/gi, '_').replace(/[^a-z0-9]/gi, '_');

export const downloadBlob = (blob: Blob, name?: string) => {

	const url = window.URL.createObjectURL(blob);
	const link = document.createElement('a');

	link.href = url;

	if (name) {
		link.download = name;
	}

	link.click();

	window.URL.revokeObjectURL(url);
};

export const pickLocalFiles = async (opts?: { multiple?: boolean, accept?: string[] }) => {

	const input = document.createElement('input');
	input.type = 'file';
	input.accept = opts?.accept?.join(',') || '';
	input.multiple = opts?.multiple || false;

	const filePromise = new Promise<FileList | null>((resolve) => {
		input.addEventListener('change', () => resolve(input.files));
		input.addEventListener('cancel', () => resolve(null));
	});

	input.click();
	const files = await filePromise;
	input.remove();

	return files;
};
