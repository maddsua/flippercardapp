<script setup lang="ts">
import GenericDropdown from '@/components/App/Inputs/GenericDropdown.vue';
import { defaultLang, useLanguage } from '@/intl';
import { useStorage } from '@/storage/storage';
import { onMounted, ref, watch } from 'vue';

const store = useStorage();
const lang = ref(useLanguage());

const languageOptions = [
	{
		value: defaultLang,
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

onMounted(() => {

	if (!languageOptions.some(opt => opt.value === lang.value)) {
		lang.value = defaultLang;
	}

	watch(() => lang.value, async (value) => {

		if (value === defaultLang) {
			store.preferences.language.clear();
		} else {
			store.preferences.language.store(value);
		}

		window.location.reload();
	});

});

</script>

<template>
	<GenericDropdown :options="languageOptions" v-model="lang" />
</template>
