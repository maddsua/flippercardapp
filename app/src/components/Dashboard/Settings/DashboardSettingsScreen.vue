<script setup lang="ts">
import AppUiHeader from '@/components/App/Layout/AppUiHeader.vue';
import DashboardMenuSection from '../DashboardMenuSection.vue';
import DashboardLanguageSettings from './DashboardLanguageSettings.vue';
import DashboardPlayModeSettings from './DashboardPlayModeSettings.vue';
import DashboardAccountSettings from './DashboardAccountSettings.vue';

</script>

<template>

	<AppUiHeader>
		<template v-slot:title>
			Dashboard
		</template>
	</AppUiHeader>

	<DashboardMenuSection>

		<template v-slot:title>
			App language
		</template>

		<template v-slot:content>
			<DashboardLanguageSettings />
		</template>

	</DashboardMenuSection>

	<DashboardMenuSection>

		<template v-slot:title>
			Play mode
		</template>

		<template v-slot:content>
			<DashboardPlayModeSettings />
		</template>

	</DashboardMenuSection>

	<DashboardMenuSection>

		<template v-slot:title>
			Account
		</template>

		<template v-slot:content>
			<DashboardAccountSettings />
		</template>

	</DashboardMenuSection>

	<DashboardMenuSection>

		<template v-slot:title>
			App info
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
