<script setup lang="ts">
import { getAppInfo } from '@/app';
import { computed } from 'vue';

const info = computed(() => {

	const info = getAppInfo();

	return {
		version: info.version || 'unknown',
		buildTime: info.buildTime?.toISOString() || 'unknown time',
		distribution: info.mode,
		platform: info.platform || 'unknown platform',
	};
});

</script>

<template>
	<div class="app-info">
		<p>
			Version: {{ info.version }} built at {{ info.buildTime }}; {{ info.platform }}; distribution: {{ info.distribution }}
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
