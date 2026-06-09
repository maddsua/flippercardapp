<script setup lang="ts">
import { reactive, watch } from 'vue';
import { intl, useLanguage } from '@/intl';
import { useStorage } from '@/storage/storage';
import GenericToggle from '@/components/App/Inputs/GenericToggle.vue';

const store = useStorage();
const lang = useLanguage();

const state = reactive({
	showNavigation: store.preferences.playMode.showNavigation.load(),
	disableCardRotation: store.preferences.playMode.disableCardRotation.load(),
});

watch(() => state.showNavigation, value => store.preferences.playMode.showNavigation.store(value));
watch(() => state.disableCardRotation, value => store.preferences.playMode.disableCardRotation.store(value));

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

		<GenericToggle v-model="state.disableCardRotation">
			{{ intl(lang, {
				en: 'Disable card rotation',
				de: 'Kartendrehung deaktivieren',
				uk: 'Не нахиляти картки'
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