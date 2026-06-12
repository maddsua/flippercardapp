<script setup lang="ts">
const emit = defineEmits<{
	(e: 'remove'): void;
	(e: 'up'): void;
	(e: 'down'): void;
}>();
</script>

<template>
	<div class="card-node-harness">

		<div class="wrapper">
			<div class="header">
				<div class="title">
					<slot name="title">
						[Title]
					</slot>
				</div>
			</div>
			<div class="content">
				<slot name="content">
					[Content]
				</slot>
			</div>
		</div>

		<div class="actions">
			<button type="button" class="remove" title="Remove node" @click="emit('remove')"></button>
			<button type="button" class="up" title="Remove node" @click="emit('up')"></button>
			<button type="button" class="down" title="Remove node" @click="emit('down')"></button>
		</div>

	</div>
</template>

<style lang="scss" scoped>
	.card-node-harness {
		display: flex;
		flex-flow: row nowrap;
		gap: 0.25rem;
		padding: 0.5rem;

		.wrapper {
			display: flex;
			flex-direction: column;
			border: 2px solid var(--app-theme-powder-trail);
			background-color: var(--app-theme-deep-lavender);
			border-radius: 0.25rem;
		}

		.header {
			border-bottom: 1px solid var(--app-theme-powder-trail);
			padding: 0.125rem 0.25rem;
			
			.title {
				font-size: 0.75rem;
				font-weight: 600;
				color: white;
			}
		}

		.content {
			display: flex;
			flex-direction: column;
			gap: 0.5rem;
			padding: 0.5rem;
			background-color: var(--app-theme-carbon);
		}

		.actions {
			flex-shrink: 0;
			width: fit-content;

			button {
				display: block;
				width: 1.5rem;
				height: 1.5rem;
				border: none;
				outline: none;
				mask-type: alpha;
				mask-position: center;
				mask-repeat: no-repeat;
				mask-size: contain;
				opacity: 0.8;

				&:hover {
					cursor: pointer;
					opacity: 1;
				}

				&.remove {
					background-color: var(--app-theme-blood-red);
					mask-image: url(/src/assets/icons/delete-mask.svg);
				}

				&.up {
					background-color: var(--app-theme-snow-white);
					mask-image: url(/src/assets/icons/arrow-bracket-mask.svg);
					transform: rotate(90deg);
				}

				&.down {
					background-color: var(--app-theme-snow-white);
					mask-image: url(/src/assets/icons/arrow-bracket-mask.svg);
					transform: rotate(-90deg);
				}
			}
		}
	}
</style>
