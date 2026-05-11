<script setup lang="ts">
import { computed } from 'vue';
import { intl, useLanguage } from '../../intl';

const lang = useLanguage();

const placehoder = computed(() => intl(lang, {
	en: 'Name, topic or a keyword',
	de: 'Name, Thema oder ein Schlüsselwort'
}));

</script>

<template>
	<input type="text" :placeholder="placehoder" />
</template>

<style lang="scss" scoped>
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
</style>
