<template>
	<div class="collection-action">
		<slot>
			[Actions]
		</slot>
	</div>
</template>

<style lang="scss" scoped>

	.collection-action {
		display: flex;
		flex-flow: row nowrap;
		justify-content: center;
		gap: 1rem;
	}

</style>
