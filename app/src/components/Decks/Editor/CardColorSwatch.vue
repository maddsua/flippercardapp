<script setup lang="ts">

const props = defineProps<{
	label: string;
	size?: 'normal' | 'small';
}>();

const model = defineModel<string | null>();

const colorOptions: string[] = [
	'var(--app-theme-snow-white)',
	'var(--app-theme-carbon)',
	'var(--app-theme-midnight)',
	'var(--app-theme-bishop)',
	'var(--app-theme-sapphire)',
	'var(--app-theme-rich-mint)',
	'var(--app-theme-sporty-yellow)',
	'var(--app-theme-spooky-orange)',
];

</script>

<template>
	<div class="color-swatch" :class="{ [`size-${size}`]:!! size }">

		<div class="header">
			<div class="label">
				{{ label }}
			</div>
			<button type="reset" @click="model = null">Reset</button>
		</div>

		<div class="options">
			<button v-for="color of colorOptions" type="button" class="color"
				:style="{ backgroundColor: color }"
				:class="{ selected: model === color }"
				@click="model = color"></button>
		</div>
	</div>
</template>

<style lang="scss" scoped>
	.color-swatch {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;

		.header {
			display: flex;
			flex-flow: row nowrap;
			gap: 1rem;
			align-items: center;
			justify-content: space-between;

			.label {
				color: var(--app-theme-mysterious-white);
				font-size: 0.75rem;
				font-weight: 300;
			}

			button[type=reset] {
				display: block;
				outline: none;
				border: 1px solid var(--app-theme-powder-trail);
				border-radius: 0.25rem;
				padding: 0.125rem 0.25rem;
				background: none;

				&:hover {
					cursor: pointer;
					background-color: var(--app-theme-powder-trail);
				}
			}
		}

		.options {
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

		&.size-small {
			.options {
				button {
					width: 1.75rem;
					height: 1.75rem;
				}
			}
		}
	}
</style>
