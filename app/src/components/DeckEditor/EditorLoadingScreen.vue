<script setup lang="ts">
import LoadingMessage from '../App/LoadingMessage.vue';
</script>

<template>
	<div class="editor-loading-screen">
		<LoadingMessage />
	</div>
</template>

<style lang="scss" scoped>
	.editor-loading-screen {
		position: absolute;
		left: 0;
		top: 0;
		width: 100%;
		height: 100%;
		backdrop-filter: blur(4px);
		background-color: rgba(0, 0, 0, 0.25);
		display: flex;
		align-items: center;
		justify-content: center;
	}
</style>
