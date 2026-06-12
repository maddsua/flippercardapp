<script setup lang="ts">
import EditableNodeHarness from './EditableNodeHarness.vue';

const model = defineModel<string>();

</script>

<template>
	<EditableNodeHarness>
		<template v-slot:title>
			Title
		</template>

		<template v-slot:content>
			<input type="text" placeholder="Title" v-model="model" />
		</template>
	</EditableNodeHarness>
</template>

<style lang="scss" scoped>
	input {
		display: block;
		width: 15rem;
		padding: 0.25rem;
		max-width: 100%;
		background-color: unset;
		border: none;
		outline: none;
		color: var(--app-theme-snow-white);
		font-weight: 600;
		font-size: 0.85rem;

		&::placeholder {
			color: var(--app-theme-mysterious-white);
			font-weight: 400;
		}
	}
</style>
