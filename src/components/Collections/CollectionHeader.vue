<template>
	<header>
		<h1>
			<slot name="title">
				[Title]
			</slot>
		</h1>
		<p v-if="$slots.summary">
			<slot name="summary" />
		</p>
	</header>
</template>

<style lang="scss" scoped>
	header {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;

		h1 {
			font-size: 2rem;
			padding: 0;
			margin: 0;
			color: var(--app-theme-snow-white);
			font-weight: 400;
		}

		p {
			font-size: 0.9rem;
			padding: 0;
			margin: 0;
			color: var(--app-theme-mysterious-white);
			font-weight: 300;
		}
	}
</style>
