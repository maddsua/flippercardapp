<template>
	<div class="input-row">
		<slot>
			[Inputs]
		</slot>
	</div>
</template>

<style lang="scss" scoped>
	.input-row {
		display: flex;
		flex-flow: row nowrap;
		gap: 0.75rem;
		width: 100%;
		align-items: center;
	}
</style>
