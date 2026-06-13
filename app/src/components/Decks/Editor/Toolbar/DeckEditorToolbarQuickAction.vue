<script setup lang="ts">
const props = defineProps<{
	label: string;
	icon: 'play' | 'publish' | 'delete';
	disabled?: boolean;
}>();
</script>

<template>
	<button type="button" :title="label" :class="{ [`icon-${icon}`]: !!icon }" :disabled="disabled"></button>
</template>

<style lang="scss" scoped>
	button {
		display: block;
		width: 1.5rem;
		height: 1.5rem;
		background-color: var(--app-theme-mysterious-white);
		outline: none;
		border: none;
		mask-type: alpha;
		mask-size: contain;
		mask-repeat: no-repeat;
		mask-position: center;
		opacity: 0.8;
		transition: all 150ms ease;

		&.icon-play {
			mask-image: url(/src/assets/icons/play-mask.svg);

			&:hover {
				background-color: var(--app-theme-irish-green);
			}
		}

		&.icon-publish {
			mask-image: url(/src/assets/icons/send-mask.svg);

			&:hover {
				background-color: var(--app-theme-deep-lavender);
			}
		}

		&.icon-delete {
			mask-image: url(/src/assets/icons/delete-mask.svg);

			&:hover {
				background-color: var(--app-theme-blood-red);
			}
		}

		&:hover {
			cursor: pointer;
			opacity: 1;
		}

		&:disabled {
			pointer-events: none;
			opacity: 0.3;
			filter: saturate(0);
		}
	}
</style>
