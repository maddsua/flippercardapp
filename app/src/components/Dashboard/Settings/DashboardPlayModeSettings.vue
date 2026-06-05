<script setup lang="ts">
import GenericToggle from '@/components/App/Inputs/GenericToggle.vue';
import { useStorage } from '@/storage/storage';
import { reactive, watch } from 'vue';

const store = useStorage();

const state = reactive({
	showNavigation: store.preferences.playModeShowNavigation.load(),
});

watch(() => state.showNavigation, async (value) => {
	store.preferences.playModeShowNavigation.store(value);
});

</script>

<template>
	<div class="preferences-group">
		<GenericToggle v-model="state.showNavigation" label="Show navigation buttons in Play Mode" />
	</div>
</template>

<style lang="scss" scoped>
	.preferences-group {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}
</style>