<template>
	<div class="error-message">
		<div v-if="$slots.title" class="message-title">
			<slot name="title" />
		</div>
		<div class="message-details">
			<slot>
				[Message]
			</slot>
		</div>
	</div>
</template>

<style lang="scss" scoped>
	.error-message {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		color: var(--app-theme-snow-white);
		padding: 0.5rem 1rem;
		background-color: var(--app-theme-bishop);
		border-radius: 0.5rem;
		border: 1px solid var(--app-theme-powder-trail);

		.message-title {
			font-size: 0.85rem;
			font-weight: 600;
		}

		.message-details {
			font-size: 0.75rem;
		}
	}
</style>
