import { nextTick, reactive } from "vue";
import { isInteractive } from "./dom";

export interface VirtualKeyboardState {
	isOpen: boolean;
}

export const detectVirtualKeyboard = (): VirtualKeyboardState => {

	let viewBox = getViewbox();

	const state = reactive<VirtualKeyboardState>({ isOpen: false });

	window.addEventListener('resize', () => {

		const resized = getViewbox();

		const interactiveFocused = isInteractive(document.activeElement);
		const detected = !!interactiveFocused && (resized.y / viewBox.y) < 0.75;

		if (detected !== state.isOpen) {

			if (detected) {
				console.debug('Virtual keyboard is open');
			} else {
				console.debug('Virtual keyboard is closed');
			}

			nextTick(() => state.isOpen = detected);
		}

		if (!interactiveFocused && (resized.y > viewBox.y || resized.x !== viewBox.x)) {
			viewBox = resized;
		}

	});

	return state;
};

const getViewbox = () => ({ x: window.innerWidth, y: window.innerHeight });
