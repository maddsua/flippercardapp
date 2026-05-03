<script setup lang="ts">
const props = defineProps<{
	icon?: 'home' | 'retry';
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
		border: 1px solid transparent;
		border-radius: 0.75rem;
		background-color: var(--app-color-blueish-dark);
		outline: none;
		transition: all 200ms ease;
		color: var(--app-theme-snow-white);

		&.filled {
			background-color: var(--app-accent-blue);
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
		}

		&:hover {
			cursor: pointer;
			background-color: var(--app-accent-light-blue);
		}
	}
</style>
