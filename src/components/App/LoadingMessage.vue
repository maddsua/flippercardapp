<template>
	<div class="loading-message">
		<span>
			<slot>Loading...</slot>
		</span>
	</div>
</template>

<style lang="scss" scoped>
	.loading-message {
		display: flex;
		flex-flow: row nowrap;
		gap: 1rem;
		align-items: center;
		justify-content: center;

		span {
			font-size: 0.85rem;
		}

		&::before {
			content: "";
			width: 1.25rem;
			height: 1.25rem;
			background-color: white;
			mask-type: alpha;
			mask-size: contain;
			mask-repeat: no-repeat;
			mask-position: center;
			mask-image: url(/src/assets/animated/ring-resize-white.svg);
		}
	}
</style>
