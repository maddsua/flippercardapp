<script setup lang="ts">
import { onMounted, reactive } from 'vue';
import { useClient, type Pagination } from '../../../api';
import type { CardDeckMetadata, CollectionMetadata } from '../../../api_models';
import { useRouter } from 'vue-router';
import { genericPageState, pageControls } from '../../../dataloader';
import AppUiHeader from '../../App/AppUiHeader.vue';
import LoadingMessage from '../../App/LoadingMessage.vue';
import ErrorMessage from '../../App/ErrorMessage.vue';
import GenericButton from '../../App/GenericButton.vue';
import CollectionEntry from './Collections/CollectionEntry.vue';

const client = useClient();
const router = useRouter();

interface LazyLoadedCollection extends CollectionMetadata {
	decks?: CardDeckMetadata[] | null;
};

const state = reactive(genericPageState<LazyLoadedCollection>());

const { more: loadMore } = pageControls(state, async (pagination: Pagination) => {
	return client.collections.list(pagination);
});

onMounted(loadMore);

const createNewCollection = () => {
	router.push('/app/dashboard/content/collections/new')
};

const manageCollection = (id: string) => {
	router.push(`/app/dashboard/content/collection/${id}`);
};

</script>

<template>

	<AppUiHeader backHref="/app/dashboard">
		<template v-slot:title>
			Content dashboard
		</template>
		<template v-slot:summary>
			Manage cards and collections
		</template>
	</AppUiHeader>

	<LoadingMessage v-if="!state.ready">
		Loading collections...
	</LoadingMessage>

	<ErrorMessage v-else-if="state.error">
		<template v-slot:message>
			Unable to load collections
		</template>
		<template v-slot:details>
			{{ state.error }}
		</template>
	</ErrorMessage>

	<div v-else class="collection-list">

		<div class="actions">
			<GenericButton theme="blue" variant="thin" @click="createNewCollection">
				+ Add collection
			</GenericButton>
		</div>

		<div class="entries">

			<template v-if="state.entries.length">
				<CollectionEntry v-for="entry of state.entries" :key="entry.id" :entry="entry" @manage="manageCollection(entry.id)" />
			</template>

			<template v-else>
				<div class="no-entries-message">
					No collections
				</div>
			</template>

		</div>

		<div v-if="state.has_next" class="actions">
			<GenericButton theme="green" variant="thin" @click="loadMore">
				Load more
			</GenericButton>
		</div>

	</div>

</template>

<style lang="scss" scoped>
	.collection-list {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;

		.actions {
			display: flex;
			flex-flow: row nowrap;
			justify-content: center;
			align-items: center;
		}

		.entries {
			display: flex;
			flex-direction: column;
			gap: 1rem;

			.no-entries-message {
				display: flex;
				flex-flow: row nowrap;
				justify-content: center;
				font-size: 0.75rem;
			}
		}
	}
</style>
