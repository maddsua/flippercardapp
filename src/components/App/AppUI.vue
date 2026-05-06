<script setup lang="ts">
//	todo: add navigation shit
</script>

<template>
	<div class="app-ui">
		<div class="app-content">
			<slot>
				[App content]
			</slot>
		</div>
	</div>
</template>

<style lang="scss" scoped>
	.app-ui {
		position: relative;
		display: flex;
		flex-direction: column;
		align-items: center;
		width: 100%;
		height: 100%;
		height: 100vh;
		height: 100dvh;

		.app-content {
			width: 100%;
			flex-grow: 1;
			max-width: 40rem;
			overflow: hidden auto;
		}
	}
</style>
