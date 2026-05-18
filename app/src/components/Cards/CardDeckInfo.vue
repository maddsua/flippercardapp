<script setup lang="ts">

const props = defineProps<{
	labels: string[]
	size: number;
	index: number;
	isMarked?: boolean;
}>();

const emit = defineEmits<{
	(e: 'toggleMarked'): void;
}>();

</script>

<template>
	<div class="deck-info">

		<div class="progress">
			<div v-for="idx of size" class="marker" :class="{ filled: idx <= index + 1 }"></div>
		</div>

		<div class="details-row">

			<div class="summary">
				<template v-for="(item, idx) of labels">
					<template v-if="idx > 0">
						<hr />
					</template>
					<span>
						{{ item }}
					</span>
				</template>
			</div>

			<div class="actions">
				<button type="button" class="save" :class="{ active: isMarked }" @click="emit('toggleMarked')">
					<template v-if="isMarked">
						Unmark
					</template>
					<template v-else>
						Mark
					</template>
				</button>
			</div>

		</div>

	</div>
</template>

<style lang="scss" scoped>
	.deck-info {
		position: absolute;
		left: 0.5rem;
		top: 0.5rem;
		right: 0.5rem;
		z-index: 10;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;

		.progress {
			display: flex;
			flex-flow: row nowrap;
			gap: 0.2rem;

			.marker {
				flex-grow: 1;
				height: 2px;
				background-color: rgba(255, 255, 255, 0.35);

				&.filled {
					background-color: white;
				}
			}
		}

		.details-row {
			display: flex;
			flex-flow: row nowrap;
			gap: 1rem;
			padding: 0 0.5rem;
			align-items: center;
			justify-content: space-between;
		}

		.summary {
			display: flex;
			flex-flow: row nowrap;
			align-items: center;
			gap: 0.5rem;
			min-width: 0;
			overflow: hidden;

			span {
				display: block;
				font-size: 0.85rem;
				white-space: nowrap;
			}

			hr {
				display: block;
				width: 3px;
				height: 3px;
				border-radius: 100%;
				background-color: white;
				border: none;
				outline: none;
				margin: 0;
				padding: 0;
				flex-shrink: 0;
			}
		}

		.actions {
			display: flex;
			flex-flow: row nowrap;
			align-items: center;
			justify-content: end;
			flex-shrink: 0;

			button {
				display: flex;
				flex-flow: row nowrap;
				align-items: center;
				gap: 0.25rem;
				color: var(--app-theme-snow-white);
				font-weight: 600;
				border: none;
				outline: none;
				border-radius: 0.25rem;
				padding: 0.25rem 0.5rem;
				font-size: 0.75rem;

				&:hover {
					cursor: pointer;
				}

				&.save {
					background-color: var(--app-theme-irish-green);

					&::before {
						content: "";
						display: block;
						width: 0.85rem;
						height: 0.85rem;
						mask-type: alpha;
						mask-size: contain;
						mask-position: center;
						mask-repeat: no-repeat;
						mask-image: url(/src/assets/icons/star-mask.svg);
						background-color: var(--app-theme-snow-white);
					}

					&.active {
						background-color: var(--app-theme-spooky-orange);

						&::before {
							mask-image: url(/src/assets/icons/star-filled-mask.svg);
						}
					}
				}
			}
		}
	}
</style>
