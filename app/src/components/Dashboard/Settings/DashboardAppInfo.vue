<script setup lang="ts">
import appManifest from '@/../public/manifest.json';
import { computed } from 'vue';

const appBuildTime = computed(() => {

	const buildTs = import.meta.env.VITE_APP_BUILD_TS;
	if (buildTs) {
		try {
			return new Date(buildTs).toISOString();
		} catch (_) {}
	}

	return 'unknown time';
});

const appVersion = import.meta.env.VITE_APP_VERSION || 'unknown';
const appPlatform = import.meta.env.VITE_APP_PLATFORM || 'unknown platform';

const appDistribution = computed(() => {

	const queries = [appManifest.display, ...appManifest.display_override];

	for (const query of queries) {
		try {
			if (window.matchMedia(`(display-mode: ${query})`).matches) {
				return 'PWA';
			}
		} catch (_) {
			return 'Web-Limited';
		}
	}

	return 'Web';
});

</script>

<template>
	<div class="app-info">
		<p>
			Version: {{ appVersion }} built at {{ appBuildTime }}; {{ appPlatform }}; distribution: {{ appDistribution }}
		</p>
		<p>
			Designed by maddsua. Provided by MWS via Railway Corp.
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
	}
</style>
