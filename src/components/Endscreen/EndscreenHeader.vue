<script setup lang="ts">

</script>

<template>
	<header>
		<div class="prize-icon"></div>
		<div class="headline">
			<h1>
				Deck complete!
			</h1>
			<h2 v-if="$slots.summary">
				<slot name="summary" />
			</h2>
		</div>
	</header>
</template>

<style lang="scss" scoped>
	header {
		display: flex;
		flex-direction: column;
		gap: 3rem;
		align-items: center;

		.prize-icon {
			display: block;
			padding: 1rem;
			background-color: #4B4AF9;
			background: linear-gradient(280deg,#4b4af9 0%, #3576ff 100%);
			border-radius: 100%;
			overflow: hidden;

			&:after {
				content: "";
				display: block;
				width: 4rem;
				height: 4rem;
				background-image: url(/src/assets/icons/trophy-color-gold.svg);
				background-position: center;
				background-repeat: no-repeat;
				background-size: contain;
			}
		}

		.headline {
			display: flex;
			flex-direction: column;
			gap: 2rem;
			align-items: center;
		}

		h1, h2 {
			line-height: 1rem;
			margin: 0;
			padding: 0;
			font-weight: 500;
		}

		h1 {
			font-size: 2.25rem;
			color: var(--app-color-white);
		}

		h2 {
			font-size: 1.75rem;
			color: var(--app-accent-light-blue);
		}
	}
</style>
