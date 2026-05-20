<script setup lang="ts">

const props = defineProps<{
	size: number;
	activeIdx: number;
}>();

const emit = defineEmits<{
	(e: 'select', idx: number): void;
	(e: 'remove', idx: number): void;
	(e: 'duplicate', idx: number): void;
	(e: 'add'): void;
}>();

</script>

<template>
	<div class="card-list">
		<div v-for="idx in size" class="card-item-tile" :class="{ selected: idx-1 === activeIdx }" @click="emit('select', idx-1)">
			<button type="button" class="duplicate" title="Remove card" @click.self.stop="emit('duplicate', idx-1)"></button>
			<div class="label">
				{{ idx }}
			</div>
			<button type="button" class="remove" title="Remove card" @click.self.stop="emit('remove', idx-1)"></button>
		</div>
		<button type="button" class="add" title="Add card" @click="emit('add')">+Add card</button>
	</div>
</template>

<style lang="scss" scoped>
	.card-list {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		user-select: none;
		overflow: hidden auto;
		scrollbar-width: thin;
		padding-right: 0.5rem;

		.card-item-tile {
			position: relative;
			display: flex;
			flex-flow: column;
			justify-content: center;
			align-items: center;
			height: 8rem;
			width: 5rem;
			background-color: var(--app-theme-mysterious-white);
			border-radius: 0.25rem;
			border: 2px solid transparent;
			flex-shrink: 0;

			.label {
				font-size: 2.5rem;
				color: var(--app-theme-carbon);
			}

			&:hover {
				cursor: pointer;
				border-color: var(--app-theme-deep-lavender);
			}

			&.selected {
				border-color: var(--app-theme-sky-blue);
			}

			button {
				display: block;
				position: absolute;
				width: 1.25rem;
				height: 1.25rem;
				border: none;
				outline: none;
				background-color: var(--app-theme-midnight);
				mask-type: alpha;
				mask-position: center;
				mask-repeat: no-repeat;
				mask-size: contain;
				transition: all 150ms ease;

				&:hover {
					cursor: pointer;
				}

				&.duplicate {
					left: 0.25rem;
					top: 0.25rem;
					mask-image: url(/src/assets/icons/copy-mask.svg);

					&:hover {
						background-color: var(--app-theme-deep-lavender);
					}
				}

				&.remove {
					right: 0.25rem;
					top: 0.25rem;
					mask-image: url(/src/assets/icons/delete-mask.svg);
		
					&:hover {
						background-color: var(--app-theme-blood-red);
					}
				}
			}
		}

		button.add {
			display: block;
			text-align: center;
			padding: 0.5rem;
			font-size: 0.65rem;
			font-weight: 600;
			background-color: var(--app-theme-sky-blue);
			color: var(--app-theme-snow-white);
			border: none;
			outline: none;
			border-radius: 0.25rem;

			&:hover {
				cursor: pointer;
				background-color: var(--app-theme-deep-lavender);
			}
		}
	}
</style>
