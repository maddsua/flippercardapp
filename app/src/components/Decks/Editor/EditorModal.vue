<script setup lang="ts">

import EditorScreenOverlay from './EditorScreenOverlay.vue';

const props = defineProps<{
	title: string;
	variant?: 'wide' | 'narrow' | 'compact';
}>();

const emit = defineEmits<{
	(e: 'close'): void;
}>();
</script>

<template>
	<EditorScreenOverlay>

		<div class="modal-container" @click.self="emit('close')">

			<div class="modal-window" :class="{ [`variant-${variant}`]: !!variant }">

				<div class="modal-status-bar">
					<div class="modal-title">
						{{ title }}
					</div>
					<div class="modal-actions">
						<button type="button" class="close" @click="emit('close')" title="Close modal"></button>
					</div>
				</div>
	
				<div class="modal-content">
					<slot>[Content]</slot>
				</div>
	
			</div>

		</div>

	</EditorScreenOverlay>
</template>

<style lang="scss" scoped>
	.modal-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		width: 100%;
		height: 100%;
		padding: 2rem 4rem;

		.modal-window {
			display: flex;
			flex-flow: column;
			gap: 2rem;
			width: 100%;
			height: 100%;
			min-height: 0;
			background-color: var(--app-theme-midnight);
			padding: 1rem;
			border-radius: 0.5rem;

			&.variant-narrow, &.variant-compact {
				max-width: 30rem;
			}

			&.variant-compact {
				height: unset;
			}

			.modal-status-bar {
				display: flex;
				flex-flow: row nowrap;
				align-items: center;
				gap: 2rem;

				.modal-title {
					font-size: 1.25rem;
					font-weight: 300;
					min-width: 0;
					flex-grow: 1;
				}

				.modal-actions {
					display: flex;
					flex-flow: row nowrap;
					align-items: center;
					gap: 1rem;
					flex-shrink: 0;

					.close {
						display: block;
						border: none;
						outline: none;
						background-color: white;
						width: 2rem;
						height: 2rem;
						mask-type: alpha;
						mask-position: center;
						mask-repeat: no-repeat;
						mask-size: contain;
						mask-image: url(/src/assets/icons/cross-cut-mask.svg);

						&:hover {
							cursor: pointer;
						}
					}
				}
			}

			.modal-content {
				display: block;
				height: 100%;
				min-height: 0;
			}
		}
	}
</style>
