<template>
	<div class="endscreen-actions">
		<slot>
			[Actions]
		</slot>
	</div>
</template>

<style lang="scss" scoped>
	.endscreen-actions {
		display: flex;
		flex-flow: row nowrap;
		gap: 1rem;
	}
</style>
