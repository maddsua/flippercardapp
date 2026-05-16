<script setup lang="ts">
import { ref } from 'vue';
import { intl, useLanguage } from '../../intl';

import AppUI from '../App/AppUI.vue';
import AppUiHeader from '../App/AppUiHeader.vue';
import RecommendedSection from './RecommendedSection.vue';
import SearchSection from './SearchSection.vue';

const lang = useLanguage();

const searchActive = ref(false);

</script>

<template>
	<AppUI>
		<AppUiHeader>
			<template v-slot:title>
				{{ intl(lang, {
					en: 'Find collections',
					uk: 'Пошук карток',
					de: 'Nach Kärtchen suchen'
				}) }}
			</template>
		</AppUiHeader>

		<SearchSection @active="state => searchActive = state" />
		
		<KeepAlive>
			<RecommendedSection v-if="!searchActive" />
		</KeepAlive>

	</AppUI>
</template>
