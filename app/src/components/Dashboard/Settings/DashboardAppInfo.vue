<script setup lang="ts">
import { getAppInfo, type AppSource } from '@/app';
import { computed } from 'vue';

const info = computed(() => {

	const info = getAppInfo();

	return {
		version: info.version || 'unknown',
		buildTime: info.buildTime?.toISOString() || 'unknown time',
		distribution: info.mode,
		platform: info.platform || 'unknown platform',
		sourceURL: formatRepoUrl(info.source),
	};
});

const formatRepoUrl = (src: AppSource) => {

	if (!src.vcs || !src.repo) {
		return null;
	}

	switch (src.vcs.toLowerCase().trim()) {
		case 'github':
			return { href: `https://github.com/${src.repo.trim()}`, title: 'GitHub' };
		default:
			return null;
	}
};

</script>

<template>
	<div class="app-info">
		<p>
			Version: {{ info.version }} built at {{ info.buildTime }}; {{ info.platform }}; distribution: {{ info.distribution }}
		</p>
		<p>
			Designed by maddsua. Provided by MWS via Railway Corp.
		</p>
		<p v-if="info.sourceURL">
			Source code at <a :href="info.sourceURL.href" target="_blank">{{ info.sourceURL.title }}</a>
		</p>
	</div>
</template>

<style lang="scss" scoped>
	.app-info {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
		font-size: 0.75rem;

		p {
			margin: 0;
			padding: 0;
		}

		a {
			color: var(--app-theme-deep-lavender);
		}
	}
</style>
