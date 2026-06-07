<template>
	<div class="card-text-box">
		<slot>
			[Card text content]
		</slot>
	</div>
</template>

<style lang="scss" scoped>
	.card-text-box {
		display: block;
		text-align: start;
		width: 100%;
	}
</style>
