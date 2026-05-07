<script setup lang="ts">
const emit = defineEmits<{
	(e: 'done', confirmed: boolean): void;
}>();
</script>

<template>
	<div class="ui-prompt">
		<div class="container">
			<div class="title">
				<slot>
					[Title]
				</slot>
			</div>
			<div class="actions">
				<button type="button" @click="emit('done', false)">Cancel</button>
				<button type="button" @click="emit('done', true)">Proceed</button>
			</div>
		</div>
	</div>
</template>

<style lang="scss" scoped>
	.ui-prompt {
		position: fixed;
		top: 0;
		left: 0;
		z-index: 100;
		background-color: rgba(0, 0, 0, 0.8);
		backdrop-filter: blur(2px);
		width: 100%;
		height: 100%;

		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;

		.container {
			display: flex;
			flex-direction: column;
			gap: 1rem;
			width: 100%;
			max-width: 20rem;
			padding: 1rem;
			border: 1px solid var(--app-theme-powder-trail);
			border-radius: 0.25rem;
		}

		.title {
			font-size: 1.125rem;
		}

		.actions {
			display: flex;
			flex-flow: row nowrap;
			border-radius: 0.25rem;
			overflow: hidden;
			gap: 0.125rem;
			flex-grow: 1;
		}

		button {
			flex-shrink: 0;
			flex-grow: 1;
			display: block;
			text-align: center;
			padding: 0.75rem 1rem;
			font-size: 0.75rem;
			font-weight: 600;
			color: var(--app-theme-midnight);
			background-color: var(--app-theme-mysterious-white);
			border: none;
			outline: none;
			transition: all 200ms ease;

			&:last-of-type {
				color: var(--app-theme-snow-white);
				background-color: var(--app-theme-bishop);
			}

			&:hover {
				cursor: pointer;
			}
		}
	}
</style>
