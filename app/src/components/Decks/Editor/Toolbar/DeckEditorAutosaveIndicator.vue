<script setup lang="ts">
const props = defineProps<{
	changed?: boolean;
	changesSaved?: boolean;
}>();
</script>

<template>
	<div class="autosave-state" :class="{ changed, saved: changesSaved }">
		<template v-if="changesSaved">
			Changes saved locally
		</template>
		<template v-else-if="changed">
			Saving changes...
		</template>
		<template v-else>
			No local changes
		</template>
	</div>
</template>

<style lang="scss" scoped>
	.autosave-state {
		display: flex;
		flex-flow: row nowrap;
		gap: 0.25rem;
		align-items: center;
		flex-shrink: 0;
		font-size: 0.55rem;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;

		&::before {
			content: "";
			display: block;
			width: 0.65rem;
			height: 0.65rem;
			background-color: white;
			mask-type: alpha;
			mask-size: contain;
			mask-repeat: no-repeat;
			mask-position: center;
			mask-image: url(/src/assets/icons/broom-mask.svg);
		}

		&.changed::before {
			mask-image: url(/src/assets/icons/refresh-mask.svg);
		}

		&.saved::before {
			mask-image: url(/src/assets/icons/check-mask.svg);
		}
	}
</style>
