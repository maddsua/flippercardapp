<template>
	<div class="input-group">
		<slot>
			[Inputs]
		</slot>
	</div>
</template>

<style lang="scss" scoped>
	.input-group {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
		width: 100%;
	}
</style>
