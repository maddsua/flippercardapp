<template>
	<div class="endscreen-stats">
		<slot>
			[Stats row]
		</slot>
	</div>
</template>

<style lang="scss" scoped>

	@use '@/media.scss';

	.endscreen-stats {
		display: grid;
		grid-template-columns: 1fr 1fr 1fr;
		gap: 1rem;

		@include media.phone {
			gap: 0.5rem;
		}
	}
</style>
