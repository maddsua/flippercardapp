<template>
	<div class="fullscreen-message">
		<slot>
			[Message content]
		</slot>
	</div>
</template>

<style lang="scss" scoped>
	.fullscreen-message {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 2rem;
		width: 100%;
		height: 100%;

		font-size: 0.85rem;
	}
</style>
