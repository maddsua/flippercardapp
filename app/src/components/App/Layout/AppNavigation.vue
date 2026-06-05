<script setup lang="ts">
import { useRoute } from 'vue-router';
import AppViewButton from './AppViewButton.vue';
import { computed } from 'vue';

const route = useRoute();

const appView = computed((): string | null => typeof route.meta.app_view === 'string' ? route.meta.app_view : null);

</script>

<template>
	<div class="app-navigation">
		<AppViewButton href="/" :active="appView === 'home'" icon="home" />
		<AppViewButton href="/starred" :active="appView === 'starred'" icon="star" />
		<AppViewButton href="/collections/discover" :active="appView === 'discover'" icon="search" />
		<AppViewButton href="/dashboard" :active="appView === 'menu'" icon="menu" />
	</div>
</template>

<style lang="scss" scoped>
	.app-navigation {
		display: flex;
		flex-flow: row nowrap;
		align-items: center;
		justify-content: space-between;
		padding: 0.5rem 1rem;
		width: 100%;
		flex-shrink: 0;
		border-radius: 1rem;
		background-color: var(--app-theme-ghostly-glow);
	}
</style>
