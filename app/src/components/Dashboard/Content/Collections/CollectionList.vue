<script setup lang="ts">
const emits = defineEmits<{
	(e: 'addCollection'): void;
	(e: 'loadMore'): void;
}>();
</script>

<template>
	<div class="collection-list">

		<div v-if="$slots.actions_before" class="actions">
			<slot name="actions_before" />
		</div>

		<div class="entries">
			<slot>
				<div class="no-entries-message">
					No collections
				</div>
			</slot>
		</div>

		<div v-if="$slots.actions_after" class="actions">
			<slot name="actions_after" />
		</div>

	</div>
</template>

<style lang="scss" scoped>

	.collection-list {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;

		.actions {
			display: flex;
			flex-flow: row nowrap;
			justify-content: center;
			align-items: center;
		}

		.entries {
			display: flex;
			flex-direction: column;
			gap: 1rem;

			.no-entries-message {
				display: flex;
				flex-flow: row nowrap;
				justify-content: center;
				font-size: 0.75rem;
			}
		}
	}

</style>
