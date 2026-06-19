<script setup lang="ts">
import type { CardContentFace } from '@/content';
import CardFace from '@/components/Cards/CardFace.vue';

const props = defineProps<{
	face: CardContentFace;
	size?: 'small' | 'normal';
	active?: boolean;
	interactive?: boolean;
	controls?: boolean;
	label?: string | number;
	hidden?: boolean;
}>();

const emit = defineEmits<{
	(e: 'remove'): void;
	(e: 'duplicate'): void;
	(e: 'moveUp'): void;
	(e: 'moveDown'): void;
}>();

</script>

<template>
	<div class="card-thumbnail" :class="{ outline: active, interactive, [`size-${size}`]: !!size, hidden }">

		<div v-if="interactive" class="controls-layer" :class="{ active: active }">

			<div v-if="controls" class="col controls">
				<button type="button" class="move-up" title="Move card up" @click.self.stop="emit('moveUp')"></button>
				<button type="button" class="move-down" title="Move card down" @click.self.stop="emit('moveDown')"></button>
			</div>

			<template v-if="label">

				<div class="label-large">
					{{ label }}
				</div>

				<div class="col label-small">
					{{ label }}
				</div>

			</template>

			<div v-if="controls" class="col controls">
				<button type="button" class="remove" title="Remove card" @click.self.stop="emit('remove')"></button>
				<button type="button" class="duplicate" title="Remove card" @click.self.stop="emit('duplicate')"></button>
			</div>

		</div>

		<div class="preview-canvas">
			<CardFace :entry="face" />
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
		transition: all 150ms ease;

		&.size-small {
			height: 5rem;
			width: 3.125rem;
		}

		&.hidden {
			opacity: 0;
		}

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
					position: relative;
					opacity: 0;
					z-index: 2;
				}
			}

			.label-small {
				position: absolute;
				top: 0.25rem;
				left: 0.25rem;
				display: flex;
				width: 1.125rem;
				height: 1.125rem;
				align-items: center;
				justify-content: center;
				overflow: visible;
				background-color: var(--app-theme-kinda-white);
				color: var(--app-theme-carbon);
				font-weight: 600;
				font-size: 0.65rem;
				border-radius: 100%;
				z-index: 1;
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

				&.move-up {
					mask-image: url(/src/assets/icons/arrow-bracket-mask.svg);
					transform: rotate(90deg);
					mask-size: 1.5rem;

					&:hover {
						background-color: var(--app-theme-sky-blue);
					}
				}

				&.move-down {
					mask-image: url(/src/assets/icons/arrow-bracket-mask.svg);
					transform: rotate(-90deg);
					mask-size: 1.5rem;

					&:hover {
						background-color: var(--app-theme-sky-blue);
					}
				}
			}

			.label-large {
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

				.label-large {
					opacity: 1;
					color: var(--app-theme-snow-white) !important;
				}

				.col.label-small {
					opacity: 0;
				}
			}

			&:hover {

				.label-large {
					color: var(--app-theme-sporty-yellow);
				}

				.col.controls {
					opacity: 1;
				}
			}
		}

		&.interactive:hover {
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
