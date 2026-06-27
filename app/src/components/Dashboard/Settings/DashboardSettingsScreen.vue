<script setup lang="ts">
import AppUiHeader from '@/components/App/Layout/AppUiHeader.vue';
import DashboardMenuSection from '../DashboardMenuSection.vue';
import DashboardAccountSettings from './DashboardAccountSettings.vue';
import DashboardAppInfo from './DashboardAppInfo.vue';
import { useLanguage, intl, defaultLang } from '@/intl';
import { onMounted, reactive, watch } from 'vue';
import { useStorage } from '@/storage/storage';
import { appCanShareData } from '@/share';
import GenericDropdown from '@/components/App/Inputs/GenericDropdown.vue';
import DashboardSettingsGroup from './DashboardSettingsGroup.vue';
import GenericToggle from '@/components/App/Inputs/GenericToggle.vue';

const store = useStorage();

const state = reactive({
	lang: useLanguage(),
	showNavigation: store.preferences.playMode.showNavigation.load(),
	disableCardRotation: store.preferences.playMode.disableCardRotation.load(),
	shareLinksOnly: store.preferences.sharing.linkOnly.load(),
});

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

	if (!languageOptions.some(opt => opt.value === state.lang)) {
		state.lang = defaultLang;
	}

	watch(() => state.lang, async (value) => {

		if (value === defaultLang) {
			store.preferences.language.clear();
		} else {
			store.preferences.language.store(value);
		}

		window.location.reload();
	});

	watch(() => state.showNavigation, value => store.preferences.playMode.showNavigation.store(value));
	watch(() => state.disableCardRotation, value => store.preferences.playMode.disableCardRotation.store(value));
	watch(() => state.shareLinksOnly, value => store.preferences.sharing.linkOnly.store(value));

});

</script>

<template>

	<AppUiHeader backHref="/">
		<template v-slot:title>
			{{ intl(state.lang, {
				en: 'Dashboard',
				de: 'Einstellungen',
				uk: 'Налаштування'
			}) }}
		</template>
	</AppUiHeader>

	<DashboardMenuSection>

		<template v-slot:title>
			{{ intl(state.lang, {
				en: 'App language',
				de: 'Sprache der App',
				uk: 'Мова застосунку'
			}) }}
		</template>

		<GenericDropdown :options="languageOptions" v-model="state.lang" />

	</DashboardMenuSection>

	<DashboardMenuSection>

		<template v-slot:title>
			{{ intl(state.lang, {
				en: 'Play mode',
				de: 'Spielmodus',
				uk: 'Режим гри'
			}) }}
		</template>

		<DashboardSettingsGroup>

			<GenericToggle v-model="state.showNavigation">
				{{ intl(state.lang, {
					en: 'Show navigation buttons in Play Mode',
					de: 'Navigationsschaltflächen im Spielmodus anzeigen',
					uk: 'Показувати кнопки навігації в режимі гри'
				}) }}
			</GenericToggle>
		
			<GenericToggle v-model="state.disableCardRotation">
				{{ intl(state.lang, {
					en: 'Disable card rotation',
					de: 'Kartendrehung deaktivieren',
					uk: 'Не нахиляти картки'
				}) }}
			</GenericToggle>

		</DashboardSettingsGroup>

	</DashboardMenuSection>

	<DashboardMenuSection v-if="appCanShareData()">

		<template v-slot:title>
			{{ intl(state.lang, {
				en: 'Sharing',
				de: 'Teilen',
				uk: 'Поширення'
			}) }}
		</template>

		<DashboardSettingsGroup>

			<GenericToggle v-model="state.shareLinksOnly">
				{{ intl(state.lang, {
					en: 'Exclude link titles',
					de: 'Link-Titel ausschließen',
					uk: 'Не підписувати посилання'
				}) }}
			</GenericToggle>
		
		</DashboardSettingsGroup>

	</DashboardMenuSection>

	<DashboardMenuSection>

		<template v-slot:title>
			{{ intl(state.lang, {
				en: 'Account',
				de: 'Konto',
				uk: 'Акаунт'
			}) }}
		</template>

		<DashboardAccountSettings />

	</DashboardMenuSection>

	<DashboardMenuSection>

		<template v-slot:title>
			{{ intl(state.lang, {
				en: 'App info',
				de: 'Angaben der App',
				uk: 'Інформація про застосунок'
			}) }}
		</template>

		<DashboardAppInfo />

	</DashboardMenuSection>

</template>
