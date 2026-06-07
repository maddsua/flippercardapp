<script setup lang="ts">
import type { CardContentFace } from '@/content';
import CardFace from '@/components/Cards/CardFace.vue';

const props = defineProps<{
	card: CardContentFace;
	active?: boolean;
	showControls?: boolean;
	label?: string | number;
}>();

const emit = defineEmits<{
	(e: 'remove'): void;
	(e: 'duplicate'): void;
}>();

</script>

<template>
	<div class="card-thumbnail" :class="{ outline: active }">

		<div class="controls-layer" :class="{ active: active }">

			<div class="label-small">
				{{ label ?? '?' }}
			</div>

			<div class="col label-large">
				{{ label ?? '?' }}
			</div>

			<div v-if="showControls" class="col controls">
				<button type="button" class="remove" title="Remove card" @click.self.stop="emit('remove')"></button>
				<button type="button" class="duplicate" title="Remove card" @click.self.stop="emit('duplicate')"></button>
			</div>

		</div>

		<div class="preview-canvas">
			<CardFace :entry="card" />
		</div>

	</div>
</template>

<style lang="scss" scoped>
	.card-thumbnail {
		position: relative;
		height: 8rem;
		width: 5rem;
		border-radius: 0.25rem;
		border: 2px solid transparent;
		flex-shrink: 0;
		overflow: hidden;

		.preview-canvas {
			position: relative;
			width: 100%;
			height: 100%;
			container-type: inline-size;
			z-index: 0;
			pointer-events: none;
		}

		.controls-layer {
			position: absolute;
			width: 100%;
			height: 100%;
			top: 0;
			left: 0;
			display: flex;
			flex-flow: row nowrap;
			justify-content: space-between;
			padding: 0.25rem;
			z-index: 1;

			.col {
				display: flex;
				flex-direction: column;
				justify-content: space-between;
				gap: 0.5rem;
				z-index: 1;
				transition: opacity 150ms ease;

				&.controls {
					opacity: 0;
				}
			}

			.label-large {
				display: flex;
				width: 1.125rem;
				height: 1.125rem;
				align-items: center;
				justify-content: center;
				overflow: visible;
				background-color: var(--app-theme-snow-white);
				color: var(--app-theme-carbon);
				font-weight: 600;
				font-size: 0.75rem;
				border-radius: 100%;
			}

			button {
				display: block;
				width: 1.25rem;
				height: 1.25rem;
				border: none;
				outline: none;
				background-color: var(--app-theme-snow-white);
				mask-type: alpha;
				mask-position: center;
				mask-repeat: no-repeat;
				mask-size: contain;
				transition: all 150ms ease;
				z-index: 2;

				&:hover {
					cursor: pointer;
				}

				&.duplicate {
					mask-image: url(/src/assets/icons/copy-mask.svg);

					&:hover {
						background-color: var(--app-theme-irish-green);
					}
				}

				&.remove {
					mask-image: url(/src/assets/icons/delete-mask.svg);

					&:hover {
						background-color: var(--app-theme-blood-red);
					}
				}
			}

			.label-small {
				position: absolute;
				width: 100%;
				height: 100%;
				top: 0;
				left: 0;
				display: flex;
				align-items: center;
				justify-content: center;
				font-size: 2.5rem;
				color: var(--app-theme-snow-white);
				z-index: 0;
				transition: opacity 150ms ease;
				opacity: 0;
			}

			&.active, &:hover {
				background-color: rgba(0, 0, 0, 0.4);
				backdrop-filter: blur(2px);

				.label-small {
					opacity: 1;
					color: var(--app-theme-snow-white) !important;
				}

				.col.label-large {
					opacity: 0;
				}
			}

			&:hover {

				.label-small {
					color: var(--app-theme-sporty-yellow);
				}

				.col.controls {
					opacity: 1;
				}
			}
		}

		&:hover {
			cursor: pointer;
			border-color: var(--app-theme-sporty-yellow);
			background-color: var(--app-theme-sporty-yellow);
		}

		&.outline {
			border-color: var(--app-theme-sky-blue);
			background-color: var(--app-theme-sky-blue);
		}
	}
</style>
