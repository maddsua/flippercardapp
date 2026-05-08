<template>
	<div class="collection-container">
		<slot>
			[Collection content]
		</slot>
	</div>
</template>

<style lang="scss" scoped>
	.collection-container {
		display: flex;
		flex-direction: column;
		gap: 2rem;
		position: relative;
		width: 100%;
		height: 100%;
		min-height: 0;
		overflow: hidden auto;
		scrollbar-width: thin;
		padding: 2rem;
	}
</style>
