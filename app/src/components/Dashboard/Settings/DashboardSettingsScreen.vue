<script setup lang="ts">
import AppUiHeader from '@/components/App/Layout/AppUiHeader.vue';
import DashboardMenuSection from '../DashboardMenuSection.vue';
import DashboardLanguageSettings from './DashboardLanguageSettings.vue';
import DashboardPlayModeSettings from './DashboardPlayModeSettings.vue';
import DashboardAccountSettings from './DashboardAccountSettings.vue';
import { useLanguage, intl } from '@/intl';

const lang = useLanguage();

</script>

<template>

	<AppUiHeader backHref="/">
		<template v-slot:title>
			{{ intl(lang, {
				en: 'Dashboard',
				de: 'Einstellungen',
				uk: 'Налаштування'
			}) }}
		</template>
	</AppUiHeader>

	<DashboardMenuSection>

		<template v-slot:title>
			{{ intl(lang, {
				en: 'App language',
				de: 'Sprache der App',
				uk: 'Мова застосунку'
			}) }}
		</template>

		<template v-slot:content>
			<DashboardLanguageSettings />
		</template>

	</DashboardMenuSection>

	<DashboardMenuSection>

		<template v-slot:title>
			{{ intl(lang, {
				en: 'Play mode',
				de: 'Spielmodus',
				uk: 'Режим гри'
			}) }}
		</template>

		<template v-slot:content>
			<DashboardPlayModeSettings />
		</template>

	</DashboardMenuSection>

	<DashboardMenuSection>

		<template v-slot:title>
			{{ intl(lang, {
				en: 'Account',
				de: 'Konto',
				uk: 'Акаунт'
			}) }}
		</template>

		<template v-slot:content>
			<DashboardAccountSettings />
		</template>

	</DashboardMenuSection>

	<DashboardMenuSection>

		<template v-slot:title>
			{{ intl(lang, {
				en: 'App info',
				de: 'Angaben der App',
				uk: 'Інформація про застосунок'
			}) }}
		</template>

		<template v-slot:content>
			<p class="service-notice">
				Designed by maddsua. Provided by MWS via Railway Corp.
			</p>
		</template>

	</DashboardMenuSection>

</template>

<style lang="scss" scoped>
	.service-notice {
		font-size: 0.75rem;
	}
</style>
