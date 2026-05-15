<template>
	<div class="header">
		<div class="overscript">
			<slot name="overscript">
				[Overscript]
			</slot>
		</div>
		<div class="title">
			<slot name="title">
				[Title]
			</slot>
		</div>
	</div>	
</template>

<style lang="scss" scoped>
	.header {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;

		.overscript {
			text-transform: uppercase;
			color: var(--app-theme-deep-lavender);
			font-size: 0.75rem;
		}

		.title {
			font-size: 300;
			font-size: 1.75rem;
		}
	}
</style>
