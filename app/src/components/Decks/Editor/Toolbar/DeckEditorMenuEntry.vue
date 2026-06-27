<script setup lang="ts">

const props = defineProps<{
	label: string;
	icon?: 'play' | 'publish' | 'delete' | 'exit' | 'history' | 'cross' | 'broom' | 'info' | 'io' | 'undo' | 'redo' | 'keyboard' | 'add-title' | 'add-text' | 'add-image' | 'add-poll';
	disabled?: boolean;
}>();

</script>

<template>
	<button class="menu-entry" type="button" :title="label" :disabled="disabled">
		<span v-if="icon" class="entry-icon" :class="icon"></span>
		<span>
			{{ label }}
		</span>
	</button>
</template>

<style lang="scss" scoped>

	button.menu-entry {
		display: flex;
		flex-flow: row nowrap;
		gap: 0.25rem;
		align-items: center;
		background: unset;
		border: none;
		outline: none;
		text-align: start;
		padding: 0.35rem 0.5rem;
		font-size: 0.65rem;

		&:hover {
			background-color: var(--app-theme-powder-trail);
			cursor: pointer;
		}

		&:disabled {
			pointer-events: none;
			opacity: 0.6;
			filter: saturate(0);
		}

		span {
			display: block;
			flex-shrink: 0;
		}

		.entry-icon {
			display: block;
			width: 0.8rem;
			height: 0.8rem;
			background-color: var(--app-theme-mysterious-white);
			mask-type: alpha;
			mask-size: contain;
			mask-repeat: no-repeat;
			mask-position: center;

			&.play {
				mask-image: url(/src/assets/icons/play-mask.svg);
				background-color: var(--app-theme-irish-green);
			}

			&.publish {
				mask-image: url(/src/assets/icons/send-mask.svg);
				background-color: var(--app-theme-irish-green);
			}

			&.delete {
				mask-image: url(/src/assets/icons/delete-mask.svg);
				background-color: var(--app-theme-blood-red);
			}

			&.exit {
				mask-image: url(/src/assets/icons/exit-mask.svg);
			}

			&.history {
				mask-image: url(/src/assets/icons/history-mask.svg);
			}

			&.cross {
				mask-image: url(/src/assets/icons/cross-cut-mask.svg);
			}

			&.broom {
				mask-image: url(/src/assets/icons/broom-mask.svg);
			}

			&.info {
				mask-image: url(/src/assets/icons/info-pixel-mask.svg);
			}

			&.io {
				mask-image: url(/src/assets/icons/import-export-mask.svg);
			}

			&.undo {
				mask-image: url(/src/assets/icons/undo-mask.svg);
			}

			&.redo {
				mask-image: url(/src/assets/icons/undo-mask.svg);
				transform: rotate(180deg);
			}

			&.keyboard {
				mask-image: url(/src/assets/icons/keyboard-mask.svg);
			}

			&.add-title {
				mask-image: url(/src/assets/icons/add-title-color.svg);
			}

			&.add-text {
				mask-image: url(/src/assets/icons/add-text-color.svg);
			}

			&.add-image {
				mask-image: url(/src/assets/icons/add-image-color.svg);
			}

			&.add-poll {
				mask-image: url(/src/assets/icons/add-poll-color.svg);
			}
		}
	}

</style>
