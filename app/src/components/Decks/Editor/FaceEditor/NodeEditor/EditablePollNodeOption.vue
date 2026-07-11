<script setup lang="ts">

const props = defineProps<{
	checked?: boolean;
}>();

const model = defineModel<string>();
const emit = defineEmits<{
	(e: 'check'): void;
	(e: 'remove'): void;
}>();

</script>

<template>
	<div class="poll-option-editor">
		<button type="button" class="chech-marker" :class="{ checked }" @click="emit('check')"></button>
		<input type="text" placeholder="Option value" v-model="model" />
		<button type="button" class="remove" @click="emit('remove')"></button>
	</div>
</template>

<style lang="scss" scoped>

	.poll-option-editor {
		display: flex;
		flex-flow: row nowrap;
		gap: 0.5rem;
		align-items: center;
		padding: 0.25rem;
		background-color: var(--app-theme-ghostly-glow);
		border-radius: 0.25rem;

		input {
			display: block;
			padding: 0.25rem;
			background-color: unset;
			border: none;
			outline: none;
			color: var(--app-theme-snow-white);
			font-weight: 600;
			font-size: 0.85rem;
			flex-grow: 1;

			&::placeholder {
				color: var(--app-theme-mysterious-white);
				font-weight: 400;
			}
		}

		button {
			display: block;
			width: 1rem;
			height: 1rem;
			border: none;
			outline: none;
			background: none;
			mask-type: alpha;
			mask-position: center;
			mask-repeat: no-repeat;
			mask-size: contain;
			opacity: 0.8;
			flex-grow: 0;

			&:hover {
				cursor: pointer;
				opacity: 1;
			}

			&.remove {
				background-color: var(--app-theme-blood-red);
				mask-image: url(/src/assets/icons/delete-mask.svg);
			}

			&.chech-marker {
				border: 2px solid white;
				border-radius: 100%;

				&.checked {
					border-color: transparent;
					background-color: var(--app-theme-sporty-yellow);
					opacity: 1;
					cursor: default;
					mask-type: alpha;
					mask-image: url(/src/assets/icons/star-filled-mask.svg);
					mask-size: contain;
					mask-repeat: no-repeat;
					mask-position: center;
				}
			}
		}
	}
</style>
