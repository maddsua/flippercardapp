<script setup lang="ts">

const props = defineProps<{
	type?: 'text' | 'password';
	placeholder?: string;
	variant?: 'normal' | 'borderless';
	multiline?: boolean;
	disabled?: boolean;
}>();

const model = defineModel<string>();

</script>

<template>
	<textarea v-if="multiline" :class="{ [`variant-${variant}`]: !!variant }" :placeholder="placeholder" v-model="model" :disabled="disabled"></textarea>
	<input v-else :type="type || 'text'" :class="{ [`variant-${variant}`]: !!variant }" :placeholder="placeholder" v-model="model" :disabled="disabled" />
</template>

<style lang="scss" scoped>

	input, textarea {
		display: block;
		width: 100%;
		font-size: 0.85rem;
		padding: 0.5rem 0.75rem;
		border-radius: 0.5rem;
		outline: none;
		border: 1px solid transparent;
		background-color: var(--app-theme-ghostly-glow);
		color: var(--app-theme-snow-white);

		&.variant-borderless {
			background-color: unset;
			border-radius: unset;
			border: none;
			border-bottom: 2px solid var(--app-theme-ghostly-glow);
		}

		&:focus {
			border-color: var(--app-theme-sky-blue);
		}

		&:active {
			border-color: var(--app-theme-deep-lavender);
		}

		&::placeholder {
			color: var(--app-theme-snow-white);
			opacity: 0.5;
		}

		&:disabled {
			pointer-events: none;
			filter: saturate(0);
			opacity: 0.75;
		}
	}

	textarea {
		resize: none;
		height: 5em;
	}

</style>
