<template>
	<div class="menu-section">
		<div class="title">
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
	.menu-section {
		display: flex;
		flex-direction: column;
		gap: 1rem;

		.title {
			font-size: 0.85rem;
			font-weight: 600;
		}

		.content {
			display: flex;
			flex-direction: column;
			gap: 0.75rem;

			@media (orientation: landscape) {
				padding-left: 1rem;
			}
		}
	}
</style>
