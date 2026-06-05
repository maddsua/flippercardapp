<template>
	<div class="error-message">
		<slot>
			[Message]
		</slot>
	</div>
</template>

<style lang="scss" scoped>
	.error-message {
		color: var(--app-theme-blood-red);
		font-weight: 400;
		font-size: 0.85rem;
	}
</style>
