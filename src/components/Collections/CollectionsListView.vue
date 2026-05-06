<script setup lang="ts">
import { onMounted, reactive } from 'vue';
import type { CardCollection } from '../../content';
import { useCollectionProvider } from '../../content.loaders';
import { useRouter } from 'vue-router';
import CollectionList from './CollectionList.vue';
import CollectionListEntry from './CollectionListEntry.vue';
import FullscreenMessage from '../App/FullscreenMessage.vue';
import ErrorMessage from '../App/ErrorMessage.vue';
import CollectionHeader from './CollectionHeader.vue';
import CollectionContainer from './CollectionContainer.vue';
import CollectionBreak from './CollectionBreak.vue';
import CollectionEndlistAction from './CollectionEndlistAction.vue';

const router = useRouter();

const state = reactive({
	collections: null as CardCollection[] | null,
	error: null as string | null,
});

onMounted(async () => {

	const { data, error } = await useCollectionProvider().collections();
	if (!data || error) {
		state.error = error?.message || 'Unable to load collections';
		return;
	}

	state.collections = data;
});

const openCollection = (id: string) => {
	router.push(`/app/collection/${id}`);
};

const openExplore = () => {
	router.push('/app/explore');
};

</script>

<template>
	<CollectionContainer>

		<CollectionHeader>

			<template v-slot:title>
				Collections
			</template>

			<template v-slot:summary>
				Your card collections
			</template>

		</CollectionHeader>

		<CollectionList v-if="state.collections?.length">
			<CollectionListEntry v-for="item of state.collections" :title="item.name" @click="openCollection(item.id)" />
		</CollectionList>

		<FullscreenMessage v-else>

			<ErrorMessage v-if="state.error">

				<template v-slot:message>
					Unable to display collections
				</template>
				
				<template v-slot:details>
					{{ state.error }}
				</template>

			</ErrorMessage>

			<p v-else>
				You haven't added any collections yet!
			</p>
			
		</FullscreenMessage>

		<CollectionBreak />

		<CollectionEndlistAction @action="openExplore">
			Explore cards
		</CollectionEndlistAction>

	</CollectionContainer>
</template>
