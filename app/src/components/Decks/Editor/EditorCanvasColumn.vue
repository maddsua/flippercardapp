<template>
	<div class="canvas-column">
		<div class="header">
			<slot name="title">
				[Title]
			</slot>
		</div>
		<div class="content">
			<slot name="content">
				[Content]
			</slot>
		</div>
	</div>
</template>

<style lang="scss" scoped>
	.canvas-column {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		min-height: 0;
		height: 100%;

		.header {
			display: flex;
			align-items: center;
			justify-content: center;
			font-size: 0.75rem;
			font-weight: 600;
		}

		.content {
			display: flex;
			flex-direction: column;
			min-height: 0;
			height: 100%;
			overflow: hidden auto;
			scrollbar-width: thin;
		}
	}
</style>
