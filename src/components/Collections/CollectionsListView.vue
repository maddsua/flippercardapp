<script setup lang="ts">
import { onMounted, reactive } from 'vue';
import type { CardCollection } from '../Cards/content';
import { sampleProvider } from '../../data/sample';


const state = reactive({
	collections: null as CardCollection[] | null,
	error: null as string | null,
});

onMounted(async () => {

	const { data, error } = await sampleProvider.collections();
	if (!data || error) {
		state.error = error?.message || 'Unable to load collections';
		return;
	}

	state.collections = data;
});

//	todo: make look nice

</script>

<template>
	<div class="collections-list-view">

		<p>
			Collections
		</p>

		<ul v-if="state.collections?.length">
			<li v-for="item of state.collections">
				<RouterLink :to="`/collection/${item.id}`">
					{{ item.name }}
				</RouterLink>
			</li>
		</ul>
		<div v-else class="message">
			No collections available
		</div>
	</div>
</template>

<style lang="scss" scoped>

</style>
