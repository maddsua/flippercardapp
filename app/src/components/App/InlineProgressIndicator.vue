<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
	title: string;
	total: number;
	done: number;
}>();

const barWidth = computed(() => Math.min(props.done / (props.total || 1), 1) * 100);

</script>

<template>
	<div class="inline-progress-indicator">
		<div class="title">
			{{ title }}
		</div>
		<div class="progress-bar">
			<div class="fill" :style="{ width: `${barWidth}%` }"></div>
		</div>
	</div>
</template>

<style lang="scss" scoped>
	.inline-progress-indicator {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;

		.title {
			font-size: 0.65rem;
		}

		.progress-bar {
			display: block;
			position: relative;
			width: 100%;
			background-color: var(--app-theme-powder-trail);
			border-radius: 4px;
			overflow: hidden;

			.fill {
				height: 2px;
				background-color: var(--app-theme-spooky-orange);
				transition: all 100ms ease;
			}
		}
	}
</style>
