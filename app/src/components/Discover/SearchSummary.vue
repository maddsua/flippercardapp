<template>
	<div class="searh-summary">
		<div class="summary">
			<slot>
				[Summary]
			</slot>
		</div>
	</div>
</template>

<style lang="scss" scoped>
	.searh-summary {
		display: flex;
		justify-content: center;

		.summary {
			font-size: 0.85rem;
			color: var(--app-theme-mysterious-white);
		}
	}
</style>
