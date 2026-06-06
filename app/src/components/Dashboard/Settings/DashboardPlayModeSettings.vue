<script setup lang="ts">
import { reactive, watch } from 'vue';
import { intl, useLanguage } from '@/intl';
import { useStorage } from '@/storage/storage';
import GenericToggle from '@/components/App/Inputs/GenericToggle.vue';

const store = useStorage();
const lang = useLanguage();

const state = reactive({
	showNavigation: store.preferences.playModeShowNavigation.load(),
});

watch(() => state.showNavigation, async (value) => {
	store.preferences.playModeShowNavigation.store(value);
});

</script>

<template>
	<div class="preferences-group">
		<GenericToggle v-model="state.showNavigation">
			{{ intl(lang, {
				en: 'Show navigation buttons in Play Mode',
				de: 'Navigationsschaltflächen im Spielmodus anzeigen',
				uk: 'Показувати кнопки навігації в режимі гри'
			}) }}
		</GenericToggle>
	</div>
</template>

<style lang="scss" scoped>
	.preferences-group {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}
</style>