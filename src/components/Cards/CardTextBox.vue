<script setup lang="ts">

</script>

<template>
	<div class="card-text-box">
		<slot>
			[Card text content]
		</slot>
	</div>
</template>

<style lang="scss" scoped>

	//	todo: add typography

	.card-text-box {
		display: inline-block;
		font-size: 1.5em;
	}

</style>
