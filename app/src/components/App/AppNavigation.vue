<script setup lang="ts">
import { useRoute } from 'vue-router';
import AppViewButton from './AppViewButton.vue';
import { computed } from 'vue';

const route = useRoute();

const appView = computed((): string | null => typeof route.meta.app_view === 'string' ? route.meta.app_view : null);

</script>

<template>
	<div class="app-navigation">
		<AppViewButton href="/app" :active="appView === 'home'" icon="home" />
		<AppViewButton href="/app/starred" :active="appView === 'starred'" icon="star" />
		<AppViewButton href="/app/discover" :active="appView === 'discover'" icon="search" />
		<AppViewButton href="/app/dashboard" :active="appView === 'menu'" icon="menu" />
	</div>
</template>

<style lang="scss" scoped>
	.app-navigation {
		display: flex;
		flex-flow: row nowrap;
		align-items: center;
		justify-content: space-between;
		padding: 1rem;
		width: 100%;
		flex-shrink: 0;
		border-radius: 1rem;
		background-color: var(--app-theme-ghostly-glow);
	}
</style>
