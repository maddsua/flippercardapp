<template>
	<div class="central-message">
		<slot>
			[Message content]
		</slot>
	</div>
</template>

<style lang="scss" scoped>
	.central-message {
		display: flex;
		justify-content: center;
		font-size: 0.8rem;
		color: var(--app-theme-mysterious-white);
	}
</style>
