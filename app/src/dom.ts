
export const isContainerOutsideClick = (container: HTMLElement | null, target: EventTarget | null): boolean => {

	//	skip non-dom generated events
	if (!target || !(target instanceof HTMLElement)) {
		return false;
	}

	// skip if the container is not yet ready
	if (!container || !document.contains(container)) {
		return false;
	}

	//	skipt if target was already removed, like buttons in a pop-up
	if (!document.contains(target)) {
		return false;
	}

	return !container.contains(target);
};

export const isInteractive = (element: Element | null) => {
	return !!element && (element instanceof HTMLInputElement || element instanceof HTMLTextAreaElement);
};
