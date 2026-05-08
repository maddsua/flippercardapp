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
		font-size: 0.85rem;
	}
</style>
