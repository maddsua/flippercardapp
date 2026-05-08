<template>
	<span class="skeleton">
		<slot />
	</span>
</template>

<style lang="scss" scoped>
span.skeleton {
	display: block;
	position: relative;
	color: transparent;
	overflow: hidden;
	flex-grow: 1;
	border-radius: 1rem;
	background-color: rgba(128, 128, 128, 0.25);

	&::before {
		content: '';
		z-index: 10;
		position: absolute;
		left: 0;
		top: 0;
		right: 0;
		bottom: 0;
		background-color: var(--app-island-focus-extra);
	}

	&::after {
		content: '';
		position: absolute;
		top: 0;
		right: 0;
		bottom: 0;
		left: 0;
		z-index: 11;
		transform: translateX(-100%);
		background-image: linear-gradient(90deg, rgba(#fff, 0) 0, rgba(#fff, 0.2) 20%, rgba(#fff, 0.5) 60%, rgba(#fff, 0));
		animation: shimmer 2.55s infinite;
	}

	@keyframes shimmer {
		100% {
			transform: translateX(100%);
		}
	}
}
</style>
