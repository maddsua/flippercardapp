<script setup lang="ts">
import { computed } from 'vue';
import { intl, useLanguage } from '@/intl';

const model = defineModel<string>();

const lang = useLanguage();

const placehoder = computed(() => intl(lang, {
	en: 'Name, topic or a keyword',
	de: 'Name, Thema oder ein Schlüsselwort'
}));

const captureInput = (event: InputEvent) => {
	const { value } = event.target as HTMLInputElement;
	model.value = value;
};

const resetInput = () => {
	model.value = undefined;
};

</script>

<template>
	<div class="searchbar">
		<input type="text" :placeholder="placehoder" :value="model" @input="captureInput" />
		<button v-if="model?.length" type="reset" @click="resetInput"></button>
	</div>
</template>

<style lang="scss" scoped>

	.searchbar {
		display: block;
		width: 100%;
		position: relative;

		input {
			display: block;
			width: 100%;
			font-size: 1.125rem;
			padding: 0.5rem 1rem;
			border-radius: 1rem;
			outline: none;
			border: 1px solid transparent;
			background-color: var(--app-theme-ghostly-glow);
			color: var(--app-theme-snow-white);
	
			&:focus {
				border-color: var(--app-theme-sky-blue);
			}
	
			&:active {
				border-color: var(--app-theme-deep-lavender);
			}
	
			&::placeholder {
				color: var(--app-theme-mysterious-white);
			}
		}

		button[type=reset] {
			position: absolute;
			right: 0.5rem;
			top: 0.4rem;
			display: block;
			width: 1.75rem;
			height: 1.75rem;
			z-index: 1;
			border: none;
			outline: none;
			background: unset;
			background-repeat: no-repeat;
			background-size: contain;
			background-image: url(/src/assets/icons/cross-cut-mask.svg);
			opacity: 0.7;

			&:hover {
				cursor: pointer;
				opacity: 1;
			}
		}
	}

</style>
