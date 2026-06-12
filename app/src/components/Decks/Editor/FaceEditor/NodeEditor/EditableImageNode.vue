<script setup lang="ts">
import { computed, onMounted, reactive, watch } from 'vue';
import { useClient } from '@/api';
import type { ImageMetadata } from '@/api_models';
import { pickLocalFiles } from '@/files';
import EditableNodeHarness from './EditableNodeHarness.vue';
import GenericButton from '@/components/App/Inputs/GenericButton.vue';

const model = defineModel<string | null | undefined>();

const state = reactive({
	data: null as ImageMetadata | null,
	name: null as string | null,
	error: null as string | null,
});

const uploading = computed((): boolean => !!(state.name && !state.data && !state.error));

const client = useClient();

const uploadFile = async () => {

	state.error = null;
	state.data = null;
	state.name = null;

	const files = await pickLocalFiles({ accept: ['.jpg', '.jpeg', '.png', '.gif', '.webp']});
	if (!files?.length) {
		return;
	}

	const uploadFile = files[0];

	state.name = uploadFile.name;

	const { data, error } = await client.images.upload(uploadFile);
	if (!data || error) {
		state.error = error?.message || 'Upload failed';
		return;
	}

	state.name = data.source_name;
	state.data = data;
	model.value = data.id;
};

const fetchMetadata = async () => {

	if (!model.value || model.value === state.data?.id) {
		return
	}

	state.data = null;
	state.error = null;
	state.name = 'Remote image';

	const { data, error } = await client.images.metadata(model.value);
	if (!data || error) {
		state.error = error?.message || 'Unable to load metadata';
		return;
	}

	state.name = data.source_name;
	state.data = data;
};

onMounted(fetchMetadata);
watch(model, fetchMetadata);

</script>

<template>
	<EditableNodeHarness>
		<template v-slot:title>
			Image
		</template>

		<template v-slot:content>
			<div class="input">

				<template v-if="!state.data">
					<GenericButton variant="thin" :disabled="uploading" :spinner="uploading" @click="uploadFile">
						Pick a file
					</GenericButton>
				</template>

				<template v-else>
					<GenericButton variant="thin" theme="orange" :disabled="uploading" :spinner="uploading" @click="uploadFile">
						Change file
					</GenericButton>
				</template>

				<div class="status" :class="{ error: !!state.error }">

					<template v-if="state.error">
						{{ state.error }}
					</template>

					<template v-else-if="uploading">
						Uploading {{ state.name }}
					</template>

					<template v-else-if="state.data">
						{{ state.name }}
					</template>

					<template v-else>
						No files selected
					</template>
					
				</div>
			</div>
		</template>
	</EditableNodeHarness>
</template>

<style lang="scss" scoped>
	.input {
		display: flex;
		flex-flow: row nowrap;
		align-items: center;
		gap: 1rem;
		user-select: none;
		max-width: 100%;

		.status {
			color: var(--app-theme-mysterious-white);
			font-size: 0.75rem;
			min-width: 0;
			overflow: hidden;
			text-overflow: ellipsis;

			&.error {
				color: var(--app-theme-blood-red);
			}
		}
	}
</style>
