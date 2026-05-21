<template>
	<div class="content-list">
		<slot>
			[Content]
		</slot>
	</div>
</template>

<style lang="scss" scoped>
	.content-list {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}
</style>
