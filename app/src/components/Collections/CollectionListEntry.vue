<script setup lang="ts">
const props = defineProps<{
	title: string;
	summary?: string;
	starrable?: boolean;
	starred?: boolean;
}>();
</script>

<template>
	<button type="button">
		<div class="header">
			<div class="title">
				{{ title }}
			</div>
			<div v-if="summary" class="summary">
				{{ summary }}
			</div>
		</div>
		<div class="stats">
			<div v-if="starrable || starred" class="icon star" :class="{ starred }"></div>
		</div>
	</button>
</template>

<style lang="scss" scoped>
	button {
		display: flex;
		flex-flow: row nowrap;
		gap: 1rem;
		padding: 1rem 2rem;
		border-radius: 1rem;
		border: none;
		outline: none;
		transition: all 150ms ease;
		text-align: start;
		text-align: unset;

		background-color: var(--app-theme-irish-green);
		color: var(--app-theme-snow-white);

		@media (orientation: portrait) {
			padding: 1rem 1.5rem;
		}

		.header {
			flex-grow: 1;
			display: flex;
			flex-direction: column;
			align-items: start;
			gap: 0.25rem;

			.title {
				font-size: 1rem;
				font-weight: 600;
				white-space: nowrap;
				overflow: hidden;
				text-overflow: ellipsis;
				width: 100%;
			}

			.summary {
				font-size: 0.85rem;
				font-weight: 400;
				white-space: nowrap;
				overflow: hidden;
				text-overflow: ellipsis;
				width: 100%;
			}
		}

		.stats {
			flex-shrink: 0;
			display: flex;
			flex-flow: row nowrap;
			gap: 0.5rem;
			align-items: center;

			.icon {
				display: block;
				width: 1.25rem;
				height: 1.25rem;
				mask-type: alpha;
				mask-size: contain;
				mask-position: center;
				mask-repeat: no-repeat;

				&.star {
					background-color: white;
					mask-image: url(/src/assets/icons/star-mask.svg);

					&.starred {
						mask-image: url(/src/assets/icons/star-filled-mask.svg);
					}
				}
			}
		}

		&:hover {
			cursor: pointer;
			background-color: var(--app-theme-spooky-orange);
		}
	}
</style>
