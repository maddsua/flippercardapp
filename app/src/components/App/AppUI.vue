<script setup lang="ts">
import AppNavigation from './AppNavigation.vue';
</script>

<template>
	<div class="app-ui">
		<div class="app-viewport">
			<div class="app-container">
				<slot>
					[App content]
				</slot>
			</div>
			<AppNavigation />
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

		.app-viewport {
			display: flex;
			flex-direction: column;
			height: 100%;
			width: 100%;

			@media (orientation: landscape) {
				max-width: 40rem;
			}
		}

		.app-container {
			display: flex;
			flex-direction: column;
			gap: 2rem;
			position: relative;
			min-height: 0;
			overflow: hidden auto;
			scrollbar-width: thin;
			padding: 2rem;
			flex-grow: 1;
		}
	}
</style>
