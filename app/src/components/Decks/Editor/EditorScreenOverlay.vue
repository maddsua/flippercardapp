<template>
	<div class="editor-screen-overlay">
		<slot>
			[Overlay]
		</slot>
	</div>
</template>

<style lang="scss" scoped>
	.editor-screen-overlay {
		position: absolute;
		left: 0;
		top: 0;
		z-index: 200;
		width: 100%;
		height: 100%;
		backdrop-filter: blur(4px);
		background-color: rgba(0, 0, 0, 0.35);
		display: flex;
		flex-direction: column;
		gap: 1rem;
		align-items: center;
		justify-content: center;
	}
</style>
