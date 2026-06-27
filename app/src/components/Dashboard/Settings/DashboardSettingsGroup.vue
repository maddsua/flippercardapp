<template>
	<div class="preferences-group">
		<slot>
			[Group]
		</slot>
	</div>
</template>

<style lang="scss" scoped>
	.preferences-group {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}
</style>

