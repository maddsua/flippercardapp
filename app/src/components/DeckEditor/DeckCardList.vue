<script setup lang="ts">
import { ref, watch } from 'vue';
import type { CardContentFace } from '../../content';
import CardFace from '../Cards/CardFace.vue';

const props = defineProps<{
	list: CardContentFace[];
	pointer: number;
}>();

const emit = defineEmits<{
	(e: 'select', idx: number): void;
	(e: 'remove', idx: number): void;
	(e: 'duplicate', idx: number): void;
	(e: 'add'): void;
}>();

const scrollableRef = ref<HTMLElement | null>(null);

watch(() => props.list.length, (length, oldLength) => {
	if (length > oldLength && scrollableRef.value) {
		setTimeout(() => scrollableRef.value!.scrollTop = scrollableRef.value!.scrollHeight, 100);
	}
});

</script>

<template>
	<div class="card-list" ref="scrollableRef">
		<div v-for="(card, idx) of list" class="card-item-tile" :class="{ selected: idx === pointer }" @click="emit('select', idx)">
			<div class="controls-layer">
				<div class="label">
						{{ idx + 1 }}
				</div>
				<div class="col index">
					<div class="index">
						{{ idx + 1 }}
					</div>
				</div>
				<div class="col controls">
					<button type="button" class="remove" title="Remove card" @click.self.stop="emit('remove', idx)"></button>
					<button type="button" class="duplicate" title="Remove card" @click.self.stop="emit('duplicate', idx)"></button>
				</div>
			</div>
			<div class="preview-canvas">
				<CardFace :entry="card" />
			</div>
		</div>
		<button type="button" class="add" title="Add card" @click="emit('add')">+Add card</button>
	</div>
</template>

<style lang="scss" scoped>
	.card-list {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		user-select: none;
		overflow: hidden auto;
		scrollbar-width: thin;
		padding-right: 0.5rem;
		scroll-behavior: smooth;

		.card-item-tile {
			position: relative;
			display: flex;
			flex-flow: column;
			justify-content: center;
			align-items: center;
			height: 8rem;
			width: 5rem;
			border-radius: 0.25rem;
			border: 2px solid transparent;
			flex-shrink: 0;

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

				.index {
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

				.label {
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

				&:hover {
					background-color: rgba(0, 0, 0, 0.4);
					backdrop-filter: blur(2px);

					.col.controls, .label {
						opacity: 1;
					}

					.col.index {
						opacity: 0;
					}
				}
			}

			&:hover {
				cursor: pointer;
				border-color: var(--app-theme-deep-lavender);
			}

			&.selected {
				border-color: var(--app-theme-sky-blue);
				background-color: var(--app-theme-sky-blue);
			}
		}

		button.add {
			display: block;
			text-align: center;
			padding: 0.5rem;
			font-size: 0.65rem;
			font-weight: 600;
			background-color: var(--app-theme-sky-blue);
			color: var(--app-theme-snow-white);
			border: none;
			outline: none;
			border-radius: 0.25rem;

			&:hover {
				cursor: pointer;
				background-color: var(--app-theme-deep-lavender);
			}
		}
	}
</style>
