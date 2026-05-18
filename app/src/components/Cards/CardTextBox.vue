<template>
	<div class="card-text-box">
		<slot>
			[Card text content]
		</slot>
	</div>
</template>

<style lang="scss" scoped>
	.card-text-box {
		display: inline-block;
		width: 100%;
		flex-grow: 1;
	}
</style>
