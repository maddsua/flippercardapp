<script setup lang="ts">
const props = defineProps<{
	icon?: 'home' | 'retry' | 'finish';
	filled?: boolean;
}>();
</script>

<template>
	<button type="button" :class="{ filled }">
		<div v-if="icon" class="icon" :class="[icon]"></div>
		<span>
			<slot>
				[Button label]
			</slot>
		</span>
	</button>
</template>

<style lang="scss" scoped>

	@use '../../media.scss';

	button {
		display: flex;
		flex-flow: row nowrap;
		gap: 0.5rem;
		align-items: center;
		justify-content: center;
		padding: 1rem;
		font-size: 1rem;
		font-weight: 600;
		flex-grow: 1;
		border-radius: 0.75rem;
		color: var(--app-theme-snow-white);
		background-color: var(--app-theme-midnight-glow);
		border: none;
		outline: none;
		transition: all 200ms ease;

		&.filled {
			background-color: var(--app-theme-sky-blue);
		}

		.icon {
			display: block;
			width: 1.25rem;
			height: 1.25rem;
			background-size: contain;
			background-repeat: no-repeat;
			background-position: center;
			flex-shrink: 0;

			&.home {
				background-image: url(/src/assets/icons/home-mask.svg);
			}

			&.retry {
				background-image: url(/src/assets/icons/retry-mask.svg);
			}

			&.finish {
				background-image: url(/src/assets/icons/finish-mask.svg);
			}
		}

		@include media.non-sticky-hover {
			background-color: var(--app-theme-deep-lavender);
		}
	}
</style>
