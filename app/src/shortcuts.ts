
interface ShortcutAction {
	ctrl?: boolean;
	shift?: boolean;
	key: string;
	prepreq?: () => boolean;
	action: () => void;
};

export class Shortcuts {

	private readonly entries: ShortcutAction[];

	private withCtrl = false;
	private withShift = false;

	constructor (actions?: ShortcutAction[]) {
		this.entries = actions?.map(item => ({ ...item, key: item.key.toLowerCase() })) || [];
	}

	private onKeyDown = (event: KeyboardEvent) => {

		const key = event.key.toLowerCase();

		switch (key) {
			case 'shift':
				this.withShift = true;
				return;
			case 'control':
				this.withCtrl = true;
				return;
		}

		const shortcut = this.entries.find(entry =>
			entry.key === key &&
			!!entry.ctrl === this.withCtrl &&
			!!entry.shift === this.withShift &&
			(!entry.prepreq || entry.prepreq()));

		if (!shortcut) {
			return;
		}

		event.preventDefault();
		event.stopImmediatePropagation();
		event.stopPropagation();
		shortcut.action();
	};

	private onKeyUp = (event: KeyboardEvent) => {

		const key = event.key.toLowerCase();

		switch (key) {
			case 'shift':
				this.withShift = false;
				return;
			case 'control':
				this.withCtrl = false;
				return;
		}
	};

	register = () => {
		window.addEventListener('keydown', this.onKeyDown);
		window.addEventListener('keyup', this.onKeyUp);
	};

	unregister = () => {
		window.removeEventListener('keydown', this.onKeyDown);
		window.removeEventListener('keyup', this.onKeyUp);
	};
};
