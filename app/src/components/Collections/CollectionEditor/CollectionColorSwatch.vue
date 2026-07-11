<script setup lang="ts">

const model = defineModel<string | null>();

const colorOptions: string[] = [
	'var(--app-theme-malachite)',
	'var(--app-theme-tangerine)',
	'var(--app-theme-electric-violet)',
	'var(--app-theme-fuchsia-rose)',
	'var(--app-theme-cobalt)',
	'var(--app-theme-spooky-orange)',
];

</script>

<template>
	<div class="color-swatch">

		<button type="button" class="color null"
			:class="{ selected: !model }"
			@click="model = null"></button>

		<button v-for="color of colorOptions" type="button" class="color"
			:style="{ backgroundColor: color }"
			:class="{ selected: model === color }"
			@click="model = color"></button>
	</div>
</template>

<style lang="scss" scoped>
	.color-swatch {
		display: flex;
		flex-flow: row nowrap;
		gap: 0.25rem;
		align-items: center;

		button {
			position: relative;
			display: block;
			background: unset;
			border: 2px solid transparent;
			outline: hidden;
			width: 2rem;
			height: 2rem;
			border-radius: 0.5rem;

			&.null {
				border: 1px solid var(--app-theme-mysterious-white);

				&:hover, &.selected {
					border-width: 2px;
				}

				&::after {
					content: "";
					display: block;
					position: absolute;
					top: 0;
					left: 0;
					width: 100%;
					height: 100%;
					background-size: contain;
					background-repeat: no-repeat;
					background-position: center;
					background-image: url(/src/assets/icons/cross-cut-mask.svg);
					opacity: 0.5;
				}
			}

			&:hover {
				cursor: pointer;
				border-color: var(--app-theme-mysterious-white);
			}

			&.selected {
				cursor: pointer;
				border-color: var(--app-theme-snow-white);
			}
		}
	}
</style>
