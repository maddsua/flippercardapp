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
	}
</style>
