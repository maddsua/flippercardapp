<script setup lang="ts">
const props = defineProps<{
	messages: string[];
}>();
</script>

<template>
	<div class="warn-list">
		<div v-for="item of messages" class="message">
			<span>
				{{ item }}
			</span>
		</div>
	</div>
</template>

<style lang="scss" scoped>
	.warn-list {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		max-width: 20rem;

		.message {
			display: flex;
			flex-flow: row nowrap;
			gap: 0.5rem;
			align-items: center;
			font-size: 0.65rem;
			color: var(--app-theme-snow-white);
			min-width: 0;
			overflow: hidden;
			text-overflow: ellipsis;
			white-space: nowrap;

			&::before {
				content: "WARN:";
				color: orange;
				font-weight: bold;
			}
		}
	}
</style>
