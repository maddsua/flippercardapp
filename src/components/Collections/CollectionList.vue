<template>
	<div class="collection-list">
		<slot>
			[Content]
		</slot>
	</div>
</template>

<style lang="scss" scoped>
	.collection-list {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		flex-grow: 1;
	}
</style>
