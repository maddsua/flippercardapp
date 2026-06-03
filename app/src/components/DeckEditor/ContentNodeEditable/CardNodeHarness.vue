<script setup lang="ts">
const emit = defineEmits<{
	(e: 'remove'): void;
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

			button.remove {
				display: block;
				width: 1.5rem;
				height: 1.5rem;
				border: none;
				outline: none;
				background-color: var(--app-theme-blood-red);
				mask-type: alpha;
				mask-position: center;
				mask-repeat: no-repeat;
				mask-size: contain;
				mask-image: url(/src/assets/icons/delete-mask.svg);
				opacity: 0.8;
	
				&:hover {
					cursor: pointer;
					opacity: 1;
				}
			}
		}
	}
</style>
