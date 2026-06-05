<script setup lang="ts">
import GenericDropdown from '@/components/App/Inputs/GenericDropdown.vue';
import { useLanguage } from '@/intl';
import { useStorage } from '@/storage/storage';
import { ref, watch } from 'vue';

const store = useStorage();
const lang = ref(useLanguage());

const languageOptions = [
	{
		value: 'en',
		label: 'English'
	},
	{
		value: 'de',
		label: 'Deutsch'
	},
	{
		value: 'uk',
		label: 'Українська'
	},
];

watch(() => lang.value, async (value) => {
	store.preferences.language.store(value);
	window.location.reload();
});

</script>

<template>
	<GenericDropdown :options="languageOptions" v-model="lang" />
</template>
