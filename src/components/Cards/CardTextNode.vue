<script setup lang="ts">
//	todo: make stylable
</script>

<template>
	<span class="card-text-node">
		<slot>
			[Text]
		</slot>
	</span>
</template>

<style lang="scss" scoped>
	.card-text-node {
		display: inline;
		font-weight: 500;
		font-size: 7.5cqw;
		line-height: 1.5em;
	}
</style>
